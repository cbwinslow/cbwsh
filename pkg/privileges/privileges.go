// Package privileges provides privilege checking and elevation for cbwsh.
package privileges

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// PrivilegeLevel represents the level of privileges.
type PrivilegeLevel int

const (
	// PrivilegeLevelUser is normal user privileges.
	PrivilegeLevelUser PrivilegeLevel = iota
	// PrivilegeLevelRoot is root/administrator privileges.
	PrivilegeLevelRoot
)

// String returns the string representation of the privilege level.
func (p PrivilegeLevel) String() string {
	switch p {
	case PrivilegeLevelUser:
		return "user"
	case PrivilegeLevelRoot:
		return "root"
	default:
		return "unknown"
	}
}

// ElevationMethod represents how to elevate privileges.
type ElevationMethod int

const (
	// ElevationMethodSudo uses sudo for elevation.
	ElevationMethodSudo ElevationMethod = iota
	// ElevationMethodSu uses su for elevation.
	ElevationMethodSu
	// ElevationMethodPkexec uses polkit for elevation.
	ElevationMethodPkexec
	// ElevationMethodDoas uses doas (OpenBSD) for elevation.
	ElevationMethodDoas
)

// String returns the string representation of the elevation method.
func (m ElevationMethod) String() string {
	switch m {
	case ElevationMethodSudo:
		return "sudo"
	case ElevationMethodSu:
		return "su"
	case ElevationMethodPkexec:
		return "pkexec"
	case ElevationMethodDoas:
		return "doas"
	default:
		return "unknown"
	}
}

// UserInfo holds information about the current user.
type UserInfo struct {
	UID          int
	GID          int
	Username     string
	Name         string
	HomeDir      string
	Shell        string
	Groups       []string
	IsRoot       bool
	EffectiveUID int
	EffectiveGID int
}

// Manager provides privilege management functionality.
type Manager struct {
	mu              sync.RWMutex
	userInfo        *UserInfo
	elevationMethod ElevationMethod
}

// NewManager creates a new privilege manager.
func NewManager() *Manager {
	m := &Manager{
		elevationMethod: ElevationMethodSudo,
	}
	m.refreshUserInfo()
	return m
}

// refreshUserInfo refreshes the user information.
func (m *Manager) refreshUserInfo() {
	m.mu.Lock()
	defer m.mu.Unlock()

	currentUser, err := user.Current()
	if err != nil {
		return
	}

	uid, _ := strconv.Atoi(currentUser.Uid)
	gid, _ := strconv.Atoi(currentUser.Gid)

	m.userInfo = &UserInfo{
		UID:          uid,
		GID:          gid,
		Username:     currentUser.Username,
		Name:         currentUser.Name,
		HomeDir:      currentUser.HomeDir,
		IsRoot:       uid == 0,
		EffectiveUID: syscall.Geteuid(),
		EffectiveGID: syscall.Getegid(),
	}

	// Get user's shell
	if shell := os.Getenv("SHELL"); shell != "" {
		m.userInfo.Shell = shell
	}

	// Get user's groups
	groups, err := currentUser.GroupIds()
	if err == nil {
		m.userInfo.Groups = groups
	}
}

// GetUserInfo returns the current user information.
func (m *Manager) GetUserInfo() *UserInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.userInfo == nil {
		return nil
	}

	// Return a copy
	info := *m.userInfo
	groups := make([]string, len(m.userInfo.Groups))
	copy(groups, m.userInfo.Groups)
	info.Groups = groups
	return &info
}

// IsRoot returns whether the current process is running as root.
func (m *Manager) IsRoot() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.userInfo != nil && m.userInfo.IsRoot
}

// IsEffectiveRoot returns whether the effective user is root.
func (m *Manager) IsEffectiveRoot() bool {
	return syscall.Geteuid() == 0
}

// GetCurrentLevel returns the current privilege level.
func (m *Manager) GetCurrentLevel() PrivilegeLevel {
	if m.IsEffectiveRoot() {
		return PrivilegeLevelRoot
	}
	return PrivilegeLevelUser
}

// SetElevationMethod sets the method to use for privilege elevation.
func (m *Manager) SetElevationMethod(method ElevationMethod) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.elevationMethod = method
}

// GetElevationMethod returns the current elevation method.
func (m *Manager) GetElevationMethod() ElevationMethod {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.elevationMethod
}

// CanElevate checks if privileges can be elevated using the current method.
func (m *Manager) CanElevate() bool {
	m.mu.RLock()
	method := m.elevationMethod
	m.mu.RUnlock()

	switch method {
	case ElevationMethodSudo:
		return commandExists("sudo")
	case ElevationMethodSu:
		return commandExists("su")
	case ElevationMethodPkexec:
		return commandExists("pkexec")
	case ElevationMethodDoas:
		return commandExists("doas")
	default:
		return false
	}
}

// ElevatedCommand wraps a command for elevated execution.
func (m *Manager) ElevatedCommand(ctx context.Context, command string) (*exec.Cmd, error) {
	m.mu.RLock()
	method := m.elevationMethod
	m.mu.RUnlock()

	var cmd *exec.Cmd
	switch method {
	case ElevationMethodSudo:
		cmd = exec.CommandContext(ctx, "sudo", "-S", "bash", "-c", command)
	case ElevationMethodSu:
		cmd = exec.CommandContext(ctx, "su", "-c", command)
	case ElevationMethodPkexec:
		cmd = exec.CommandContext(ctx, "pkexec", "bash", "-c", command)
	case ElevationMethodDoas:
		cmd = exec.CommandContext(ctx, "doas", "bash", "-c", command)
	default:
		return nil, errors.New("unknown elevation method")
	}

	return cmd, nil
}

// ExecuteElevated executes a command with elevated privileges.
func (m *Manager) ExecuteElevated(ctx context.Context, command string) (string, error) {
	cmd, err := m.ElevatedCommand(ctx, command)
	if err != nil {
		return "", err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("elevated command failed: %w", err)
	}

	return string(output), nil
}

// CheckSudoAccess checks if the user has sudo access without a password.
func (m *Manager) CheckSudoAccess() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}

// InGroup checks if the user is in the specified group.
func (m *Manager) InGroup(groupName string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.userInfo == nil {
		return false
	}

	// Try to find group by name
	grp, err := user.LookupGroup(groupName)
	if err != nil {
		return false
	}

	for _, gid := range m.userInfo.Groups {
		if gid == grp.Gid {
			return true
		}
	}

	return false
}

// GetCapabilities returns the capabilities of the current process.
func (m *Manager) GetCapabilities() []string {
	var caps []string

	// Check common capabilities
	if m.IsEffectiveRoot() {
		caps = append(caps, "CAP_ALL")
		return caps
	}

	// Check specific capabilities by testing operations
	// This is a simplified check

	// Check if we can bind to low ports (CAP_NET_BIND_SERVICE)
	if m.CheckSudoAccess() {
		caps = append(caps, "SUDO_ACCESS")
	}

	// Check if user is in sudo/wheel group
	if m.InGroup("sudo") || m.InGroup("wheel") || m.InGroup("admin") {
		caps = append(caps, "ADMIN_GROUP")
	}

	return caps
}

// RequiresElevation checks if a command requires elevated privileges.
func RequiresElevation(command string) bool {
	// Commands that typically require root
	rootCommands := []string{
		"mount", "umount", "fdisk", "mkfs", "iptables", "ip6tables",
		"systemctl", "service", "journalctl", "apt", "apt-get", "yum",
		"dnf", "pacman", "apk", "pkg", "useradd", "userdel", "usermod",
		"groupadd", "groupdel", "groupmod", "chown", "chmod", "passwd",
		"visudo", "shutdown", "reboot", "poweroff", "halt", "init",
		"modprobe", "insmod", "rmmod", "lsmod", "dmesg", "sysctl",
		"crontab", "at", "dpkg", "rpm",
	}

	// Parse command to get the base command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return false
	}

	baseCmd := parts[0]
	// Remove path if present
	if idx := strings.LastIndex(baseCmd, "/"); idx >= 0 {
		baseCmd = baseCmd[idx+1:]
	}

	for _, rc := range rootCommands {
		if baseCmd == rc {
			return true
		}
	}

	// Check for paths that require root
	for _, arg := range parts {
		if strings.HasPrefix(arg, "/etc/") ||
			strings.HasPrefix(arg, "/usr/") ||
			strings.HasPrefix(arg, "/var/") ||
			strings.HasPrefix(arg, "/sys/") ||
			strings.HasPrefix(arg, "/proc/") ||
			strings.HasPrefix(arg, "/boot/") {
			return true
		}
	}

	return false
}

// FilePermissions holds file permission information.
type FilePermissions struct {
	Path       string
	Mode       os.FileMode
	OwnerUID   int
	OwnerGID   int
	Readable   bool
	Writable   bool
	Executable bool
}

// CheckFilePermissions checks file permissions for the current user.
func CheckFilePermissions(path string) (*FilePermissions, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, errors.New("failed to get file stat")
	}

	perms := &FilePermissions{
		Path:     path,
		Mode:     info.Mode(),
		OwnerUID: int(stat.Uid),
		OwnerGID: int(stat.Gid),
	}

	// Check readable
	_, err = os.Open(path)
	perms.Readable = err == nil

	// Check writable
	if info.IsDir() {
		// For directories, try to create a temp file
		tmpPath := path + "/.cbwsh_perm_check"
		f, err := os.Create(tmpPath)
		if err == nil {
			f.Close()
			os.Remove(tmpPath)
			perms.Writable = true
		}
	} else {
		f, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err == nil {
			f.Close()
			perms.Writable = true
		}
	}

	// Check executable
	perms.Executable = info.Mode()&0o111 != 0

	return perms, nil
}

// DropPrivileges drops root privileges to the specified user.
func DropPrivileges(uid, gid int) error {
	// Set supplementary groups
	if err := syscall.Setgroups([]int{gid}); err != nil {
		return fmt.Errorf("failed to set groups: %w", err)
	}

	// Set GID first (must be done before UID)
	if err := syscall.Setgid(gid); err != nil {
		return fmt.Errorf("failed to set GID: %w", err)
	}

	// Set UID
	if err := syscall.Setuid(uid); err != nil {
		return fmt.Errorf("failed to set UID: %w", err)
	}

	// Verify we can't regain privileges
	if err := syscall.Setuid(0); err == nil {
		return errors.New("failed to drop privileges permanently")
	}

	return nil
}

// commandExists checks if a command exists in PATH.
func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetSecureEnvironment returns a cleaned environment for elevated commands.
func GetSecureEnvironment() []string {
	// Only keep essential environment variables
	safeVars := []string{
		"PATH",
		"HOME",
		"USER",
		"LOGNAME",
		"SHELL",
		"TERM",
		"LANG",
		"LC_ALL",
	}

	var result []string
	for _, v := range safeVars {
		if val := os.Getenv(v); val != "" {
			result = append(result, fmt.Sprintf("%s=%s", v, val))
		}
	}

	return result
}
