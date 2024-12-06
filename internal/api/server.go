// This file manages a webserver that listens for requests:
//   - files in directory
//   - converted files
//   - progress of conversion
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"blockbuffer/internal/io"
	opts "blockbuffer/internal/settings"
	store "blockbuffer/internal/store"
	types "blockbuffer/internal/types"
)

type Config struct {
	AutoConvert       *bool `json:"autoConvert,omitempty"`
	DeleteAfter       *bool `json:"deleteAfter,omitempty"`
	OverwriteExisting *bool `json:"overwriteExisting,omitempty"`
}

func isDevServer() bool {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		io.Logf("Failed to determine executable path: %v", io.Error, err)
	}

	// Check if the path indicates a temporary binary (used by `go run`)
	tempDir := os.TempDir()
	return strings.HasPrefix(execPath, tempDir)
}

func startNuxtDev() {
	cmd := exec.Command("yarn", "dev")
	cmd.Dir = "./blockbuffer-fe"

	// listen to stdout and stderr and pass to stdout with prefix "NUXT: "
	cmd.Stdout = types.Writer()
	cmd.Stderr = types.Writer()

	io.Log("Starting Nuxt.js development server...")
	if err := cmd.Start(); err != nil {
		io.Logf("Failed to start Nuxt.js: %v", io.Fatal, err)
	}

	// Wait for the process to exit or handle as a background process
	go func() {
		err := cmd.Wait()
		if err != nil {
			io.Logf("Nuxt.js process exited with error: %v", io.Fatal, err)
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
	if *opts.Headless {
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
	io.Logf("Server listening on port: %s", io.Info, strconv.Itoa(*opts.Port))
	io.Panicf(
		http.ListenAndServe(*opts.ListenAddr+":"+strconv.Itoa(*opts.Port), nil),
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
	if r.Method == "GET" {
		io.SuccessJSON(w, map[string]interface{}{
			"autoConvert":    opts.AutoConvert,
			"deleteAfter":    opts.DeleteAfter,
			"ignoreExisting": opts.OverwriteExisting,
		})
	}

	if r.Method == "POST" {
		var config Config = Config{}
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			io.ErrorJSON(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		if config.AutoConvert != nil {
			opts.AutoConvert = config.AutoConvert
		}
		if config.DeleteAfter != nil {
			opts.DeleteAfter = config.DeleteAfter
		}
		if config.OverwriteExisting != nil {
			opts.OverwriteExisting = config.OverwriteExisting
		}
		io.SuccessJSON(w, "success")
	}
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fileArray := []types.File{}
	for _, file := range store.FileList {
		fileArray = append(fileArray, file)
	}
	io.SuccessJSON(w, fileArray)
}
