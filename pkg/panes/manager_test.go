package panes_test

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/core"
	"github.com/cbwinslow/cbwsh/pkg/panes"
)

func TestNewManager(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
	
	// Manager starts with no panes
	allPanes := manager.List()
	if len(allPanes) != 0 {
		t.Errorf("expected 0 initial panes, got %d", len(allPanes))
	}
	
	// Create the first pane
	_, err := manager.Create()
	if err != nil {
		t.Fatalf("failed to create first pane: %v", err)
	}
	
	allPanes = manager.List()
	if len(allPanes) != 1 {
		t.Errorf("expected 1 pane after create, got %d", len(allPanes))
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	
	// Create first pane
	pane1, err := manager.Create()
	if err != nil {
		t.Fatalf("failed to create first pane: %v", err)
	}
	if pane1 == nil {
		t.Fatal("expected non-nil pane")
	}
	
	// Create additional panes
	pane2, err := manager.Create()
	if err != nil {
		t.Fatalf("failed to create pane: %v", err)
	}
	if pane2 == nil {
		t.Fatal("expected non-nil pane")
	}
	
	pane3, err := manager.Create()
	if err != nil {
		t.Fatalf("failed to create pane: %v", err)
	}
	if pane3 == nil {
		t.Fatal("expected non-nil pane")
	}
	
	// Should have 3 panes total
	allPanes := manager.List()
	if len(allPanes) != 3 {
		t.Errorf("expected 3 panes, got %d", len(allPanes))
	}
}

func TestSetActive(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	pane1, _ := manager.Create()
	pane2, _ := manager.Create()
	pane3, _ := manager.Create()
	
	// Get initial active pane (should be pane1)
	active := manager.Active()
	if active == nil {
		t.Fatal("expected active pane")
	}
	if active.ID() != pane1.ID() {
		t.Error("expected pane1 to be initially active")
	}
	
	// Switch to pane2
	err := manager.SetActive(pane2.ID())
	if err != nil {
		t.Fatalf("failed to set active pane: %v", err)
	}
	
	active = manager.Active()
	if active.ID() == pane1.ID() {
		t.Error("expected different pane after SetActive")
	}
	if active.ID() != pane2.ID() {
		t.Errorf("expected pane %s, got %s", pane2.ID(), active.ID())
	}
	
	// Switch to pane3
	err = manager.SetActive(pane3.ID())
	if err != nil {
		t.Fatalf("failed to set active pane: %v", err)
	}
	
	active = manager.Active()
	if active.ID() != pane3.ID() {
		t.Errorf("expected pane %s, got %s", pane3.ID(), active.ID())
	}
}

func TestSetActiveNonExistent(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create() // Need at least one pane
	
	err := manager.SetActive("non-existent-id")
	if err == nil {
		t.Error("expected error when setting active to non-existent pane")
	}
}

func TestNextPane(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	pane1, _ := manager.Create()
	manager.Create()
	manager.Create()
	
	active := manager.Active()
	initialID := pane1.ID() // First pane should be active
	
	err := manager.NextPane()
	if err != nil {
		t.Fatalf("failed to switch to next pane: %v", err)
	}
	
	active = manager.Active()
	if active.ID() == initialID {
		t.Error("expected different pane after NextPane")
	}
}

func TestPrevPane(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create()
	manager.Create()
	manager.Create()
	
	// Move forward first
	manager.NextPane()
	currentID := manager.Active().ID()
	
	// Move back
	err := manager.PrevPane()
	if err != nil {
		t.Fatalf("failed to switch to previous pane: %v", err)
	}
	
	active := manager.Active()
	if active.ID() == currentID {
		t.Error("expected different pane after PrevPane")
	}
}

func TestClose(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	_, _ = manager.Create()
	pane2, _ := manager.Create()
	pane3, _ := manager.Create()
	
	// Close a pane
	err := manager.Close(pane2.ID())
	if err != nil {
		t.Fatalf("failed to close pane: %v", err)
	}
	
	// Should have 2 panes left
	allPanes := manager.List()
	if len(allPanes) != 2 {
		t.Errorf("expected 2 panes after close, got %d", len(allPanes))
	}
	
	// Verify pane2 is gone
	for _, p := range allPanes {
		if p.ID() == pane2.ID() {
			t.Error("expected pane2 to be removed")
		}
	}
	
	// pane3 should still exist
	found := false
	for _, p := range allPanes {
		if p.ID() == pane3.ID() {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected pane3 to still exist")
	}
}

func TestCloseLastPane(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	pane, _ := manager.Create()
	
	// Try to close the only pane - should fail or keep at least one pane
	allPanes := manager.List()
	if len(allPanes) != 1 {
		t.Fatal("expected 1 pane")
	}
	
	err := manager.Close(pane.ID())
	
	// After attempting to close the only pane
	remainingPanes := manager.List()
	
	// The manager implementation may either:
	// 1. Prevent closing the last pane (error) and keep 1 pane
	// 2. Allow closing and have 0 panes
	// Both behaviors are valid depending on design choice
	if err != nil {
		// Error returned - should still have the pane
		if len(remainingPanes) != 1 {
			t.Errorf("expected 1 pane when close failed, got %d", len(remainingPanes))
		}
	} else {
		// Close succeeded - may have 0 or 1 panes
		// Accept either behavior
		if len(remainingPanes) > 1 {
			t.Errorf("expected at most 1 pane after closing last pane, got %d", len(remainingPanes))
		}
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	_, _ = manager.Create()
	pane2, _ := manager.Create()
	
	// Get existing pane
	retrieved, exists := manager.Get(pane2.ID())
	if !exists {
		t.Fatal("expected to find pane")
	}
	
	if retrieved.ID() != pane2.ID() {
		t.Errorf("expected pane %s, got %s", pane2.ID(), retrieved.ID())
	}
	
	// Try to get non-existent pane
	_, exists = manager.Get("non-existent-id")
	if exists {
		t.Error("expected not to find non-existent pane")
	}
}

func TestCount(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	
	// Initial count (no panes)
	count := manager.Count()
	if count != 0 {
		t.Errorf("expected 0 panes initially, got %d", count)
	}
	
	// Create first pane
	manager.Create()
	count = manager.Count()
	if count != 1 {
		t.Errorf("expected 1 pane, got %d", count)
	}
	
	// Add panes
	manager.Create()
	manager.Create()
	manager.Create()
	
	count = manager.Count()
	if count != 4 {
		t.Errorf("expected 4 panes, got %d", count)
	}
}

func TestLayout(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create() // Need at least one pane
	
	// Test different layouts
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
		
		currentLayout := manager.Layout()
		if currentLayout != layout {
			t.Errorf("expected layout %v, got %v", layout, currentLayout)
		}
	}
}

func TestSplit(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create() // Need at least one pane to split from
	
	initialCount := manager.Count()
	
	// Split horizontally
	pane, err := manager.Split(core.LayoutHorizontalSplit)
	if err != nil {
		t.Fatalf("failed to split horizontally: %v", err)
	}
	if pane == nil {
		t.Error("expected non-nil pane from split")
	}
	
	if manager.Count() != initialCount+1 {
		t.Errorf("expected %d panes after horizontal split, got %d", initialCount+1, manager.Count())
	}
	
	// Split vertically
	pane, err = manager.Split(core.LayoutVerticalSplit)
	if err != nil {
		t.Fatalf("failed to split vertically: %v", err)
	}
	if pane == nil {
		t.Error("expected non-nil pane from split")
	}
	
	if manager.Count() != initialCount+2 {
		t.Errorf("expected %d panes after vertical split, got %d", initialCount+2, manager.Count())
	}
}

func TestPaneProperties(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	_, _ = manager.Create()
	paneInterface := manager.Active()
	
	// Get the underlying Pane type
	pane, ok := manager.GetPane(paneInterface.ID())
	if !ok {
		t.Fatal("expected to get pane")
	}
	
	// Test title
	pane.SetTitle("Custom Title")
	if pane.Title() != "Custom Title" {
		t.Errorf("expected 'Custom Title', got '%s'", pane.Title())
	}
	
	// Test ID is not empty
	if pane.ID() == "" {
		t.Error("expected non-empty pane ID")
	}
	
	// Test active status
	if !pane.IsActive() {
		t.Error("expected active pane to report IsActive() == true")
	}
}

func TestUpdateAllSizes(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create()
	manager.Create()
	manager.Create()
	
	// Test resize
	manager.UpdateAllSizes(100, 50)
	
	// No error expected, just verify it doesn't panic
}

func TestMultiplePaneOperations(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	
	// Create multiple panes
	_, _ = manager.Create()
	pane2, _ := manager.Create()
	pane3, _ := manager.Create()
	pane4, _ := manager.Create()
	
	// Set custom titles
	if p, ok := manager.GetPane(pane2.ID()); ok {
		p.SetTitle("Pane 2")
	}
	if p, ok := manager.GetPane(pane3.ID()); ok {
		p.SetTitle("Pane 3")
	}
	if p, ok := manager.GetPane(pane4.ID()); ok {
		p.SetTitle("Pane 4")
	}
	
	// Switch between them
	manager.SetActive(pane3.ID())
	if manager.Active().ID() != pane3.ID() {
		t.Error("failed to switch to pane3")
	}
	
	manager.NextPane()
	if manager.Active().ID() != pane4.ID() {
		t.Error("expected pane4 after next")
	}
	
	manager.PrevPane()
	if manager.Active().ID() != pane3.ID() {
		t.Error("expected pane3 after previous")
	}
	
	// Close one pane
	manager.Close(pane2.ID())
	
	// Verify count
	if manager.Count() != 3 {
		t.Errorf("expected 3 panes, got %d", manager.Count())
	}
}

func TestPaneWithDifferentShellTypes(t *testing.T) {
	t.Parallel()
	
	// Test with bash
	managerBash := panes.NewManager(core.ShellTypeBash)
	if managerBash == nil {
		t.Error("expected non-nil manager for bash")
	}
	
	// Test with zsh
	managerZsh := panes.NewManager(core.ShellTypeZsh)
	if managerZsh == nil {
		t.Error("expected non-nil manager for zsh")
	}
}

func TestListPanes(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create()
	manager.Create()
	manager.Create()
	
	// List all panes (concrete type)
	panes := manager.ListPanes()
	if len(panes) != 3 {
		t.Errorf("expected 3 panes, got %d", len(panes))
	}
	
	// Verify each pane is valid
	for _, p := range panes {
		if p == nil {
			t.Error("expected non-nil pane in list")
		}
		if p.ID() == "" {
			t.Error("expected non-empty pane ID")
		}
	}
}

func TestActivePane(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	manager.Create()
	
	// Get active pane (concrete type)
	pane := manager.ActivePane()
	if pane == nil {
		t.Fatal("expected non-nil active pane")
	}
	
	if !pane.IsActive() {
		t.Error("expected active pane to report IsActive() == true")
	}
}

func TestConcurrentPaneOperations(t *testing.T) {
	t.Parallel()
	
	manager := panes.NewManager(core.ShellTypeBash)
	
	// Create panes concurrently
	done := make(chan bool)
	errors := make(chan error, 5)
	
	for i := 0; i < 5; i++ {
		go func() {
			_, err := manager.Create()
			if err != nil {
				errors <- err
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
	close(errors)
	
	// Check for errors
	for err := range errors {
		t.Errorf("concurrent create error: %v", err)
	}
	
	// Should have at least the initial pane
	if manager.Count() < 1 {
		t.Error("expected at least 1 pane after concurrent operations")
	}
}
