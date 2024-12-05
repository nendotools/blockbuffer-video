package main

import (
	"os"

	// import internal package
	api "blockbuffer/internal/api"
	fs "blockbuffer/internal/filesystem"
	"blockbuffer/internal/io"
	opts "blockbuffer/internal/settings"
)

func main() {
	// Ensure output directory exists
	if _, err := os.Stat(*opts.OutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(*opts.OutputDir, 0755)
		if err != nil {
			io.Logf("Failed to create output directory: %v", io.Fatal, err)
		}
	}

	// Ensure upload directory exists
	if _, err := os.Stat(*opts.WatchDir); os.IsNotExist(err) {
		err := os.MkdirAll(*opts.WatchDir, 0755)
		if err != nil {
			io.Logf("Failed to create upload directory: %v", io.Fatal, err)
		}
	}

	if _, err := os.Stat(*opts.UploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(*opts.UploadDir, 0755)
		if err != nil {
			io.Logf("Failed to create upload directory: %v", io.Fatal, err)
		}
	}

	// Scan input directory and queue files for conversion
	go fs.ScanAndQueueFiles(*opts.WatchDir, *opts.OutputDir)

	// Start watching the directory
	go fs.WatchDirectory(*opts.WatchDir, *opts.OutputDir)

	// Check the queue and process files
	go fs.ProcessQueue()

	// Start the server
	api.StartServer()
}
