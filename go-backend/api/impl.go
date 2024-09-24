package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body to 100MB
	r.Body = http.MaxBytesReader(w, r.Body, 100*1024*1024+1024) // Add a little extra to account for multipart overhead

	// Parse the multipart form data
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart/form-data", http.StatusBadRequest)
		return
	}

	// Process each part of the multipart form
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break // No more parts
		}
		if err != nil {
			http.Error(w, "Error reading multipart data", http.StatusInternalServerError)
			return
		}

		if part.FileName() == "" {
			// Skip if not a file
			continue
		}

		// Securely join the upload directory and the filename to prevent directory traversal
		filename := filepath.Base(part.FileName())
		dstPath := filepath.Join("./uploads", filename)

		// Create a destination file
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Unable to create the file on server", http.StatusInternalServerError)

			log.Println(err)

			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination file in chunks
		written, err := io.Copy(dst, part)
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			return
		}

		// Check if the uploaded file size exceeds 100MB
		if written > 100*1024*1024 {
			http.Error(w, "File size exceeds 100MB limit", http.StatusBadRequest)
			// Remove the partially uploaded file
			os.Remove(dstPath)
			return
		}
	}

	log.Println("UploadFile OK")
	// Respond to the client
	fmt.Fprintln(w, "File uploaded successfully.")
}

func (Server) Thingy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "File uploaded successfully.")
}
