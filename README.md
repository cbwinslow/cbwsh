# cbwsh - Custom Bubble Tea Shell

A modern, modular terminal shell built with the complete [Bubble Tea](https://github.com/charmbracelet/bubbletea) ecosystem. Features rich TUI components, animations, syntax highlighting, AI integration, and more.

## Features

- 🐚 **Multi-shell support**: Execute bash and zsh commands
- 📊 **Multiple panes**: Split terminal with different layouts (single, horizontal, vertical, grid)
- 🔌 **Plugin system**: Extensible architecture with command, UI, hook, and formatter plugins
- 🔐 **Secrets manager**: Encrypted storage using AES-256-GCM with Argon2id key derivation
- 📡 **SSH manager**: Manage and connect to SSH hosts with key authentication
- 🤖 **AI tools**: Integrated AI agents for command suggestions, explanations, and error fixes
- 📈 **Progress bars**: Beautiful progress indicators using Bubble Tea
- 📝 **Markdown rendering**: Rich markdown display with glamour
- ✨ **Animations**: Smooth spring-physics animations with harmonica
- 🎨 **Syntax highlighting**: Shell command highlighting with chroma
- ⌨️ **Autocompletion**: Smart command, file, and history completion
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
    ├── secrets/        # Encrypted secrets storage
    ├── ssh/            # SSH connection management
    ├── ai/             # AI agents and tools
    └── ui/
        ├── progress/   # Progress bar component
        ├── markdown/   # Markdown renderer
        ├── animation/  # Harmonica animations
        ├── autocomplete/ # Autocompletion
        ├── highlight/  # Syntax highlighting
        └── styles/     # UI styles and themes
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
  provider: none
  enable_suggestions: true

keybindings:
  quit: ctrl+q
  help: ctrl+?
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

