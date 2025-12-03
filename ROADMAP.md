# cbwsh Implementation Roadmap

This document outlines the implementation plan for enabling the features listed in [TODO.md](TODO.md).

## Implementation Phases

### Phase 1: Core Infrastructure (Foundation)
*Priority: 🔴 Critical*
*Timeline: Sprint 1-2*

These foundational features are required for other features to build upon.

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| Natural language to command translation | `pkg/ai/nlp` | 🚧 Scaffold | AI agent |
| Context-aware suggestions | `pkg/ai/context` | 🚧 Scaffold | Shell executor, history |
| Error fix suggestions | `pkg/ai/errorfix` | 🚧 Scaffold | AI agent |
| Command palette | `pkg/ui/palette` | 🚧 Scaffold | UI, keybindings |
| Debug mode logging | `pkg/logging` | ✅ Exists | - |

### Phase 2: Developer Experience (DX)
*Priority: 🔴 Critical + 🟠 High*
*Timeline: Sprint 3-4*

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| Git status integration | `pkg/vcs/git` | 🚧 Scaffold | - |
| Git branch in prompt/status bar | `pkg/ui/prompt` | 🚧 Scaffold | Git integration |
| Command duration display | `pkg/ui/duration` | 🚧 Scaffold | Shell executor |
| Exit code visualization | `pkg/ui/exitcode` | 🚧 Scaffold | Shell executor |
| Model switching at runtime | `pkg/ai/models` | 🚧 Scaffold | AI manager |

### Phase 3: UI/UX Enhancements
*Priority: 🔴 Critical + 🟠 High*
*Timeline: Sprint 5-6*

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| Toast notifications | `pkg/ui/notifications` | 🚧 Scaffold | - |
| Theme hot-reloading | `pkg/ui/themes` | 🚧 Scaffold | Config |
| Cross-platform clipboard | `pkg/clipboard` | 🚧 Scaffold | - |
| Floating panes | `pkg/panes` | ⏳ Planned | Pane manager |
| Command completion for git/docker | `pkg/ui/autocomplete` | 🚧 Extend | Autocomplete |

### Phase 4: Advanced AI Features
*Priority: 🟠 High*
*Timeline: Sprint 7-8*

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| Command chain generation | `pkg/ai/chains` | ⏳ Planned | AI agent |
| Shell script generation | `pkg/ai/scriptgen` | ⏳ Planned | AI agent |
| Code review for scripts | `pkg/ai/review` | ⏳ Planned | AI agent |
| Inline diff preview | `pkg/ui/diff` | ⏳ Planned | - |
| File editing agent | `pkg/ai/agents/file` | ⏳ Planned | AI agent |

### Phase 5: SSH & Remote
*Priority: 🟠 High*
*Timeline: Sprint 9-10*

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| SSH config import | `pkg/ssh` | ⏳ Planned | SSH manager |
| Connection bookmarks | `pkg/ssh` | ⏳ Planned | Config |
| Multi-hop SSH | `pkg/ssh` | ⏳ Planned | SSH manager |

### Phase 6: Security Enhancements
*Priority: 🔴 Critical + 🟠 High*
*Timeline: Sprint 11-12*

| Feature | Package | Status | Dependencies |
|---------|---------|--------|--------------|
| 1Password integration | `pkg/secrets/providers` | ⏳ Planned | Secrets manager |
| Bitwarden integration | `pkg/secrets/providers` | ⏳ Planned | Secrets manager |
| TOTP/HOTP generator | `pkg/auth/totp` | ⏳ Planned | - |
| Touch ID for sudo (macOS) | `pkg/privileges` | ⏳ Planned | Privileges |

---

## Package Structure

```
pkg/
├── ai/
│   ├── agent.go          # ✅ Existing
│   ├── a2a.go            # ✅ Existing
│   ├── nlp/              # 🚧 NEW: Natural language processing
│   │   ├── translator.go # Natural language to command
│   │   └── translator_test.go
│   ├── context/          # 🚧 NEW: Context-aware suggestions
│   │   ├── analyzer.go   # Context analysis
│   │   └── analyzer_test.go
│   ├── errorfix/         # 🚧 NEW: Error fix suggestions
│   │   ├── fixer.go      # Error fixing agent
│   │   └── fixer_test.go
│   ├── models/           # 🚧 NEW: Model management
│   │   ├── switcher.go   # Runtime model switching
│   │   └── switcher_test.go
│   ├── chains/           # ⏳ Command chain generation
│   ├── scriptgen/        # ⏳ Script generation
│   ├── review/           # ⏳ Code review
│   └── agents/           # ⏳ Specialized agents
│       ├── file/         # File editing agent
│       ├── git/          # Git agent
│       └── devops/       # DevOps agent
├── vcs/
│   └── git/              # 🚧 NEW: Git integration
│       ├── status.go     # Git status
│       ├── branch.go     # Branch operations
│       └── git_test.go
├── clipboard/            # 🚧 NEW: Cross-platform clipboard
│   ├── clipboard.go
│   └── clipboard_test.go
└── ui/
    ├── palette/          # 🚧 NEW: Command palette
    │   ├── palette.go
    │   └── palette_test.go
    ├── notifications/    # 🚧 NEW: Toast notifications
    │   ├── toast.go
    │   └── toast_test.go
    ├── prompt/           # 🚧 NEW: Prompt customization
    │   ├── prompt.go     # Starship-like prompt
    │   └── prompt_test.go
    ├── duration/         # 🚧 NEW: Command duration
    │   ├── tracker.go
    │   └── tracker_test.go
    ├── exitcode/         # 🚧 NEW: Exit code visualization
    │   ├── display.go
    │   └── display_test.go
    ├── themes/           # 🚧 NEW: Theme management
    │   ├── loader.go     # Theme hot-reloading
    │   ├── themes.go     # Theme definitions
    │   └── themes_test.go
    └── diff/             # ⏳ Diff viewer
        └── viewer.go
```

---

## Feature Implementation Details

### 1. Natural Language to Command Translation (`pkg/ai/nlp`)

**Purpose**: Convert natural language descriptions to shell commands.

**Example**:
- Input: "find large files over 100MB"
- Output: `find . -type f -size +100M`

**Implementation**:
```go
type Translator interface {
    Translate(ctx context.Context, description string) (*TranslationResult, error)
    TranslateWithContext(ctx context.Context, description string, context *ShellContext) (*TranslationResult, error)
}

type TranslationResult struct {
    Command     string
    Explanation string
    Confidence  float64
    Alternatives []string
}
```

### 2. Context-Aware Suggestions (`pkg/ai/context`)

**Purpose**: Provide intelligent suggestions based on current directory, git status, and command history.

**Context Sources**:
- Current working directory
- Git repository status
- Recent command history
- Project type (package.json, go.mod, etc.)
- Environment variables

**Implementation**:
```go
type ContextAnalyzer interface {
    Analyze(cwd string) (*Context, error)
    SuggestCommands(ctx *Context, partialInput string) ([]Suggestion, error)
}

type Context struct {
    CWD           string
    GitInfo       *GitInfo
    ProjectType   ProjectType
    RecentCommands []string
    EnvVars       map[string]string
}
```

### 3. Command Palette (`pkg/ui/palette`)

**Purpose**: VSCode-style command palette (Ctrl+P) for quick access to all features.

**Implementation**:
```go
type CommandPalette struct {
    commands []PaletteCommand
    filter   string
    selected int
}

type PaletteCommand struct {
    ID          string
    Name        string
    Description string
    Shortcut    string
    Action      func() tea.Cmd
}
```

### 4. Git Integration (`pkg/vcs/git`)

**Purpose**: Native git integration for status, branch info, and operations.

**Implementation**:
```go
type GitRepo interface {
    Status() (*GitStatus, error)
    CurrentBranch() (string, error)
    Branches() ([]string, error)
    Diff() (string, error)
    Log(n int) ([]Commit, error)
}

type GitStatus struct {
    Branch       string
    Ahead        int
    Behind       int
    Staged       []string
    Modified     []string
    Untracked    []string
    HasConflicts bool
}
```

### 5. Toast Notifications (`pkg/ui/notifications`)

**Purpose**: Non-intrusive notifications for long-running command completion.

**Implementation**:
```go
type NotificationManager interface {
    Show(notification Notification)
    Dismiss(id string)
    DismissAll()
    List() []Notification
}

type Notification struct {
    ID        string
    Type      NotificationType
    Title     string
    Message   string
    Duration  time.Duration
    Action    func()
}
```

### 6. Model Switching (`pkg/ai/models`)

**Purpose**: Switch between AI models at runtime.

**Implementation**:
```go
type ModelSwitcher interface {
    ListModels() []ModelInfo
    CurrentModel() ModelInfo
    Switch(modelID string) error
    SwitchProvider(provider AIProvider) error
}

type ModelInfo struct {
    ID          string
    Provider    AIProvider
    Name        string
    Description string
    MaxTokens   int
    Available   bool
}
```

### 7. Theme Hot-Reloading (`pkg/ui/themes`)

**Purpose**: Reload themes without restarting the application.

**Implementation**:
```go
type ThemeManager interface {
    Load(name string) (*Theme, error)
    Reload() error
    Watch() error
    Current() *Theme
    List() []string
}

type Theme struct {
    Name       string
    Colors     ColorScheme
    Styles     StyleSet
}
```

### 8. Clipboard Support (`pkg/clipboard`)

**Purpose**: Cross-platform clipboard operations.

**Implementation**:
```go
type Clipboard interface {
    Read() (string, error)
    Write(text string) error
    WriteImage(img image.Image) error
    ReadImage() (image.Image, error)
    History() []ClipboardEntry
}
```

---

## Testing Strategy

Each new package should include:
1. Unit tests for core functionality
2. Integration tests where applicable
3. Mock implementations for testing dependent packages

Example test structure:
```go
func TestTranslator_Translate(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantCommand string
        wantErr     bool
    }{
        {
            name:        "find large files",
            input:       "find files larger than 100MB",
            wantCommand: "find . -type f -size +100M",
            wantErr:     false,
        },
    }
    // ...
}
```

---

## Configuration Updates

New configuration sections needed:

```yaml
# AI NLP Settings
ai:
  nlp:
    enabled: true
    min_confidence: 0.7
    show_alternatives: true

# Notifications
notifications:
  enabled: true
  position: "top-right"
  duration: 5s
  sound: false

# Git Integration
git:
  show_in_prompt: true
  show_in_status_bar: true
  show_ahead_behind: true

# Theme Settings
themes:
  hot_reload: true
  watch_interval: 1s

# Clipboard
clipboard:
  history_size: 100
  sync_enabled: false
```

---

## Migration Guide

For existing users upgrading to new versions:

1. Configuration files are backward compatible
2. New features are opt-in by default
3. Default keybindings remain unchanged
4. New keybindings documented in CHANGELOG

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on:
- Code style and formatting
- Testing requirements
- Pull request process
- Feature request process

---

*Last updated: December 2024*
