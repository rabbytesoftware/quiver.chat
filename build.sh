#!/bin/bash

# Quiver Chat - Cross-platform build script
# This script builds executables for Windows, Linux, and macOS

set -e

echo "🚀 Building Quiver Chat for multiple platforms..."

# Clean previous builds
rm -rf dist/
mkdir -p dist/

# App info
APP_NAME="quiver-chat"
VERSION="1.0.0"

# Build for different platforms
echo "📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/${APP_NAME}-linux-amd64 .

echo "📦 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/${APP_NAME}-linux-arm64 .

echo "📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/${APP_NAME}-macos-amd64 .

echo "📦 Building for macOS (arm64 - Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/${APP_NAME}-macos-arm64 .

echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/${APP_NAME}-windows-amd64.exe .

echo "📦 Building for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/${APP_NAME}-windows-arm64.exe .

echo ""
echo "✅ Build completed! Binaries available in dist/ directory:"
echo ""
ls -la dist/

echo ""
echo "📋 Usage instructions:"
echo "1. Copy the appropriate binary for your platform"
echo "2. Run the binary: ./${APP_NAME}-[platform]"
echo "3. Open your browser and navigate to http://localhost:8080"
echo "4. Choose a nickname and start chatting!"
echo ""
echo "🌐 No external dependencies required - just run and go!" 