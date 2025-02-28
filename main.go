package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	inputFile  = "/acme/acme.json"
	outputDir  = "/extracted-certs"
	htmlFile   = "./index.html"
	serverPort = ":8080"
)

// ACMEData represents the structure of Traefik's acme.json file
type ACMEData struct {
	LetsEncrypt struct {
		Certificates []struct {
			Domain struct {
				Main string   `json:"main"`
				SANs []string `json:"sans,omitempty"`
			} `json:"domain"`
			Certificate string `json:"certificate"`
			Key         string `json:"key"`
		} `json:"Certificates"`
	} `json:"letsencrypt"`
}

// CertInfo stores information about a certificate
type CertInfo struct {
	Domain    string
	Issuer    string
	NotBefore time.Time
	NotAfter  time.Time
	SANs      []string
	Files     []string
}

// Global map to store certificate information
var certInfoMap = make(map[string]*CertInfo)

// extractCerts reads Traefik's acme.json file and extracts certificates
func extractCerts() error {
	log.Println("Starting certificate extraction...")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read acme.json file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read acme.json: %w", err)
	}

	// Parse JSON
	var acmeData ACMEData
	if err := json.Unmarshal(data, &acmeData); err != nil {
		return fmt.Errorf("failed to parse acme.json: %w", err)
	}

	// Extract certificates
	for _, cert := range acmeData.LetsEncrypt.Certificates {
		domain := cert.Domain.Main
		log.Printf("Processing certificate for domain: %s", domain)

		if cert.Certificate == "" || cert.Key == "" {
			log.Printf("Certificate or key for domain %s is empty", domain)
			continue
		}

		// Decode certificate
		decodedCert, err := base64.StdEncoding.DecodeString(cert.Certificate)
		if err != nil {
			log.Printf("Failed to decode certificate for domain %s: %v", domain, err)
			continue
		}

		// Extract certificate information
		certInfo, err := parseCertificate(decodedCert)
		if err != nil {
			log.Printf("Failed to parse certificate for domain %s: %v", domain, err)
		} else {
			certInfo.Domain = domain
			certInfo.SANs = append(certInfo.SANs, cert.Domain.SANs...)
			certInfoMap[domain] = certInfo
			log.Printf("Certificate info for %s: Valid from %s to %s",
				domain,
				certInfo.NotBefore.Format("2006-01-02"),
				certInfo.NotAfter.Format("2006-01-02"))
		}

		// Write certificate file
		certPath := filepath.Join(outputDir, fmt.Sprintf("%s.crt", domain))
		if err := os.WriteFile(certPath, decodedCert, 0644); err != nil {
			log.Printf("Failed to write certificate file for domain %s: %v", domain, err)
			continue
		}

		if certInfo != nil {
			certInfo.Files = append(certInfo.Files, fmt.Sprintf("%s.crt", domain))
		}

		// Decode key
		decodedKey, err := base64.StdEncoding.DecodeString(cert.Key)
		if err != nil {
			log.Printf("Failed to decode key for domain %s: %v", domain, err)
			continue
		}

		// Write key file
		keyPath := filepath.Join(outputDir, fmt.Sprintf("%s.key", domain))
		if err := os.WriteFile(keyPath, decodedKey, 0600); err != nil {
			log.Printf("Failed to write key file for domain %s: %v", domain, err)
			continue
		}

		if certInfo != nil {
			certInfo.Files = append(certInfo.Files, fmt.Sprintf("%s.key", domain))
		}

		log.Printf("Successfully extracted certificate and key for %s", domain)
	}

	log.Println("Certificate extraction completed")
	return nil
}

// parseCertificate extracts information from a certificate
func parseCertificate(certData []byte) (*CertInfo, error) {
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	info := &CertInfo{
		Issuer:    cert.Issuer.CommonName,
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
		SANs:      cert.DNSNames,
	}

	return info, nil
}

// listFiles handles the display of the list of certificate files
func listFiles(w http.ResponseWriter, r *http.Request) {
	// If root path, serve the HTML file
	if r.URL.Path == "/" {
		htmlContent, err := os.ReadFile(htmlFile)
		if err != nil {
			http.Error(w, "Unable to read HTML template", http.StatusInternalServerError)
			return
		}

		// Generate HTML for certificate cards with detailed information
		var certsHTML strings.Builder

		if len(certInfoMap) == 0 {
			certsHTML.WriteString(`
			<div class="no-certs">
				<i class="fas fa-certificate"></i>
				<p>No certificates found</p>
			</div>
			`)
		} else {
			// Process each certificate
			for domain, certInfo := range certInfoMap {
				// Calculate days until expiration
				daysLeft := int(certInfo.NotAfter.Sub(time.Now()).Hours() / 24)

				// Determine status based on expiration
				statusClass := "status-valid"
				statusIcon := "fa-check-circle"
				if daysLeft < 30 {
					statusClass = "status-warning"
					statusIcon = "fa-exclamation-circle"
				}
				if daysLeft < 7 {
					statusClass = "status-danger"
					statusIcon = "fa-exclamation-triangle"
				}

				certsHTML.WriteString(fmt.Sprintf(`
				<div class="cert-card">
					<div class="cert-header">
						<div class="cert-domain">%s</div>
						<button class="cert-toggle">
							<i class="fas fa-chevron-down"></i>
						</button>
					</div>
					<div class="cert-content">
						<div class="cert-details">
							<div class="status-badge %s">
								<i class="fas %s"></i>
								<span class="days-remaining">%d days</span> remaining
							</div>
							
							<div class="cert-info-grid">
								<div class="cert-info-item">
									<div class="cert-info-label">Valid From</div>
									<div class="cert-info-value">%s</div>
								</div>
								<div class="cert-info-item">
									<div class="cert-info-label">Valid Until</div>
									<div class="cert-info-value">%s</div>
								</div>
								<div class="cert-info-item">
									<div class="cert-info-label">Issuer</div>
									<div class="cert-info-value">%s</div>
								</div>
								<div class="cert-info-item">
									<div class="cert-info-label">Total SANs</div>
									<div class="cert-info-value">%d</div>
								</div>
							</div>
				`,
					domain,
					statusClass,
					statusIcon,
					daysLeft,
					certInfo.NotBefore.Format("2006-01-02"),
					certInfo.NotAfter.Format("2006-01-02"),
					certInfo.Issuer,
					len(certInfo.SANs),
				))

				// Add Subject Alternative Names if any
				if len(certInfo.SANs) > 0 {
					certsHTML.WriteString(`
					<div class="sans-list">
						<div class="sans-title">Subject Alternative Names</div>
						<ul>
					`)
					for _, san := range certInfo.SANs {
						certsHTML.WriteString(fmt.Sprintf("<li>%s</li>", san))
					}
					certsHTML.WriteString(`
						</ul>
					</div>
					`)
				}

				// Add download links
				certsHTML.WriteString(`
					<div class="file-list">
						<div class="file-list-title">Files</div>
						<div class="file-items">
				`)
				for _, file := range certInfo.Files {
					fileIcon := "fa-file-alt"
					if strings.HasSuffix(file, ".key") {
						fileIcon = "fa-key"
					} else if strings.HasSuffix(file, ".crt") {
						fileIcon = "fa-certificate"
					}

					certsHTML.WriteString(fmt.Sprintf(`
						<a href="/certs/%s" download class="file-item">
							<i class="fas %s"></i>
							%s
						</a>
					`, file, fileIcon, file))
				}
				certsHTML.WriteString(`
						</div>
					</div>
				</div>
			</div>
				`)
			}
		}

		// Replace the placeholder in the HTML template
		htmlStr := strings.Replace(string(htmlContent), "{{CERT_GROUPS}}", certsHTML.String(), 1)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlStr)
		return
	}

	// Handle other paths
	http.NotFound(w, r)
}

// serveFile handles serving the certificate files
func serveFile(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(outputDir, filepath.Base(r.URL.Path))
	http.ServeFile(w, r, filePath)
}

func startWebServer() {
	http.HandleFunc("/", listFiles)
	http.HandleFunc("/certs/", serveFile)

	log.Printf("Starting web server on %s", serverPort)
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	// Extract certificates on startup
	if err := extractCerts(); err != nil {
		log.Printf("Initial certificate extraction failed: %v", err)
	}

	// Start periodic extraction in a goroutine
	go func() {
		ticker := time.NewTicker(8 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := extractCerts(); err != nil {
				log.Printf("Periodic certificate extraction failed: %v", err)
			}
		}
	}()

	// Start web server
	startWebServer()
}
