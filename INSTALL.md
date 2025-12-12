# cbwsh Installation Guide

This guide provides detailed installation instructions for cbwsh on various platforms.

## Table of Contents

- [System Requirements](#system-requirements)
- [Quick Install](#quick-install)
  - [Linux / macOS](#linux--macos)
  - [Windows](#windows)
- [Install Methods](#install-methods)
  - [Using Installation Script](#using-installation-script)
  - [Using Go](#using-go)
  - [Build from Source](#build-from-source)
  - [Using Package Managers](#using-package-managers)
- [Post-Installation](#post-installation)
- [Uninstalling](#uninstalling)
- [Troubleshooting](#troubleshooting)

## System Requirements

### Minimum Requirements
- **Operating System**: Linux, macOS, FreeBSD, or Windows 10/11
- **Terminal**: Any modern terminal emulator with Unicode support
- **Memory**: 50MB RAM minimum
- **Disk Space**: 25MB for binary and configuration

### Recommended
- **Terminal**: A terminal with true color support (e.g., iTerm2, Alacritty, Windows Terminal)
- **Shell**: bash or zsh installed for command execution
- **Go**: Go 1.24+ if building from source

### Optional Dependencies
- **Ollama**: For local AI features (https://ollama.ai)
- **Git**: For version control integration
- **GPG/Age**: For advanced secrets encryption

## Quick Install

### Linux / macOS

#### Using curl
```bash
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

#### Using wget
```bash
wget -qO- https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

The script will:
1. Detect your operating system and architecture
2. Download the appropriate binary
3. Install cbwsh to `/usr/local/bin` (may require sudo)
4. Create a default configuration file at `~/.cbwsh/config.yaml`
5. Verify the installation

### Windows

#### Using PowerShell (Run as Administrator)
```powershell
iwr https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1 | iex
```

Or download and run manually:
```powershell
# Download the script
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1" -OutFile "install.ps1"

# Run the installer
.\install.ps1 -AddToPath
```

## Install Methods

### Using Installation Script

The installation scripts support several options for customization.

#### Linux/macOS Options

```bash
# Install specific version
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --version v1.0.0

# Install to custom location (no sudo required)
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --prefix ~/.local/bin

# Install without sudo
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --no-sudo --prefix ~/.local/bin

# Get help
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --help
```

**Script Options:**
- `--version <version>` - Install a specific version (default: latest)
- `--prefix <path>` - Install to custom location (default: /usr/local/bin)
- `--no-sudo` - Don't use sudo for installation
- `--help` - Show help message

#### Windows Options

```powershell
# Install specific version
.\install.ps1 -Version v1.0.0

# Install to custom location
.\install.ps1 -Prefix "C:\Tools\cbwsh"

# Install and add to PATH
.\install.ps1 -AddToPath

# Get help
Get-Help .\install.ps1 -Full
```

**Script Parameters:**
- `-Version` - Install a specific version (default: latest)
- `-Prefix` - Install to custom location
- `-AddToPath` - Automatically add installation directory to PATH

### Using Go

If you have Go 1.24 or later installed:

```bash
go install github.com/cbwinslow/cbwsh@latest
```

This will install cbwsh to your `$GOPATH/bin` directory (typically `~/go/bin`).

**Make sure your GOPATH/bin is in your PATH:**

```bash
# Add to your ~/.bashrc, ~/.zshrc, or equivalent
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Build from Source

#### Prerequisites
- Git
- Go 1.24 or later
- Make (optional, but recommended)

#### Using Make

```bash
# Clone the repository
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh

# Build and install (installs to /usr/local/bin)
make install

# Or install to custom location
make install PREFIX=~/.local

# View all available make targets
make help
```

#### Without Make

```bash
# Clone the repository
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh

# Build the binary
go build -o cbwsh .

# Install manually
sudo cp cbwsh /usr/local/bin/

# Or install to user directory (no sudo)
mkdir -p ~/.local/bin
cp cbwsh ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

#### Build Options

You can customize the build with various flags:

```bash
# Build for different platform
GOOS=linux GOARCH=amd64 go build -o cbwsh .

# Build with version information
go build -ldflags "-X main.version=1.0.0" -o cbwsh .

# Build without debug symbols (smaller binary)
go build -ldflags="-s -w" -o cbwsh .

# Cross-compile for all platforms
make cross-compile
```

### Using Package Managers

#### Homebrew (macOS and Linux)

```bash
# Add the tap (coming soon)
brew tap cbwinslow/cbwsh

# Install
brew install cbwsh
```

#### APT (Debian/Ubuntu)

```bash
# Coming soon
# sudo apt install cbwsh
```

#### Snap (Linux)

```bash
# Coming soon
# sudo snap install cbwsh
```

#### Chocolatey (Windows)

```powershell
# Coming soon
# choco install cbwsh
```

## Post-Installation

### Verify Installation

Check that cbwsh is installed correctly:

```bash
# Check if cbwsh is in PATH
which cbwsh
# or on Windows: where cbwsh

# Run cbwsh
cbwsh
```

You should see the cbwsh shell interface. Press `Ctrl+Q` to quit.

### Configuration

cbwsh creates a default configuration file at `~/.cbwsh/config.yaml` on first run or during installation.

#### Review Configuration

```bash
# View configuration
cat ~/.cbwsh/config.yaml

# Edit configuration
nano ~/.cbwsh/config.yaml
# or use your preferred editor
```

#### Configuration Locations

cbwsh looks for configuration in the following order:
1. `~/.cbwsh/config.yaml` (user configuration)
2. `/etc/cbwsh/config.yaml` (system-wide configuration)
3. Default built-in configuration

#### Key Configuration Options

```yaml
shell:
  default_shell: bash        # or zsh
  history_size: 10000
  history_path: ~/.cbwsh/history

ui:
  theme: default            # default, dracula, nord
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true

ai:
  provider: none            # none, ollama, openai, anthropic, gemini
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: false
```

See [README.md](README.md#configuration) for complete configuration documentation.

### Set Up AI Features (Optional)

To use AI features with Ollama (local, privacy-focused):

```bash
# Install Ollama
curl https://ollama.ai/install.sh | sh

# Pull a model
ollama pull llama2
# or a smaller model: ollama pull phi3

# Update cbwsh configuration
cat >> ~/.cbwsh/config.yaml << 'EOF'
ai:
  provider: ollama
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true
  monitoring_interval: 30
EOF

# Start cbwsh
cbwsh
```

Press `Ctrl+M` to toggle the AI monitor pane.

### Adding to Shell Profile (Optional)

To make cbwsh your default shell or add it to your shell profile:

#### As Login Shell

```bash
# Find cbwsh path
which cbwsh

# Add to /etc/shells (requires sudo)
echo "$(which cbwsh)" | sudo tee -a /etc/shells

# Change default shell
chsh -s $(which cbwsh)
```

**Note:** cbwsh is designed as an interactive shell and may not work correctly as a login shell for all use cases. Test thoroughly before setting as default.

#### As Shell Alias

Add to your `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
# Create alias
alias csh='cbwsh'

# Or replace common shell commands
alias sh='cbwsh'
```

## Uninstalling

### Using Installation Script

If you installed using the installation script:

```bash
# Remove binary
sudo rm /usr/local/bin/cbwsh
# or if installed to custom location:
rm ~/.local/bin/cbwsh

# Remove configuration (optional)
rm -rf ~/.cbwsh
```

### Installed via Go

```bash
# Remove binary
rm $(go env GOPATH)/bin/cbwsh

# Remove configuration (optional)
rm -rf ~/.cbwsh
```

### Installed via Make

```bash
# From the source directory
make uninstall

# Or manually
sudo rm /usr/local/bin/cbwsh
```

### Clean Uninstall

To completely remove cbwsh including all configuration:

```bash
# Remove binary
sudo rm /usr/local/bin/cbwsh  # or your installation location

# Remove configuration and data
rm -rf ~/.cbwsh

# If you added cbwsh to /etc/shells, remove it (requires sudo)
sudo sed -i '/cbwsh/d' /etc/shells
```

## Troubleshooting

### Installation Issues

#### "Permission denied" during installation

**Solution:** Use sudo or install to a user directory:

```bash
# Option 1: Use sudo
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | sudo bash

# Option 2: Install to user directory
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --no-sudo --prefix ~/.local/bin

# Make sure ~/.local/bin is in PATH
export PATH="$HOME/.local/bin:$PATH"
```

#### "Command not found" after installation

**Solution:** Add the installation directory to your PATH:

```bash
# Check where cbwsh was installed
which cbwsh
ls -la /usr/local/bin/cbwsh
ls -la ~/.local/bin/cbwsh

# Add to PATH in your shell configuration
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Download fails or times out

**Solutions:**
1. Check your internet connection
2. Try using a VPN if the download is blocked
3. Download the binary manually from GitHub releases:
   ```bash
   # Visit: https://github.com/cbwinslow/cbwsh/releases
   # Download appropriate binary for your platform
   # Extract and install manually
   ```

#### "No suitable binary found" error

**Solution:** Your platform may not be supported by pre-built binaries. Try building from source:

```bash
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
go build -o cbwsh .
./cbwsh
```

### Runtime Issues

#### cbwsh crashes on startup

**Solutions:**
1. Check terminal compatibility:
   ```bash
   echo $TERM
   # Should output something like: xterm-256color
   ```

2. Reset configuration:
   ```bash
   mv ~/.cbwsh/config.yaml ~/.cbwsh/config.yaml.bak
   cbwsh  # Will create new default config
   ```

3. Run with debug logging (coming soon):
   ```bash
   cbwsh --log-level debug
   ```

#### Display issues or garbled characters

**Solutions:**
1. Ensure your terminal supports Unicode:
   ```bash
   locale | grep UTF-8
   ```

2. Use a modern terminal emulator:
   - Linux: Alacritty, Kitty, GNOME Terminal
   - macOS: iTerm2, Alacritty
   - Windows: Windows Terminal, Alacritty

3. Check terminal size:
   ```bash
   tput cols  # Should be at least 80
   tput lines # Should be at least 24
   ```

#### AI features not working

**Solutions:**
1. Verify Ollama is running:
   ```bash
   curl http://localhost:11434/api/tags
   ```

2. Check model is installed:
   ```bash
   ollama list
   ```

3. Verify configuration:
   ```bash
   grep -A 5 "ai:" ~/.cbwsh/config.yaml
   ```

### Getting Help

If you encounter issues not covered here:

1. **Check existing issues:** https://github.com/cbwinslow/cbwsh/issues
2. **Documentation:** https://github.com/cbwinslow/cbwsh
3. **Open an issue:** Include:
   - Operating system and version
   - Installation method
   - Error messages
   - Configuration file (redact sensitive info)
   - Steps to reproduce

---

**Next Steps:** See [USAGE.md](USAGE.md) for a complete guide on using cbwsh.
