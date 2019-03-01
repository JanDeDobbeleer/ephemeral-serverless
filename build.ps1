# Clean the build artifacts
Remove-Item -Force -Recurse ./bin -ErrorAction SilentlyContinue
Remove-Item -Force -Recurse ./vendor -ErrorAction SilentlyContinue
Remove-Item -Force Gopkg.lock -ErrorAction SilentlyContinue
# Ensure we have all dependencies
dep ensure -v
# Build for Linux
$env:GOOS = "linux"
go build -ldflags="-s -w" -o bin/twitter main.go
# Reset to Windows for safety
$env:GOOS = "windows"
