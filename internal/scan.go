// This file handles scanning of directories and files. It also
// handles managing the conversion queue.
package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/u2takey/go-utils/uuid"
)

type FileStatus string

const (
	New        FileStatus = "new"
	Queued     FileStatus = "queued"
	Processing FileStatus = "processing"
	Completed  FileStatus = "completed"
	Cancelled  FileStatus = "cancelled"
	Failed     FileStatus = "failed"
	Deleted    FileStatus = "deleted"
)

type File struct {
	ID       string     `json:"id"`
	FilePath string     `json:"filePath"`
	Status   FileStatus `json:"status"`
	Progress float32    `json:"progress"`
	Duration float64    `json:"duration"`
}

var FileListMutex = &sync.Mutex{}
var fileList = make(map[string]File)

// isVideoFile checks if a file is a supported video format (case-insensitive)
func isVideoFile(filePath string) bool {
	// Convert file extension to lowercase for case-insensitive comparison
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv"
}

func ScanAndQueueFiles(inputDir string, outputDir string) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && isVideoFile(file.Name()) {
			inputFile := file.Name()
			var filePath = inputDir + "/" + inputFile
			var totalDuration = PollFile(filePath)

			outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_dnxhr.mov"
			outputPath := filepath.Join(outputDir, outputFile)
			file := File{
				ID:       uuid.NewUUID(),
				FilePath: filePath,
				Duration: totalDuration,
				Status:   Queued,
				Progress: 0,
			}
			FileListMutex.Lock()
			fileList[file.ID] = file
			fmt.Println("file: ", file)
			FileListMutex.Unlock()

			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				fmt.Printf("Queueing file for conversion: %s\n", inputFile)
				fileQueue <- file
			} else {
				fmt.Printf("Output file already exists: %s\n", outputFile)
				file.Status = Completed
				file.Progress = 100
				FileListMutex.Lock()
				fileList[file.ID] = file
				FileListMutex.Unlock()
			}
		}
	}
}

// WatchDirectory watches the directory for new video files and triggers conversion to DNxHR
func WatchDirectory(inputDir string, outputDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Watching directory: %s\n", inputDir)

	// Watch for events in the directory
	for {
		select {
		case event := <-watcher.Events:
			// Create is triggered when a new file is created AND not Rename
			if event.Op.Has(fsnotify.Create) {
				// When a new file is created, process it if it's a video file
				if isVideoFile(event.Name) {
					fmt.Printf("Detected new video: %s\n", event.Name)
					var totalDuration = PollFile(event.Name)
					file := File{
						ID:       uuid.NewUUID(),
						FilePath: event.Name,
						Status:   Queued,
						Duration: totalDuration,
						Progress: 0,
					}
					fileQueue <- file
					FileListMutex.Lock()
					fileList[file.ID] = file
					FileListMutex.Unlock()
					BroadcastMessage(Message{
						MessageType: CreateFile,
						Data:        map[string]File{file.ID: file},
					})

					// remove file from skip list
					delete(skipList, file.ID)
				}
			}
			if event.Op.Has(fsnotify.Rename) || event.Op.Has(fsnotify.Remove) {
				fmt.Printf("Detected removed video: %s\n", event.Name)
				// add file to skip list
				// search for the file path in the fileList.FilePath
				for _, file := range fileList {
					fmt.Println(file.FilePath == event.Name, " file: ", file)
					if file.FilePath == event.Name {
						fmt.Println("canceling conversion: ", file.ID)
						CancelConversion(file.ID)
						skipList[file.ID] = true
						BroadcastMessage(Message{
							MessageType: DeleteFile,
							Data:        map[string]File{file.ID: file},
						})
						delete(fileList, file.ID)
						break
					}
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("Error:", err)
		}
	}
}
