package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary git repository for testing.
func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		cleanup()
		t.Skipf("git init failed (git may not be installed): %v", err)
	}

	// Configure git for test
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	return tmpDir, cleanup
}

func TestRepo_IsRepo(t *testing.T) {
	// Test with git repo
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	repo := NewRepo(tmpDir)
	if !repo.IsRepo() {
		t.Error("IsRepo() should return true for git repository")
	}

	// Test with non-git directory
	nonGitDir, err := os.MkdirTemp("", "non-git")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(nonGitDir)

	repo2 := NewRepo(nonGitDir)
	if repo2.IsRepo() {
		t.Error("IsRepo() should return false for non-git directory")
	}
}

func TestRepo_Status(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	repo := NewRepo(tmpDir)
	status, err := repo.Status()
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}

	if !status.IsRepo {
		t.Error("Status().IsRepo should be true")
	}

	if status.Branch == "" {
		t.Error("Status().Branch should not be empty")
	}
}

func TestRepo_Status_NotGitRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "non-git")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	repo := NewRepo(tmpDir)
	_, err = repo.Status()
	if err != ErrNotGitRepo {
		t.Errorf("Status() error = %v, want %v", err, ErrNotGitRepo)
	}
}

func TestRepo_CurrentBranch(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	repo := NewRepo(tmpDir)
	branch, err := repo.CurrentBranch()
	if err != nil {
		t.Fatalf("CurrentBranch() error = %v", err)
	}

	// Branch should be "master" or "main" depending on git version
	if branch != "master" && branch != "main" {
		t.Errorf("CurrentBranch() = %s, want master or main", branch)
	}
}

func TestRepo_Branches(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	// Create a new branch
	cmd = exec.Command("git", "branch", "feature")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	repo := NewRepo(tmpDir)
	branches, err := repo.Branches()
	if err != nil {
		t.Fatalf("Branches() error = %v", err)
	}

	if len(branches) < 2 {
		t.Errorf("Branches() len = %d, want >= 2", len(branches))
	}

	// Find feature branch
	found := false
	for _, b := range branches {
		if b.Name == "feature" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Branches() should include 'feature' branch")
	}
}

func TestRepo_Log(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create commits
	for i := 0; i < 3; i++ {
		testFile := filepath.Join(tmpDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("content "+string(rune('0'+i))), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		cmd := exec.Command("git", "add", ".")
		cmd.Dir = tmpDir
		_ = cmd.Run()

		cmd = exec.Command("git", "commit", "-m", "commit "+string(rune('0'+i)))
		cmd.Dir = tmpDir
		_ = cmd.Run()
	}

	repo := NewRepo(tmpDir)
	commits, err := repo.Log(5)
	if err != nil {
		t.Fatalf("Log() error = %v", err)
	}

	if len(commits) != 3 {
		t.Errorf("Log() len = %d, want 3", len(commits))
	}

	if commits[0].Hash == "" {
		t.Error("Log()[0].Hash should not be empty")
	}

	if commits[0].ShortHash == "" {
		t.Error("Log()[0].ShortHash should not be empty")
	}
}

func TestRepo_Tags(t *testing.T) {
	tmpDir, cleanup := setupTestRepo(t)
	defer cleanup()

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	// Create tags
	cmd = exec.Command("git", "tag", "v1.0.0")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	cmd = exec.Command("git", "tag", "v1.1.0")
	cmd.Dir = tmpDir
	_ = cmd.Run()

	repo := NewRepo(tmpDir)
	tags, err := repo.Tags()
	if err != nil {
		t.Fatalf("Tags() error = %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("Tags() len = %d, want 2", len(tags))
	}
}

func TestFormatStatus(t *testing.T) {
	tests := []struct {
		name   string
		status *Status
		want   string
	}{
		{
			name:   "nil status",
			status: nil,
			want:   "",
		},
		{
			name:   "not a repo",
			status: &Status{IsRepo: false},
			want:   "",
		},
		{
			name: "clean repo",
			status: &Status{
				IsRepo:  true,
				Branch:  "main",
				IsClean: true,
			},
			want: "main",
		},
		{
			name: "with ahead/behind",
			status: &Status{
				IsRepo:  true,
				Branch:  "main",
				Ahead:   2,
				Behind:  1,
				IsClean: true,
			},
			want: "main ↑2 ↓1",
		},
		{
			name: "with changes",
			status: &Status{
				IsRepo:    true,
				Branch:    "feature",
				Staged:    1,
				Modified:  2,
				Untracked: 3,
			},
			want: "feature +1 !2 ?3",
		},
		{
			name: "merging",
			status: &Status{
				IsRepo:    true,
				Branch:    "main",
				IsMerging: true,
			},
			want: "main MERGING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatStatus(tt.status)
			if got != tt.want {
				t.Errorf("FormatStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatStatusCompact(t *testing.T) {
	tests := []struct {
		name   string
		status *Status
		want   string
	}{
		{
			name:   "nil status",
			status: nil,
			want:   "",
		},
		{
			name: "clean repo",
			status: &Status{
				IsRepo:  true,
				Branch:  "main",
				IsClean: true,
			},
			want: "main",
		},
		{
			name: "dirty repo",
			status: &Status{
				IsRepo:   true,
				Branch:   "main",
				Modified: 1,
				IsClean:  false,
			},
			want: "main *",
		},
		{
			name: "ahead",
			status: &Status{
				IsRepo:  true,
				Branch:  "main",
				Ahead:   3,
				IsClean: true,
			},
			want: "main ↑3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatStatusCompact(tt.status)
			if got != tt.want {
				t.Errorf("FormatStatusCompact() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseTrackInfo(t *testing.T) {
	tests := []struct {
		input      string
		wantAhead  int
		wantBehind int
	}{
		{"[ahead 1]", 1, 0},
		{"[behind 2]", 0, 2},
		{"[ahead 3, behind 4]", 3, 4},
		{"", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ahead, behind := parseTrackInfo(tt.input)
			if ahead != tt.wantAhead {
				t.Errorf("parseTrackInfo() ahead = %d, want %d", ahead, tt.wantAhead)
			}
			if behind != tt.wantBehind {
				t.Errorf("parseTrackInfo() behind = %d, want %d", behind, tt.wantBehind)
			}
		})
	}
}
