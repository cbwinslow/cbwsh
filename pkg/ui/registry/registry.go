// Package registry provides a component registry system for cbwsh UI components,
// similar to shadcn/ui. Components can be discovered, browsed, and added to projects.
package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ComponentMetadata describes a UI component in the registry.
type ComponentMetadata struct {
	// Name is the component name (e.g., "button", "card")
	Name string `json:"name"`
	// DisplayName is the human-readable name
	DisplayName string `json:"display_name"`
	// Description describes what the component does
	Description string `json:"description"`
	// Category groups components (e.g., "interactive", "layout", "display")
	Category string `json:"category"`
	// Dependencies lists other components this depends on
	Dependencies []string `json:"dependencies,omitempty"`
	// Files lists the files that make up this component
	Files []string `json:"files"`
	// Example provides example usage code
	Example string `json:"example,omitempty"`
	// Version is the component version
	Version string `json:"version"`
	// Tags for searching and filtering
	Tags []string `json:"tags,omitempty"`
}

// Registry manages the component registry.
type Registry struct {
	components map[string]*ComponentMetadata
	categories map[string][]*ComponentMetadata
}

// New creates a new registry.
func New() *Registry {
	return &Registry{
		components: make(map[string]*ComponentMetadata),
		categories: make(map[string][]*ComponentMetadata),
	}
}

// Register adds a component to the registry.
func (r *Registry) Register(component *ComponentMetadata) error {
	if component.Name == "" {
		return fmt.Errorf("component name cannot be empty")
	}
	
	r.components[component.Name] = component
	
	// Add to category index
	if component.Category != "" {
		r.categories[component.Category] = append(r.categories[component.Category], component)
	}
	
	return nil
}

// Get retrieves a component by name.
func (r *Registry) Get(name string) (*ComponentMetadata, error) {
	component, ok := r.components[name]
	if !ok {
		return nil, fmt.Errorf("component %q not found", name)
	}
	return component, nil
}

// List returns all registered components.
func (r *Registry) List() []*ComponentMetadata {
	components := make([]*ComponentMetadata, 0, len(r.components))
	for _, c := range r.components {
		components = append(components, c)
	}
	
	// Sort by name
	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})
	
	return components
}

// ListByCategory returns all components in a category.
func (r *Registry) ListByCategory(category string) []*ComponentMetadata {
	components := r.categories[category]
	
	// Sort by name
	sorted := make([]*ComponentMetadata, len(components))
	copy(sorted, components)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	
	return sorted
}

// Categories returns all available categories.
func (r *Registry) Categories() []string {
	categories := make([]string, 0, len(r.categories))
	for cat := range r.categories {
		categories = append(categories, cat)
	}
	sort.Strings(categories)
	return categories
}

// Search finds components matching a query.
func (r *Registry) Search(query string) []*ComponentMetadata {
	matches := make([]*ComponentMetadata, 0)
	
	for _, component := range r.components {
		if matchesQuery(component, query) {
			matches = append(matches, component)
		}
	}
	
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Name < matches[j].Name
	})
	
	return matches
}

func matchesQuery(component *ComponentMetadata, query string) bool {
	// Simple substring matching for now
	// Could be enhanced with fuzzy matching or more sophisticated algorithms
	query = toLower(query)
	
	if contains(toLower(component.Name), query) {
		return true
	}
	if contains(toLower(component.DisplayName), query) {
		return true
	}
	if contains(toLower(component.Description), query) {
		return true
	}
	if contains(toLower(component.Category), query) {
		return true
	}
	
	for _, tag := range component.Tags {
		if contains(toLower(tag), query) {
			return true
		}
	}
	
	return false
}

// LoadFromFile loads component metadata from a JSON file.
func (r *Registry) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read registry file: %w", err)
	}
	
	var components []*ComponentMetadata
	if err := json.Unmarshal(data, &components); err != nil {
		return fmt.Errorf("failed to parse registry file: %w", err)
	}
	
	for _, component := range components {
		if err := r.Register(component); err != nil {
			return fmt.Errorf("failed to register component %q: %w", component.Name, err)
		}
	}
	
	return nil
}

// SaveToFile saves the registry to a JSON file.
func (r *Registry) SaveToFile(path string) error {
	components := r.List()
	
	data, err := json.MarshalIndent(components, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}
	
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry file: %w", err)
	}
	
	return nil
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c + ('a' - 'A')
		}
		result[i] = c
	}
	return string(result)
}

// DefaultRegistry returns a registry with all built-in components registered.
func DefaultRegistry() *Registry {
	r := New()
	
	// Register all built-in components
	components := []*ComponentMetadata{
		{
			Name:        "button",
			DisplayName: "Button",
			Description: "A versatile button component with multiple variants and sizes",
			Category:    "interactive",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"interactive", "action", "form"},
			Example: `btn := components.NewButton("Save Changes")
btn.Variant = components.ButtonPrimary
btn.Size = components.ButtonSizeMedium
btn.Icon = "💾"
output := btn.Render()`,
		},
		{
			Name:        "card",
			DisplayName: "Card",
			Description: "A container component for grouping related content with title, content, and footer",
			Category:    "layout",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"container", "layout", "card"},
			Example: `card := components.NewCard("User Profile", "John Doe\njohn@example.com")
card.Footer = "Last updated: 2 minutes ago"
card.Elevated = true
output := card.Render()`,
		},
		{
			Name:        "badge",
			DisplayName: "Badge",
			Description: "A small label component for status indicators and counts",
			Category:    "display",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"label", "status", "indicator"},
			Example: `badge := components.NewBadge("NEW")
badge.Variant = components.BadgeInfo
output := badge.Render()`,
		},
		{
			Name:        "stack",
			DisplayName: "Stack",
			Description: "A layout component for organizing items vertically or horizontally",
			Category:    "layout",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"layout", "container", "flex"},
			Example: `stack := components.NewStack(components.StackVertical)
stack.Spacing = 8
stack.Add("Item 1")
stack.Add("Item 2")
output := stack.Render()`,
		},
		{
			Name:        "status-indicator",
			DisplayName: "Status Indicator",
			Description: "A component for displaying status with icon and label",
			Category:    "display",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"status", "indicator", "icon"},
			Example: `status := components.NewStatusIndicator("Server Running", components.StatusTypeSuccess)
status.ShowIcon = true
output := status.Render()`,
		},
		{
			Name:        "divider",
			DisplayName: "Divider",
			Description: "A visual separator between content sections",
			Category:    "layout",
			Files:       []string{"pkg/ui/components/components.go"},
			Version:     "1.0.0",
			Tags:        []string{"separator", "divider", "line"},
			Example: `divider := components.NewDivider(50)
divider.Char = "─"
output := divider.Render()`,
		},
	}
	
	for _, component := range components {
		_ = r.Register(component)
	}
	
	return r
}
