// This file handles scanning of directories and files. It also
// handles managing the conversion queue.
package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/u2takey/go-utils/uuid"

	api "blockbuffer/internal/api"
	"blockbuffer/internal/io"
	store "blockbuffer/internal/store"
	types "blockbuffer/internal/types"
)

var skipList = make(map[string]bool)

// isVideoFile checks if a file is a supported video format (case-insensitive)
func isVideoFile(filePath string) bool {
	// Convert file extension to lowercase for case-insensitive comparison
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv"
}

func ScanAndQueueFiles(inputDir string, outputDir string) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		io.Logf("Error reading directory: %v", io.Error, err)
	}

	for _, file := range files {
		if !file.IsDir() && isVideoFile(file.Name()) {
			inputFile := file.Name()
			var filePath = inputDir + "/" + inputFile
			var totalDuration = PollFile(filePath)

			outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + "_dnxhr.mov"
			outputPath := filepath.Join(outputDir, outputFile)
			file := types.File{
				ID:       uuid.NewUUID(),
				FilePath: filePath,
				Duration: totalDuration,
				Status:   types.Queued,
				Progress: 0,
			}
			store.UpdateFile(file)

			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				io.Logf("Queueing file for conversion: %s", io.Info, inputFile)
				store.FileQueue <- file
			} else {
				io.Logf("Output file already exists: %s", io.Info, outputFile)
				file.Status = types.Completed
				file.Progress = 100
				store.UpdateFile(file)
			}
		}
	}
}

// WatchDirectory watches the directory for new video files and triggers conversion to DNxHR
func WatchDirectory(inputDir string, outputDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		io.Logf("Error creating watcher: %v", io.Fatal, err)
	}
	defer watcher.Close()

	err = watcher.Add(inputDir)
	if err != nil {
		io.Logf("Error adding directory to watcher: %v", io.Fatal, err)
	}
	io.Logf("Watching directory: %s", io.Info, inputDir)

	// Watch for events in the directory
	for {
		select {
		case event := <-watcher.Events:
			// Create is triggered when a new file is created AND not Rename
			if event.Op.Has(fsnotify.Create) {
				// When a new file is created, process it if it's a video file
				if isVideoFile(event.Name) {
					io.Logf("Detected new video: %s", io.Info, event.Name)
					var totalDuration = PollFile(event.Name)
					file := types.File{
						ID:       uuid.NewUUID(),
						FilePath: event.Name,
						Status:   types.Queued,
						Duration: totalDuration,
						Progress: 0,
					}
					store.FileQueue <- file
					store.UpdateFile(file)
					api.BroadcastMessage(types.Message{
						MessageType: types.CreateFile,
						MustSend:    true,
						Data:        map[string]types.File{file.ID: file},
					})

					// remove file from skip list
					delete(skipList, file.ID)
				}
			}
			if event.Op.Has(fsnotify.Rename) || event.Op.Has(fsnotify.Remove) {
				io.Logf("Detected renamed/removed video: %s", io.Info, event.Name)
				// add file to skip list
				// search for the file path in the fileList.FilePath
				for _, file := range store.FileList {
					if file.FilePath == event.Name {
						io.Logf("canceling conversion: %s", io.Info, file.ID)
						CancelConversion(file.ID)
						skipList[file.ID] = true
						if file.Status == types.CompleteDeleted {
							io.Logf("Skipping UI notification: %s", io.Info, file.ID)
							break
						}
						api.BroadcastMessage(types.Message{
							MessageType: types.DeleteFile,
							MustSend:    true,
							Data:        map[string]types.File{file.ID: file},
						})
						delete(store.FileList, file.ID)
						break
					}
				}
			}
		case err := <-watcher.Errors:
			io.Logf("Error in watcher: %v", io.Error, err)
		}
	}
}
