@echo off
REM Build the Go application directly from the current folder
REM This assumes Go is installed and on your PATH

go build -o booking-system.exe .

echo Build complete: booking-system.exe
