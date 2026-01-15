package tokens

import (
	"testing"
	"time"
)

func TestDefaultSpacing(t *testing.T) {
	spacing := DefaultSpacing()

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"Base", spacing.Base, 4},
		{"XXS", spacing.XXS, 1},
		{"XS", spacing.XS, 2},
		{"SM", spacing.SM, 4},
		{"MD", spacing.MD, 8},
		{"LG", spacing.LG, 12},
		{"XL", spacing.XL, 16},
		{"XXL", spacing.XXL, 24},
		{"XXXL", spacing.XXXL, 32},
		{"Huge", spacing.Huge, 48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Spacing.%s = %d, want %d", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestDefaultTypography(t *testing.T) {
	typography := DefaultTypography()

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"FontXS", typography.FontXS, "10px"},
		{"FontSM", typography.FontSM, "12px"},
		{"FontBase", typography.FontBase, "14px"},
		{"FontMD", typography.FontMD, "16px"},
		{"FontLG", typography.FontLG, "18px"},
		{"FontXL", typography.FontXL, "24px"},
		{"FontXXL", typography.FontXXL, "32px"},
		{"FontHuge", typography.FontHuge, "48px"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Typography.%s = %s, want %s", tt.name, tt.value, tt.expected)
			}
		})
	}

	if typography.LineHeightNormal != 1.5 {
		t.Errorf("LineHeightNormal = %f, want 1.5", typography.LineHeightNormal)
	}
}

func TestDefaultBorders(t *testing.T) {
	borders := DefaultBorders()

	if borders.WidthThin != 1 {
		t.Errorf("WidthThin = %d, want 1", borders.WidthThin)
	}

	if borders.RadiusSM != 4 {
		t.Errorf("RadiusSM = %d, want 4", borders.RadiusSM)
	}

	if borders.RadiusFull != 9999 {
		t.Errorf("RadiusFull = %d, want 9999", borders.RadiusFull)
	}
}

func TestDefaultAnimation(t *testing.T) {
	animation := DefaultAnimation()

	tests := []struct {
		name     string
		value    time.Duration
		expected time.Duration
	}{
		{"DurationInstant", animation.DurationInstant, 50 * time.Millisecond},
		{"DurationFast", animation.DurationFast, 150 * time.Millisecond},
		{"DurationNormal", animation.DurationNormal, 250 * time.Millisecond},
		{"DurationSlow", animation.DurationSlow, 500 * time.Millisecond},
		{"DurationSlower", animation.DurationSlower, 750 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("Animation.%s = %v, want %v", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestDefaultZIndex(t *testing.T) {
	zIndex := DefaultZIndex()

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{"Base", zIndex.Base, 0},
		{"Dropdown", zIndex.Dropdown, 1000},
		{"Modal", zIndex.Modal, 1300},
		{"Toast", zIndex.Toast, 1600},
		{"Overlay", zIndex.Overlay, 1700},
		{"Debug", zIndex.Debug, 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("ZIndex.%s = %d, want %d", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tokens := Default()

	if tokens.Spacing.Base != 4 {
		t.Errorf("Expected Base spacing to be 4, got %d", tokens.Spacing.Base)
	}

	if tokens.Typography.FontBase != "14px" {
		t.Errorf("Expected FontBase to be 14px, got %s", tokens.Typography.FontBase)
	}

	if tokens.Borders.RadiusMD != 8 {
		t.Errorf("Expected RadiusMD to be 8, got %d", tokens.Borders.RadiusMD)
	}

	if tokens.Animation.DurationNormal != 250*time.Millisecond {
		t.Errorf("Expected DurationNormal to be 250ms, got %v", tokens.Animation.DurationNormal)
	}

	if tokens.ZIndex.Modal != 1300 {
		t.Errorf("Expected Modal z-index to be 1300, got %d", tokens.ZIndex.Modal)
	}
}

func TestApplySpacing(t *testing.T) {
	spacing := DefaultSpacing()
	styles := ApplySpacing(spacing)

	if len(styles) == 0 {
		t.Error("ApplySpacing returned no styles")
	}

	expectedKeys := []string{
		"padding-xs", "padding-sm", "padding-md", "padding-lg", "padding-xl", "padding-xxl",
		"margin-xs", "margin-sm", "margin-md", "margin-lg", "margin-xl", "margin-xxl",
	}

	for _, key := range expectedKeys {
		if _, ok := styles[key]; !ok {
			t.Errorf("Expected style key %s not found", key)
		}
	}
}

func TestDefaultSemanticColors(t *testing.T) {
	colors := DefaultSemanticColors()

	if colors.InteractiveDefault == "" {
		t.Error("InteractiveDefault color should not be empty")
	}

	if colors.TextPrimary == "" {
		t.Error("TextPrimary color should not be empty")
	}

	if colors.BgPrimary == "" {
		t.Error("BgPrimary color should not be empty")
	}

	if colors.StatusSuccess == "" {
		t.Error("StatusSuccess color should not be empty")
	}
}

func TestDefaultComponentTokens(t *testing.T) {
	spacing := DefaultSpacing()
	typography := DefaultTypography()
	borders := DefaultBorders()
	shadows := DefaultShadows()

	tokens := DefaultComponentTokens(spacing, typography, borders, shadows)

	if tokens.ButtonPaddingX != spacing.MD {
		t.Errorf("ButtonPaddingX = %d, want %d", tokens.ButtonPaddingX, spacing.MD)
	}

	if tokens.InputFontSize != typography.FontBase {
		t.Errorf("InputFontSize = %s, want %s", tokens.InputFontSize, typography.FontBase)
	}

	if tokens.CardBorderRadius != borders.RadiusMD {
		t.Errorf("CardBorderRadius = %d, want %d", tokens.CardBorderRadius, borders.RadiusMD)
	}

	if tokens.ModalMaxWidth != 600 {
		t.Errorf("ModalMaxWidth = %d, want 600", tokens.ModalMaxWidth)
	}
}

func TestSpacingConsistency(t *testing.T) {
	spacing := DefaultSpacing()

	// Check that spacing follows the 4px grid
	if spacing.SM%spacing.Base != 0 {
		t.Error("SM spacing should be a multiple of Base")
	}
	if spacing.MD%spacing.Base != 0 {
		t.Error("MD spacing should be a multiple of Base")
	}
	if spacing.LG%spacing.Base != 0 {
		t.Error("LG spacing should be a multiple of Base")
	}
}

func TestZIndexOrdering(t *testing.T) {
	zIndex := DefaultZIndex()

	// Verify z-index values are in ascending order
	if zIndex.Base >= zIndex.Dropdown {
		t.Error("Base should be less than Dropdown")
	}
	if zIndex.Dropdown >= zIndex.Modal {
		t.Error("Dropdown should be less than Modal")
	}
	if zIndex.Modal >= zIndex.Toast {
		t.Error("Modal should be less than Toast")
	}
	if zIndex.Toast >= zIndex.Overlay {
		t.Error("Toast should be less than Overlay")
	}
}
