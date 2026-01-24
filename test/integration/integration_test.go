// Package integration provides integration tests for cbwsh components.
package integration_test

import (
"testing"
"time"

"github.com/cbwinslow/cbwsh/pkg/ai"
"github.com/cbwinslow/cbwsh/pkg/core"
"github.com/cbwinslow/cbwsh/pkg/panes"
"github.com/cbwinslow/cbwsh/pkg/secrets"
"github.com/cbwinslow/cbwsh/pkg/shell"
"github.com/cbwinslow/cbwsh/pkg/ssh"
)

// TestShellExecutorMultiplePanes tests shell execution across multiple panes
func TestShellExecutorMultiplePanes(t *testing.T) {
t.Parallel()

manager := panes.NewManager(core.ShellTypeBash)
pane1, _ := manager.Create()
pane2, _ := manager.Create()
_, _ = manager.Create()

if manager.Count() != 3 {
t.Errorf("expected 3 panes, got %d", manager.Count())
}

manager.SetActive(pane1.ID())
if manager.Active().ID() != pane1.ID() {
t.Error("failed to activate pane1")
}

manager.NextPane()
if manager.Active().ID() != pane2.ID() {
t.Error("expected pane2 after next")
}
}

// TestSecretsWithSSH tests secrets manager integration with SSH
func TestSecretsWithSSH(t *testing.T) {
t.Parallel()

tempDir := t.TempDir()
secretsMgr := secrets.NewManager(tempDir + "/secrets.enc")
err := secretsMgr.Initialize("test-password")
if err != nil {
t.Fatalf("failed to initialize secrets: %v", err)
}

err = secretsMgr.Store("ssh-password", []byte("secret-password"))
if err != nil {
t.Fatalf("failed to store secret: %v", err)
}

password, err := secretsMgr.Retrieve("ssh-password")
if err != nil {
t.Fatalf("failed to retrieve secret: %v", err)
}

if string(password) != "secret-password" {
t.Errorf("expected 'secret-password', got %s", password)
}

sshMgr := ssh.NewManager(tempDir+"/hosts", 5*time.Second)
host := core.SSHHost{
Name: "secure-host",
Host: "example.com",
Port: 22,
User: "testuser",
}

err = sshMgr.SaveHost(host)
if err != nil {
t.Fatalf("failed to save host: %v", err)
}

hosts, err := sshMgr.ListSavedHosts()
if err != nil {
t.Fatalf("failed to list hosts: %v", err)
}

found := false
for _, h := range hosts {
if h.Name == "secure-host" {
found = true
break
}
}

if !found {
t.Error("host not found in saved hosts")
}
}

// TestShellHistoryPersistence tests shell history across sessions
func TestShellHistoryPersistence(t *testing.T) {
t.Parallel()

tempFile := t.TempDir() + "/history"
history1 := shell.NewHistory(100, tempFile)
history1.Add("command1")
history1.Add("command2")
history1.Add("command3")

err := history1.Save()
if err != nil {
t.Fatalf("failed to save history: %v", err)
}

history2 := shell.NewHistory(100, tempFile)
err = history2.Load()
if err != nil {
t.Fatalf("failed to load history: %v", err)
}

all := history2.All()
if len(all) != 3 {
t.Errorf("expected 3 commands, got %d", len(all))
}
}

// TestAIAgentManager tests AI agent management
func TestAIAgentManager(t *testing.T) {
t.Parallel()

manager := ai.NewManager()
agent1 := ai.NewAgent("agent1", core.AIProviderLocal, "", "")
agent2 := ai.NewAgent("agent2", core.AIProviderLocal, "", "")

manager.RegisterAgent(agent1)
manager.RegisterAgent(agent2)

agents := manager.ListAgents()
if len(agents) != 2 {
t.Errorf("expected 2 agents, got %d", len(agents))
}

retrieved, exists := manager.GetAgent("agent1")
if !exists {
t.Fatal("expected to find agent1")
}
if retrieved.Name() != "agent1" {
t.Errorf("expected 'agent1', got %s", retrieved.Name())
}
}

// TestPaneLayoutManagement tests pane layout switching
func TestPaneLayoutManagement(t *testing.T) {
t.Parallel()

manager := panes.NewManager(core.ShellTypeBash)
manager.Create()
manager.Create()
manager.Create()

layouts := []core.PaneLayout{
core.LayoutSingle,
core.LayoutHorizontalSplit,
core.LayoutVerticalSplit,
core.LayoutGrid,
}

for _, layout := range layouts {
err := manager.SetLayout(layout)
if err != nil {
t.Errorf("failed to set layout %v: %v", layout, err)
}

current := manager.Layout()
if current != layout {
t.Errorf("expected layout %v, got %v", layout, current)
}
}

manager.UpdateAllSizes(120, 40)

if manager.Count() != 3 {
t.Errorf("expected 3 panes after layout changes, got %d", manager.Count())
}
}
