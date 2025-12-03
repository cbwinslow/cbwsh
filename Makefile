# cbwsh Makefile
# ================
# Build and install automation for cbwsh
#
# Usage:
#   make              - Build the binary
#   make install      - Install to /usr/local/bin
#   make uninstall    - Remove installed binary
#   make test         - Run tests
#   make lint         - Run linter
#   make clean        - Clean build artifacts
#
# Variables:
#   PREFIX            - Installation prefix (default: /usr/local)
#   DESTDIR           - Destination root for packaging
#   VERSION           - Version string override
#

# Project metadata
BINARY_NAME := cbwsh
PACKAGE := github.com/cbwinslow/cbwsh
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Installation paths
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin
MANDIR := $(PREFIX)/share/man/man1
CONFIG_DIR := $(HOME)/.cbwsh

# Go settings
GO ?= go
GOFLAGS ?=
CGO_ENABLED ?= 0

# Build settings
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(BUILD_DATE)

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m

.PHONY: all build install uninstall clean test lint fmt vet deps help \
        install-config cross-compile release version

# Default target
all: build

# Help message
help:
	@echo "$(CYAN)cbwsh Makefile$(NC)"
	@echo ""
	@echo "$(GREEN)Build Targets:$(NC)"
	@echo "  make              Build the binary"
	@echo "  make build        Build the binary"
	@echo "  make cross        Cross-compile for all platforms"
	@echo ""
	@echo "$(GREEN)Install Targets:$(NC)"
	@echo "  make install      Install to $(PREFIX)/bin"
	@echo "  make uninstall    Remove installed binary"
	@echo "  make install-config Create default configuration"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  make test         Run tests"
	@echo "  make test-v       Run tests with verbose output"
	@echo "  make test-race    Run tests with race detector"
	@echo "  make coverage     Run tests with coverage"
	@echo "  make lint         Run golangci-lint"
	@echo "  make fmt          Format code"
	@echo "  make vet          Run go vet"
	@echo "  make deps         Download dependencies"
	@echo ""
	@echo "$(GREEN)Other:$(NC)"
	@echo "  make clean        Clean build artifacts"
	@echo "  make version      Show version info"
	@echo "  make help         Show this help"
	@echo ""
	@echo "$(YELLOW)Variables:$(NC)"
	@echo "  PREFIX=$(PREFIX)"
	@echo "  VERSION=$(VERSION)"
	@echo "  CGO_ENABLED=$(CGO_ENABLED)"

# Build the binary
build:
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "$(GREEN)Build complete: ./$(BINARY_NAME)$(NC)"

# Build with debug info
build-debug:
	@echo "$(CYAN)Building $(BINARY_NAME) (debug)...$(NC)"
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .
	@echo "$(GREEN)Debug build complete: ./$(BINARY_NAME)$(NC)"

# Install the binary
install: build
	@echo "$(CYAN)Installing $(BINARY_NAME) to $(DESTDIR)$(BINDIR)...$(NC)"
	@mkdir -p $(DESTDIR)$(BINDIR)
	@install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	@echo "$(GREEN)Installed to $(DESTDIR)$(BINDIR)/$(BINARY_NAME)$(NC)"

# Install with config
install-all: install install-config
	@echo "$(GREEN)Full installation complete$(NC)"

# Install default configuration
install-config:
	@echo "$(CYAN)Creating default configuration...$(NC)"
	@mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/config.yaml ]; then \
		echo "# cbwsh Configuration" > $(CONFIG_DIR)/config.yaml; \
		echo "# https://github.com/cbwinslow/cbwsh" >> $(CONFIG_DIR)/config.yaml; \
		echo "" >> $(CONFIG_DIR)/config.yaml; \
		echo "shell:" >> $(CONFIG_DIR)/config.yaml; \
		echo "  default_shell: bash" >> $(CONFIG_DIR)/config.yaml; \
		echo "  history_size: 10000" >> $(CONFIG_DIR)/config.yaml; \
		echo "" >> $(CONFIG_DIR)/config.yaml; \
		echo "ui:" >> $(CONFIG_DIR)/config.yaml; \
		echo "  theme: default" >> $(CONFIG_DIR)/config.yaml; \
		echo "  layout: single" >> $(CONFIG_DIR)/config.yaml; \
		echo "  show_status_bar: true" >> $(CONFIG_DIR)/config.yaml; \
		echo "  enable_animations: true" >> $(CONFIG_DIR)/config.yaml; \
		echo "  syntax_highlighting: true" >> $(CONFIG_DIR)/config.yaml; \
		echo "" >> $(CONFIG_DIR)/config.yaml; \
		echo "ai:" >> $(CONFIG_DIR)/config.yaml; \
		echo "  provider: none" >> $(CONFIG_DIR)/config.yaml; \
		echo "  api_key: \"\"" >> $(CONFIG_DIR)/config.yaml; \
		echo "  model: \"\"" >> $(CONFIG_DIR)/config.yaml; \
		echo "  enable_suggestions: false" >> $(CONFIG_DIR)/config.yaml; \
		echo "$(GREEN)Created $(CONFIG_DIR)/config.yaml$(NC)"; \
	else \
		echo "$(YELLOW)Configuration already exists$(NC)"; \
	fi

# Uninstall the binary
uninstall:
	@echo "$(CYAN)Removing $(BINARY_NAME)...$(NC)"
	@rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
	@echo "$(GREEN)Uninstalled $(BINARY_NAME)$(NC)"

# Clean build artifacts
clean:
	@echo "$(CYAN)Cleaning...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f coverage.txt coverage.html
	@rm -rf dist/
	@$(GO) clean
	@echo "$(GREEN)Clean complete$(NC)"

# Download dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

# Run tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	$(GO) test ./...

# Run tests with verbose output
test-v:
	@echo "$(CYAN)Running tests (verbose)...$(NC)"
	$(GO) test -v ./...

# Run tests with race detector
test-race:
	@echo "$(CYAN)Running tests with race detector...$(NC)"
	$(GO) test -race ./...

# Run tests with coverage
coverage:
	@echo "$(CYAN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

# Run linter
lint:
	@echo "$(CYAN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not found, installing...$(NC)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

# Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

# Run go vet
vet:
	@echo "$(CYAN)Running vet...$(NC)"
	$(GO) vet ./...

# Cross-compile for multiple platforms
cross-compile: clean
	@echo "$(CYAN)Cross-compiling for multiple platforms...$(NC)"
	@mkdir -p dist
	
	# Linux
	@echo "  Building linux/amd64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_linux_amd64 .
	@echo "  Building linux/arm64..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_linux_arm64 .
	@echo "  Building linux/386..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_linux_386 .
	
	# macOS
	@echo "  Building darwin/amd64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_darwin_amd64 .
	@echo "  Building darwin/arm64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_darwin_arm64 .
	
	# Windows
	@echo "  Building windows/amd64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_windows_amd64.exe .
	@echo "  Building windows/arm64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_windows_arm64.exe .
	
	# FreeBSD
	@echo "  Building freebsd/amd64..."
	@CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_freebsd_amd64 .
	
	@echo "$(GREEN)Cross-compilation complete. Binaries in dist/$(NC)"

# Show version information
version:
	@echo "$(CYAN)Version Information$(NC)"
	@echo "  Binary:  $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Date:    $(BUILD_DATE)"
	@echo "  Go:      $(shell $(GO) version)"

# Run the application
run: build
	@./$(BINARY_NAME)

# Development mode - watch for changes and rebuild
dev:
	@echo "$(CYAN)Starting development mode...$(NC)"
	@echo "$(YELLOW)Install 'watchexec' or 'entr' for auto-rebuild$(NC)"
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -r -e go make build; \
	elif command -v entr >/dev/null 2>&1; then \
		find . -name '*.go' | entr -r make build; \
	else \
		echo "No file watcher found. Running single build."; \
		make build; \
	fi

# Docker build
docker:
	@echo "$(CYAN)Building Docker image...$(NC)"
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "$(GREEN)Docker image built: $(BINARY_NAME):$(VERSION)$(NC)"

# Generate checksums for releases
checksums:
	@echo "$(CYAN)Generating checksums...$(NC)"
	@cd dist && sha256sum * > checksums.txt
	@echo "$(GREEN)Checksums saved to dist/checksums.txt$(NC)"

# Security audit
audit:
	@echo "$(CYAN)Running security audit...$(NC)"
	$(GO) list -json -m all | docker run --rm -i sonatypecommunity/nancy:latest sleuth || true
	@echo "$(GREEN)Audit complete$(NC)"

# Update dependencies
update-deps:
	@echo "$(CYAN)Updating dependencies...$(NC)"
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies updated$(NC)"
