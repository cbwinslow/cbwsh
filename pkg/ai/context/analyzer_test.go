package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzer_DetectProjectType(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  []string
		wantType    ProjectType
		wantFileLen int
	}{
		{
			name:        "go project",
			setupFiles:  []string{"go.mod"},
			wantType:    ProjectTypeGo,
			wantFileLen: 1,
		},
		{
			name:        "node project",
			setupFiles:  []string{"package.json"},
			wantType:    ProjectTypeNode,
			wantFileLen: 1,
		},
		{
			name:        "python project",
			setupFiles:  []string{"requirements.txt"},
			wantType:    ProjectTypePython,
			wantFileLen: 1,
		},
		{
			name:        "rust project",
			setupFiles:  []string{"Cargo.toml"},
			wantType:    ProjectTypeRust,
			wantFileLen: 1,
		},
		{
			name:        "docker project",
			setupFiles:  []string{"Dockerfile"},
			wantType:    ProjectTypeDocker,
			wantFileLen: 1,
		},
		{
			name:        "make project",
			setupFiles:  []string{"Makefile"},
			wantType:    ProjectTypeMake,
			wantFileLen: 1,
		},
		{
			name:        "unknown project",
			setupFiles:  []string{"random.txt"},
			wantType:    ProjectTypeUnknown,
			wantFileLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir, err := os.MkdirTemp("", "context-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for _, file := range tt.setupFiles {
				path := filepath.Join(tmpDir, file)
				if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			analyzer := NewAnalyzer()
			projectType, files := analyzer.detectProjectType(tmpDir)

			if projectType != tt.wantType {
				t.Errorf("detectProjectType() type = %v, want %v", projectType, tt.wantType)
			}

			if len(files) != tt.wantFileLen {
				t.Errorf("detectProjectType() files len = %v, want %v", len(files), tt.wantFileLen)
			}
		})
	}
}

func TestAnalyzer_AddRecentCommand(t *testing.T) {
	analyzer := NewAnalyzer()

	// Add some commands
	commands := []string{"ls", "cd /tmp", "git status", "make build"}
	for _, cmd := range commands {
		analyzer.AddRecentCommand(cmd)
	}

	recent := analyzer.getRecentCommands()
	if len(recent) != len(commands) {
		t.Errorf("getRecentCommands() len = %v, want %v", len(recent), len(commands))
	}

	// Verify order
	for i, cmd := range commands {
		if recent[i] != cmd {
			t.Errorf("getRecentCommands()[%d] = %v, want %v", i, recent[i], cmd)
		}
	}
}

func TestAnalyzer_Analyze(t *testing.T) {
	// Create temp directory with a go project
	tmpDir, err := os.MkdirTemp("", "context-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0o644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	analyzer := NewAnalyzer()
	ctx, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("Analyze() error = %v", err)
	}

	if ctx.CWD != tmpDir {
		t.Errorf("Analyze() CWD = %v, want %v", ctx.CWD, tmpDir)
	}

	if ctx.ProjectType != ProjectTypeGo {
		t.Errorf("Analyze() ProjectType = %v, want %v", ctx.ProjectType, ProjectTypeGo)
	}
}

func TestAnalyzer_Cache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "context-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	analyzer := NewAnalyzer()

	// First call
	ctx1, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("First Analyze() error = %v", err)
	}

	// Second call should return cached result
	ctx2, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("Second Analyze() error = %v", err)
	}

	// Should be the same instance
	if ctx1 != ctx2 {
		t.Error("Analyze() should return cached result")
	}

	// Clear cache
	analyzer.ClearCache()

	// Next call should be different instance
	ctx3, err := analyzer.Analyze(tmpDir)
	if err != nil {
		t.Fatalf("Third Analyze() error = %v", err)
	}

	if ctx1 == ctx3 {
		t.Error("Analyze() should return new result after cache clear")
	}
}

func TestAnalyzer_SuggestCommands(t *testing.T) {
	analyzer := NewAnalyzer()

	// Test git suggestions
	gitInfo := &GitInfo{
		IsRepo:                true,
		Branch:                "main",
		HasUncommittedChanges: true,
	}

	ctx := &Context{
		CWD:     "/tmp",
		GitInfo: gitInfo,
	}

	suggestions, err := analyzer.SuggestCommands(ctx, "")
	if err != nil {
		t.Fatalf("SuggestCommands() error = %v", err)
	}

	// Should have git status suggestion
	found := false
	for _, s := range suggestions {
		if s.Command == "git status" {
			found = true
			break
		}
	}

	if !found {
		t.Error("SuggestCommands() should include 'git status' for repo with uncommitted changes")
	}
}

func TestAnalyzer_ProjectSuggestions(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		projectType ProjectType
		wantCommand string
	}{
		{ProjectTypeGo, "go build"},
		{ProjectTypeNode, "npm install"},
		{ProjectTypePython, "pip install -r requirements.txt"},
		{ProjectTypeRust, "cargo build"},
		{ProjectTypeDocker, "docker-compose up"},
		{ProjectTypeMake, "make"},
	}

	for _, tt := range tests {
		t.Run(string(tt.projectType), func(t *testing.T) {
			suggestions := analyzer.projectSuggestions(tt.projectType)

			found := false
			for _, s := range suggestions {
				if s.Command == tt.wantCommand {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("projectSuggestions(%s) should include %s", tt.projectType, tt.wantCommand)
			}
		})
	}
}

func TestGitInfo(t *testing.T) {
	info := &GitInfo{
		IsRepo:                true,
		Branch:                "main",
		Ahead:                 2,
		Behind:                1,
		HasUncommittedChanges: true,
		HasUntrackedFiles:     true,
		HasConflicts:          false,
		RemoteURL:             "git@github.com:user/repo.git",
	}

	if !info.IsRepo {
		t.Error("IsRepo should be true")
	}

	if info.Branch != "main" {
		t.Errorf("Branch = %v, want %v", info.Branch, "main")
	}

	if info.Ahead != 2 {
		t.Errorf("Ahead = %v, want %v", info.Ahead, 2)
	}

	if info.Behind != 1 {
		t.Errorf("Behind = %v, want %v", info.Behind, 1)
	}
}
