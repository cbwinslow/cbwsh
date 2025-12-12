# cbwsh Usage Guide

This guide covers everything you need to know to use cbwsh effectively.

## Table of Contents

- [Getting Started](#getting-started)
- [Basic Usage](#basic-usage)
- [Key Bindings](#key-bindings)
- [Features](#features)
  - [Pane Management](#pane-management)
  - [Command History](#command-history)
  - [Autocompletion](#autocompletion)
  - [AI Features](#ai-features)
  - [Job Control](#job-control)
  - [Secrets Management](#secrets-management)
  - [SSH Management](#ssh-management)
- [Configuration](#configuration)
- [Tips and Tricks](#tips-and-tricks)
- [Examples](#examples)

## Getting Started

### Starting cbwsh

Simply run the `cbwsh` command:

```bash
cbwsh
```

You'll see the cbwsh interface with:
- A header showing "🐚 cbwsh"
- An empty output area
- A command prompt at the bottom
- A status bar (if enabled)

### First Commands

Try these basic commands to get familiar:

```bash
# List files
ls -la

# Check current directory
pwd

# Change directory
cd ~

# View environment
env

# Get help
help
# Or press Ctrl+?
```

### Exiting cbwsh

There are several ways to exit:
- Press `Ctrl+Q`
- Type `exit` or `quit` and press Enter
- Press `Ctrl+C` then `Ctrl+Q`

## Basic Usage

### Running Commands

cbwsh supports all standard shell commands:

```bash
# File operations
ls -la
cat file.txt
touch newfile.txt
rm oldfile.txt

# Directory navigation
cd /path/to/directory
cd ..
cd ~

# Process management
ps aux | grep process
top
kill PID

# Piping and redirection
cat file.txt | grep pattern
echo "text" > file.txt
command 2>&1 | tee output.log
```

### Built-in Commands

cbwsh provides several built-in commands:

| Command | Description | Example |
|---------|-------------|---------|
| `cd <dir>` | Change directory | `cd /home/user` |
| `clear` | Clear output | `clear` |
| `exit` | Exit shell | `exit` |
| `quit` | Exit shell | `quit` |
| `help` | Show help | `help` |
| `jobs` | List background jobs | `jobs` |
| `whoami` | Show current user | `whoami` |

## Key Bindings

### Global Keys

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+Q` | Quit | Exit cbwsh |
| `Ctrl+C` | Cancel | Cancel current command |
| `Ctrl+?` / `F1` | Help | Show/hide help screen |
| `Enter` | Execute | Run the current command |
| `Ctrl+L` | Clear | Clear the output screen |

### Navigation Keys

| Key | Action | Description |
|-----|--------|-------------|
| `↑` | History up | Previous command in history |
| `↓` | History down | Next command in history |
| `Tab` | Autocomplete | Complete command/file/path |
| `←` / `→` | Move cursor | Navigate in input line |

### Pane Management

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+N` | New pane | Create a new shell pane |
| `Ctrl+W` | Close pane | Close the active pane |
| `Ctrl+]` | Next pane | Switch to next pane |
| `Ctrl+[` | Previous pane | Switch to previous pane |
| `Ctrl+\` | Split vertical | Split current pane vertically |
| `Ctrl+-` | Split horizontal | Split current pane horizontally |

### UI Controls

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+B` | Toggle sidebar | Show/hide sidebar |
| `Ctrl+M` | Toggle AI monitor | Show/hide AI monitor pane |
| `F10` / `Alt+M` | Toggle menu | Show/hide menu bar |
| `Alt+F` | File menu | Open File menu |
| `Alt+E` | Edit menu | Open Edit menu |
| `Alt+V` | View menu | Open View menu |
| `Alt+H` | Help menu | Open Help menu |

### AI Features

| Key | Action | Description |
|-----|--------|-------------|
| `Ctrl+A` | AI assist | Toggle AI assist mode |
| `Ctrl+M` | AI monitor | Toggle AI monitor pane |

## Features

### Pane Management

cbwsh supports multiple panes for working on different tasks simultaneously.

#### Creating Panes

```bash
# Press Ctrl+N to create a new pane
# The new pane starts in the current directory
```

#### Splitting Panes

```bash
# Press Ctrl+\ for vertical split
# Press Ctrl+- for horizontal split
```

#### Switching Between Panes

```bash
# Press Ctrl+] for next pane
# Press Ctrl+[ for previous pane
```

#### Closing Panes

```bash
# Press Ctrl+W to close the active pane
# Note: Can't close the last pane
```

### Command History

cbwsh maintains a persistent command history.

#### Navigating History

```bash
# Press ↑ to go to previous command
# Press ↓ to go to next command
# Press Enter to execute the current command
```

#### History Features

- **Persistent**: History is saved between sessions
- **Searchable**: Use ↑/↓ to find previous commands
- **Configurable**: Set `history_size` in config to control history length

#### History Configuration

```yaml
shell:
  history_size: 10000
  history_path: ~/.cbwsh/history
```

### Autocompletion

cbwsh provides intelligent autocompletion for:
- Commands
- File paths
- Directory names
- Environment variables

#### Using Autocompletion

```bash
# Start typing and press Tab
ls /ho<Tab>         # Completes to /home/
cd ~/Doc<Tab>       # Completes to ~/Documents/
echo $PA<Tab>       # Completes to $PATH

# Multiple matches show suggestions
# Use ↑/↓ to select, Tab to complete
```

### AI Features

cbwsh integrates AI capabilities for enhanced productivity.

#### Setting Up AI

1. **Install Ollama** (recommended for privacy):
```bash
curl https://ollama.ai/install.sh | sh
ollama pull llama2
```

2. **Configure cbwsh**:
```yaml
ai:
  provider: ollama
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true
  monitoring_interval: 30
```

3. **Start cbwsh** and press `Ctrl+M` to open AI monitor

#### AI Assist Mode

Press `Ctrl+A` to toggle AI assist mode:

```bash
# In AI mode, ask questions naturally
"How do I find large files?"
# AI suggests: find . -type f -size +100M

"compress all .log files"
# AI suggests: tar -czf logs.tar.gz *.log
```

#### AI Monitor

Press `Ctrl+M` to toggle the AI monitor pane:

**Features:**
- **Activity tracking**: Monitors your command usage
- **Error analysis**: Provides fixes for failed commands
- **Pattern detection**: Suggests aliases for repeated commands
- **Contextual tips**: Offers relevant advice based on your workflow

**Example workflow:**
```bash
# Run a command that fails
npm install
# Error: Permission denied

# AI monitor automatically suggests:
# "Try: sudo npm install"
# Or: "Use nvm to avoid permission issues"
```

### Job Control

Manage background jobs from cbwsh.

#### Viewing Jobs

```bash
# List all background jobs
jobs

# Output example:
# [1] Running    sleep 100 &
# [2] Stopped    vim file.txt
```

#### Managing Jobs

```bash
# Start a background job (from underlying shell)
sleep 100 &

# Stop a job: Ctrl+Z (in the underlying shell)
# Resume a job: fg (in the underlying shell)
# Kill a job: kill <PID>
```

### Secrets Management

cbwsh includes encrypted secrets storage.

#### Storing Secrets

```bash
# Coming soon - secrets management commands
```

#### Encryption Backends

- **AES-256-GCM**: Built-in, fast encryption
- **age**: Modern encryption tool
- **GPG**: Traditional PGP encryption

#### Configuration

```yaml
secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM
  key_derivation: argon2id
```

### SSH Management

Connect to remote hosts with cbwsh's SSH manager.

#### Quick SSH

```bash
# Regular SSH (uses system ssh command)
ssh user@hostname
```

#### SSH Manager Features

Coming soon:
- Saved host configurations
- Key management
- Port forwarding
- Connection profiles

## Configuration

### Configuration File Location

cbwsh looks for configuration in:
1. `~/.cbwsh/config.yaml` (user config)
2. `/etc/cbwsh/config.yaml` (system config)
3. Built-in defaults

### Creating Configuration

```bash
# Create config directory
mkdir -p ~/.cbwsh

# Create configuration file
cat > ~/.cbwsh/config.yaml << 'EOF'
shell:
  default_shell: bash
  history_size: 10000

ui:
  theme: default
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true

ai:
  provider: none
  enable_monitoring: false
EOF
```

### Configuration Options

#### Shell Settings

```yaml
shell:
  default_shell: bash          # or zsh
  history_size: 10000         # number of commands to remember
  history_path: ~/.cbwsh/history
```

#### UI Settings

```yaml
ui:
  theme: default              # default, dracula, nord
  layout: single              # single, horizontal, vertical, grid
  show_status_bar: true
  enable_animations: true
  syntax_highlighting: true
```

#### AI Settings

```yaml
ai:
  provider: ollama            # none, ollama, openai, anthropic, gemini
  api_key: ""                 # not needed for Ollama
  model: ""
  enable_suggestions: true
  
  # Ollama-specific
  ollama_url: http://localhost:11434
  ollama_model: llama2
  enable_monitoring: true
  monitoring_interval: 30     # seconds between recommendations
```

#### Secrets Settings

```yaml
secrets:
  store_path: ~/.cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM
  key_derivation: argon2id
```

#### Key Bindings

```yaml
keybindings:
  quit: ctrl+q
  help: ctrl+?
  ai_assist: ctrl+a
  # Add custom bindings as needed
```

## Tips and Tricks

### Performance Tips

1. **Disable animations** for faster performance:
```yaml
ui:
  enable_animations: false
```

2. **Reduce history size** if startup is slow:
```yaml
shell:
  history_size: 1000
```

3. **Disable AI monitoring** when not needed:
```yaml
ai:
  enable_monitoring: false
```

### Workflow Tips

1. **Use panes** for parallel tasks:
   - One pane for editing
   - One pane for testing
   - One pane for monitoring

2. **Enable AI monitor** for learning:
   - See better ways to run commands
   - Get automatic error fixes
   - Learn from usage patterns

3. **Use command history** effectively:
   - Press ↑ to find recent commands
   - Modify and re-run commands
   - Build complex commands incrementally

### Customization Tips

1. **Choose a theme** that suits your terminal:
```yaml
ui:
  theme: dracula  # or nord, default
```

2. **Customize the status bar**:
```yaml
ui:
  show_status_bar: true
```

3. **Set your preferred shell**:
```yaml
shell:
  default_shell: zsh
```

## Examples

### Example 1: Basic Workflow

```bash
# Start cbwsh
cbwsh

# Check current directory
pwd

# List files
ls -la

# Create a new file
touch test.txt

# Edit the file (opens in default editor)
nano test.txt

# View the file
cat test.txt

# Remove the file
rm test.txt

# Exit
exit
```

### Example 2: Multi-Pane Development

```bash
# Start cbwsh
cbwsh

# Split vertically for editor and terminal
# Press Ctrl+\

# In left pane: edit code
vim main.go

# Switch to right pane: Press Ctrl+]
# In right pane: watch for changes
watch -n 2 go build

# Split right pane horizontally: Press Ctrl+-
# In bottom right: run tests
go test -v ./...
```

### Example 3: Using AI Features

```bash
# Start cbwsh with AI enabled
cbwsh

# Press Ctrl+M to open AI monitor

# Run some commands - AI monitors and learns
git status
git add .
git commit -m "update"

# Make a typo
gti push
# AI suggests: "Did you mean: git push"

# Ask AI for help: Press Ctrl+A
"how do I undo the last commit"
# AI suggests: git reset --soft HEAD~1
```

### Example 4: Job Control

```bash
# Start a long-running process in background
sleep 300 &

# Check running jobs
jobs

# Continue working on other tasks
ls
pwd
cd /tmp

# Check jobs again
jobs

# Kill the background job if needed
kill %1
```

### Example 5: Working with History

```bash
# Run several commands
ls -la
cd /tmp
pwd
ls -la

# Press ↑ to go back to "ls -la"
# Press ↑ again to go back to "pwd"
# Press ↓ to go forward to "ls -la"

# Modify and execute
# Changes "ls -la" to "ls -lah"
# Press Enter to execute
```

---

## Getting Help

### In-App Help

- Press `Ctrl+?` or `F1` to show help screen
- Press `Alt+H` to open Help menu

### Online Resources

- **Documentation**: https://github.com/cbwinslow/cbwsh
- **Issues**: https://github.com/cbwinslow/cbwsh/issues
- **Discussions**: https://github.com/cbwinslow/cbwsh/discussions

### Reporting Issues

When reporting issues, include:
1. Operating system and version
2. cbwsh version
3. Configuration file (redact sensitive info)
4. Steps to reproduce
5. Expected vs actual behavior

---

**Next Steps:**
- Read the [Installation Guide](INSTALL.md) for installation options
- Check [README.md](README.md) for feature overview
- See [TODO.md](TODO.md) for upcoming features
