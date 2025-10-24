# GoThink Makefile

# Variables
BINARY_NAME=gothink
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_MACOS=$(BINARY_NAME)_macos

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
BUILD_FLAGS=-ldflags "-s -w"

.PHONY: all build clean test deps run help

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) -v .

# Build for different platforms
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_UNIX) -v ./...

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_WINDOWS) -v ./...

build-macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_MACOS) -v ./...

# Build all platforms
build-all: build-linux build-windows build-macos

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_MACOS)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

# Run in development mode
dev:
	$(GOCMD) run main.go

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
install:
	$(GOGET) -d -v ./...

# Update dependencies
update:
	$(GOMOD) download
	$(GOMOD) tidy

# Docker commands
docker-build:
	docker build -t gothink .

docker-run:
	docker run -p 8080:8080 gothink

# Help
help:
	@echo "Available targets:"
	@echo "  all          - Run tests and build"
	@echo "  build        - Build the application"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  build-macos  - Build for macOS"
	@echo "  build-all    - Build for all platforms"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  deps         - Download dependencies"
	@echo "  run          - Build and run the application"
	@echo "  dev          - Run in development mode"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  install      - Install dependencies"
	@echo "  update       - Update dependencies"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message"
