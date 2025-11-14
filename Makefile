.PHONY: build build-all test test-coverage test-race lint fmt clean install help

# Binary name
BINARY_NAME=git-issue

# Build for current platform
build:
	@echo "Building for current platform..."
	go build -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

# Cross-compile for all supported platforms
build-all:
	@echo "Building for all platforms..."
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64
	@echo "Cross-compilation complete"
	@ls -lh $(BINARY_NAME)-*

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run tests with coverage report
test-coverage-report:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install to /usr/local/bin
install: build
	@echo "Installing to /usr/local/bin..."
	sudo mv $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-all          - Cross-compile for macOS (ARM64/AMD64) and Linux (AMD64)"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage"
	@echo "  test-coverage-report - Generate HTML coverage report"
	@echo "  lint               - Run golangci-lint"
	@echo "  fmt                - Format code with go fmt"
	@echo "  clean              - Remove build artifacts"
	@echo "  install            - Build and install to /usr/local/bin"
	@echo "  help               - Show this help message"
