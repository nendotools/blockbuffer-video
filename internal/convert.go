// This file handles the conversion of files to DNxHR format and
// communicates the progress of the conversion to the user.
package internal

import (
	"encoding/json"
	"fmt"
	"log"
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

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var Cmd *exec.Cmd

// map file ID to Conversion command
type Conversion struct {
	inFile  string
	outFile string
	cmd     *exec.Cmd
}

func PollFile(inputFile string) float64 {
	a, err := ffmpeg.Probe(inputFile)
	CheckError(err)
	err = json.Unmarshal([]byte(a), &InputProbeData)
	CheckError(err)
	totalDuration, err := ProbeDuration(InputProbeData)
	CheckError(err)
	return totalDuration
}

var ConversionMap = make(map[string]Conversion)

func Ternary(condition bool, a string, b string) string {
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
			fmt.Printf("Error stating file %s: %v\n", filePath, err)
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
		for file := range fileQueue {
			// if fileis in skip list, skip it
			if !skipList[file.ID] {
				go convertToDNxHR(file, *OutputDir)
			} else {
				fmt.Printf("Skipping file: %s\n", file.FilePath)
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
func convertToDNxHR(inputFile File, outputDir string) {
	for !*AutoConvert {
		time.Sleep(2 * time.Second)
	}

	conv <- 1
	if !waitForFileReady(inputFile.FilePath) {
		fmt.Printf("File %s is not ready to be processed\n", inputFile.ID)
		fileQueue <- inputFile
		<-conv
		return
	}

	// Prepare paths
	outputFile := strings.TrimSuffix(filepath.Base(inputFile.FilePath), filepath.Ext(inputFile.FilePath)) + "_dnxhr.mov"
	outputPath := filepath.Join(outputDir, outputFile)

	// probe input file
	if len(InputProbeData.Streams) == 0 {
		a, err := ffmpeg.Probe(inputFile.FilePath)
		CheckError(err)
		err = json.Unmarshal([]byte(a), &InputProbeData)
		CheckError(err)
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
	if *IgnoreExisting {
		ffmpegArgs["y"] = nil
	}

	// set resolution
	if inputWidth > inputHeight && inputHeight > 1080 {
		ffmpegArgs["vf"] = "scale=-2:1080"
	} else if inputHeight >= inputWidth && inputWidth > 1080 {
		ffmpegArgs["vf"] = "scale=1080:-2"
	}

	convertWithProgress(inputFile.ID, inputFile.FilePath, outputPath, ffmpegArgs)
	updateProgress(inputFile.ID, 100, true)
	if *DeleteAfter {
		_, err := os.Stat(inputFile.FilePath)
		if err == nil {
			os.Remove(inputFile.FilePath)
		}
	}

	fmt.Printf("Successfully converted: %s -> %s\n", inputFile.FilePath, outputPath)
	<-conv
}

// convertWithProgress uses the ffmpeg `-progress` option with a unix-domain socket to report progress
func convertWithProgress(fileId string, inFileName string, outFileName string, ffmpegArgs ffmpeg.KwArgs) {
	var err error

	// get duration of video (3 seconds if preview mode)
	totalDuration, err := ProbeDuration(InputProbeData)
	CheckError(err)

	fmt.Println("preparing conversion")

	fmt.Printf("Processing file: %s\n", inFileName)
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
		fmt.Printf("Error converting file: %v\n", err)
	}

	defer CancelConversion(fileId)
	fmt.Printf("Successfully queued file: %s -> %s\n", inFileName, outFileName)
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
			log.Fatal("accept error:", err)
		}
		buf := make([]byte, 16)
		data := ""
		fmt.Print("\033[s") // save the cursor position
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
				fmt.Printf("\033[2K\rProgress: %.2f%%", cp*100)
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
	file, ok := fileList[fileId]
	if !ok {
		return
	}
	var completionStatus = Ternary(*DeleteAfter, string(CompleteDeleted), string(Completed))

	FileListMutex.Lock()
	fileList[fileId] = File{
		ID:       fileId,
		FilePath: file.FilePath,
		Status:   FileStatus(Ternary(progress == 100, completionStatus, Ternary(progress == -1, string(Failed), string(Processing)))),
		Duration: file.Duration,
		Progress: progress,
	}
	FileListMutex.Unlock()

	BroadcastMessage(Message{
		MessageType: UpdateFile,
		MustSend:    mustSend,
		Data: map[string]File{
			fileId: fileList[fileId],
		},
	})
}

func CancelConversion(fileId string) {
	if conv, ok := ConversionMap[fileId]; ok {
		fmt.Println("cmd found for: ", fileId)
		if err := conv.cmd.Process.Signal(os.Interrupt); err == nil {
			if _, err := os.Stat(conv.outFile); !os.IsNotExist(err) {
				fmt.Println("Removing incomplete file: ", conv.outFile)
				if err := os.Remove(conv.outFile); err != nil {
					fmt.Println("Error deleting file: ", err)
				}
			}
		}
		delete(ConversionMap, fileId)
	}
}
