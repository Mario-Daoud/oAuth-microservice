package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0) // 0 = unlimited size - change to put apply file size limit
	if err != nil {
		http.Error(w, "Unable to parse form file", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]

	storageDir := "storage"
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening uploaded file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := filepath.Join(storageDir, fileHeader.Filename)

		if _, err := os.Stat(storageDir); os.IsNotExist(err) {
			if err := os.MkdirAll(storageDir, 0777); err != nil {
				http.Error(w, "Error creating directory", http.StatusInternalServerError)
				return
			}
		}

		out, err := os.Create(filename)
		if err != nil {
			http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		// buffer for handling bigger files too
		const chunkSize = 8192 // 8KB chunk size
		buffer := make([]byte, chunkSize)

		// copy file content to server in chunks
		for {
			// read chunk from the input file
			n, err := file.Read(buffer)
			if err == io.EOF {
				break // end of file reached
			}
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			// write chunk to output file
			chunk := buffer[:n]
			_, err = out.Write(chunk)
			if err != nil {
				http.Error(w, "Error writing file", http.StatusInternalServerError)
				return
			}
		}

		fmt.Fprintf(w, "File %s uploaded successfully\n", fileHeader.Filename)
	}
}
