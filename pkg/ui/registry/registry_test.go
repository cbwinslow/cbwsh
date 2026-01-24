package registry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New() returned nil")
	}
	if r.components == nil {
		t.Error("components map is nil")
	}
	if r.categories == nil {
		t.Error("categories map is nil")
	}
}

func TestRegister(t *testing.T) {
	r := New()
	
	component := &ComponentMetadata{
		Name:        "test-button",
		DisplayName: "Test Button",
		Description: "A test button component",
		Category:    "interactive",
		Version:     "1.0.0",
	}
	
	err := r.Register(component)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}
	
	// Test duplicate registration (should not add duplicate to category)
	err = r.Register(component)
	if err != nil {
		t.Errorf("Register() should allow duplicate registration: %v", err)
	}
	
	// Verify component is in category only once
	categoryComponents := r.ListByCategory("interactive")
	if len(categoryComponents) != 1 {
		t.Errorf("ListByCategory() returned %d components after duplicate registration, want 1", len(categoryComponents))
	}
	
	// Test empty name
	emptyComponent := &ComponentMetadata{
		Name: "",
	}
	err = r.Register(emptyComponent)
	if err == nil {
		t.Error("Register() should fail for empty name")
	}
}

func TestGet(t *testing.T) {
	r := New()
	
	component := &ComponentMetadata{
		Name:        "test-card",
		DisplayName: "Test Card",
		Description: "A test card component",
		Category:    "layout",
		Version:     "1.0.0",
	}
	
	_ = r.Register(component)
	
	// Test successful get
	got, err := r.Get("test-card")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}
	if got.Name != component.Name {
		t.Errorf("Get() = %v, want %v", got.Name, component.Name)
	}
	
	// Test not found
	_, err = r.Get("nonexistent")
	if err == nil {
		t.Error("Get() should fail for nonexistent component")
	}
}

func TestList(t *testing.T) {
	r := New()
	
	components := []*ComponentMetadata{
		{Name: "button", DisplayName: "Button", Category: "interactive", Version: "1.0.0"},
		{Name: "card", DisplayName: "Card", Category: "layout", Version: "1.0.0"},
		{Name: "badge", DisplayName: "Badge", Category: "display", Version: "1.0.0"},
	}
	
	for _, c := range components {
		_ = r.Register(c)
	}
	
	list := r.List()
	if len(list) != 3 {
		t.Errorf("List() returned %d components, want 3", len(list))
	}
	
	// Check if sorted by name
	if list[0].Name != "badge" || list[1].Name != "button" || list[2].Name != "card" {
		t.Error("List() is not sorted by name")
	}
}

func TestListByCategory(t *testing.T) {
	r := New()
	
	components := []*ComponentMetadata{
		{Name: "button", DisplayName: "Button", Category: "interactive", Version: "1.0.0"},
		{Name: "card", DisplayName: "Card", Category: "layout", Version: "1.0.0"},
		{Name: "badge", DisplayName: "Badge", Category: "display", Version: "1.0.0"},
		{Name: "stack", DisplayName: "Stack", Category: "layout", Version: "1.0.0"},
	}
	
	for _, c := range components {
		_ = r.Register(c)
	}
	
	layoutComponents := r.ListByCategory("layout")
	if len(layoutComponents) != 2 {
		t.Errorf("ListByCategory('layout') returned %d components, want 2", len(layoutComponents))
	}
	
	// Check if sorted
	if layoutComponents[0].Name != "card" || layoutComponents[1].Name != "stack" {
		t.Error("ListByCategory() is not sorted by name")
	}
}

func TestCategories(t *testing.T) {
	r := New()
	
	components := []*ComponentMetadata{
		{Name: "button", Category: "interactive", Version: "1.0.0"},
		{Name: "card", Category: "layout", Version: "1.0.0"},
		{Name: "badge", Category: "display", Version: "1.0.0"},
	}
	
	for _, c := range components {
		_ = r.Register(c)
	}
	
	categories := r.Categories()
	if len(categories) != 3 {
		t.Errorf("Categories() returned %d categories, want 3", len(categories))
	}
	
	// Check if sorted
	if categories[0] != "display" || categories[1] != "interactive" || categories[2] != "layout" {
		t.Error("Categories() is not sorted")
	}
}

func TestSearch(t *testing.T) {
	r := New()
	
	components := []*ComponentMetadata{
		{
			Name:        "button",
			DisplayName: "Button",
			Description: "A button for interactions",
			Category:    "interactive",
			Tags:        []string{"action", "click"},
			Version:     "1.0.0",
		},
		{
			Name:        "card",
			DisplayName: "Card",
			Description: "A container for content",
			Category:    "layout",
			Tags:        []string{"container", "box"},
			Version:     "1.0.0",
		},
	}
	
	for _, c := range components {
		_ = r.Register(c)
	}
	
	tests := []struct {
		query    string
		wantName string
		wantLen  int
	}{
		{"button", "button", 1},
		{"Button", "button", 1},
		{"BUTTON", "button", 1},
		{"interaction", "button", 1},
		{"container", "card", 1},
		{"layout", "card", 1},
		{"click", "button", 1},
		{"nonexistent", "", 0},
	}
	
	for _, tt := range tests {
		results := r.Search(tt.query)
		if len(results) != tt.wantLen {
			t.Errorf("Search(%q) returned %d results, want %d", tt.query, len(results), tt.wantLen)
		}
		if tt.wantLen > 0 && results[0].Name != tt.wantName {
			t.Errorf("Search(%q) returned component %q, want %q", tt.query, results[0].Name, tt.wantName)
		}
	}
}

func TestSaveAndLoadFromFile(t *testing.T) {
	r := New()
	
	components := []*ComponentMetadata{
		{Name: "button", DisplayName: "Button", Category: "interactive", Version: "1.0.0"},
		{Name: "card", DisplayName: "Card", Category: "layout", Version: "1.0.0"},
	}
	
	for _, c := range components {
		_ = r.Register(c)
	}
	
	// Create temp directory
	tmpDir := t.TempDir()
	registryFile := filepath.Join(tmpDir, "registry.json")
	
	// Save
	err := r.SaveToFile(registryFile)
	if err != nil {
		t.Fatalf("SaveToFile() failed: %v", err)
	}
	
	// Verify file exists
	if _, err := os.Stat(registryFile); os.IsNotExist(err) {
		t.Fatal("Registry file was not created")
	}
	
	// Load into new registry
	r2 := New()
	err = r2.LoadFromFile(registryFile)
	if err != nil {
		t.Fatalf("LoadFromFile() failed: %v", err)
	}
	
	// Verify components
	list := r2.List()
	if len(list) != 2 {
		t.Errorf("Loaded registry has %d components, want 2", len(list))
	}
	
	button, err := r2.Get("button")
	if err != nil {
		t.Errorf("Get('button') failed: %v", err)
	}
	if button.DisplayName != "Button" {
		t.Errorf("button.DisplayName = %q, want %q", button.DisplayName, "Button")
	}
}

func TestDefaultRegistry(t *testing.T) {
	r := DefaultRegistry()
	
	// Check that we have the expected built-in components
	expectedComponents := []string{"button", "card", "badge", "stack", "status-indicator", "divider"}
	
	for _, name := range expectedComponents {
		component, err := r.Get(name)
		if err != nil {
			t.Errorf("DefaultRegistry() missing component %q: %v", name, err)
		}
		if component.Name != name {
			t.Errorf("Component name = %q, want %q", component.Name, name)
		}
		if component.Version == "" {
			t.Errorf("Component %q has no version", name)
		}
		if component.Category == "" {
			t.Errorf("Component %q has no category", name)
		}
	}
	
	// Check categories
	categories := r.Categories()
	if len(categories) < 2 {
		t.Errorf("DefaultRegistry() has %d categories, want at least 2", len(categories))
	}
}

func TestMatchesQuery(t *testing.T) {
	component := &ComponentMetadata{
		Name:        "test-button",
		DisplayName: "Test Button",
		Description: "A test button component",
		Category:    "interactive",
		Tags:        []string{"action", "click"},
	}
	
	tests := []struct {
		query string
		want  bool
	}{
		{"test", true},
		{"button", true},
		{"Button", true},
		{"BUTTON", true},
		{"interactive", true},
		{"action", true},
		{"click", true},
		{"nonexistent", false},
		{"", true}, // Empty query matches everything (contains empty string)
	}
	
	for _, tt := range tests {
		got := matchesQuery(component, tt.query)
		if got != tt.want {
			t.Errorf("matchesQuery(%q) = %v, want %v", tt.query, got, tt.want)
		}
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test strings.ToLower (standard library)
	tests := []struct {
		input string
		want  string
	}{
		{"HELLO", "hello"},
		{"Hello", "hello"},
		{"hello", "hello"},
		{"Hello123", "hello123"},
		{"", ""},
	}
	
	for _, tt := range tests {
		got := strings.ToLower(tt.input)
		if got != tt.want {
			t.Errorf("strings.ToLower(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
	
	// Test strings.Contains (standard library)
	containsTests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "hello", true},
		{"hello world", "world", true},
		{"hello world", "o w", true},
		{"hello world", "xyz", false},
		{"", "", true},
		{"hello", "", true},
	}
	
	for _, tt := range containsTests {
		got := strings.Contains(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("strings.Contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}
