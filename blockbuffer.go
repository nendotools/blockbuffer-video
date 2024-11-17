package main

import (
	"fmt"
	"log"
	"net/http"
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
	// call function from interal/scan.go
	go i.ScanAndQueueFiles(i.WatchDir, i.OutputDir)

	// Start watching the directory
	go i.WatchDirectory(i.WatchDir, i.OutputDir)

	// Check the queue and process files
	go i.ProcessQueue()

	fs := http.FileServer(http.Dir("public"))
	http.HandleFunc("/api", apiHandler)
	http.Handle("/", fs)

	fmt.Println("Server listening on port 3000")
	log.Panic(
		http.ListenAndServe(":3000", nil),
	)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// respond to the request with json formatted data
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "hello world"}`))
}
