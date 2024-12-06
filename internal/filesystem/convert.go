// This file handles the conversion of files to DNxHR format and
// communicates the progress of the conversion to the user.
package filesystem

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	api "blockbuffer/internal/api"
	io "blockbuffer/internal/io"
	opts "blockbuffer/internal/settings"
	store "blockbuffer/internal/store"
	types "blockbuffer/internal/types"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// var conv chan int // conv is a channel to limit the number of concurrent conversions
var conv = make(chan int, *opts.MaxConcurrent)
var Cmd *exec.Cmd

// map file ID to Conversion command
type Conversion struct {
	inFile  string
	outFile string
	cmd     *exec.Cmd
}

func PollFile(inputFile string) float64 {
	a, err := ffmpeg.Probe(inputFile)
	io.CheckError(err)
	err = json.Unmarshal([]byte(a), &InputProbeData)
	io.CheckError(err)
	totalDuration, err := ProbeDuration(InputProbeData)
	io.CheckError(err)
	return totalDuration
}

var ConversionMap = make(map[string]Conversion)

func Ternary(condition bool, a any, b any) any {
	if condition {
		return a
	}
	return b
}

func waitForFileReady(filePath string) bool {
	const checkInterval = 5 * time.Second
	const maxChecks = 10

	var lastSize int64 = -1
	for i := 0; i < maxChecks; i++ {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			io.Logf("Error stating file %s: %v", io.Error, filePath, err)
			return false
		}

		currentSize := fileInfo.Size()
		if currentSize == lastSize {
			return true
		}

		lastSize = currentSize
		time.Sleep(checkInterval)
	}

	return false
}

func ProcessQueue() {
	for {
		for file := range store.FileQueue {
			// if fileis in skip list, skip it
			if !skipList[file.ID] {
				go convertToDNxHR(file, *opts.OutputDir)
			} else {
				io.Logf("Skipping file: %s", io.Info, file.FilePath)
				delete(skipList, file.FilePath) // skipped files are removed from the skip list
			}
		}

		// sleep before retrying processQueue
		time.Sleep(15 * time.Second)
	}
}

type probeData struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

var InputProbeData = probeData{}

// convertToDNxHR runs FFmpeg to convert the video to DNxHR
func convertToDNxHR(inputFile types.File, outputDir string) {
	for !*opts.AutoConvert {
		time.Sleep(2 * time.Second)
	}

	conv <- 1
	if !waitForFileReady(inputFile.FilePath) {
		io.Logf("File %s is not ready to be processed", io.Info, inputFile.ID)
		store.FileQueue <- inputFile
		<-conv
		return
	}

	// Prepare paths
	outputFile := strings.TrimSuffix(filepath.Base(inputFile.FilePath), filepath.Ext(inputFile.FilePath)) + "_dnxhr.mov"
	outputPath := filepath.Join(outputDir, outputFile)
	// if file exists and not overwriting, skip conversion
	if _, err := os.Stat(outputPath); err == nil && !*opts.OverwriteExisting {
		io.Logf("Skipping existing file: %s", io.Info, outputFile)
		updateProgress(inputFile.ID, -10, true)
		<-conv
		return
	}

	// probe input file
	if len(InputProbeData.Streams) == 0 {
		a, err := ffmpeg.Probe(inputFile.FilePath)
		io.CheckError(err)
		err = json.Unmarshal([]byte(a), &InputProbeData)
		io.CheckError(err)
	}

	var inputHeight = 0
	var inputWidth = 0

	// find width / height of video stream
	if inputWidth == 0 {
		for _, stream := range InputProbeData.Streams {
			if stream.Width != 0 && stream.Height != 0 {
				inputWidth = stream.Width
				inputHeight = stream.Height
				break
			}
		}
	}

	// create ffmpeg args
	ffmpegArgs := ffmpeg.KwArgs{}
	ffmpegArgs["c:v"] = "dnxhd"
	ffmpegArgs["profile:v"] = "dnxhr_hq"
	ffmpegArgs["pix_fmt"] = "yuv420p"
	ffmpegArgs["c:a"] = "pcm_s16le"

	// set resolution
	if inputWidth > inputHeight && inputHeight > 1080 {
		ffmpegArgs["vf"] = "scale=-2:1080"
	} else if inputHeight >= inputWidth && inputWidth > 1080 {
		ffmpegArgs["vf"] = "scale=1080:-2"
	}

	convertWithProgress(inputFile.ID, inputFile.FilePath, outputPath, ffmpegArgs)
	updateProgress(inputFile.ID, 100, true)
	if *opts.DeleteAfter {
		_, err := os.Stat(inputFile.FilePath)
		if err == nil {
			os.Remove(inputFile.FilePath)
		}
	}

	io.Logf("Successfully converted: %s -> %s", io.Info, inputFile.FilePath, outputPath)
	<-conv
}

// convertWithProgress uses the ffmpeg `-progress` option with a unix-domain socket to report progress
func convertWithProgress(fileId string, inFileName string, outFileName string, ffmpegArgs ffmpeg.KwArgs) {
	var err error

	// get duration of video (3 seconds if preview mode)
	totalDuration, err := ProbeDuration(InputProbeData)
	io.CheckError(err)

	io.Logf("Processing file: %s", io.Info, inFileName)
	Cmd = ffmpeg.Input(inFileName).
		Output(outFileName, ffmpegArgs).
		GlobalArgs("-progress", "unix://"+TempSock(totalDuration, fileId)).
		OverWriteOutput().
		Silent(true).
		Compile()
	ConversionMap[fileId] = Conversion{
		inFile:  inFileName,
		outFile: outFileName,
		cmd:     Cmd,
	}

	err = Cmd.Run()
	if err != nil {
		io.Logf("Error converting file: %s: %v", io.Error, inFileName, err)
	}

	defer CancelConversion(fileId)
	io.Logf("Successfully queued file: %s -> %s", io.Info, inFileName, outFileName)
}

func ProbeDuration(data probeData) (float64, error) {
	f, err := strconv.ParseFloat(data.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

func TempSock(totalDuration float64, fileId string) string {
	// serve
	sockFileName := path.Join(os.TempDir(), fmt.Sprintf("%d_sock", rand.Int()))
	l, err := net.Listen("unix", sockFileName)
	if err != nil {
		updateProgress(fileId, -1, true)
		panic(err)
	}

	go func() {
		re := regexp.MustCompile(`out_time_ms=(\d+)`)
		fd, err := l.Accept()
		if err != nil {
			io.Logf("Error accepting connection: %v", io.Fatal, err)
		}
		buf := make([]byte, 16)
		data := ""
		for {
			_, err := fd.Read(buf)
			if err != nil {
				return
			}
			data += string(buf)
			a := re.FindAllStringSubmatch(data, -1)
			cp := 0.00
			if len(a) > 0 && len(a[len(a)-1]) > 0 {
				c, _ := strconv.Atoi(a[len(a)-1][len(a[len(a)-1])-1])
				cp = float64(c) / totalDuration / 1000000
			}
			if strings.Contains(data, "progress=end") {
				cp = 1.00
				updateProgress(fileId, 100, true)
				break
			}
			if cp > 0.00 && cp <= 1.00 {
				updateProgress(fileId, float32(cp*100), false)
				if cp == 1.00 {
					l.Close()
					break
				}
			}
		}
	}()

	return sockFileName
}

func updateProgress(fileId string, progress float32, mustSend bool) {
	file, ok := store.FileList[fileId]
	if !ok {
		return
	}
	status := file.Status
	switch progress {
	case -10:
		status = types.Rejected
		break
	case -1:
		status = types.Failed
		break
	case 100:
		status = Ternary(*opts.DeleteAfter, types.CompleteDeleted, types.Completed).(types.FileStatus)
		break
	default:
		status = types.Processing
		break
	}

	store.FileListMutex.Lock()
	store.FileList[fileId] = types.File{
		ID:       fileId,
		FilePath: file.FilePath,
		Status:   status,
		Duration: file.Duration,
		Progress: progress,
	}
	store.FileListMutex.Unlock()

	api.BroadcastMessage(types.Message{
		MessageType: types.UpdateFile,
		MustSend:    mustSend,
		Data: map[string]types.File{
			fileId: store.FileList[fileId],
		},
	})
}

func CancelConversion(fileId string) {
	if conv, ok := ConversionMap[fileId]; ok {
		io.Logf("Cancelling conversion: %s", io.Info, fileId)
		if err := conv.cmd.Process.Signal(os.Interrupt); err == nil {
			if _, err := os.Stat(conv.outFile); !os.IsNotExist(err) {
				io.Logf("Removing incomplete file: %s", io.Info, conv.outFile)
				if err := os.Remove(conv.outFile); err != nil {
					io.Logf("Error deleting file: %v", io.Error, err)
				}
			}
		}
		delete(ConversionMap, fileId)
	}
}
