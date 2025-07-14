# Quiver Chat - Cross-platform Makefile
# This Makefile builds executables for Windows, Linux, and macOS

# Variables
APP_NAME := quiver-chat
VERSION := 27.7.1
DIST_DIR := dist
FRONTEND_DIR := frontend

# Go build flags
LDFLAGS := -ldflags="-s -w"

# Platforms and architectures
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

# Default target
.PHONY: all
all: clean frontend build

# Clean previous builds
.PHONY: clean
clean:
	@echo "🧹 Cleaning previous builds..."
	@rm -rf $(DIST_DIR)/
	@mkdir -p $(DIST_DIR)/

# Build the Next.js frontend
.PHONY: frontend
frontend:
	@echo "🎨 Building Next.js frontend..."
	@cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "✅ Frontend build completed!"

# Build for all platforms
.PHONY: build
build: $(PLATFORMS)

# Build for specific platforms
.PHONY: linux/amd64
linux/amd64:
	@echo "📦 Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .

.PHONY: linux/arm64
linux/arm64:
	@echo "📦 Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .

.PHONY: darwin/amd64
darwin/amd64:
	@echo "📦 Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-macos-amd64 .

.PHONY: darwin/arm64
darwin/arm64:
	@echo "📦 Building for macOS (arm64 - Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-macos-arm64 .

.PHONY: windows/amd64
windows/amd64:
	@echo "📦 Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .

.PHONY: windows/arm64
windows/arm64:
	@echo "📦 Building for Windows (arm64)..."
	@GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/$(APP_NAME)-windows-arm64.exe .

# Show build results
.PHONY: list
list:
	@echo ""
	@echo "✅ Build completed! Binaries available in $(DIST_DIR)/ directory:"
	@echo ""
	@ls -la $(DIST_DIR)/
	@echo ""
	@echo "📋 Usage instructions:"
	@echo "1. Copy the appropriate binary for your platform"
	@echo "2. Run the binary: ./$(APP_NAME)-[platform]"
	@echo "3. Open your browser and navigate to http://localhost:8080"
	@echo "4. Choose a nickname and start chatting!"
	@echo ""
	@echo "🌐 No external dependencies required - just run and go!"

# Development targets
.PHONY: dev
dev:
	@echo "🚀 Starting development server..."
	@go run .

.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v -race ./...

.PHONY: coverage
coverage:
	@echo "📊 Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...

# Release targets
.PHONY: release-archives
release-archives: build
	@echo "📦 Creating release archives..."
	@cd $(DIST_DIR) && \
		tar -czf $(APP_NAME)-linux-amd64.tar.gz $(APP_NAME)-linux-amd64 && \
		tar -czf $(APP_NAME)-linux-arm64.tar.gz $(APP_NAME)-linux-arm64 && \
		tar -czf $(APP_NAME)-macos-amd64.tar.gz $(APP_NAME)-macos-amd64 && \
		tar -czf $(APP_NAME)-macos-arm64.tar.gz $(APP_NAME)-macos-arm64 && \
		zip $(APP_NAME)-windows-amd64.zip $(APP_NAME)-windows-amd64.exe && \
		zip $(APP_NAME)-windows-arm64.zip $(APP_NAME)-windows-arm64.exe

# Help target
.PHONY: help
help:
	@echo "🚀 Quiver Chat Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build everything (clean + frontend + all platforms)"
	@echo "  clean            - Clean previous builds"
	@echo "  frontend         - Build Next.js frontend only"
	@echo "  build            - Build for all platforms"
	@echo "  linux/amd64     - Build for Linux x64"
	@echo "  linux/arm64     - Build for Linux ARM64"
	@echo "  darwin/amd64    - Build for macOS Intel"
	@echo "  darwin/arm64    - Build for macOS Apple Silicon"
	@echo "  windows/amd64   - Build for Windows x64"
	@echo "  windows/arm64   - Build for Windows ARM64"
	@echo "  list             - Show build results"
	@echo "  dev              - Start development server"
	@echo "  test             - Run tests"
	@echo "  coverage         - Run tests with coverage"
	@echo "  release-archives - Create release archives"
	@echo "  help             - Show this help message" 