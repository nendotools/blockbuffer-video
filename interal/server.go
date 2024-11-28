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
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func isDevServer() bool {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to determine executable path: %v", err)
	}

	// Check if the path indicates a temporary binary (used by `go run`)
	tempDir := os.TempDir()
	return strings.HasPrefix(execPath, tempDir)
}

func startNuxtDev() {
	cmd := exec.Command("yarn", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "./blockbuffer-fe"

	log.Println("Starting Nuxt.js development server...")
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Nuxt.js: %v", err)
	}

	// Wait for the process to exit or handle as a background process
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("Nuxt.js process exited with error: %v", err)
		}
	}()
}

// Proxy requests to Nuxt dev server
func proxyToNuxtDev(w http.ResponseWriter, r *http.Request) {
	target, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Nuxt dev server is not available", http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

func StartServer() {
	if isDevServer() {
		startNuxtDev()
		http.HandleFunc("/", proxyToNuxtDev)
	} else {
		fs := http.FileServer(http.Dir("public"))
		http.Handle("/", fs)
	}

	http.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(apiHandler)))
	fmt.Println("Server listening on port 8080")
	log.Panic(
		http.ListenAndServe(":8080", nil),
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
	fileArray := []File{}
	for _, file := range fileList {
		fileArray = append(fileArray, file)
	}
	jsonData.Encode(fileArray)
}
