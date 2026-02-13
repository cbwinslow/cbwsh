#!/bin/bash
#
# cbwsh Setup Script
# ===================
# Comprehensive setup for cbwsh shell environment
#
# This script:
# - Creates proper directory structure following XDG Base Directory specification
# - Sets up configuration files
# - Installs shell integration
# - Configures logging
# - Provides upgrade mechanism
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/setup.sh | bash
#   or
#   ./setup.sh [options]
#
# Options:
#   --dev                Install development version
#   --config-dir <path>  Custom config directory (default: ~/.config/cbwsh)
#   --data-dir <path>    Custom data directory (default: ~/.local/share/cbwsh)
#   --no-shell-integration  Skip shell integration setup
#   --help               Show this help message

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Default values - Follow XDG Base Directory specification
XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

CONFIG_DIR="${XDG_CONFIG_HOME}/cbwsh"
DATA_DIR="${XDG_DATA_HOME}/cbwsh"
CACHE_DIR="${XDG_CACHE_HOME}/cbwsh"
STATE_DIR="${XDG_STATE_HOME}/cbwsh"
LOG_DIR="${STATE_DIR}/logs"
PLUGIN_DIR="${DATA_DIR}/plugins"
THEME_DIR="${DATA_DIR}/themes"
HISTORY_FILE="${STATE_DIR}/history"
SECRETS_FILE="${DATA_DIR}/secrets.enc"

DEV_MODE="false"
SKIP_SHELL_INTEGRATION="false"

# Logging functions
info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

success() {
    echo -e "${GREEN}[✓]${NC} $*"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $*"
}

error() {
    echo -e "${RED}[ERROR]${NC} $*"
    exit 1
}

step() {
    echo -e "${CYAN}[STEP]${NC} $*"
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
                                        
  Setup & Configuration
EOF
    echo -e "${NC}"
}

# Show help
show_help() {
    cat << EOF
cbwsh Setup Script

This script sets up cbwsh with proper directory structure and configuration.

Usage: $0 [options]

Options:
  --dev                      Install development version
  --config-dir <path>        Custom config directory (default: $CONFIG_DIR)
  --data-dir <path>          Custom data directory (default: $DATA_DIR)
  --no-shell-integration     Skip shell integration setup
  --help                     Show this help message

Directory Structure (XDG compliant):
  Config:  $CONFIG_DIR
  Data:    $DATA_DIR
  Cache:   $CACHE_DIR
  State:   $STATE_DIR
  Logs:    $LOG_DIR

Examples:
  $0                          # Standard setup
  $0 --dev                    # Development setup
  $0 --config-dir ~/.cbwsh    # Custom config directory

EOF
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --dev)
                DEV_MODE="true"
                shift
                ;;
            --config-dir)
                CONFIG_DIR="$2"
                shift 2
                ;;
            --data-dir)
                DATA_DIR="$2"
                shift 2
                ;;
            --no-shell-integration)
                SKIP_SHELL_INTEGRATION="true"
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

# Create directory structure
create_directories() {
    step "Creating directory structure..."
    
    # XDG directories
    mkdir -p "$CONFIG_DIR"
    mkdir -p "$DATA_DIR"
    mkdir -p "$CACHE_DIR"
    mkdir -p "$STATE_DIR"
    mkdir -p "$LOG_DIR"
    mkdir -p "$PLUGIN_DIR"
    mkdir -p "$THEME_DIR"
    
    # Create subdirectories
    mkdir -p "$DATA_DIR/sessions"
    mkdir -p "$DATA_DIR/backups"
    mkdir -p "$CACHE_DIR/downloads"
    mkdir -p "$CACHE_DIR/temp"
    
    success "Directory structure created"
    
    info "Directories created:"
    info "  Config:  $CONFIG_DIR"
    info "  Data:    $DATA_DIR"
    info "  Cache:   $CACHE_DIR"
    info "  State:   $STATE_DIR"
    info "  Logs:    $LOG_DIR"
}

# Create default configuration
create_default_config() {
    step "Creating default configuration..."
    
    local config_file="$CONFIG_DIR/config.yaml"
    
    if [ -f "$config_file" ]; then
        warn "Configuration file already exists at $config_file"
        read -p "Do you want to backup and overwrite it? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            local backup_file="$DATA_DIR/backups/config.yaml.$(date +%Y%m%d_%H%M%S).bak"
            cp "$config_file" "$backup_file"
            success "Backed up existing config to $backup_file"
        else
            info "Keeping existing configuration"
            return 0
        fi
    fi
    
    cat > "$config_file" << 'EOF'
# cbwsh Configuration
# https://github.com/cbwinslow/cbwsh

# Shell settings
shell:
  default_shell: bash
  history_size: 10000
  history_path: ~/.local/state/cbwsh/history
  aliases:
    ll: ls -lah
    gs: git status
    gd: git diff
  environment:
    EDITOR: vim
    PAGER: less

# UI settings
ui:
  theme: default
  layout: single
  show_status_bar: true
  show_menu_bar: false
  enable_animations: true
  syntax_highlighting: true

# AI settings (disabled by default)
ai:
  provider: none
  enable_suggestions: false
  enable_monitoring: false

# SSH settings
ssh:
  default_key_path: ~/.ssh/id_rsa
  known_hosts_path: ~/.ssh/known_hosts
  connect_timeout: 30
  keep_alive_interval: 60

# Secrets settings
secrets:
  store_path: ~/.local/share/cbwsh/secrets.enc
  encryption_algorithm: AES-256-GCM

# Keybindings
keybindings:
  quit: ctrl+q
  help: ctrl+?
  ai_assist: ctrl+a
  ai_monitor: ctrl+m
  new_pane: ctrl+n
  close_pane: ctrl+w
  next_pane: ctrl+]
  prev_pane: ctrl+[
  split_vertical: ctrl+\
  split_horizontal: ctrl+-
  clear_screen: ctrl+l
  command_palette: ctrl+p

# Plugins
plugins:
  enabled: true
  auto_load: true
  directory: ~/.local/share/cbwsh/plugins

# Logging
logging:
  enabled: true
  level: info
  path: ~/.local/state/cbwsh/logs/cbwsh.log
  max_size: 10
  max_backups: 3
  max_age: 30

# Advanced settings
advanced:
  shell_integration: true
  cd_on_exit: false
  save_session: true
  restore_session: true
  max_output_lines: 10000
  tab_width: 4
EOF
    
    success "Created default configuration at $config_file"
}

# Setup logging
setup_logging() {
    step "Setting up logging..."
    
    # Create log rotation configuration
    local logrotate_conf="$CONFIG_DIR/logrotate.conf"
    
    cat > "$logrotate_conf" << EOF
# cbwsh log rotation configuration

$LOG_DIR/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 $(whoami) $(id -gn)
}
EOF
    
    # Create initial log file
    touch "$LOG_DIR/cbwsh.log"
    
    success "Logging configured at $LOG_DIR"
}

# Setup shell integration
setup_shell_integration() {
    if [ "$SKIP_SHELL_INTEGRATION" = "true" ]; then
        info "Skipping shell integration setup"
        return 0
    fi
    
    step "Setting up shell integration..."
    
    # Detect current shell
    local current_shell=$(basename "$SHELL")
    
    case "$current_shell" in
        bash)
            setup_bash_integration
            ;;
        zsh)
            setup_zsh_integration
            ;;
        *)
            warn "Unsupported shell: $current_shell"
            warn "Shell integration skipped"
            return 0
            ;;
    esac
    
    success "Shell integration configured for $current_shell"
}

# Setup bash integration
setup_bash_integration() {
    local bashrc="$HOME/.bashrc"
    local integration_marker="# cbwsh integration"
    
    if grep -q "$integration_marker" "$bashrc" 2>/dev/null; then
        info "Bash integration already configured"
        return 0
    fi
    
    cat >> "$bashrc" << 'EOF'

# cbwsh integration
if [ -f ~/.config/cbwsh/shell-integration.bash ]; then
    source ~/.config/cbwsh/shell-integration.bash
fi

# cbwsh alias (optional)
# alias sh='cbwsh'
EOF
    
    # Create shell integration file
    cat > "$CONFIG_DIR/shell-integration.bash" << 'EOF'
# cbwsh shell integration for bash

# Add cbwsh to PATH if installed
if [ -d "$HOME/.local/bin" ]; then
    export PATH="$HOME/.local/bin:$PATH"
fi

# Set XDG directories
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
export XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
export XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

# cbwsh environment
export CBWSH_CONFIG_DIR="$XDG_CONFIG_HOME/cbwsh"
export CBWSH_DATA_DIR="$XDG_DATA_HOME/cbwsh"
export CBWSH_LOG_DIR="$XDG_STATE_HOME/cbwsh/logs"

# Optional: Set cbwsh as default shell for new terminals
# export SHELL="$(which cbwsh)"
EOF
    
    info "Bash integration added to $bashrc"
}

# Setup zsh integration
setup_zsh_integration() {
    local zshrc="$HOME/.zshrc"
    local integration_marker="# cbwsh integration"
    
    if grep -q "$integration_marker" "$zshrc" 2>/dev/null; then
        info "Zsh integration already configured"
        return 0
    fi
    
    cat >> "$zshrc" << 'EOF'

# cbwsh integration
if [ -f ~/.config/cbwsh/shell-integration.zsh ]; then
    source ~/.config/cbwsh/shell-integration.zsh
fi

# cbwsh alias (optional)
# alias sh='cbwsh'
EOF
    
    # Create shell integration file
    cat > "$CONFIG_DIR/shell-integration.zsh" << 'EOF'
# cbwsh shell integration for zsh

# Add cbwsh to PATH if installed
if [ -d "$HOME/.local/bin" ]; then
    path=("$HOME/.local/bin" $path)
    export PATH
fi

# Set XDG directories
export XDG_CONFIG_HOME="${XDG_CONFIG_HOME:-$HOME/.config}"
export XDG_DATA_HOME="${XDG_DATA_HOME:-$HOME/.local/share}"
export XDG_CACHE_HOME="${XDG_CACHE_HOME:-$HOME/.cache}"
export XDG_STATE_HOME="${XDG_STATE_HOME:-$HOME/.local/state}"

# cbwsh environment
export CBWSH_CONFIG_DIR="$XDG_CONFIG_HOME/cbwsh"
export CBWSH_DATA_DIR="$XDG_DATA_HOME/cbwsh"
export CBWSH_LOG_DIR="$XDG_STATE_HOME/cbwsh/logs"

# Optional: Set cbwsh as default shell for new terminals
# export SHELL="$(which cbwsh)"
EOF
    
    info "Zsh integration added to $zshrc"
}

# Create upgrade script
create_upgrade_script() {
    step "Creating upgrade script..."
    
    local upgrade_script="$DATA_DIR/upgrade.sh"
    
    cat > "$upgrade_script" << 'EOF'
#!/bin/bash
# cbwsh upgrade script

set -euo pipefail

CYAN='\033[0;36m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${CYAN}Upgrading cbwsh...${NC}"

# Determine installation location
CBWSH_BIN=$(command -v cbwsh || echo "/usr/local/bin/cbwsh")
INSTALL_DIR=$(dirname "$CBWSH_BIN")

# Backup current version
if [ -f "$CBWSH_BIN" ]; then
    BACKUP_DIR="$HOME/.local/share/cbwsh/backups"
    mkdir -p "$BACKUP_DIR"
    BACKUP_FILE="$BACKUP_DIR/cbwsh.$(date +%Y%m%d_%H%M%S).bak"
    cp "$CBWSH_BIN" "$BACKUP_FILE"
    echo -e "${GREEN}Backed up current version to $BACKUP_FILE${NC}"
fi

# Download and install latest version
curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash -s -- --prefix "$INSTALL_DIR"

echo -e "${GREEN}cbwsh upgraded successfully!${NC}"
echo "Run 'cbwsh --version' to see the new version"
EOF
    
    chmod +x "$upgrade_script"
    success "Created upgrade script at $upgrade_script"
    info "Run '$upgrade_script' to upgrade cbwsh in the future"
}

# Create README
create_readme() {
    step "Creating README..."
    
    local readme_file="$CONFIG_DIR/README.md"
    
    cat > "$readme_file" << EOF
# cbwsh Configuration

This directory contains your cbwsh configuration following the XDG Base Directory specification.

## Directory Structure

- \`$CONFIG_DIR\` - Configuration files
- \`$DATA_DIR\` - Application data (plugins, themes, secrets)
- \`$CACHE_DIR\` - Cache files
- \`$STATE_DIR\` - State files (history, logs)

## Key Files

- \`config.yaml\` - Main configuration file
- \`shell-integration.bash/zsh\` - Shell integration scripts
- \`logrotate.conf\` - Log rotation configuration

## Data Locations

- Logs: \`$LOG_DIR\`
- History: \`$HISTORY_FILE\`
- Plugins: \`$PLUGIN_DIR\`
- Themes: \`$THEME_DIR\`
- Secrets: \`$SECRETS_FILE\`

## Upgrading

Run the upgrade script:
\`\`\`bash
$DATA_DIR/upgrade.sh
\`\`\`

## Configuration

Edit \`config.yaml\` to customize cbwsh:
\`\`\`bash
\${EDITOR:-vi} $CONFIG_DIR/config.yaml
\`\`\`

## Documentation

- GitHub: https://github.com/cbwinslow/cbwsh
- Usage Guide: https://github.com/cbwinslow/cbwsh/blob/main/USAGE.md
- Install Guide: https://github.com/cbwinslow/cbwsh/blob/main/INSTALL.md

## Support

- Issues: https://github.com/cbwinslow/cbwsh/issues
- Discussions: https://github.com/cbwinslow/cbwsh/discussions
EOF
    
    success "Created README at $readme_file"
}

# Create .gitignore for config directory
create_gitignore() {
    step "Creating .gitignore..."
    
    local gitignore_file="$CONFIG_DIR/.gitignore"
    
    cat > "$gitignore_file" << 'EOF'
# cbwsh gitignore

# Logs
logs/
*.log

# Secrets
secrets.enc
*.key
*.pem

# Cache
cache/
*.cache

# Session data
sessions/
*.session

# Backups
backups/
*.bak

# Temporary files
*.tmp
*.swp
*~

# OS files
.DS_Store
Thumbs.db
EOF
    
    success "Created .gitignore at $gitignore_file"
}

# Validate installation
validate_setup() {
    step "Validating setup..."
    
    local errors=0
    
    # Check directories
    for dir in "$CONFIG_DIR" "$DATA_DIR" "$CACHE_DIR" "$STATE_DIR" "$LOG_DIR"; do
        if [ ! -d "$dir" ]; then
            error "Directory not found: $dir"
            ((errors++))
        fi
    done
    
    # Check config file
    if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
        error "Configuration file not found"
        ((errors++))
    fi
    
    # Check if cbwsh is installed
    if ! command -v cbwsh &> /dev/null; then
        warn "cbwsh binary not found in PATH"
        info "You may need to install cbwsh first:"
        info "  curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash"
    fi
    
    if [ $errors -eq 0 ]; then
        success "Setup validation completed successfully"
        return 0
    else
        error "Setup validation failed with $errors errors"
        return 1
    fi
}

# Show summary
show_summary() {
    echo
    echo -e "${BOLD}${GREEN}Setup Complete!${NC}"
    echo
    echo -e "${CYAN}Next Steps:${NC}"
    echo
    echo "1. Install cbwsh binary (if not already installed):"
    echo -e "   ${YELLOW}curl -fsSL https://raw.githubusercontent.com/cbwinslow/cbwsh/main/install.sh | bash${NC}"
    echo
    echo "2. Reload your shell configuration:"
    echo -e "   ${YELLOW}source ~/.bashrc${NC}  # or source ~/.zshrc"
    echo
    echo "3. Configure AI features (optional):"
    echo -e "   ${YELLOW}\${EDITOR:-vi} $CONFIG_DIR/config.yaml${NC}"
    echo
    echo "4. Start cbwsh:"
    echo -e "   ${YELLOW}cbwsh${NC}"
    echo
    echo -e "${CYAN}Configuration:${NC}"
    echo "  Config: $CONFIG_DIR/config.yaml"
    echo "  Logs:   $LOG_DIR"
    echo "  Data:   $DATA_DIR"
    echo
    echo -e "${CYAN}Documentation:${NC}"
    echo "  https://github.com/cbwinslow/cbwsh"
    echo
}

# Main installation flow
main() {
    show_banner
    parse_args "$@"
    
    info "Starting cbwsh setup..."
    echo
    
    create_directories
    create_default_config
    setup_logging
    setup_shell_integration
    create_upgrade_script
    create_readme
    create_gitignore
    validate_setup
    
    show_summary
}

# Run main function
main "$@"
