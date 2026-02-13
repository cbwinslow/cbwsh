# cbwsh Shell Requirements & Dependencies

Complete reference for all requirements needed to build, install, and run cbwsh.

## Table of Contents

1. [System Requirements](#system-requirements)
2. [Build Dependencies](#build-dependencies)
3. [Runtime Dependencies](#runtime-dependencies)
4. [Bubble Tea Ecosystem](#bubble-tea-ecosystem)
5. [Optional Dependencies](#optional-dependencies)
6. [Platform-Specific Requirements](#platform-specific-requirements)
7. [Feature Compatibility Matrix](#feature-compatibility-matrix)
8. [Plugin System Requirements](#plugin-system-requirements)

## System Requirements

### Minimum Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **OS** | Linux 3.10+, macOS 10.13+, Windows 10, FreeBSD 12+ | Latest stable version |
| **CPU** | x86_64, ARM64 | Any modern CPU |
| **RAM** | 50 MB | 100 MB |
| **Disk Space** | 25 MB | 100 MB (with plugins) |
| **Terminal** | Any with Unicode support | True color support (24-bit) |

### Supported Platforms

✅ **Fully Supported:**
- Linux (x86_64, ARM64, ARM, 386)
- macOS (Intel, Apple Silicon)
- FreeBSD (x86_64)
- Windows 10/11 (x86_64, ARM64)

⚠️ **Experimental:**
- OpenBSD
- NetBSD
- Solaris/illumos

## Build Dependencies

### Required for Building

```bash
# Minimum Go version
Go 1.24 or later

# Build tools
make
git

# Optional but recommended
golangci-lint  # For linting
goreleaser     # For releases
```

### Installation Commands

**Ubuntu/Debian:**
```bash
# Install Go
sudo apt-get update
sudo apt-get install golang-go make git

# Verify Go version
go version  # Should be 1.24+
```

**macOS (Homebrew):**
```bash
brew install go make git
brew install golangci-lint  # Optional
```

**From source:**
```bash
# Download Go
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

## Runtime Dependencies

### Core Dependencies

cbwsh has **no external runtime dependencies** for basic functionality. Everything is statically compiled into the binary.

### System Components (Recommended)

| Component | Purpose | Installation |
|-----------|---------|--------------|
| **bash or zsh** | Command execution | Usually pre-installed |
| **terminal emulator** | Running cbwsh | iTerm2, Alacritty, Windows Terminal |
| **UTF-8 locale** | Unicode support | Usually configured |

### Verifying Runtime Environment

```bash
# Check shell
echo $SHELL

# Check locale
locale

# Check terminal capabilities
echo $TERM
tput colors  # Should output 256 or more
```

## Bubble Tea Ecosystem

cbwsh is built on the excellent Charm Bubble Tea ecosystem.

### Core Libraries

All these dependencies are managed by Go modules and compiled into the binary:

```go
// TUI Framework
github.com/charmbracelet/bubbletea     // v1.3.10
  - Terminal UI framework
  - Event handling and rendering
  - Cross-platform terminal support

// UI Components
github.com/charmbracelet/bubbles       // v0.21.1
  - List, table, progress bars
  - Text input, viewport
  - Spinners and other widgets

// Styling
github.com/charmbracelet/lipgloss      // v1.1.1
  - CSS-like styling for terminal
  - Layout and positioning
  - Color and text formatting

// Markdown Rendering
github.com/charmbracelet/glamour       // v0.10.0
  - Markdown rendering in terminal
  - Syntax highlighting for code blocks
  - Theme support

// Animations
github.com/charmbracelet/harmonica     // v0.2.0
  - Spring physics animations
  - Smooth transitions
  - Easing functions
```

### Additional Charm Tools

```go
// Terminal Utilities
github.com/muesli/termenv             // v0.16.0
  - Terminal capability detection
  - Color profile detection
  
github.com/muesli/reflow              // v0.3.0
  - Text wrapping and formatting
  
// Syntax Highlighting
github.com/alecthomas/chroma/v2       // v2.20.0
  - Code syntax highlighting
  - Multiple language support
```

### Complete Dependency Tree

Run this to see all dependencies:

```bash
cd cbwsh
go mod graph | grep charmbracelet
```

Current full tree:

```
bubbletea
├── ansi
├── colorprofile
├── x/ansi
├── x/cellbuf
├── x/term
└── x/input

bubbles
├── list
├── table
├── progress
├── textinput
├── textarea
└── viewport

lipgloss
├── color
└── border

glamour
├── chroma (syntax highlighting)
└── goldmark (markdown parsing)

harmonica (spring animations)
```

## Optional Dependencies

### AI Features

| Provider | Requirement | Installation |
|----------|-------------|--------------|
| **Ollama** (Recommended) | Local server | `curl https://ollama.ai/install.sh \| sh` |
| **OpenAI** | API key | Sign up at platform.openai.com |
| **Anthropic** | API key | Sign up at console.anthropic.com |
| **Google Gemini** | API key | Sign up at makersuite.google.com |

**Ollama Model Installation:**
```bash
# Install Ollama
curl https://ollama.ai/install.sh | sh

# Download models
ollama pull llama2          # General purpose (3.8GB)
ollama pull codellama       # Code-focused (3.8GB)
ollama pull mistral         # Fast & efficient (4.1GB)
ollama pull deepseek-coder  # Advanced code (6.7GB)

# Verify installation
ollama list
```

### Git Integration

```bash
# Git is optional but recommended for VCS features
git --version  # Should be 2.0+
```

### SSH Features

```bash
# SSH client (usually pre-installed)
ssh -V

# Generate SSH keys if needed
ssh-keygen -t ed25519 -C "your_email@example.com"
```

### Secrets Encryption

**Default (AES-256-GCM):**
- No external dependencies, built-in

**Optional backends:**
```bash
# Age encryption
brew install age      # macOS
apt install age       # Debian/Ubuntu

# GPG encryption
brew install gnupg    # macOS
apt install gnupg     # Debian/Ubuntu
```

## Platform-Specific Requirements

### Linux

**Additional packages for full functionality:**
```bash
# Debian/Ubuntu
sudo apt-get install -y \
    ca-certificates \
    curl \
    git \
    ssh-client

# Arch Linux
sudo pacman -S ca-certificates curl git openssh

# Fedora/RHEL
sudo dnf install ca-certificates curl git openssh-clients
```

### macOS

**Xcode Command Line Tools:**
```bash
xcode-select --install
```

**Recommended terminal:**
- iTerm2: https://iterm2.com/
- Alacritty: `brew install alacritty`
- Warp: https://www.warp.dev/

### Windows

**Requirements:**
- Windows Terminal (recommended): Install from Microsoft Store
- PowerShell 7+ or WSL2

**WSL2 Setup (Recommended):**
```powershell
# Enable WSL
wsl --install

# Install Ubuntu
wsl --install -d Ubuntu

# Run cbwsh in WSL
wsl
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

**Native Windows:**
```powershell
# Install using install.ps1
iwr https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1 | iex
```

### FreeBSD

```bash
# Install packages
pkg install go git gmake

# Build from source
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
gmake build
sudo gmake install
```

## Feature Compatibility Matrix

| Feature | Linux | macOS | Windows | FreeBSD | Notes |
|---------|-------|-------|---------|---------|-------|
| **Core Shell** | ✅ | ✅ | ✅ | ✅ | Full support |
| **Multi-pane** | ✅ | ✅ | ✅ | ✅ | Full support |
| **Syntax Highlighting** | ✅ | ✅ | ✅ | ✅ | Full support |
| **Themes** | ✅ | ✅ | ✅ | ✅ | Full support |
| **AI Integration** | ✅ | ✅ | ✅ | ✅ | Requires AI provider |
| **SSH** | ✅ | ✅ | ⚠️ | ✅ | Windows: Use WSL |
| **Job Control** | ✅ | ✅ | ⚠️ | ✅ | Windows: Limited |
| **File Permissions** | ✅ | ✅ | ⚠️ | ✅ | Windows: Limited |
| **Plugins** | ✅ | ✅ | ✅ | ✅ | Full support |
| **Shell Integration** | ✅ | ✅ | ⚠️ | ✅ | Windows: PowerShell only |

**Legend:**
- ✅ Full support
- ⚠️ Partial support or limitations
- ❌ Not supported

## Plugin System Requirements

### Plugin Development

**Required:**
- Go 1.24+
- cbwsh plugin SDK

**Plugin structure:**
```go
package main

import "github.com/cbwinslow/cbwsh/pkg/plugins"

func main() {
    plugins.Register(&MyPlugin{})
}

type MyPlugin struct{}

func (p *MyPlugin) Name() string { return "my-plugin" }
func (p *MyPlugin) Version() string { return "1.0.0" }
func (p *MyPlugin) Init() error { return nil }
```

### Plugin Installation

```bash
# Install plugin
cbwsh plugin install <plugin-url>

# List installed plugins
cbwsh plugin list

# Enable/disable
cbwsh plugin enable <name>
cbwsh plugin disable <name>
```

### Plugin Locations

Following XDG Base Directory specification:

```bash
# System plugins
/usr/local/share/cbwsh/plugins/

# User plugins
~/.local/share/cbwsh/plugins/

# Custom location (set in config)
$CBWSH_PLUGIN_DIR
```

## Verification Checklist

Use this checklist to verify your environment:

```bash
# ✓ System
uname -a
echo $TERM

# ✓ Go (for building)
go version

# ✓ Runtime
bash --version
which cbwsh

# ✓ AI (optional)
ollama --version
ollama list

# ✓ Git (optional)
git --version

# ✓ SSH (optional)
ssh -V

# ✓ Directories
echo $XDG_CONFIG_HOME
echo $XDG_DATA_HOME
echo $XDG_CACHE_HOME

# ✓ cbwsh
cbwsh --version
cbwsh --help
```

## Troubleshooting

### Common Issues

**Terminal doesn't support colors:**
```bash
export TERM=xterm-256color
# Or use a modern terminal emulator
```

**Go version too old:**
```bash
# Download newer Go version
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**cbwsh not in PATH:**
```bash
# Add to PATH
export PATH="$HOME/.local/bin:$PATH"

# Make permanent
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## Getting Help

- **Documentation**: https://github.com/cbwinslow/cbwsh
- **Issues**: https://github.com/cbwinslow/cbwsh/issues
- **Discussions**: https://github.com/cbwinslow/cbwsh/discussions
- **Discord**: [Community Server]

## Version History

| Version | Go Min | Notable Changes |
|---------|--------|-----------------|
| v1.0.0 | 1.24 | Initial release |
| main | 1.24 | Development branch |

---

**Last Updated**: 2026-02-13
**cbwsh Version**: dev
**Go Version**: 1.24+
