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
		body {
			font-family: Arial, sans-serif;
			background-color: #121212;
			color: #e0e0e0;
			margin: 0;
			padding: 20px;
		}
		h1 {
			color: #ffffff;
			text-align: center;
			margin-bottom: 30px;
		}
		.details {
			background: linear-gradient(135deg, #434343 0%, #000000 100%);
			border-radius: 8px;
			margin-bottom: 20px;
			padding: 15px;
		}
		summary {
			font-size: 18px;
			font-weight: bold;
			cursor: pointer;
			color: #1e88e5;
		}
		ul {
			list-style-type: none;
			padding: 0;
			margin: 10px 0 0 0;
		}
		li {
			margin: 5px 0;
		}
		a {
			text-decoration: none;
			color: #81d4fa;
		}
		a:hover {
			text-decoration: underline;
		}
	</style>
</head>
<body>
	<h1>Certificates</h1>`)

	for group, files := range groups {
		fmt.Fprintf(w, `<details class="details"><summary>%s</summary><ul>`, group)
		for _, file := range files {
			fmt.Fprintf(w, `<li><a href="/certs/%s" download>%s</a></li>`, file, file)
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
