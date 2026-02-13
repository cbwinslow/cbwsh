# cbwsh Visual Guide & UI Documentation

Complete visual reference for cbwsh's user interface, themes, and layout options.

## Table of Contents

1. [UI Overview](#ui-overview)
2. [Terminal Requirements](#terminal-requirements)
3. [Layout Options](#layout-options)
4. [Themes](#themes)
5. [Components](#components)
6. [AI Integration UI](#ai-integration-ui)
7. [Capturing Screenshots](#capturing-screenshots)

## UI Overview

cbwsh provides a modern, bubbletea-powered terminal user interface with multiple panes, syntax highlighting, and rich visual feedback.

### Main Interface

```
┌────────────────────────────────────────────────────────────────────────────┐
│ File  Edit  View  Help                                          [cbwsh]    │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│  ~/projects/myapp                                                          │
│  $ ls -la                                                                  │
│  total 24                                                                  │
│  drwxr-xr-x  6 user user 4096 Feb 13 21:00 .                             │
│  drwxr-xr-x  3 user user 4096 Feb 13 20:55 ..                            │
│  -rw-r--r--  1 user user   45 Feb 13 21:00 README.md                     │
│  drwxr-xr-x  2 user user 4096 Feb 13 21:00 src                           │
│                                                                            │
│  $ █                                                                       │
│                                                                            │
├────────────────────────────────────────────────────────────────────────────┤
│ Ctrl+Q: Quit | Ctrl+?: Help | Ctrl+A: AI | Ctrl+M: Monitor | [  main  ]  │
└────────────────────────────────────────────────────────────────────────────┘
```

### Key Visual Elements

1. **Menu Bar** (optional) - Top menu for mouse users
2. **Command Area** - Where you type and see output
3. **Input Line** - Current command being typed (with syntax highlighting)
4. **Status Bar** - Shows keybindings and git branch
5. **Pane Borders** - Visual separation for multi-pane layouts

## Terminal Requirements

### Recommended Terminals

For the best visual experience, use a terminal with:

✅ **True Color (24-bit) Support:**
- iTerm2 (macOS)
- Alacritty (Cross-platform)
- Windows Terminal
- Warp
- Kitty
- WezTerm

✅ **Font Requirements:**
- Nerd Font (recommended) or any monospace font with good Unicode support
- Ligature support (optional, for prettier code)

**Recommended Fonts:**
- FiraCode Nerd Font
- JetBrains Mono Nerd Font
- Cascadia Code
- Source Code Pro

### Testing Your Terminal

```bash
# Test color support
echo $COLORTERM  # Should show 'truecolor' or '24bit'

# Test 256 colors
for i in {0..255}; do echo -en "\e[48;5;${i}m ${i} \e[0m"; done; echo

# Test true color
awk 'BEGIN{
    s="/\\/\\/\\/\\/\\"; s=s s s s s s s s;
    for (colnum = 0; colnum<77; colnum++) {
        r = 255-(colnum*255/76);
        g = (colnum*510/76);
        b = (colnum*255/76);
        if (g>255) g = 510-g;
        printf "\033[48;2;%d;%d;%dm", r,g,b;
        printf "\033[38;2;%d;%d;%dm", 255-r,255-g,255-b;
        printf "%s\033[0m", substr(s,colnum+1,1);
    }
    printf "\n";
}'
```

## Layout Options

### 1. Single Pane (Default)

Full terminal dedicated to one shell session.

```
┌────────────────────────────────────────┐
│                                        │
│                                        │
│           Single Pane                  │
│           Full Screen                  │
│                                        │
│                                        │
│                                        │
└────────────────────────────────────────┘
```

**Config:**
```yaml
ui:
  layout: single
```

**Use case:** Simple, focused work

### 2. Horizontal Split

Two panes side by side.

```
┌──────────────────┬──────────────────┐
│                  │                  │
│                  │                  │
│   Left Pane      │   Right Pane     │
│                  │                  │
│                  │                  │
│                  │                  │
│                  │                  │
└──────────────────┴──────────────────┘
```

**Config:**
```yaml
ui:
  layout: horizontal
```

**Use case:** Code on left, tests on right

**Keybindings:**
- `Ctrl+\` - Create horizontal split
- `Ctrl+]` - Switch to next pane
- `Ctrl+[` - Switch to previous pane

### 3. Vertical Split

Two panes stacked vertically.

```
┌────────────────────────────────────────┐
│                                        │
│           Top Pane                     │
│                                        │
├────────────────────────────────────────┤
│                                        │
│         Bottom Pane                    │
│                                        │
└────────────────────────────────────────┘
```

**Config:**
```yaml
ui:
  layout: vertical
```

**Use case:** Edit above, run below

**Keybindings:**
- `Ctrl+-` - Create vertical split
- `Ctrl+]` - Switch to next pane
- `Ctrl+[` - Switch to previous pane

### 4. Grid Layout

Four panes in a 2x2 grid.

```
┌──────────────────┬──────────────────┐
│                  │                  │
│   Top-Left       │   Top-Right      │
│                  │                  │
├──────────────────┼──────────────────┤
│                  │                  │
│  Bottom-Left     │  Bottom-Right    │
│                  │                  │
└──────────────────┴──────────────────┘
```

**Config:**
```yaml
ui:
  layout: grid
```

**Use case:** Full development environment
- Top-left: Editor
- Top-right: Build output
- Bottom-left: Tests
- Bottom-right: Server logs

### 5. With AI Monitor

Any layout can add an AI monitor pane on the right.

```
┌────────────────────────┬──────────────┐
│                        │ AI Monitor   │
│                        │              │
│    Main Pane           │ Activity:    │
│                        │ • git status │
│                        │ • npm build  │
│                        │              │
│                        │ Tips:        │
│                        │ Try 'git log'│
└────────────────────────┴──────────────┘
```

**Toggle:** `Ctrl+M`

## Themes

cbwsh includes multiple built-in themes with different color schemes.

### Default Theme

Clean, professional theme with subtle colors.

**Colors:**
- Background: Terminal default
- Foreground: Terminal default
- Accent: Blue (#5B9FFF)
- Success: Green (#7FD962)
- Warning: Yellow (#FFCA27)
- Error: Red (#FF5A67)
- Comment: Gray (#6B7280)

**Visual Style:**
```
$ echo "Hello World"
  ^^^^^              (Blue - command)
        ^^^^^^^^^^^^^ (Green - string)

$ ls -la
  ^^                 (Blue - command)
     ^^^             (Yellow - flag)
```

### Dracula Theme

Popular dark theme with purple/pink accents.

**Config:**
```yaml
ui:
  theme: dracula
```

**Colors:**
- Background: #282A36
- Foreground: #F8F8F2
- Accent: #BD93F9 (Purple)
- Success: #50FA7B (Green)
- Warning: #F1FA8C (Yellow)
- Error: #FF5555 (Red)
- Cyan: #8BE9FD
- Pink: #FF79C6

### Nord Theme

Arctic-inspired cool palette.

**Config:**
```yaml
ui:
  theme: nord
```

**Colors:**
- Background: #2E3440
- Foreground: #D8DEE9
- Accent: #88C0D0 (Frost Blue)
- Success: #A3BE8C (Green)
- Warning: #EBCB8B (Yellow)
- Error: #BF616A (Red)

### Tokyo Night Theme

Vibrant, modern dark theme inspired by Tokyo at night.

**Config:**
```yaml
ui:
  theme: tokyo-night
```

**Colors:**
- Background: #1A1B26
- Foreground: #A9B1D6
- Accent: #7AA2F7 (Blue)
- Success: #9ECE6A (Green)
- Warning: #E0AF68 (Yellow)
- Error: #F7768E (Red)
- Purple: #BB9AF7
- Cyan: #7DCFFF

### Gruvbox Theme

Retro groove color scheme with warm contrast.

**Config:**
```yaml
ui:
  theme: gruvbox
```

**Colors:**
- Background: #282828
- Foreground: #EBDBB2
- Accent: #458588 (Blue)
- Success: #B8BB26 (Green)
- Warning: #FABD2F (Yellow)
- Error: #FB4934 (Red)
- Orange: #FE8019
- Purple: #D3869B

## Components

### Syntax Highlighting

Commands are highlighted in real-time as you type:

```bash
# Command (blue)
git commit -m "message"
    ^^^^^^             # Subcommand (cyan)
           ^^          # Flag (yellow)
              ^^^^^^^^^ # String (green)

# Variables (purple)
$HOME/projects
^^^^^

# Paths (white/default)
/usr/local/bin

# Operators (orange)
> output.txt
^

# Numbers (magenta)
seq 1 10
    ^ ^^
```

### Autocomplete

Live suggestions as you type:

```
$ git che█
        ├─ checkout
        ├─ cherry
        └─ cherry-pick

↑/↓: Navigate  Tab: Complete  Esc: Cancel
```

### Progress Bars

For long-running operations:

```
Installing dependencies...
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓░░░░░░░░░░ 67% (120/180) [2m15s]

Building project...
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ 100% (50/50) [45s]
```

### Notifications

Toast-style notifications:

```
┌────────────────────────────┐
│ ✓ Build completed          │
│   Time: 2m 15s             │
└────────────────────────────┘

┌────────────────────────────┐
│ ⚠ Warning: disk space low  │
│   Available: 2.3 GB        │
└────────────────────────────┘

┌────────────────────────────┐
│ ✗ Error: test failed       │
│   See logs for details     │
└────────────────────────────┘
```

### Menu Bar

Mouse-accessible menu:

```
┌─────┬──────┬──────┬──────┐
│File │ Edit │ View │ Help │
└─────┴──────┴──────┴──────┘
  │
  ├─ New Pane      (Ctrl+N)
  ├─ Close Pane    (Ctrl+W)
  ├─ Split Horiz   (Ctrl+\)
  ├─ Split Vert    (Ctrl+-)
  ├─────────────
  └─ Quit          (Ctrl+Q)
```

## AI Integration UI

### AI Monitor Pane

Real-time AI analysis and recommendations:

```
┌─────────────────────────────────────┬──────────────────────┐
│ $ npm test                          │  AI Monitor         │
│ Running tests...                    │                     │
│                                     │ Activity            │
│ Test Suites: 2 passed, 2 total     │ ✓ Tests passing     │
│ Tests:       15 passed, 15 total   │ ✓ Build successful  │
│ Snapshots:   0 total                │                     │
│ Time:        3.142s                 │ Recommendations     │
│                                     │ • Add more tests    │
│ $ █                                 │ • Update deps       │
│                                     │                     │
│                                     │ Tips                │
│                                     │ Try 'npm audit' to  │
│                                     │ check security      │
│                                     │                     │
│                                     │ Last updated: 2s ago│
└─────────────────────────────────────┴──────────────────────┘
```

**Toggle:** Press `Ctrl+M`

### AI Chat Pane

Interactive conversation with AI:

```
┌─────────────────────────────────────────────────────────────────────┐
│                         AI Assistant Chat                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│ You: How do I find large files?                                    │
│                                                                     │
│ AI: You can find large files using the `find` command:             │
│                                                                     │
│     find . -type f -size +100M                                     │
│                                                                     │
│     This finds files larger than 100MB in the current directory    │
│     and subdirectories.                                            │
│                                                                     │
│     Or use `du` to see directory sizes:                            │
│                                                                     │
│     du -h --max-depth=1 | sort -hr                                │
│                                                                     │
│ You: █                                                              │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│ Ctrl+Enter: Send | Ctrl+L: Clear | Esc: Close                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Open:** Press `Ctrl+A`

### AI Suggestions

Inline suggestions based on context:

```
$ git commit 
           █

  💡 AI Suggests:
     git commit -m "Add new feature"
     git commit --amend
     git commit -a -m "Update files"

Tab: Accept  ↑/↓: Navigate  Esc: Dismiss
```

## Capturing Screenshots

### Using Built-in Tool

cbwsh includes a screenshot utility:

```bash
# Take screenshot of current view
cbwsh screenshot

# Save to specific file
cbwsh screenshot -o ~/screenshots/cbwsh-demo.png

# Capture specific pane
cbwsh screenshot --pane 1

# Capture with theme
cbwsh screenshot --theme dracula
```

### Manual Screenshot Methods

**macOS:**
```bash
# Full window
Cmd+Shift+4, then Space, then click window

# Selected area
Cmd+Shift+4, then drag
```

**Linux:**
```bash
# Using gnome-screenshot
gnome-screenshot -w  # Window
gnome-screenshot -a  # Area

# Using scrot
scrot -s
```

**Windows:**
```bash
# Windows Terminal
Right-click > Export Text
# Or use Snipping Tool
```

### ASCII Cinema Recording

Record terminal sessions as animated SVG/GIF:

```bash
# Install asciinema
npm install -g asciinema

# Record session
asciinema rec cbwsh-demo.cast

# Play back
asciinema play cbwsh-demo.cast

# Upload and share
asciinema upload cbwsh-demo.cast
```

### VHS - Programmatic Screenshots

Create reproducible terminal demonstrations:

```bash
# Install VHS
go install github.com/charmbracelet/vhs@latest

# Create demo script
cat > demo.tape << 'EOF'
Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1200
Set Height 600
Set Theme "Dracula"

Type "cbwsh"
Enter
Sleep 2s

Type "echo 'Hello cbwsh!'"
Enter
Sleep 1s

Type "ls -la"
Enter
Sleep 2s

Ctrl+Q
EOF

# Generate GIF
vhs demo.tape
```

## UI Customization

### Custom Themes

Create your own theme:

```yaml
# ~/.config/cbwsh/themes/mytheme.yaml

name: My Custom Theme
colors:
  background: "#1E1E1E"
  foreground: "#D4D4D4"
  cursor: "#AEAFAD"
  
  # ANSI colors
  black: "#000000"
  red: "#CD3131"
  green: "#0DBC79"
  yellow: "#E5E510"
  blue: "#2472C8"
  magenta: "#BC3FBC"
  cyan: "#11A8CD"
  white: "#E5E5E5"
  
  # Bright colors
  bright_black: "#666666"
  bright_red: "#F14C4C"
  bright_green: "#23D18B"
  bright_yellow: "#F5F543"
  bright_blue: "#3B8EEA"
  bright_magenta: "#D670D6"
  bright_cyan: "#29B8DB"
  bright_white: "#FFFFFF"
  
  # UI elements
  border: "#404040"
  selection: "#264F78"
  highlight: "#4E4E4E"
```

Use your theme:

```yaml
# config.yaml
ui:
  theme: mytheme
  theme_path: ~/.config/cbwsh/themes
```

### Layout Presets

Save custom layouts:

```yaml
# ~/.config/cbwsh/layouts/dev.yaml

name: Development Layout
layout:
  type: custom
  panes:
    - position: { x: 0, y: 0, width: 50%, height: 70% }
      command: "vim"
    - position: { x: 50%, y: 0, width: 50%, height: 35% }
      command: "npm run watch"
    - position: { x: 50%, y: 35%, width: 50%, height: 35% }
      command: "npm test -- --watch"
    - position: { x: 0, y: 70%, width: 100%, height: 30% }
      command: ""  # Interactive shell
```

Load preset:

```bash
cbwsh --layout dev
```

## Tips for Best Visual Experience

1. **Use a modern terminal** with true color support
2. **Install a Nerd Font** for better icons and symbols
3. **Enable animations** in config for smooth transitions
4. **Use the AI monitor** for live assistance
5. **Experiment with themes** to find what works for you
6. **Try different layouts** for different tasks
7. **Customize keybindings** to match your workflow

## Visual Examples Repository

For actual screenshots and recordings, visit:

https://github.com/cbwinslow/cbwsh/tree/main/screenshots

Contents:
- Theme comparisons
- Layout examples
- Feature demonstrations
- Animation samples
- Multi-pane workflows

---

**Note:** Visual elements are rendered using the Bubble Tea ecosystem with:
- lipgloss for styling
- bubbles for components
- glamour for markdown
- harmonica for animations

All rendering is done in the terminal using ANSI escape sequences and Unicode characters.
