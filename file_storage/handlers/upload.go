package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse data
	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Unable to parse form file", http.StatusBadRequest)
		return
	}

	// get file from form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get form from file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	path := filepath.Join("storage", filename)

	// create file
	out, err := os.Create(path)
	if err != nil {
		http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// copy file in chunks
	chunksize := 8192
	buffer := make([]byte, chunksize) // 8kb buffer for chunks
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break // stop when reached end of file
		}
		if err != nil {
			http.Error(w, "Error trying to read file", http.StatusInternalServerError)
			return
		}

		// write chunk to output
		_, err = out.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Error trying to write to file", http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "File %s uploaded successfully\n", filename)
}
