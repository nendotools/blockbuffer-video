package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	appIO "blockbuffer/internal/io"
	opts "blockbuffer/internal/settings"
)

func HandleUpload(w http.ResponseWriter, r *http.Request, deferMove bool) {
	if r.Method != http.MethodPost {
		appIO.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		appIO.ErrorJSON(w, "Failed to get file", http.StatusBadRequest)
		return
	}

	_, tempPath, err := writeFile(file, *header)
	if err != nil {
		appIO.ErrorJSON(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	if !deferMove {
		moveFile(tempPath, filepath.Join(*opts.WatchDir, header.Filename))
	}
	fmt.Fprintf(w, "File uploaded: %s\n", header.Filename)
}

func HandleUploadMultipleFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		appIO.ErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		appIO.ErrorJSON(w, "Failed to get files", http.StatusBadRequest)
		return
	}

	tempFiles := []string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if part.FormName() != "files" || part.FileName() == "" {
			continue
		}

		_, tempPath, err := writeFile(part, *&multipart.FileHeader{Filename: part.FileName()})
		if err != nil {
			appIO.ErrorJSON(w, "Failed to write file", http.StatusInternalServerError)
			return
		}
		tempFiles = append(tempFiles, tempPath)
	}

	filenames := []string{}
	for _, tempPath := range tempFiles {
		filename := filepath.Base(tempPath)
		filenames = append(filenames, filename)
		moveFile(tempPath, filepath.Join(*opts.WatchDir, filename))
	}

	jsonData := json.NewEncoder(w)
	jsonData.Encode(filenames)
	return
}

func writeFile(reader io.Reader, header multipart.FileHeader) (*os.File, string, error) {
	tempPath := filepath.Join(*opts.UploadDir, header.Filename)
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return nil, "", err
	}

	defer tempFile.Close()

	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return nil, "", err
	}

	return tempFile, tempPath, nil
}

func moveFile(tempPath string, watchFilePath string) error {
	err := os.Rename(tempPath, watchFilePath)
	if err != nil {
		return err
	}
	return nil
}
