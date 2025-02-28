# Build for standard Linux (like Debian, Ubuntu) with amd64 architecture
Write-Host "Building for Linux amd64..."
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o tce.amd64

Write-Host "Building for Linux arm64..."
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -o tce.arm64


$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:GOARM = ""


Write-Host "Build completed and Saved"
