package themes

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_New(t *testing.T) {
	m := NewManager("")

	if m == nil {
		t.Fatal("NewManager() returned nil")
	}

	// Should have built-in themes
	themes := m.List()
	if len(themes) == 0 {
		t.Error("NewManager() should have built-in themes")
	}
}

func TestManager_Load(t *testing.T) {
	m := NewManager("")

	// Load built-in theme
	theme, err := m.Load("default")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if theme.Name != "default" {
		t.Errorf("Load() theme.Name = %s, want default", theme.Name)
	}

	// Load non-existent theme
	_, err = m.Load("nonexistent")
	if err != ErrThemeNotFound {
		t.Errorf("Load() error = %v, want ErrThemeNotFound", err)
	}
}

func TestManager_Current(t *testing.T) {
	m := NewManager("")

	// Default current theme
	theme := m.Current()
	if theme == nil {
		t.Fatal("Current() returned nil")
	}

	if theme.Name != "default" {
		t.Errorf("Current() theme.Name = %s, want default", theme.Name)
	}
}

func TestManager_SetCurrent(t *testing.T) {
	m := NewManager("")

	err := m.SetCurrent("dracula")
	if err != nil {
		t.Fatalf("SetCurrent() error = %v", err)
	}

	theme := m.Current()
	if theme.Name != "dracula" {
		t.Errorf("Current() after SetCurrent = %s, want dracula", theme.Name)
	}

	// Non-existent theme
	err = m.SetCurrent("nonexistent")
	if err != ErrThemeNotFound {
		t.Errorf("SetCurrent() error = %v, want ErrThemeNotFound", err)
	}
}

func TestManager_List(t *testing.T) {
	m := NewManager("")

	themes := m.List()

	expectedThemes := []string{"default", "dracula", "nord", "monokai", "one-dark", "gruvbox", "catppuccin"}
	for _, expected := range expectedThemes {
		found := false
		for _, name := range themes {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("List() missing theme: %s", expected)
		}
	}
}

func TestManager_ListThemes(t *testing.T) {
	m := NewManager("")

	themes := m.ListThemes()
	if len(themes) == 0 {
		t.Error("ListThemes() returned empty list")
	}
}

func TestManager_Register(t *testing.T) {
	m := NewManager("")

	customTheme := &Theme{
		Name:        "custom",
		Description: "Custom theme",
		IsDark:      true,
		Colors: ColorScheme{
			Background: "#000000",
			Foreground: "#ffffff",
		},
	}

	m.Register(customTheme)

	theme, err := m.Load("custom")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if theme.Description != "Custom theme" {
		t.Errorf("Load() theme.Description = %s, want Custom theme", theme.Description)
	}
}

func TestManager_Unregister(t *testing.T) {
	m := NewManager("")

	m.Unregister("monokai")

	_, err := m.Load("monokai")
	if err != ErrThemeNotFound {
		t.Errorf("Load() after Unregister error = %v, want ErrThemeNotFound", err)
	}
}

func TestManager_LoadFromFile(t *testing.T) {
	// Create temp directory with theme file
	tmpDir, err := os.MkdirTemp("", "themes-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	themeContent := `
name: file-theme
description: Theme from file
is_dark: true
colors:
  background: "#1a1a1a"
  foreground: "#f0f0f0"
  primary: "#ff0000"
`

	themePath := filepath.Join(tmpDir, "file-theme.yaml")
	if err := os.WriteFile(themePath, []byte(themeContent), 0o644); err != nil {
		t.Fatalf("Failed to write theme file: %v", err)
	}

	m := NewManager(tmpDir)

	theme, err := m.Load("file-theme")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if theme.Description != "Theme from file" {
		t.Errorf("Load() theme.Description = %s, want 'Theme from file'", theme.Description)
	}

	if theme.Colors.Primary != "#ff0000" {
		t.Errorf("Load() theme.Colors.Primary = %s, want #ff0000", theme.Colors.Primary)
	}
}

func TestManager_Reload(t *testing.T) {
	m := NewManager("")

	err := m.Reload()
	if err != nil {
		t.Fatalf("Reload() error = %v", err)
	}

	// Themes should still be available
	themes := m.List()
	if len(themes) == 0 {
		t.Error("Reload() should preserve themes")
	}
}

func TestManager_OnChange(t *testing.T) {
	m := NewManager("")

	called := false
	var receivedTheme *Theme

	m.OnChange(func(theme *Theme) {
		called = true
		receivedTheme = theme
	})

	_ = m.SetCurrent("nord")

	if !called {
		t.Error("OnChange callback should be called")
	}

	if receivedTheme.Name != "nord" {
		t.Errorf("OnChange received theme.Name = %s, want nord", receivedTheme.Name)
	}
}

func TestBuiltinThemes(t *testing.T) {
	themes := []struct {
		name   string
		fn     func() *Theme
		isDark bool
	}{
		{"default", DefaultTheme, true},
		{"dracula", DraculaTheme, true},
		{"nord", NordTheme, true},
		{"solarized-dark", SolarizedDarkTheme, true},
		{"solarized-light", SolarizedLightTheme, false},
		{"monokai", MonokaiTheme, true},
		{"one-dark", OneDarkTheme, true},
		{"gruvbox", GruvboxTheme, true},
		{"catppuccin", CatppuccinTheme, true},
	}

	for _, tt := range themes {
		t.Run(tt.name, func(t *testing.T) {
			theme := tt.fn()

			if theme.Name != tt.name {
				t.Errorf("Name = %s, want %s", theme.Name, tt.name)
			}

			if theme.IsDark != tt.isDark {
				t.Errorf("IsDark = %v, want %v", theme.IsDark, tt.isDark)
			}

			// Check essential colors are set
			if theme.Colors.Background == "" {
				t.Error("Background color should be set")
			}
			if theme.Colors.Foreground == "" {
				t.Error("Foreground color should be set")
			}
			if theme.Colors.Primary == "" {
				t.Error("Primary color should be set")
			}
		})
	}
}

func TestApplyTheme(t *testing.T) {
	theme := DefaultTheme()
	styles := ApplyTheme(theme)

	expectedStyles := []string{"background", "foreground", "primary", "secondary", "success", "warning", "error", "info"}
	for _, name := range expectedStyles {
		if _, ok := styles[name]; !ok {
			t.Errorf("ApplyTheme() missing style: %s", name)
		}
	}
}

func TestToLipglossColor(t *testing.T) {
	color := ToLipglossColor("#ff0000")
	if string(color) != "#ff0000" {
		t.Errorf("ToLipglossColor() = %s, want #ff0000", string(color))
	}
}

func TestManager_Watch(t *testing.T) {
	m := NewManager("")

	err := m.Watch()
	if err != nil {
		t.Fatalf("Watch() error = %v", err)
	}

	// Should be watching
	if !m.watching {
		t.Error("Watch() should set watching to true")
	}

	// Stop watching
	m.StopWatch()

	// Should not be watching
	if m.watching {
		t.Error("StopWatch() should set watching to false")
	}
}

func TestColorScheme(t *testing.T) {
	theme := DefaultTheme()
	colors := theme.Colors

	// Verify all ANSI colors are set
	ansiColors := []string{
		colors.Black, colors.Red, colors.Green, colors.Yellow,
		colors.Blue, colors.Magenta, colors.Cyan, colors.White,
	}

	for i, c := range ansiColors {
		if c == "" {
			t.Errorf("ANSI color %d is empty", i)
		}
	}
}
