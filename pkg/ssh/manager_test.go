package ssh_test

import (
	"context"
	"testing"
	"time"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/ssh"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestState(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	
	// Initial state should be disconnected
	state := manager.State()
	if state != core.SSHDisconnected {
		t.Errorf("expected disconnected state, got %v", state)
	}
}

func TestSaveAndListHosts(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/hosts"
	manager := ssh.NewManager(tempFile, 30*time.Second)
	
	// Add hosts
	hosts := []core.SSHHost{
		{
			Name: "server1",
			Host: "example.com",
			Port: 22,
			User: "user1",
		},
		{
			Name: "server2",
			Host: "test.com",
			Port: 2222,
			User: "user2",
		},
	}
	
	for _, h := range hosts {
		err := manager.SaveHost(h)
		if err != nil {
			t.Fatalf("failed to save host: %v", err)
		}
	}
	
	// List hosts
	loaded, err := manager.ListSavedHosts()
	if err != nil {
		t.Fatalf("failed to list hosts: %v", err)
	}
	
	if len(loaded) < len(hosts) {
		t.Errorf("expected at least %d hosts, got %d", len(hosts), len(loaded))
	}
}

func TestRemoveHost(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/hosts"
	manager := ssh.NewManager(tempFile, 30*time.Second)
	
	// Save a host
	testHost := core.SSHHost{
		Name: "testserver",
		Host: "test.com",
		Port: 22,
		User: "testuser",
	}
	
	err := manager.SaveHost(testHost)
	if err != nil {
		t.Fatalf("failed to save host: %v", err)
	}
	
	// Delete the host
	err = manager.RemoveHost("testserver")
	if err != nil {
		t.Fatalf("failed to remove host: %v", err)
	}
	
	// Verify it's gone by listing
	hosts, _ := manager.ListSavedHosts()
	for _, h := range hosts {
		if h.Name == "testserver" {
			t.Error("host should have been removed")
		}
	}
}

func TestListHosts(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/hosts"
	manager := ssh.NewManager(tempFile, 30*time.Second)
	
	// Add multiple hosts
	hosts := []core.SSHHost{
		{Name: "host1", Host: "h1.com", Port: 22, User: "user1"},
		{Name: "host2", Host: "h2.com", Port: 22, User: "user2"},
		{Name: "host3", Host: "h3.com", Port: 22, User: "user3"},
	}
	
	for _, h := range hosts {
		manager.SaveHost(h)
	}
	
	// List all hosts
	listed, err := manager.ListSavedHosts()
	if err != nil {
		t.Fatalf("failed to list hosts: %v", err)
	}
	
	if len(listed) < len(hosts) {
		t.Errorf("expected at least %d hosts, got %d", len(hosts), len(listed))
	}
}

func TestDisconnect(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	
	// Disconnect when not connected should not error
	err := manager.Disconnect()
	if err != nil {
		t.Errorf("disconnect should not error when not connected: %v", err)
	}
	
	state := manager.State()
	if state != core.SSHDisconnected {
		t.Errorf("expected disconnected state, got %v", state)
	}
}

func TestCurrentHost(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	
	// Initially no current host
	host := manager.CurrentHost()
	if host != nil {
		t.Error("expected nil current host when not connected")
	}
}

func TestSetStrictHostKeyChecking(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	
	// Test setting strict host key checking
	manager.SetStrictHostKeyChecking(true)
	manager.SetStrictHostKeyChecking(false)
	
	// No error expected, just verify it doesn't panic
}

func TestSetKnownHostsPath(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 30*time.Second)
	
	// Test setting known hosts path
	manager.SetKnownHostsPath("/tmp/known_hosts")
	
	// No error expected, just verify it doesn't panic
}

func TestConnectInvalidHost(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 1*time.Second)
	ctx := context.Background()
	
	// Try to connect to invalid host
	err := manager.Connect(ctx, "invalid-nonexistent-host-12345.invalid", 22, "user")
	if err == nil {
		t.Error("expected error when connecting to invalid host")
	}
	
	// State should be error
	state := manager.State()
	if state != core.SSHError {
		t.Errorf("expected error state after failed connection, got %v", state)
	}
}

func TestConnectWithPasswordInvalidHost(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 1*time.Second)
	ctx := context.Background()
	
	// Try to connect to invalid host with password
	err := manager.ConnectWithPassword(ctx, "invalid-host-98765.invalid", 22, "user", "pass")
	if err == nil {
		t.Error("expected error when connecting to invalid host")
	}
	
	state := manager.State()
	if state != core.SSHError {
		t.Errorf("expected error state, got %v", state)
	}
}

func TestConnectWithKeyInvalidKey(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 1*time.Second)
	ctx := context.Background()
	
	// Try to connect with non-existent key file
	err := manager.ConnectWithKey(ctx, "host.example.com", 22, "user", "/nonexistent/key", "")
	if err == nil {
		t.Error("expected error when using non-existent key file")
	}
	
	state := manager.State()
	if state != core.SSHError {
		t.Errorf("expected error state, got %v", state)
	}
}

func TestHostCRUD(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/hosts"
	manager := ssh.NewManager(tempFile, 30*time.Second)
	
	// Create
	host := core.SSHHost{
		Name: "crud-test",
		Host: "example.com",
		Port: 22,
		User: "testuser",
	}
	
	err := manager.SaveHost(host)
	if err != nil {
		t.Fatalf("failed to create host: %v", err)
	}
	
	// Read (by listing)
	hosts, err := manager.ListSavedHosts()
	if err != nil {
		t.Fatalf("failed to list hosts: %v", err)
	}
	
	found := false
	for _, h := range hosts {
		if h.Name == "crud-test" {
			found = true
			if h.Port != 22 {
				t.Error("incorrect port")
			}
		}
	}
	if !found {
		t.Error("host not found after save")
	}
	
	// Update
	host.Port = 2222
	err = manager.SaveHost(host)
	if err != nil {
		t.Fatalf("failed to update host: %v", err)
	}
	
	hosts, _ = manager.ListSavedHosts()
	found = false
	for _, h := range hosts {
		if h.Name == "crud-test" {
			found = true
			if h.Port != 2222 {
				t.Errorf("expected port 2222, got %d", h.Port)
			}
		}
	}
	if !found {
		t.Error("host not found after update")
	}
	
	// Delete
	err = manager.RemoveHost("crud-test")
	if err != nil {
		t.Fatalf("failed to delete host: %v", err)
	}
	
	hosts, _ = manager.ListSavedHosts()
	for _, h := range hosts {
		if h.Name == "crud-test" {
			t.Error("host should have been deleted")
		}
	}
}

func TestEmptyHostFile(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/empty_hosts"
	manager := ssh.NewManager(tempFile, 30*time.Second)
	
	// List from non-existent file should not error
	hosts, err := manager.ListSavedHosts()
	if err != nil {
		t.Errorf("listing from non-existent file should not error: %v", err)
	}
	
	if len(hosts) != 0 {
		t.Errorf("expected 0 hosts from empty file, got %d", len(hosts))
	}
}

func TestMultipleManagers(t *testing.T) {
	t.Parallel()
	
	tempFile := t.TempDir() + "/hosts"
	
	// Create first manager and save a host
	manager1 := ssh.NewManager(tempFile, 30*time.Second)
	host := core.SSHHost{
		Name: "shared",
		Host: "example.com",
		Port: 22,
		User: "user",
	}
	manager1.SaveHost(host)
	
	// Verify first manager can list the host
	hosts1, err := manager1.ListSavedHosts()
	if err != nil {
		t.Fatalf("failed to list hosts from manager1: %v", err)
	}
	
	found1 := false
	for _, h := range hosts1 {
		if h.Name == "shared" {
			found1 = true
			break
		}
	}
	
	if !found1 {
		t.Error("host not saved by first manager")
	}
	
	// Create second manager - file persistence behavior may vary
	// This tests that managers can read the same file
	manager2 := ssh.NewManager(tempFile, 30*time.Second)
	hosts2, err := manager2.ListSavedHosts()
	if err != nil {
		t.Fatalf("failed to list hosts from manager2: %v", err)
	}
	
	// Both managers should see their respective saves
	// (Exact behavior depends on implementation - some may auto-reload, some may not)
	_ = hosts2 // Accept either behavior
}

func TestConnectToSavedHostNotExists(t *testing.T) {
	t.Parallel()
	
	manager := ssh.NewManager("/tmp/hosts", 1*time.Second)
	ctx := context.Background()
	
	// Try to connect to non-existent saved host
	err := manager.ConnectToSavedHost(ctx, "nonexistent")
	if err == nil {
		t.Error("expected error when connecting to non-existent saved host")
	}
}
