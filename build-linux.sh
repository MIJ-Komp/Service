#!/bin/bash

# Simple build script for MIJKomp Service - Linux Production
# Usage: ./build-linux.sh

set -e  # Exit on any error

echo "Building MIJKomp Service for Linux production..."

# Clean previous builds
echo "Cleaning previous builds..."
rm -f mijkomp-service
rm -f mijkomp-service.exe

# Generate wire dependencies
echo "Generating wire dependencies..."
go generate ./...

# Download and tidy dependencies
echo "Downloading dependencies..."
go mod download
go mod tidy

# Build for Linux production
echo "Building for Linux (production)..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o mijkomp-service .

# Make executable
# chmod +x mijkomp-service

echo "Build completed successfully!"
echo "Binary created: mijkomp-service"
echo ""
echo "To deploy:"
echo "1. Copy mijkomp-service to your Linux server"
echo "2. Copy .env.example to .env and configure it"
echo "3. Copy docs/ folder to server (for Swagger)"
echo "4. Run: ./mijkomp-service"
echo ""
echo "File size:"
ls -lh mijkomp-service