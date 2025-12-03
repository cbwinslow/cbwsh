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
    └── ui/
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

```bash
go install github.com/cbwinslow/cbwsh@latest
```

Or build from source:

```bash
git clone https://github.com/cbwinslow/cbwsh.git
cd cbwsh
go build -o cbwsh .
./cbwsh
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
```

## License

MIT

---

_Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the Charm ecosystem_

