// server.go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Create the uploads directory if it doesn't exist
	os.MkdirAll("./uploads", os.ModePerm)

	// Handle the upload route
	http.HandleFunc("/upload", uploadHandler)

	// Serve the HTML form at the root route
	http.HandleFunc("/", formHandler)

	// Start the server
	fmt.Println("Server listening on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

// formHandler serves the HTML form for uploading files
func formHandler(w http.ResponseWriter, r *http.Request) {
	html := `<html>
<head>
    <title>Upload File</title>
</head>
<body>
    <h1>Upload File (Max 100MB)</h1>
    <form enctype="multipart/form-data" action="/upload" method="post">
        <input type="file" name="uploadfile" />
        <input type="submit" value="Upload" />
    </form>
</body>
</html>`
	fmt.Fprint(w, html)
}

// uploadHandler handles the file upload logic
func uploadHandler(w http.ResponseWriter, r *http.Request) {
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
		if written > 1*1024*1024 {
			http.Error(w, "File size exceeds 1MB limit", http.StatusBadRequest)
			// Remove the partially uploaded file
			os.Remove(dstPath)
			return
		}
	}

	// Respond to the client
	fmt.Fprintln(w, "File uploaded successfully.")
}
