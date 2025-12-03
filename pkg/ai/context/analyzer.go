// Package context provides context-aware suggestions for cbwsh.
package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// ProjectType represents the type of project detected.
type ProjectType string

const (
	ProjectTypeUnknown ProjectType = "unknown"
	ProjectTypeGo      ProjectType = "go"
	ProjectTypeNode    ProjectType = "node"
	ProjectTypePython  ProjectType = "python"
	ProjectTypeRust    ProjectType = "rust"
	ProjectTypeJava    ProjectType = "java"
	ProjectTypeRuby    ProjectType = "ruby"
	ProjectTypeMake    ProjectType = "make"
	ProjectTypeDocker  ProjectType = "docker"
)

// GitInfo contains git repository information.
type GitInfo struct {
	// IsRepo indicates if the directory is in a git repository.
	IsRepo bool
	// Branch is the current branch name.
	Branch string
	// Ahead is the number of commits ahead of remote.
	Ahead int
	// Behind is the number of commits behind remote.
	Behind int
	// HasUncommittedChanges indicates uncommitted changes.
	HasUncommittedChanges bool
	// HasUntrackedFiles indicates untracked files.
	HasUntrackedFiles bool
	// HasConflicts indicates merge conflicts.
	HasConflicts bool
	// RemoteURL is the remote repository URL.
	RemoteURL string
}

// Context represents the current shell context.
type Context struct {
	// CWD is the current working directory.
	CWD string
	// GitInfo contains git repository information.
	GitInfo *GitInfo
	// ProjectType is the detected project type.
	ProjectType ProjectType
	// ProjectFiles are relevant project files found.
	ProjectFiles []string
	// RecentCommands are recently executed commands.
	RecentCommands []string
	// EnvVars are relevant environment variables.
	EnvVars map[string]string
}

// Suggestion represents a context-aware command suggestion.
type Suggestion struct {
	// Command is the suggested command.
	Command string
	// Description describes what the command does.
	Description string
	// Priority is the suggestion priority (higher = more relevant).
	Priority int
	// Category is the suggestion category.
	Category string
}

// Analyzer provides context analysis for the shell.
type Analyzer struct {
	mu             sync.RWMutex
	recentCommands []string
	maxRecent      int
	cache          map[string]*Context
}

// NewAnalyzer creates a new context analyzer.
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		recentCommands: make([]string, 0, 100),
		maxRecent:      100,
		cache:          make(map[string]*Context),
	}
}

// Analyze analyzes the current directory context.
func (a *Analyzer) Analyze(cwd string) (*Context, error) {
	a.mu.RLock()
	if cached, ok := a.cache[cwd]; ok {
		a.mu.RUnlock()
		return cached, nil
	}
	a.mu.RUnlock()

	ctx := &Context{
		CWD:            cwd,
		EnvVars:        make(map[string]string),
		RecentCommands: a.getRecentCommands(),
	}

	// Detect project type
	ctx.ProjectType, ctx.ProjectFiles = a.detectProjectType(cwd)

	// Get git info
	ctx.GitInfo = a.getGitInfo(cwd)

	// Get relevant environment variables
	ctx.EnvVars = a.getRelevantEnvVars()

	// Cache result
	a.mu.Lock()
	a.cache[cwd] = ctx
	a.mu.Unlock()

	return ctx, nil
}

// AddRecentCommand adds a command to the recent commands list.
func (a *Analyzer) AddRecentCommand(command string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.recentCommands = append(a.recentCommands, command)
	if len(a.recentCommands) > a.maxRecent {
		a.recentCommands = a.recentCommands[1:]
	}

	// Clear cache when commands change
	a.cache = make(map[string]*Context)
}

// SuggestCommands returns context-aware command suggestions.
func (a *Analyzer) SuggestCommands(ctx *Context, partialInput string) ([]Suggestion, error) {
	var suggestions []Suggestion

	// Add git-related suggestions if in a git repo
	if ctx.GitInfo != nil && ctx.GitInfo.IsRepo {
		suggestions = append(suggestions, a.gitSuggestions(ctx.GitInfo)...)
	}

	// Add project-specific suggestions
	suggestions = append(suggestions, a.projectSuggestions(ctx.ProjectType)...)

	// Add suggestions based on partial input
	if partialInput != "" {
		suggestions = append(suggestions, a.inputBasedSuggestions(partialInput, ctx)...)
	}

	// Sort by priority
	a.sortByPriority(suggestions)

	return suggestions, nil
}

// ClearCache clears the context cache.
func (a *Analyzer) ClearCache() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache = make(map[string]*Context)
}

func (a *Analyzer) getRecentCommands() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]string, len(a.recentCommands))
	copy(result, a.recentCommands)
	return result
}

func (a *Analyzer) detectProjectType(cwd string) (ProjectType, []string) {
	var files []string

	// Check for project-specific files
	projectIndicators := map[string]ProjectType{
		"go.mod":             ProjectTypeGo,
		"go.sum":             ProjectTypeGo,
		"package.json":       ProjectTypeNode,
		"node_modules":       ProjectTypeNode,
		"requirements.txt":   ProjectTypePython,
		"setup.py":           ProjectTypePython,
		"pyproject.toml":     ProjectTypePython,
		"Cargo.toml":         ProjectTypeRust,
		"pom.xml":            ProjectTypeJava,
		"build.gradle":       ProjectTypeJava,
		"Gemfile":            ProjectTypeRuby,
		"Makefile":           ProjectTypeMake,
		"Dockerfile":         ProjectTypeDocker,
		"docker-compose.yml": ProjectTypeDocker,
	}

	for filename, projectType := range projectIndicators {
		path := filepath.Join(cwd, filename)
		if _, err := os.Stat(path); err == nil {
			files = append(files, filename)
			return projectType, files
		}
	}

	return ProjectTypeUnknown, files
}

func (a *Analyzer) getGitInfo(cwd string) *GitInfo {
	info := &GitInfo{}

	// Check if in git repo
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = cwd
	if output, err := cmd.Output(); err == nil && strings.TrimSpace(string(output)) == "true" {
		info.IsRepo = true
	} else {
		return nil
	}

	// Get current branch
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = cwd
	if output, err := cmd.Output(); err == nil {
		info.Branch = strings.TrimSpace(string(output))
	}

	// Check for uncommitted changes
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = cwd
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}
			if line[0] == '?' && line[1] == '?' {
				info.HasUntrackedFiles = true
			} else if line[0] == 'U' || line[1] == 'U' {
				info.HasConflicts = true
			} else {
				info.HasUncommittedChanges = true
			}
		}
	}

	// Get ahead/behind
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = cwd
	if output, err := cmd.Output(); err == nil {
		parts := strings.Fields(string(output))
		if len(parts) >= 2 {
			if n, err := parsePositiveInt(parts[0]); err == nil {
				info.Ahead = n
			}
			if n, err := parsePositiveInt(parts[1]); err == nil {
				info.Behind = n
			}
		}
	}

	// Get remote URL
	cmd = exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = cwd
	if output, err := cmd.Output(); err == nil {
		info.RemoteURL = strings.TrimSpace(string(output))
	}

	return info
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, os.ErrInvalid
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func (a *Analyzer) getRelevantEnvVars() map[string]string {
	relevant := []string{
		"EDITOR", "VISUAL", "PAGER",
		"GOPATH", "GOROOT",
		"NODE_PATH", "NPM_TOKEN",
		"VIRTUAL_ENV", "PYTHONPATH",
		"DOCKER_HOST",
		"AWS_PROFILE", "AWS_REGION",
		"KUBECONFIG",
	}

	result := make(map[string]string)
	for _, name := range relevant {
		if value := os.Getenv(name); value != "" {
			result[name] = value
		}
	}
	return result
}

func (a *Analyzer) gitSuggestions(gitInfo *GitInfo) []Suggestion {
	var suggestions []Suggestion

	if gitInfo.HasUncommittedChanges {
		suggestions = append(suggestions, Suggestion{
			Command:     "git status",
			Description: "View current changes",
			Priority:    90,
			Category:    "git",
		})
		suggestions = append(suggestions, Suggestion{
			Command:     "git diff",
			Description: "View unstaged changes",
			Priority:    85,
			Category:    "git",
		})
		suggestions = append(suggestions, Suggestion{
			Command:     "git add -p",
			Description: "Interactively stage changes",
			Priority:    80,
			Category:    "git",
		})
	}

	if gitInfo.HasUntrackedFiles {
		suggestions = append(suggestions, Suggestion{
			Command:     "git add .",
			Description: "Stage all files",
			Priority:    75,
			Category:    "git",
		})
	}

	if gitInfo.HasConflicts {
		suggestions = append(suggestions, Suggestion{
			Command:     "git mergetool",
			Description: "Resolve merge conflicts",
			Priority:    100,
			Category:    "git",
		})
	}

	if gitInfo.Ahead > 0 {
		suggestions = append(suggestions, Suggestion{
			Command:     "git push",
			Description: "Push local commits",
			Priority:    70,
			Category:    "git",
		})
	}

	if gitInfo.Behind > 0 {
		suggestions = append(suggestions, Suggestion{
			Command:     "git pull",
			Description: "Pull remote changes",
			Priority:    70,
			Category:    "git",
		})
	}

	return suggestions
}

func (a *Analyzer) projectSuggestions(projectType ProjectType) []Suggestion {
	var suggestions []Suggestion

	switch projectType {
	case ProjectTypeGo:
		suggestions = append(suggestions,
			Suggestion{Command: "go build", Description: "Build the project", Priority: 60, Category: "project"},
			Suggestion{Command: "go test ./...", Description: "Run tests", Priority: 55, Category: "project"},
			Suggestion{Command: "go mod tidy", Description: "Tidy dependencies", Priority: 50, Category: "project"},
		)
	case ProjectTypeNode:
		suggestions = append(suggestions,
			Suggestion{Command: "npm install", Description: "Install dependencies", Priority: 60, Category: "project"},
			Suggestion{Command: "npm test", Description: "Run tests", Priority: 55, Category: "project"},
			Suggestion{Command: "npm run build", Description: "Build the project", Priority: 50, Category: "project"},
		)
	case ProjectTypePython:
		suggestions = append(suggestions,
			Suggestion{Command: "pip install -r requirements.txt", Description: "Install dependencies", Priority: 60, Category: "project"},
			Suggestion{Command: "python -m pytest", Description: "Run tests", Priority: 55, Category: "project"},
		)
	case ProjectTypeRust:
		suggestions = append(suggestions,
			Suggestion{Command: "cargo build", Description: "Build the project", Priority: 60, Category: "project"},
			Suggestion{Command: "cargo test", Description: "Run tests", Priority: 55, Category: "project"},
		)
	case ProjectTypeDocker:
		suggestions = append(suggestions,
			Suggestion{Command: "docker-compose up", Description: "Start containers", Priority: 60, Category: "project"},
			Suggestion{Command: "docker-compose down", Description: "Stop containers", Priority: 55, Category: "project"},
			Suggestion{Command: "docker build .", Description: "Build image", Priority: 50, Category: "project"},
		)
	case ProjectTypeMake:
		suggestions = append(suggestions,
			Suggestion{Command: "make", Description: "Run default target", Priority: 60, Category: "project"},
			Suggestion{Command: "make test", Description: "Run tests", Priority: 55, Category: "project"},
			Suggestion{Command: "make clean", Description: "Clean build artifacts", Priority: 50, Category: "project"},
		)
	}

	return suggestions
}

func (a *Analyzer) inputBasedSuggestions(input string, _ *Context) []Suggestion {
	var suggestions []Suggestion

	// Add suggestions based on common patterns
	patterns := map[string][]Suggestion{
		"git": {
			{Command: "git status", Description: "Show working tree status", Priority: 50},
			{Command: "git log --oneline", Description: "Show commit log", Priority: 45},
			{Command: "git branch", Description: "List branches", Priority: 40},
		},
		"docker": {
			{Command: "docker ps", Description: "List running containers", Priority: 50},
			{Command: "docker images", Description: "List images", Priority: 45},
			{Command: "docker logs", Description: "View container logs", Priority: 40},
		},
		"kubectl": {
			{Command: "kubectl get pods", Description: "List pods", Priority: 50},
			{Command: "kubectl get services", Description: "List services", Priority: 45},
			{Command: "kubectl logs", Description: "View pod logs", Priority: 40},
		},
	}

	for prefix, patternSuggestions := range patterns {
		if strings.HasPrefix(input, prefix) {
			suggestions = append(suggestions, patternSuggestions...)
		}
	}

	return suggestions
}

func (a *Analyzer) sortByPriority(suggestions []Suggestion) {
	// Simple bubble sort for small lists
	n := len(suggestions)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if suggestions[j].Priority < suggestions[j+1].Priority {
				suggestions[j], suggestions[j+1] = suggestions[j+1], suggestions[j]
			}
		}
	}
}
