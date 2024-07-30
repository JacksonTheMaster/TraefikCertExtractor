# Build for standard Linux (like Debian, Ubuntu) with amd64 architecture
Write-Host "Starting build for standard Linux with amd64 architecture..."
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:GOARM = ""
go build -o tce.amd64
Write-Host "Build for standard Linux with amd64 architecture completed."

$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:GOARM = ""


Write-Host "Build completed and Saved"
