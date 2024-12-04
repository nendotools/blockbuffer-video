// This file manages a webserver that listens for requests:
//   - files in directory
//   - converted files
//   - progress of conversion
package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Config struct {
	AutoConvert    *bool `json:"autoConvert,omitempty"`
	DeleteAfter    *bool `json:"deleteAfter,omitempty"`
	IgnoreExisting *bool `json:"ignoreExisting,omitempty"`
}

func SuccessJSON(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	jsonData := json.NewEncoder(w)
	var response = map[string]interface{}{"message": message}
	if data != nil {
		response = map[string]interface{}{"message": message, "data": data}
	}
	jsonData.Encode(response)
}

func ErrorJSON(w http.ResponseWriter, message string, code int) {
	// respond with json formatted error message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData := json.NewEncoder(w)
	jsonData.Encode(map[string]interface{}{"error": message, "code": code})
}

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

type filterWriter struct {
	writer io.Writer
}

func (fw *filterWriter) Write(p []byte) (n int, err error) {
	// get string from byte slice
	// filter out nuxt.js logs
	// write to writer
	s := string(p)
	if !strings.Contains(s, "WARN  Deprecation") {
		return fw.writer.Write(p)
	}
	return len(p), nil
}

func startNuxtDev() {
	cmd := exec.Command("yarn", "dev")
	cmd.Dir = "./blockbuffer-fe"

	// listen to stdout and stderr and pass to stdout with prefix "NUXT: "
	cmd.Stdout = &filterWriter{writer: log.Writer()}
	cmd.Stderr = &filterWriter{writer: log.Writer()}

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
	if *Headless {
		return
	}

	if isDevServer() {
		startNuxtDev()
		http.HandleFunc("/", proxyToNuxtDev)
	} else {
		fs := http.FileServer(http.Dir("public"))
		http.Handle("/", fs)
	}

	// Web Scocket Server
	go HandleMessages()
	http.HandleFunc("/ws", HandleSocketConnections)

	http.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(apiHandler)))
	fmt.Printf("Server listening on port %s\n", strconv.Itoa(*Port))
	log.Panic(
		http.ListenAndServe(*ListenAddr+":"+strconv.Itoa(*Port), nil),
	)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// respond to the request with json formatted data
	w.Header().Set("Content-Type", "application/json")

	// handle routing requests with router
	router := http.NewServeMux()
	router.HandleFunc("GET /config", configHandler)
	router.HandleFunc("POST /config", configHandler)
	router.HandleFunc("GET /files", filesHandler)
	router.HandleFunc("POST /upload", HandleUploadMultipleFiles)
	router.ServeHTTP(w, r)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	// handle GET
	if r.Method == "GET" {
		// respond with json formatted configuration
		w.Header().Set("Content-Type", "application/json")
		jsonData := json.NewEncoder(w)
		jsonData.Encode(map[string]interface{}{
			"autoConvert":    AutoConvert,
			"deleteAfter":    DeleteAfter,
			"ignoreExisting": IgnoreExisting,
		})
	}

	// handle POST
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		// parse request body as Config
		var config Config = Config{}
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		if config.AutoConvert != nil {
			AutoConvert = config.AutoConvert
		}
		if config.DeleteAfter != nil {
			DeleteAfter = config.DeleteAfter
		}
		if config.IgnoreExisting != nil {
			IgnoreExisting = config.IgnoreExisting
		}
		SuccessJSON(w, "success", nil)
	}
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
