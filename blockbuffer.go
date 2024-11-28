package main

import (
	"log"
	"os"

	// import internal package
	i "blockbuffer/interal"
)

func main() {

	// Ensure output directory exists
	if _, err := os.Stat(i.OutputDir); os.IsNotExist(err) {
		err := os.MkdirAll(i.OutputDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// Scan input directory and queue files for conversion
	go i.ScanAndQueueFiles(*i.WatchDir, i.OutputDir)

	// Start watching the directory
	go i.WatchDirectory(*i.WatchDir, i.OutputDir)

	// Check the queue and process files
	go i.ProcessQueue()

	// Start the server
	i.StartServer()
}
