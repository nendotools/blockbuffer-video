package main

import (
	"log"
	"os"

	// import internal package
	i "blockbuffer/internal"
)

func main() {
	// Ensure output directory exists
	if _, err := os.Stat(*i.OutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(*i.OutputDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Ensure upload directory exists
	if _, err := os.Stat(*i.WatchDir); os.IsNotExist(err) {
		err := os.MkdirAll(*i.WatchDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}

	if _, err := os.Stat(*i.UploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(*i.UploadDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
	}

	// Scan input directory and queue files for conversion
	go i.ScanAndQueueFiles(*i.WatchDir, *i.OutputDir)

	// Start watching the directory
	go i.WatchDirectory(*i.WatchDir, *i.OutputDir)

	// Check the queue and process files
	go i.ProcessQueue()

	// Start the server
	i.StartServer()
}
