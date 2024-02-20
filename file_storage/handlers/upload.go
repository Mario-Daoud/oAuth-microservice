package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(0)
    if err != nil {
        http.Error(w, "Unable to parse form file", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Unable to get form from file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    filename := handler.Filename
    path := filepath.Join("storage", filename)

    out, err := os.Create(path)
    if err != nil {
        http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
        return
    }
    defer out.Close()

    _, err = io.Copy(out, file)
    if err != nil {
        http.Error(w, "Unable to copy file contents", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "File %s uploaded successfully\n", filename)
}
