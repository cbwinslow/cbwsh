#!/bin/bash
#
# cbwsh Installation Script
# ===========================
# A modern, modular terminal shell built with Bubble Tea
#
# This script installs cbwsh on Unix-like systems (Linux, macOS, FreeBSD)
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
#   or
#   wget -qO- https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
#
# Options:
#   --version <version>  Install a specific version (default: latest)
#   --prefix <path>      Install to a custom location (default: /usr/local/bin)
#   --no-sudo            Don't use sudo for installation
#   --help               Show this help message
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values
VERSION="latest"
PREFIX="/usr/local/bin"
USE_SUDO="true"
GITHUB_REPO="cbwinslow/cbwsh"
BINARY_NAME="cbwsh"

# Global variables for temp directory (for cleanup trap)
TMP_DIR=""

# Logging functions
info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

error() {
    echo -e "${RED}[ERROR]${NC} $*"
    exit 1
}

# Show banner
show_banner() {
    echo -e "${BOLD}${BLUE}"
    cat << 'EOF'
   _____ ______          _____  _    _ 
  / ____|  _ \ \        / / __|| |  | |
 | |    | |_) \ \  /\  / / |__ | |__| |
 | |    |  _ < \ \/  \/ /\__ \|  __  |
 | |____| |_) | \  /\  / ___) | |  | |
  \_____|____/   \/  \/ |____/|_|  |_|
                                        
  Custom Bubble Tea Shell
EOF
    echo -e "${NC}"
}

# Show help
show_help() {
    cat << EOF
cbwsh Installer

Usage: $0 [options]

Options:
  --version <version>  Install a specific version (default: latest)
  --prefix <path>      Install to a custom location (default: /usr/local/bin)
  --no-sudo            Don't use sudo for installation
  --help               Show this help message

Examples:
  $0                           # Install latest version to /usr/local/bin
  $0 --version v1.0.0          # Install specific version
  $0 --prefix ~/.local/bin     # Install to custom location
  $0 --no-sudo                 # Install without sudo

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --prefix)
                PREFIX="$2"
                shift 2
                ;;
            --no-sudo)
                USE_SUDO="false"
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                error "Unknown option: $1"
                ;;
        esac
    done
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    # Normalize OS name
    case "$OS" in
        darwin)
            OS="darwin"
            ;;
        linux)
            OS="linux"
            ;;
        freebsd)
            OS="freebsd"
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac

    # Normalize architecture
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        armv7l|armv7)
            ARCH="armv7"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac

    info "Detected platform: ${OS}/${ARCH}"
}

# Check for required dependencies
check_dependencies() {
    local missing_deps=()

    # Check for curl or wget
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        missing_deps+=("curl or wget")
    fi

    # Check for tar
    if ! command -v tar &> /dev/null; then
        missing_deps+=("tar")
    fi

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        error "Missing required dependencies: ${missing_deps[*]}"
    fi
}

# Get the latest release version
get_latest_version() {
    info "Fetching latest version..."
    
    if command -v curl &> /dev/null; then
        VERSION=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        VERSION=$(wget -qO- "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
    fi

    if [[ -z "$VERSION" || "$VERSION" == "null" ]]; then
        error "Could not determine latest version. Please specify a version with --version."
    fi

    info "Latest version: $VERSION"
}

# Download the binary
download_binary() {
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    TMP_DIR=$(mktemp -d)

    info "Downloading cbwsh ${VERSION} for ${OS}/${ARCH}..."
    info "URL: $download_url"

    cd "$TMP_DIR" || error "Failed to create temp directory"

    if command -v curl &> /dev/null; then
        curl -fsSL "$download_url" -o "${BINARY_NAME}.tar.gz" || error "Download failed. Please check the version and try again."
    else
        wget -q "$download_url" -O "${BINARY_NAME}.tar.gz" || error "Download failed. Please check the version and try again."
    fi

    info "Extracting archive..."
    tar -xzf "${BINARY_NAME}.tar.gz" || error "Extraction failed"

    BINARY_PATH="${TMP_DIR}/${BINARY_NAME}"
    if [[ ! -f "$BINARY_PATH" ]]; then
        # Try common subdirectory patterns first before using find
        for subpath in "./${BINARY_NAME}" "./bin/${BINARY_NAME}" "./${BINARY_NAME}-${VERSION}/${BINARY_NAME}"; do
            if [[ -f "${TMP_DIR}/${subpath}" ]]; then
                BINARY_PATH="${TMP_DIR}/${subpath}"
                break
            fi
        done
        
        # Fall back to find if common patterns didn't work
        if [[ ! -f "$BINARY_PATH" ]]; then
            BINARY_PATH=$(find "$TMP_DIR" -name "$BINARY_NAME" -type f -print -quit 2>/dev/null)
            if [[ -z "$BINARY_PATH" || ! -f "$BINARY_PATH" ]]; then
                error "Binary not found in archive"
            fi
        fi
    fi

    success "Download complete"
}

# Install the binary
install_binary() {
    info "Installing to ${PREFIX}..."

    # Create prefix directory if it doesn't exist
    if [[ ! -d "$PREFIX" ]]; then
        if [[ "$USE_SUDO" == "true" ]] && [[ "$EUID" -ne 0 ]]; then
            sudo mkdir -p "$PREFIX" || error "Failed to create directory: $PREFIX"
        else
            mkdir -p "$PREFIX" || error "Failed to create directory: $PREFIX"
        fi
    fi

    # Copy binary to destination
    if [[ "$USE_SUDO" == "true" ]] && [[ "$EUID" -ne 0 ]]; then
        sudo cp "$BINARY_PATH" "${PREFIX}/${BINARY_NAME}" || error "Failed to copy binary"
        sudo chmod +x "${PREFIX}/${BINARY_NAME}" || error "Failed to set permissions"
    else
        cp "$BINARY_PATH" "${PREFIX}/${BINARY_NAME}" || error "Failed to copy binary"
        chmod +x "${PREFIX}/${BINARY_NAME}" || error "Failed to set permissions"
    fi

    success "cbwsh installed to ${PREFIX}/${BINARY_NAME}"
}

# Create default configuration
create_config() {
    local config_dir="${HOME}/.cbwsh"
    local config_file="${config_dir}/config.yaml"

    if [[ ! -d "$config_dir" ]]; then
        info "Creating configuration directory..."
        mkdir -p "$config_dir"
    fi

    if [[ ! -f "$config_file" ]]; then
        info "Creating default configuration..."
        cat > "$config_file" << 'EOF'
# cbwsh Configuration
# https://github.com/cbwinslow/cbwsh

shell:
  default_shell: bash
  history_size: 10000

ui:
  theme: default
  layout: single
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true

ai:
  provider: none  # Options: none, openai, anthropic, gemini, local
  api_key: ""
  model: ""
  enable_suggestions: false

secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM
  key_derivation: argon2id

keybindings:
  quit: ctrl+q
  help: ctrl+?
  ai_assist: ctrl+a
EOF
        success "Created default configuration at $config_file"
    else
        info "Configuration file already exists at $config_file"
    fi
}

# Verify installation
verify_installation() {
    info "Verifying installation..."

    if command -v "${PREFIX}/${BINARY_NAME}" &> /dev/null; then
        local version_output
        version_output=$("${PREFIX}/${BINARY_NAME}" --version 2>/dev/null || echo "installed")
        success "cbwsh is installed and ready to use!"
    else
        # Check if PREFIX is in PATH
        if [[ ":$PATH:" != *":${PREFIX}:"* ]]; then
            warn "Installation directory is not in your PATH"
            echo ""
            echo "Add the following to your shell configuration file:"
            echo ""
            echo "  export PATH=\"${PREFIX}:\$PATH\""
            echo ""
        fi
    fi
}

# Print post-installation instructions
print_instructions() {
    echo ""
    echo -e "${BOLD}${GREEN}Installation Complete!${NC}"
    echo ""
    echo "To start cbwsh, run:"
    echo ""
    echo -e "  ${BOLD}cbwsh${NC}"
    echo ""
    echo "Quick Start:"
    echo "  - Press Ctrl+? or F1 for help"
    echo "  - Press Ctrl+A to toggle AI assist mode"
    echo "  - Press Ctrl+Q to quit"
    echo ""
    echo "Configuration: ~/.cbwsh/config.yaml"
    echo ""
    echo "For more information, visit:"
    echo "  https://github.com/${GITHUB_REPO}"
    echo ""
}

# Cleanup
cleanup() {
    if [[ -n "${TMP_DIR:-}" ]] && [[ -d "$TMP_DIR" ]]; then
        rm -rf "$TMP_DIR"
    fi
}

# Set up trap for cleanup on exit
trap cleanup EXIT

# Main installation function
main() {
    show_banner
    parse_args "$@"
    check_dependencies
    detect_platform
    
    if [[ "$VERSION" == "latest" ]]; then
        get_latest_version
    fi

    download_binary
    install_binary
    create_config
    verify_installation
    print_instructions
}

# Run main function
main "$@"
