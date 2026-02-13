# cbwsh 🐚

**A Modern, AI-Powered Terminal Shell Built with Bubble Tea**

cbwsh is a next-generation terminal shell that combines the power of traditional shells (bash/zsh) with modern TUI components, AI assistance, and advanced features for developers and power users.

[![Build Status](https://github.com/cbwinslow/cbwsh/workflows/Build/badge.svg)](https://github.com/cbwinslow/cbwsh/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbwinslow/cbwsh)](https://goreportcard.com/report/github.com/cbwinslow/cbwsh)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![cbwsh Demo](https://via.placeholder.com/800x400?text=cbwsh+Terminal+Shell)

## ✨ Features

### 🎨 Modern Terminal UI
- **Multi-pane Support** - Split your terminal horizontally, vertically, or in a grid layout
- **Figma-inspired Design System** - Beautiful, consistent UI with customizable themes
- **Syntax Highlighting** - Shell commands highlighted in real-time
- **Visual Effects** - Water ripples, fluid simulations, spring animations
- **Markdown Rendering** - View documentation with rich formatting
- **Menu Bar** - File, Edit, View, and Help menus for easy navigation

### 🤖 AI Integration
- **Multiple AI Providers** - OpenAI, Anthropic Claude, Google Gemini, Local LLMs (Ollama)
- **Natural Language Commands** - Ask AI to generate commands for you
- **Error Fix Suggestions** - AI analyzes failed commands and suggests fixes
- **AI Monitor Pane** - Real-time activity tracking with contextual recommendations
- **Command Explanations** - Learn what commands do before running them
- **Smart Suggestions** - Context-aware command completion

### 🔐 Security & Secrets
- **Encrypted Secrets Storage** - AES-256-GCM encryption with Argon2id key derivation
- **SSH Key Management** - Store and manage SSH keys securely
- **API Key Storage** - Safely store API keys for various services
- **Multiple Encryption Backends** - Support for Age and GPG

### 🌐 SSH & Remote Access
- **SSH Connection Manager** - Save and manage SSH connections
- **Key-based Authentication** - Secure SSH with public key authentication
- **Port Forwarding** - Local and remote port forwarding support
- **Known Hosts Management** - Security through host key verification

### 🎯 Developer Features
- **Git Integration** - View status, branches, and diffs
- **Command History** - Persistent, searchable command history
- **Smart Autocompletion** - File paths, commands, environment variables
- **Job Control** - Manage background jobs
- **Plugin System** - Extensible architecture for custom functionality
- **Custom Themes** - Default, Dracula, Nord, and more

### ⚡ Performance
- **Lightweight** - Built with Go for speed and efficiency
- **Concurrent Panes** - Each pane runs independently
- **Optimized Rendering** - Smooth animations and responsive UI

## 📦 Installation

### Quick Install (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

### Quick Install (Windows)

```powershell
irm https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1 | iex
```

### Install with Go

```bash
go install github.com/cbwinslow/cbwsh@latest
```

### Build from Source

```bash
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
make build
sudo make install
```

For detailed installation instructions, see [INSTALL.md](INSTALL.md).

## 🚀 Quick Start

### Starting cbwsh

Simply run:

```bash
cbwsh
```

### Basic Key Bindings

| Key | Action |
|-----|--------|
| `Ctrl+Q` | Quit cbwsh |
| `Ctrl+?` | Show help |
| `Ctrl+A` | AI assist mode |
| `Ctrl+M` | Toggle AI monitor |
| `Ctrl+N` | New pane |
| `Ctrl+]` | Next pane |
| `Ctrl+\` | Split vertical |
| `Ctrl+-` | Split horizontal |
| `Tab` | Autocomplete |
| `↑/↓` | Command history |

For complete key bindings, see [USAGE.md](USAGE.md).

### Using AI Features

1. **Install Ollama** (recommended for privacy):
```bash
curl https://ollama.ai/install.sh | sh
ollama pull llama2
```

2. **Configure AI in cbwsh**:
Create `~/.cbwsh/config.yaml`:
```yaml
ai:
  provider: ollama
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true
```

3. **Enable AI Monitor**:
Press `Ctrl+M` to open the AI monitor pane

4. **Use AI Assist**:
Press `Ctrl+A` and ask: *"How do I find large files?"*

## 🎨 Multi-Pane Workflow

cbwsh supports multiple panes for efficient multitasking:

```bash
# Start cbwsh
cbwsh

# Split vertically (Ctrl+\)
# Now you have left and right panes

# In left pane: Edit code
vim main.go

# Switch to right pane (Ctrl+])
# Split horizontally (Ctrl+-)

# In right pane top: Build/test
go build && go test

# In right pane bottom: Monitor
watch -n 2 git status
```

### Pane Layouts

- **Single** - One full-screen pane (default)
- **Horizontal Split** - Two panes side by side
- **Vertical Split** - Two panes stacked
- **Grid** - Four panes in a 2x2 grid
- **Custom** - Create your own layout

## 📚 Documentation

- **[USAGE.md](USAGE.md)** - Comprehensive usage guide
- **[INSTALL.md](INSTALL.md)** - Installation instructions
- **[DESIGN_SYSTEM.md](DESIGN_SYSTEM.md)** - UI design principles
- **[INTEGRATION.md](INTEGRATION.md)** - Integration guides
- **[TODO.md](TODO.md)** - Feature roadmap
- **[AGENTS.md](AGENTS.md)** - AI agent configuration

## ⚙️ Configuration

cbwsh uses YAML configuration files. Create `~/.cbwsh/config.yaml`:

```yaml
# Shell settings
shell:
  default_shell: bash          # or zsh
  history_size: 10000
  history_path: ~/.cbwsh/history

# UI settings
ui:
  theme: default              # default, dracula, nord
  layout: single              # single, horizontal, vertical, grid
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true

# AI settings
ai:
  provider: ollama            # none, ollama, openai, anthropic, gemini
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true
  monitoring_interval: 30     # seconds

# SSH settings
ssh:
  default_user: your-username
  key_path: ~/.ssh/id_rsa
  known_hosts: ~/.ssh/known_hosts

# Secrets settings
secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM
```

## 🤝 AI Agents for Code Review

cbwsh integrates with multiple AI agents for automated code review and assistance:

- **CodeRabbit** - Automated PR reviews
- **GitHub Copilot** - Code suggestions
- **OpenAI Codex** - Code generation
- **Google Gemini** - Multimodal AI
- **Anthropic Claude** - Long-context analysis

See [AGENTS.md](AGENTS.md) for configuration details.

## 🛠️ Development

### Prerequisites

- Go 1.24 or later
- Make
- Git

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run linter
make lint

# Install locally
make install
```

### Project Structure

```
cbwsh/
├── cmd/            # Command-line interface
├── internal/       # Internal application code
│   └── app/        # Main application logic
├── pkg/            # Public packages
│   ├── ai/         # AI integration
│   ├── config/     # Configuration management
│   ├── panes/      # Pane management
│   ├── secrets/    # Secrets storage
│   ├── shell/      # Shell execution
│   ├── ssh/        # SSH management
│   └── ui/         # UI components
├── test/           # Integration tests
└── examples/       # Example applications
```

## 🎯 Use Cases

### For Developers
- Split panes for code editing, building, and testing
- AI-assisted debugging and error fixing
- Git operations with visual feedback
- SSH into remote servers with saved configurations

### For DevOps
- Monitor multiple servers in different panes
- Execute commands across multiple machines
- Track system metrics with AI analysis
- Manage secrets and credentials securely

### For Data Scientists
- Run long-running jobs with progress tracking
- AI-powered command suggestions
- Manage Python/R environments
- SSH tunnel management for Jupyter notebooks

### For System Administrators
- Multi-server management
- Automated task scheduling
- Security audit with AI insights
- Centralized secret management

## 🌟 Inspiration

cbwsh draws inspiration from:
- **[Warp](https://www.warp.dev/)** - AI-powered terminal
- **[Fig](https://fig.io/)** - Autocomplete and AI features
- **[tmux](https://github.com/tmux/tmux)** - Terminal multiplexing
- **[Zellij](https://zellij.dev/)** - Modern workspace
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - TUI framework

## 📝 License

cbwsh is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with the amazing [Bubble Tea ecosystem](https://github.com/charmbracelet):
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [harmonica](https://github.com/charmbracelet/harmonica) - Physics animations

## 🐛 Bug Reports & Feature Requests

Found a bug or have a feature request? Please open an issue on [GitHub Issues](https://github.com/cbwinslow/cbwsh/issues).

## 💬 Community

- **GitHub Discussions**: [cbwinslow/cbwsh/discussions](https://github.com/cbwinslow/cbwsh/discussions)
- **Twitter**: [@cbwinslow](https://twitter.com/cbwinslow)

## 🚀 Roadmap

See [TODO.md](TODO.md) for planned features and improvements.

---

**Made with ❤️ and ☕ by the cbwsh team**
