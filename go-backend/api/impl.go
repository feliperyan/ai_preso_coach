package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body to 200MB
	r.Body = http.MaxBytesReader(w, r.Body, 200*1024*1024+1024) // Add a little extra to account for multipart overhead

	// Parse the multipart form data
	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Invalid multipart/form-data", http.StatusBadRequest)
		return
	}

	// Not convinced I need this for loop. I'll only ever process 1 file which is the video upload.
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
		extension := strings.Trim(filepath.Ext(filename), ".")
		if extension == "" {
			http.Error(w, fmt.Sprintf("error: your file '%s' does not have an extension", filename), http.StatusInternalServerError)
			log.Printf("error: your file '%s' does not have an extension", filename)
			return
		}

		var buf bytes.Buffer
		tee := io.TeeReader(part, &buf)

		head := make([]byte, 512)
		_, err = tee.Read(head)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading the file content", http.StatusInternalServerError)
			return
		}

		contentType := http.DetectContentType(head)
		if !strings.HasPrefix(contentType, "video/") {
			http.Error(w, "The uploaded file is not a video", http.StatusBadRequest)
			log.Println("error: The uploaded file is not a video")
			return
		}

		dstPath := filepath.Join("./uploads", filename)

		// Create a destination file
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Unable to create the file on server", http.StatusInternalServerError)
			log.Println("error: ", err)
			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination file in chunks
		written, err := io.Copy(dst, io.MultiReader(&buf, part))
		if err != nil {
			http.Error(w, "Error saving the file", http.StatusInternalServerError)
			log.Println("error: ", err)

			// Remove the partially uploaded file
			os.Remove(dstPath)
			return
		}

		// Check if the uploaded file size exceeds 200MB
		if written > 200*1024*1024 {
			http.Error(w, "File size exceeds 200MB limit", http.StatusBadRequest)
			log.Println("error: File size exceeds 200MB limit", err)

			// Remove the partially uploaded file
			os.Remove(dstPath)
			return
		}

	}

	log.Println("UploadFile OK")
	fmt.Fprintln(w, "File uploaded successfully.")
}

func (Server) Thingy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Got thingy")
}
