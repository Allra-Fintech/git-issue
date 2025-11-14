.PHONY: build build-all test test-coverage test-coverage-report lint fmt clean install help

GO ?= go
BINARY_NAME ?= gi
INSTALL_DIR ?= $(HOME)/.local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS ?= -s -w -X main.version=$(VERSION)

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) (version $(VERSION))..."
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/gi
	@echo "Build complete: $(BINARY_NAME)"

# Cross-compile for all supported platforms
build-all:
	@echo "Cross-compiling binaries (version $(VERSION))..."
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 ./cmd/gi
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 ./cmd/gi
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 ./cmd/gi
	@ls -lh $(BINARY_NAME)-darwin-arm64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-linux-amd64

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -cover ./...

# Run tests with coverage report
test-coverage-report:
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running golangci-lint..."
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME) $(BINARY_NAME)-darwin-arm64 $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME)-linux-amd64
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install to INSTALL_DIR (default: ~/.local/bin)
install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p "$(INSTALL_DIR)"
	install -m 755 $(BINARY_NAME) "$(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "Installation complete"

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build                 Build for current platform"
	@echo "  build-all             Cross-compile for macOS (arm64/amd64) and Linux (amd64)"
	@echo "  test                  Run unit tests"
	@echo "  test-coverage         Run tests with coverage summary"
	@echo "  test-coverage-report  Generate HTML coverage report"
	@echo "  lint                  Run golangci-lint"
	@echo "  fmt                   Format source files"
	@echo "  clean                 Remove build artifacts"
	@echo "  install               Build and install to $(INSTALL_DIR)"
	@echo "  help                  Show this help message"
