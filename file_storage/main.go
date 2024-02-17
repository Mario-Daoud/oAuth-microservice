package main

import (
	"fmt"
	"log"
	"net/http"

    "harbor/file_storage/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	router.HandleFunc("/delete/{filename}", handlers.DeleteHandler).Methods("DELETE")

	port := fmt.Sprintf(":8080")
	log.Printf("Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
