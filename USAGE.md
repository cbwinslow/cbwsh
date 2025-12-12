# cbwsh Usage Guide

## Table of Contents
- [Getting Started](#getting-started)
- [Basic Usage](#basic-usage)
- [Command-Line Options](#command-line-options)
- [Configuration](#configuration)
- [Key Bindings](#key-bindings)
- [Built-in Commands](#built-in-commands)
- [AI Features](#ai-features)
- [Advanced Features](#advanced-features)
- [Troubleshooting](#troubleshooting)

## Getting Started

After installation, simply run:

```bash
cbwsh
```

On first run, cbwsh will create a default configuration file at `~/.cbwsh/config.yaml`.

### Quick Start

1. **Start cbwsh**:
   ```bash
   cbwsh
   ```

2. **Execute a command**:
   Type any shell command and press Enter:
   ```
   ls -la
   ```

3. **Get help**:
   Press `Ctrl+?` or `F1` to see the help screen.

4. **Exit**:
   Press `Ctrl+Q` or type `exit`.

## Basic Usage

### Running Commands

cbwsh works like a regular shell. Type commands and press Enter to execute:

```bash
# List files
ls -la

# Change directory
cd /path/to/directory

# Run any shell command
npm install
python script.py
```

### Command History

- **Up/Down arrows**: Navigate through command history
- **Ctrl+R**: Search command history (future feature)
- History is automatically saved to `~/.cbwsh/history`

### Autocompletion

Press `Tab` to:
- Complete command names
- Complete file paths
- Complete environment variables
- See available completions

## Command-Line Options

cbwsh supports several command-line flags:

### Version Information

```bash
cbwsh --version
```

Output:
```
cbwsh version dev
  commit: unknown
  built:  unknown
```

### Help

```bash
cbwsh --help
```

Shows all available flags and key bindings.

### Custom Configuration

```bash
cbwsh --config /path/to/config.yaml
```

Use a custom configuration file instead of the default `~/.cbwsh/config.yaml`.

## Configuration

cbwsh is configured via a YAML file located at `~/.cbwsh/config.yaml`.

### Default Configuration

The default configuration is created automatically:

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
  provider: none  # Options: none, openai, anthropic, gemini, ollama, local
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
```

### Configuration Options

#### Shell Configuration

```yaml
shell:
  default_shell: bash          # or zsh
  history_size: 10000          # Number of commands to keep
  history_path: ~/.cbwsh/history
  startup_commands:            # Commands to run on startup
    - echo "Welcome to cbwsh!"
  environment:                 # Additional environment variables
    CUSTOM_VAR: value
  aliases:                     # Command aliases
    ll: ls -la
    gs: git status
```

#### UI Configuration

```yaml
ui:
  theme: default               # default, dracula, nord
  layout: single               # single, horizontal, vertical, grid
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true
  highlight_theme: monokai
  markdown_theme: dark
```

#### AI Configuration

```yaml
ai:
  provider: ollama             # none, openai, anthropic, gemini, ollama, local
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true      # Enable AI shell monitoring
  monitoring_interval: 30      # Seconds between recommendations
  enable_suggestions: true
```

## Key Bindings

### Essential Keys

| Key | Action |
|-----|--------|
| `Ctrl+Q` | Quit cbwsh |
| `Ctrl+C` | Cancel current command |
| `Enter` | Execute command |
| `Ctrl+?` or `F1` | Show help |

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate command history |
| `Tab` | Autocomplete |
| `Ctrl+L` | Clear screen |

### Pane Management

| Key | Action |
|-----|--------|
| `Ctrl+N` | Create new pane |
| `Ctrl+W` | Close current pane |
| `Ctrl+]` | Next pane |
| `Ctrl+[` | Previous pane |
| `Ctrl+\` | Split vertical |
| `Ctrl+-` | Split horizontal |

### AI Features

| Key | Action |
|-----|--------|
| `Ctrl+A` | Toggle AI assist mode |
| `Ctrl+M` | Toggle AI monitor pane |

### UI Controls

| Key | Action |
|-----|--------|
| `Ctrl+B` | Toggle sidebar |
| `F10` or `Alt+M` | Toggle menu bar |

## Built-in Commands

cbwsh provides several built-in commands:

### Directory Navigation

```bash
cd /path/to/directory     # Change directory
```

### Screen Management

```bash
clear                     # Clear the screen
```

### Process Management

```bash
jobs                      # List background jobs
```

### System Information

```bash
whoami                    # Show current user
```

### Exit

```bash
exit                      # Exit cbwsh
quit                      # Exit cbwsh
```

## AI Features

cbwsh integrates AI to enhance your shell experience.

### Setting Up AI (Ollama - Recommended)

1. **Install Ollama**:
   ```bash
   curl https://ollama.ai/install.sh | sh
   ```

2. **Pull a model**:
   ```bash
   ollama pull llama2
   # Or use a smaller model
   ollama pull phi3
   ```

3. **Configure cbwsh**:
   Edit `~/.cbwsh/config.yaml`:
   ```yaml
   ai:
     provider: ollama
     ollama_url: http://localhost:11434
     ollama_model: llama2
     enable_monitoring: true
     monitoring_interval: 30
   ```

4. **Start cbwsh and enable AI monitor**:
   ```bash
   cbwsh
   ```
   Press `Ctrl+M` to toggle the AI monitor pane.

### AI Monitor

The AI monitor watches your shell activity and provides:
- Error analysis and fixes
- Command suggestions
- Performance tips
- Pattern detection
- Alias recommendations

### AI Assist Mode

Press `Ctrl+A` to enter AI assist mode where you can:
- Ask questions about commands
- Get command explanations
- Request command suggestions
- Natural language to command translation

## Advanced Features

### Multiple Panes

Create split layouts for multitasking:

1. **Vertical split**: `Ctrl+\`
2. **Horizontal split**: `Ctrl+-`
3. **Navigate panes**: `Ctrl+]` (next) or `Ctrl+[` (previous)
4. **Close pane**: `Ctrl+W`

### Aliases

Define aliases in your config:

```yaml
shell:
  aliases:
    ll: ls -la
    gs: git status
    gp: git push
    gc: git commit -m
```

Usage:
```bash
ll          # Executes: ls -la
gs          # Executes: git status
```

### Environment Variables

Set custom environment variables:

```yaml
shell:
  environment:
    EDITOR: vim
    CUSTOM_PATH: /path/to/custom
```

### Themes

Change the color theme:

```yaml
ui:
  theme: dracula    # or 'nord', 'default'
```

### Syntax Highlighting

Enable/disable syntax highlighting:

```yaml
ui:
  syntax_highlighting: true
  highlight_theme: monokai
```

## Troubleshooting

### cbwsh won't start

**Check if cbwsh is in your PATH**:
```bash
which cbwsh
```

If not found, add the installation directory to your PATH:
```bash
export PATH="/usr/local/bin:$PATH"
```

### Configuration errors

**Validate your configuration**:
```bash
cbwsh --config ~/.cbwsh/config.yaml
```

If there are errors, cbwsh will report them on startup.

**Reset to default configuration**:
```bash
rm ~/.cbwsh/config.yaml
cbwsh  # Will create a new default config
```

### Commands not working

**Check which shell is being used**:
The status bar shows the current shell type (bash/zsh).

**Check your default shell configuration**:
```yaml
shell:
  default_shell: bash  # or zsh
```

### AI features not working

**For Ollama**:
1. Check if Ollama is running:
   ```bash
   curl http://localhost:11434/api/version
   ```

2. Check if the model is available:
   ```bash
   ollama list
   ```

3. Verify configuration:
   ```yaml
   ai:
     provider: ollama
     ollama_url: http://localhost:11434
     ollama_model: llama2  # Must match installed model
   ```

### Performance issues

**Disable animations**:
```yaml
ui:
  enable_animations: false
```

**Reduce history size**:
```yaml
shell:
  history_size: 1000  # Default is 10000
```

**Disable AI monitoring**:
```yaml
ai:
  enable_monitoring: false
```

### Getting help

1. **In-app help**: Press `Ctrl+?` or `F1`
2. **Check documentation**: https://github.com/cbwinslow/cbwsh
3. **Report issues**: https://github.com/cbwinslow/cbwsh/issues

## Examples

### Basic Command Execution

```bash
# Start cbwsh
cbwsh

# Run commands
$ ls -la
$ pwd
$ echo "Hello, cbwsh!"

# Navigate history
# Press ↑ to see previous commands
# Press ↓ to go forward in history

# Exit
$ exit
```

### Using Multiple Panes

```bash
# Start cbwsh
cbwsh

# Split vertically (Ctrl+\)
# Now you have two panes side by side

# In left pane:
$ tail -f /var/log/syslog

# Switch to right pane (Ctrl+])
$ htop

# Close a pane (Ctrl+W)
```

### Working with AI

```bash
# Start cbwsh with AI enabled
cbwsh

# Toggle AI monitor (Ctrl+M)
# The AI monitor appears on the right

# Run some commands
$ git status
$ npm install
$ python script.py

# AI monitor will provide suggestions and tips

# Enter AI assist mode (Ctrl+A)
# Ask questions in natural language
```

### Custom Configuration

Create `~/.cbwsh/config.yaml`:

```yaml
shell:
  default_shell: zsh
  aliases:
    ll: ls -laGh
    gs: git status -sb
    gd: git diff
    
ui:
  theme: dracula
  syntax_highlighting: true
  
ai:
  provider: ollama
  ollama_model: codellama
  enable_monitoring: true
```

Then start cbwsh:
```bash
cbwsh --config ~/.cbwsh/config.yaml
```

## Best Practices

1. **Start with defaults**: Use the default configuration first, then customize as needed.

2. **Use aliases**: Create aliases for frequently used commands.

3. **Enable AI monitoring**: Let AI help you improve your workflow.

4. **Use panes wisely**: Split when you need to monitor multiple things.

5. **Keep history**: The larger the history, the better the autocomplete.

6. **Backup your config**: Keep a backup of your customized configuration.

7. **Update regularly**: Keep cbwsh updated for the latest features and fixes.

## Next Steps

- Explore the [README](README.md) for feature details
- Check the [TODO](TODO.md) for upcoming features
- Review the [ROADMAP](ROADMAP.md) for the project direction
- Contribute to the project on [GitHub](https://github.com/cbwinslow/cbwsh)
