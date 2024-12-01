// This file handles scanning of directories and files. It also
// handles managing the conversion queue.
package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/u2takey/go-utils/uuid"
)

type File struct {
	ID       string `json:"id"`
	FilePath string `json:"filePath"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

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
			outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_dnxhr.mov"
			outputPath := filepath.Join(outputDir, outputFile)
			file := File{
				ID:       uuid.NewUUID(),
				FilePath: inputDir + "/" + inputFile,
				Status:   "queued",
				Progress: "0",
			}
			fileList[file.ID] = file
			fmt.Println("file: ", file)

			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				fmt.Printf("Queueing file for conversion: %s\n", inputFile)
				fileQueue <- file
			} else {
				fmt.Printf("Output file already exists: %s\n", outputFile)
				file.Status = "done"
				file.Progress = "100"
				fileList[file.ID] = file
			}
		}
	}
}

// WatchDirectory watches the directory for new video files and triggers conversion to DNxHR
func WatchDirectory(inputDir, outputDir string) {
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
			if event.Op&fsnotify.Create == fsnotify.Create {
				// When a new file is created, process it if it's a video file
				if isVideoFile(event.Name) {
					fmt.Printf("Detected new video: %s\n", event.Name)
					file := File{
						ID:       uuid.NewUUID(),
						FilePath: event.Name,
						Status:   "queued",
						Progress: "0",
					}
					fileQueue <- file
					fileList[file.ID] = file
					BroadcastFiles(fileList)

					// remove file from skip list
					delete(skipList, file.ID)
				}
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf("Detected removed video: %s\n", event.Name)
				// add file to skip list
				// search for the file path in the fileList.FilePath
				for _, file := range fileList {
					if file.FilePath == event.Name {
						skipList[file.ID] = true
						break
					}
				}
			}
		case err := <-watcher.Errors:
			fmt.Println("Error:", err)
		}
	}
}
