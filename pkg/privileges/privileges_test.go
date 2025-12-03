package privileges_test

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/privileges"
)

func TestPrivilegeLevelString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level    privileges.PrivilegeLevel
		expected string
	}{
		{privileges.PrivilegeLevelUser, "user"},
		{privileges.PrivilegeLevelRoot, "root"},
		{privileges.PrivilegeLevel(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestElevationMethodString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		method   privileges.ElevationMethod
		expected string
	}{
		{privileges.ElevationMethodSudo, "sudo"},
		{privileges.ElevationMethodSu, "su"},
		{privileges.ElevationMethodPkexec, "pkexec"},
		{privileges.ElevationMethodDoas, "doas"},
		{privileges.ElevationMethod(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			result := tt.method.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewManager(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestGetUserInfo(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()
	info := manager.GetUserInfo()

	if info == nil {
		t.Fatal("expected non-nil user info")
	}

	if info.Username == "" {
		t.Error("expected non-empty username")
	}

	if info.HomeDir == "" {
		t.Error("expected non-empty home directory")
	}
}

func TestGetCurrentLevel(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()
	level := manager.GetCurrentLevel()

	// Most tests run as non-root
	if level != privileges.PrivilegeLevelUser && level != privileges.PrivilegeLevelRoot {
		t.Errorf("unexpected privilege level: %s", level)
	}
}

func TestSetElevationMethod(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()

	manager.SetElevationMethod(privileges.ElevationMethodDoas)

	if manager.GetElevationMethod() != privileges.ElevationMethodDoas {
		t.Error("expected doas elevation method")
	}
}

func TestRequiresElevation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		command  string
		requires bool
	}{
		{"ls -la", false},
		{"echo hello", false},
		{"cat /tmp/test.txt", false},
		{"apt update", true},
		{"systemctl restart nginx", true},
		{"mount /dev/sda1 /mnt", true},
		{"chmod 755 /etc/hosts", true},
		{"cat /etc/passwd", true},
		{"useradd testuser", true},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			t.Parallel()
			result := privileges.RequiresElevation(tt.command)
			if result != tt.requires {
				t.Errorf("command '%s': expected %v, got %v", tt.command, tt.requires, result)
			}
		})
	}
}

func TestCheckFilePermissions(t *testing.T) {
	t.Parallel()

	perms, err := privileges.CheckFilePermissions("/tmp")
	if err != nil {
		t.Fatalf("failed to check permissions: %v", err)
	}

	if perms.Path != "/tmp" {
		t.Errorf("expected path /tmp, got %s", perms.Path)
	}

	if !perms.Readable {
		t.Error("expected /tmp to be readable")
	}

	if !perms.Writable {
		t.Error("expected /tmp to be writable")
	}
}

func TestCheckFilePermissionsNonExistent(t *testing.T) {
	t.Parallel()

	_, err := privileges.CheckFilePermissions("/nonexistent/path")
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestGetSecureEnvironment(t *testing.T) {
	t.Parallel()

	env := privileges.GetSecureEnvironment()

	// Should contain at least PATH
	hasPath := false
	for _, v := range env {
		if len(v) >= 5 && v[:5] == "PATH=" {
			hasPath = true
			break
		}
	}

	if !hasPath {
		t.Error("expected PATH in secure environment")
	}
}

func TestInGroup(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()

	// Check for a non-existent group
	if manager.InGroup("nonexistentgroup12345") {
		t.Error("expected to not be in nonexistent group")
	}
}

func TestCanElevate(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()
	manager.SetElevationMethod(privileges.ElevationMethodSudo)

	// This should return true if sudo is installed
	// Just verify it doesn't panic
	_ = manager.CanElevate()
}

func TestIsRoot(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()

	// Most tests run as non-root
	// Just verify it doesn't panic and returns a boolean
	_ = manager.IsRoot()
	_ = manager.IsEffectiveRoot()
}

func TestGetCapabilities(t *testing.T) {
	t.Parallel()

	manager := privileges.NewManager()
	caps := manager.GetCapabilities()

	// Just verify it returns a slice without panic
	if caps == nil {
		t.Error("expected non-nil capabilities slice")
	}
}
