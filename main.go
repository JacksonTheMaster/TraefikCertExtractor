package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const outputDir = "/extracted-certs"

// groupFiles groups certificate files by their names without extensions
func groupFiles(files []os.DirEntry) map[string][]string {
	groups := make(map[string][]string)
	for _, file := range files {
		if !file.IsDir() {
			name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			groups[name] = append(groups[name], file.Name())
		}
	}
	return groups
}

// listFiles handles the display of the list of certificate files
func listFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}

	groups := groupFiles(files)

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
		.details { margin-bottom: 20px; }
		summary { font-size: 18px; font-weight: bold; cursor: pointer; }
		ul { list-style-type: none; padding: 0; }
		li { margin: 5px 0; }
		a { text-decoration: none; color: #1a73e8; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<h1>Certificates</h1>`)

	for group, files := range groups {
		fmt.Fprintf(w, `<details class="details"><summary>%s</summary><ul>`, group)
		for _, file := range files {
			fmt.Fprintf(w, `<li><a href="/certs/%s">%s</a></li>`, file, file)
		}
		fmt.Fprintln(w, `</ul></details>`)
	}

	fmt.Fprintln(w, `</body>
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
