// Package git provides git integration for cbwsh.
package git

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Common errors.
var (
	ErrNotGitRepo    = errors.New("not a git repository")
	ErrNoRemote      = errors.New("no remote configured")
	ErrCommandFailed = errors.New("git command failed")
)

// Status represents the git status.
type Status struct {
	// IsRepo indicates if the directory is in a git repository.
	IsRepo bool
	// Branch is the current branch name.
	Branch string
	// Ahead is the number of commits ahead of upstream.
	Ahead int
	// Behind is the number of commits behind upstream.
	Behind int
	// Staged is the count of staged files.
	Staged int
	// Modified is the count of modified files.
	Modified int
	// Untracked is the count of untracked files.
	Untracked int
	// Deleted is the count of deleted files.
	Deleted int
	// Conflicted is the count of conflicted files.
	Conflicted int
	// Stashed is the count of stashed entries.
	Stashed int
	// IsClean indicates no pending changes.
	IsClean bool
	// IsMerging indicates a merge in progress.
	IsMerging bool
	// IsRebasing indicates a rebase in progress.
	IsRebasing bool
	// IsCherryPicking indicates a cherry-pick in progress.
	IsCherryPicking bool
	// IsBisecting indicates a bisect in progress.
	IsBisecting bool
	// RemoteBranch is the upstream branch name.
	RemoteBranch string
	// RemoteURL is the remote origin URL.
	RemoteURL string
}

// Commit represents a git commit.
type Commit struct {
	Hash       string
	ShortHash  string
	Author     string
	AuthorDate string
	Message    string
}

// Branch represents a git branch.
type Branch struct {
	Name      string
	IsCurrent bool
	IsRemote  bool
	Upstream  string
	Ahead     int
	Behind    int
}

// Repo provides git repository operations.
type Repo struct {
	mu   sync.RWMutex
	path string
}

// NewRepo creates a new git repo interface for the given path.
func NewRepo(path string) *Repo {
	return &Repo{path: path}
}

// IsRepo returns true if the path is inside a git repository.
func (r *Repo) IsRepo() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = r.path
	output, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(output)) == "true"
}

// Status returns the current git status.
func (r *Repo) Status() (*Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return nil, ErrNotGitRepo
	}

	status := &Status{IsRepo: true}

	// Get branch name
	status.Branch = r.getBranch()

	// Get ahead/behind
	status.Ahead, status.Behind = r.getAheadBehind()

	// Get file counts
	r.getFileCounts(status)

	// Get stash count
	status.Stashed = r.getStashCount()

	// Check for operations in progress
	r.checkOperations(status)

	// Get remote info
	status.RemoteBranch = r.getUpstreamBranch()
	status.RemoteURL = r.getRemoteURL()

	// Determine if clean
	status.IsClean = status.Staged == 0 && status.Modified == 0 &&
		status.Untracked == 0 && status.Deleted == 0 && status.Conflicted == 0

	return status, nil
}

// CurrentBranch returns the current branch name.
func (r *Repo) CurrentBranch() (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return "", ErrNotGitRepo
	}

	return r.getBranch(), nil
}

// Branches returns all local branches.
func (r *Repo) Branches() ([]Branch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return nil, ErrNotGitRepo
	}

	cmd := exec.Command("git", "branch", "-vv", "--format=%(refname:short)|%(upstream:short)|%(upstream:track)")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrCommandFailed
	}

	currentBranch := r.getBranch()
	var branches []Branch

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		branch := Branch{
			Name:      parts[0],
			IsCurrent: parts[0] == currentBranch,
		}

		if len(parts) > 1 && parts[1] != "" {
			branch.Upstream = parts[1]
		}

		if len(parts) > 2 && parts[2] != "" {
			branch.Ahead, branch.Behind = parseTrackInfo(parts[2])
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// RemoteBranches returns all remote branches.
func (r *Repo) RemoteBranches() ([]Branch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return nil, ErrNotGitRepo
	}

	cmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrCommandFailed
	}

	var branches []Branch
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		branches = append(branches, Branch{
			Name:     line,
			IsRemote: true,
		})
	}

	return branches, nil
}

// Log returns the last n commits.
func (r *Repo) Log(n int) ([]Commit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return nil, ErrNotGitRepo
	}

	format := "%H|%h|%an|%ar|%s"
	cmd := exec.Command("git", "log", "-n", strconv.Itoa(n), "--format="+format)
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrCommandFailed
	}

	var commits []Commit
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}

		commits = append(commits, Commit{
			Hash:       parts[0],
			ShortHash:  parts[1],
			Author:     parts[2],
			AuthorDate: parts[3],
			Message:    parts[4],
		})
	}

	return commits, nil
}

// Diff returns the current diff.
func (r *Repo) Diff() (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return "", ErrNotGitRepo
	}

	cmd := exec.Command("git", "diff")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return "", ErrCommandFailed
	}

	return string(output), nil
}

// DiffStaged returns the staged diff.
func (r *Repo) DiffStaged() (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return "", ErrNotGitRepo
	}

	cmd := exec.Command("git", "diff", "--staged")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return "", ErrCommandFailed
	}

	return string(output), nil
}

// Tags returns all tags.
func (r *Repo) Tags() ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if !r.isRepo() {
		return nil, ErrNotGitRepo
	}

	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return nil, ErrCommandFailed
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var tags []string
	for _, line := range lines {
		if line != "" {
			tags = append(tags, line)
		}
	}

	return tags, nil
}

// Internal helpers

func (r *Repo) isRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = r.path
	output, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(output)) == "true"
}

func (r *Repo) getBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (r *Repo) getAheadBehind() (int, int) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return 0, 0
	}

	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])
	return ahead, behind
}

func (r *Repo) getFileCounts(status *Status) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		x, y := line[0], line[1]

		// Staged files (index status)
		if x == 'M' || x == 'A' || x == 'D' || x == 'R' || x == 'C' {
			status.Staged++
		}

		// Modified files (work tree status)
		if y == 'M' {
			status.Modified++
		}

		// Deleted files
		if y == 'D' {
			status.Deleted++
		}

		// Untracked files
		if x == '?' && y == '?' {
			status.Untracked++
		}

		// Conflicts
		if x == 'U' || y == 'U' {
			status.Conflicted++
		}
	}
}

func (r *Repo) getStashCount() int {
	cmd := exec.Command("git", "stash", "list")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return 0
	}
	return len(lines)
}

func (r *Repo) checkOperations(status *Status) {
	gitDir := r.getGitDir()
	if gitDir == "" {
		return
	}

	// Check for merge in progress
	cmd := exec.Command("git", "rev-parse", "-q", "--verify", "MERGE_HEAD")
	cmd.Dir = r.path
	if err := cmd.Run(); err == nil {
		status.IsMerging = true
	}

	// Check for rebase in progress
	cmd = exec.Command("git", "rev-parse", "-q", "--verify", "REBASE_HEAD")
	cmd.Dir = r.path
	if err := cmd.Run(); err == nil {
		status.IsRebasing = true
	}

	// Check for cherry-pick in progress
	cmd = exec.Command("git", "rev-parse", "-q", "--verify", "CHERRY_PICK_HEAD")
	cmd.Dir = r.path
	if err := cmd.Run(); err == nil {
		status.IsCherryPicking = true
	}
}

func (r *Repo) getGitDir() string {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (r *Repo) getUpstreamBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func (r *Repo) getRemoteURL() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = r.path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func parseTrackInfo(track string) (int, int) {
	// Parse format like "[ahead 1, behind 2]" or "[ahead 1]"
	var ahead, behind int

	if strings.Contains(track, "ahead") {
		parts := strings.Split(track, "ahead ")
		if len(parts) > 1 {
			numStr := strings.Split(parts[1], ",")[0]
			numStr = strings.TrimRight(numStr, "]")
			ahead, _ = strconv.Atoi(strings.TrimSpace(numStr))
		}
	}

	if strings.Contains(track, "behind") {
		parts := strings.Split(track, "behind ")
		if len(parts) > 1 {
			numStr := strings.TrimRight(parts[1], "]")
			behind, _ = strconv.Atoi(strings.TrimSpace(numStr))
		}
	}

	return ahead, behind
}

// FormatStatus formats the git status for display.
func FormatStatus(s *Status) string {
	if s == nil || !s.IsRepo {
		return ""
	}

	var parts []string

	// Branch
	parts = append(parts, s.Branch)

	// Ahead/behind
	if s.Ahead > 0 {
		parts = append(parts, "↑"+strconv.Itoa(s.Ahead))
	}
	if s.Behind > 0 {
		parts = append(parts, "↓"+strconv.Itoa(s.Behind))
	}

	// File status
	if s.Staged > 0 {
		parts = append(parts, "+"+strconv.Itoa(s.Staged))
	}
	if s.Modified > 0 {
		parts = append(parts, "!"+strconv.Itoa(s.Modified))
	}
	if s.Untracked > 0 {
		parts = append(parts, "?"+strconv.Itoa(s.Untracked))
	}
	if s.Conflicted > 0 {
		parts = append(parts, "✖"+strconv.Itoa(s.Conflicted))
	}

	// Stash
	if s.Stashed > 0 {
		parts = append(parts, "⚑"+strconv.Itoa(s.Stashed))
	}

	// Operations
	if s.IsMerging {
		parts = append(parts, "MERGING")
	}
	if s.IsRebasing {
		parts = append(parts, "REBASING")
	}
	if s.IsCherryPicking {
		parts = append(parts, "CHERRY-PICKING")
	}

	return strings.Join(parts, " ")
}

// FormatStatusCompact formats the git status in a compact form.
func FormatStatusCompact(s *Status) string {
	if s == nil || !s.IsRepo {
		return ""
	}

	result := s.Branch

	if s.Ahead > 0 || s.Behind > 0 {
		result += " "
		if s.Ahead > 0 {
			result += "↑" + strconv.Itoa(s.Ahead)
		}
		if s.Behind > 0 {
			result += "↓" + strconv.Itoa(s.Behind)
		}
	}

	if !s.IsClean {
		result += " *"
	}

	return result
}
