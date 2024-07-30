package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const outputDir = "./extracted-certs"

// listFiles handles the display of the list of certificate files
func listFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Certificates</title>
	<style>
		body { font-family: Arial, sans-serif; }
		h1 { color: #333; }
		ul { list-style-type: none; padding: 0; }
		li { margin: 5px 0; }
		a { text-decoration: none; color: #1a73e8; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<h1>Certificates</h1>
	<ul>`)

	for _, file := range files {
		if !file.IsDir() {
			fmt.Fprintf(w, `<li><a href="/certs/%s">%s</a></li>`, file.Name(), file.Name())
		}
	}

	fmt.Fprintln(w, `</ul>
</body>
</html>`)
}

// serveFile handles serving the certificate files
func serveFile(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(outputDir, filepath.Base(r.URL.Path))
	http.ServeFile(w, r, filePath)
}

func main() {
	http.HandleFunc("/", listFiles)
	http.HandleFunc("/certs/", serveFile)

	log.Println("Starting web server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
