// This file manages a webserver that listens for requests:
//   - files in directory
//   - converted files
//   - progress of conversion
package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(apiHandler)))
	http.Handle("/", fs)

	fmt.Println("Server listening on port 3001")
	log.Panic(
		http.ListenAndServe(":3001", nil),
	)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// respond to the request with json formatted data
	w.Header().Set("Content-Type", "application/json")

	// handle routing requests with router
	router := http.NewServeMux()
	router.HandleFunc("/files", filesHandler)
	router.ServeHTTP(w, r)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	// respond with json formatted list from fileList
	w.Header().Set("Content-Type", "application/json")
	jsonData := json.NewEncoder(w)
	jsonData.Encode(fileList)
}
