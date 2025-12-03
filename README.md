# cbwsh - Custom Bubble Tea Shell

A modern, modular terminal shell built with the complete [Bubble Tea](https://github.com/charmbracelet/bubbletea) ecosystem. Features rich TUI components, animations, syntax highlighting, AI integration, visual effects, and more.

## Features

### Core Shell Features
- 🐚 **Multi-shell support**: Execute bash and zsh commands
- 📊 **Multiple panes**: Split terminal with different layouts (single, horizontal, vertical, grid)
- 🔌 **Plugin system**: Extensible architecture with command, UI, hook, and formatter plugins
- ⌨️ **Autocompletion**: Smart command, file, environment variable, and history completion

### Security & Secrets
- 🔐 **Secrets manager**: Encrypted storage using AES-256-GCM with Argon2id key derivation
- 🔑 **Multi-backend encryption**: Support for AES, age, and GPG encryption
- 📂 **Git integration**: Sync secrets with git or yadm (Yet Another Dotfiles Manager)
- 🔒 **API key management**: Specialized management for API keys

### AI Integration
- 🤖 **AI agents**: Integrated AI agents for command suggestions, explanations, and error fixes
- 🌐 **Multiple providers**: Support for OpenAI, Anthropic, Gemini, and local LLMs
- 🗣️ **A2A Protocol**: Agent-to-agent communication for building complex AI workflows
- 💬 **AI Chat Pane**: Resizable, configurable chat interface for AI conversations
- 📝 **Markdown rendering**: AI responses rendered with full markdown support

### SSH & Remote
- 📡 **SSH manager**: Manage and connect to SSH hosts with key/password authentication
- 🔗 **Port forwarding**: Local port forwarding support
- 📋 **Host management**: Save and recall SSH host configurations

### UI Components
- 📈 **Progress bars**: Beautiful progress indicators using Bubble Tea
- 📝 **Markdown editor**: Edit markdown files with live preview
- ✨ **Visual effects**: Water waves, fluid dynamics, and particle systems using harmonica
- 🎨 **Syntax highlighting**: Shell command highlighting with chroma
- 🎭 **Themes**: Multiple color themes (default, dracula, nord)

### Terminal Features
- 📋 **Menu bar**: Standard File, Edit, View, Help menus with keyboard shortcuts
- 📝 **Command history**: Persistent history with search
- ⚙️ **Job control**: Background job management (start, stop, resume, kill)
- 📊 **Logging**: Structured logging with multiple levels and formats
- 🔑 **Privilege management**: Sudo/su integration and privilege elevation
- 🖥️ **POSIX signals**: Signal handling and process management

## Architecture

The project follows a modular architecture with complete abstraction:

```
cbwsh/
├── cmd/cbwsh/          # Main entry point
├── internal/app/       # Main application logic
└── pkg/
    ├── core/           # Core types, interfaces, and enums
    ├── config/         # Configuration management
    ├── shell/          # Shell executor (bash/zsh)
    ├── panes/          # Pane/layout management
    ├── plugins/        # Plugin system
    ├── secrets/        # Encrypted secrets storage (AES, age, GPG)
    ├── ssh/            # SSH connection management
    ├── ai/             # AI agents, A2A protocol, and tools
    ├── logging/        # Structured logging infrastructure
    ├── process/        # Job control and process management
    ├── privileges/     # Privilege checking and elevation
    ├── posix/          # POSIX signals and system calls
    └── ui/
        ├── menu/       # Menu bar component
        ├── progress/   # Progress bar component
        ├── markdown/   # Markdown renderer
        ├── animation/  # Harmonica animations
        ├── autocomplete/ # Autocompletion
        ├── highlight/  # Syntax highlighting
        ├── styles/     # UI styles and themes
        ├── effects/    # Visual effects (water, fluid, particles)
        ├── aichat/     # AI chat pane component
        └── editor/     # Markdown editor component
```

## Dependencies

Built with the complete Charm ecosystem:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Harmonica](https://github.com/charmbracelet/harmonica) - Smooth animations
- [Huh](https://github.com/charmbracelet/huh) - Form components
- [Log](https://github.com/charmbracelet/log) - Logging
- [Wish](https://github.com/charmbracelet/wish) - SSH server
- [Chroma](https://github.com/alecthomas/chroma) - Syntax highlighting

## Key Bindings

| Key | Action |
|-----|--------|
| Ctrl+Q | Quit |
| Ctrl+C | Cancel current command |
| Enter | Execute command |
| Tab | Autocomplete |
| ↑/↓ | Navigate history |
| Ctrl+L | Clear screen |
| Ctrl+N | New pane |
| Ctrl+W | Close pane |
| Ctrl+] | Next pane |
| Ctrl+[ | Previous pane |
| Ctrl+\ | Split vertical |
| Ctrl+- | Split horizontal |
| Ctrl+B | Toggle sidebar |
| Ctrl+A | AI assist mode |
| Ctrl+? | Help |

### Menu Bar
| Key | Action |
|-----|--------|
| Alt+M / F10 | Toggle menu bar |
| Alt+F | File menu |
| Alt+E | Edit menu |
| Alt+V | View menu |
| Alt+H | Help menu |
| ←/→ | Navigate menus |
| ↑/↓ | Navigate items |
| Enter | Select item |
| Escape | Close menu |

### AI Chat Pane
| Key | Action |
|-----|--------|
| Ctrl+Enter | Send message |
| Ctrl+L | Clear chat |
| Escape | Unfocus pane |
| PageUp/PageDown | Scroll chat |

### Markdown Editor
| Key | Action |
|-----|--------|
| Ctrl+S | Save |
| Ctrl+P | Toggle mode (edit/preview/split) |
| Ctrl+B | Insert bold |
| Ctrl+I | Insert italic |
| Ctrl+K | Insert link |
| Ctrl+H | Insert heading |

## Installation

### Quick Install (Recommended)

#### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash
```

#### Windows (PowerShell)

```powershell
iwr https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.ps1 | iex
```

Or download and run manually:

```powershell
.\install.ps1 -AddToPath
```

### Using Go

```bash
go install github.com/cbwinslow/cbwsh@latest
```

### Build from Source

```bash
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
make install
```

Or without make:

```bash
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
go build -o cbwsh .
./cbwsh
```

### Installation Options

#### Linux/macOS Script Options

```bash
# Install specific version
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --version v1.0.0

# Install to custom location
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --prefix ~/.local/bin

# Install without sudo
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --no-sudo --prefix ~/.local/bin
```

#### Windows PowerShell Options

```powershell
# Install specific version
.\install.ps1 -Version v1.0.0

# Install to custom location
.\install.ps1 -Prefix "C:\Tools\cbwsh"

# Install and add to PATH
.\install.ps1 -AddToPath
```

#### Makefile Targets

```bash
make                  # Build the binary
make install          # Install to /usr/local/bin
make install PREFIX=~/.local  # Install to custom prefix
make uninstall        # Remove installed binary
make test             # Run tests
make lint             # Run linter
make clean            # Clean build artifacts
make cross-compile    # Build for all platforms
make help             # Show all available targets
```

## Configuration

Configuration is stored in `~/.cbwsh/config.yaml`. Example:

```yaml
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
  provider: gemini  # Options: none, openai, anthropic, gemini, local
  api_key: ""       # Store in secrets manager for security
  model: "gemini-pro"
  enable_suggestions: true

secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM  # Options: AES-256-GCM, age, gpg
  key_derivation: argon2id

keybindings:
  quit: ctrl+q
  help: ctrl+?
  ai_assist: ctrl+a
```

## Visual Effects

cbwsh includes several visual effects powered by harmonica physics:

### Water Effect
```go
water := effects.NewWaterEffect(width, height)
water.SetColors(effects.DefaultWaterColors)
water.Update()
output := water.RenderColored()
```

### Fluid Simulation
```go
fluid := effects.NewFluidSimulation(width, height, viscosity, diffusion)
fluid.AddDensity(x, y, amount)
fluid.AddVelocity(x, y, vx, vy)
fluid.Step()
output := fluid.RenderColored()
```

### Particle System
```go
particles := effects.NewParticleSystem(width, height, maxParticles)
particles.SetColors(effects.DefaultFireColors)
particles.Emit(x, y, count, spread, speed)
particles.Update()
output := particles.RenderColored()
```

## Logging

```go
import "github.com/cbwinslow/cbwsh/pkg/logging"

// Create a logger
logger := logging.New(
    logging.WithLevel(logging.LevelDebug),
    logging.WithOutput(os.Stderr),
)

// Log messages
logger.Info("Application started")
logger.Debug("Debug message")
logger.Error("Something went wrong")

// With fields
logger.WithField("user", "john").Info("User logged in")
logger.WithFields(map[string]any{
    "action": "execute",
    "command": "ls -la",
}).Debug("Command executed")
```

## Job Control

```go
import "github.com/cbwinslow/cbwsh/pkg/process"

// Create job manager
manager := process.NewJobManager(100)

// Start a background job
job, err := manager.StartJob(ctx, "sleep 10", "/bin/bash")

// List jobs
for _, job := range manager.ListJobs() {
    fmt.Printf("%s\n", job.String())
}

// Stop a job
manager.StopJob(job.ID)

// Continue a stopped job
manager.ContinueJob(job.ID)

// Kill a job
manager.KillJob(job.ID)
```

## Privilege Management

```go
import "github.com/cbwinslow/cbwsh/pkg/privileges"

// Create privilege manager
manager := privileges.NewManager()

// Check current privileges
if manager.IsRoot() {
    fmt.Println("Running as root")
}

// Check if command requires elevation
if privileges.RequiresElevation("apt update") {
    fmt.Println("Needs sudo")
}

// Execute elevated command
output, err := manager.ExecuteElevated(ctx, "apt update")
```

## POSIX Signals

```go
import "github.com/cbwinslow/cbwsh/pkg/posix"

// Create signal manager
manager := posix.NewSignalManager()

// Register signal handler
manager.RegisterHandler(posix.SIGINT, func(sig posix.Signal) {
    fmt.Printf("Received %s\n", sig)
})

// Start handling signals
manager.Start(posix.SIGINT, posix.SIGTERM)

// Send signal to process
posix.Send(pid, posix.SIGTERM)
```

## AI Integration

### Using AI Agents
```go
// Create an AI agent
agent := ai.NewAgent("assistant", core.AIProviderGemini, apiKey, "gemini-pro")

// Query the agent
response, err := agent.Query(ctx, "How do I list files?")

// Get command suggestions
cmd, err := agent.SuggestCommand(ctx, "find large files")

// Explain a command
explanation, err := agent.ExplainCommand(ctx, "find . -size +100M")
```

### A2A Protocol
```go
// Create A2A router
router := ai.NewA2ARouter()

// Register agents
router.RegisterHandler("shell-assistant", ai.NewShellAssistant(agent))

// Send messages between agents
msg := &ai.A2AMessage{
    Type:    ai.A2AMessageTypeQuery,
    From:    "user",
    To:      "shell-assistant",
    Payload: "Suggest a command to compress files",
}
response, err := router.Send(ctx, msg)
```

## Development

```bash
# Run tests
go test -v ./...

# Run linter
golangci-lint run ./...

# Build
go build -o cbwsh .

# Or use make
make test           # Run tests
make lint           # Run linter
make build          # Build binary
make dev            # Development mode with file watching
```

## Roadmap

See [TODO.md](TODO.md) for a comprehensive list of planned features, including:

- 🤖 **AI Features**: Natural language commands, error fixes, code generation
- 🖥️ **Terminal UI**: Floating panes, tabs, themes, visual effects
- 🐚 **Shell Features**: Block-based input, structured output, job progress
- 🔒 **Security**: Password manager integrations, 2FA, hardware keys
- 📡 **SSH & Remote**: Multi-hop connections, SFTP, Mosh support
- 🎨 **Customization**: Starship-like prompts, custom themes, key bindings
- 🔧 **Integration**: Git, Docker, Kubernetes, cloud providers

We welcome contributions! Check the [TODO.md](TODO.md) for areas where you can help.

## License

MIT

---

_Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the Charm ecosystem_

