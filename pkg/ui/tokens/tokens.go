// Package tokens provides design tokens for cbwsh.
//
// Design tokens are the visual design atoms of the design system - specifically,
// they are named entities that store visual design attributes. Following Figma's
// design system principles, tokens ensure consistency across the entire UI.
//
// Token categories:
//   - Spacing: Base grid system (4px base)
//   - Typography: Font sizes, weights, line heights
//   - Borders: Border widths and radius values
//   - Shadows: Elevation and depth effects
//   - Animation: Duration and easing functions
//   - Z-index: Stacking order for layered UI
package tokens

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// SpacingScale defines a consistent spacing system based on a 4px grid.
// This provides predictable, harmonious spacing throughout the UI.
type SpacingScale struct {
	// Base is the fundamental unit (4px)
	Base int
	// XXS is 0.25x base (1px)
	XXS int
	// XS is 0.5x base (2px)
	XS int
	// SM is base (4px)
	SM int
	// MD is 2x base (8px)
	MD int
	// LG is 3x base (12px)
	LG int
	// XL is 4x base (16px)
	XL int
	// XXL is 6x base (24px)
	XXL int
	// XXXL is 8x base (32px)
	XXXL int
	// Huge is 12x base (48px)
	Huge int
}

// DefaultSpacing returns the default spacing scale.
func DefaultSpacing() SpacingScale {
	return SpacingScale{
		Base: 4,
		XXS:  1,
		XS:   2,
		SM:   4,
		MD:   8,
		LG:   12,
		XL:   16,
		XXL:  24,
		XXXL: 32,
		Huge: 48,
	}
}

// TypographyScale defines font sizes, weights, and line heights.
// Following Figma's type scale principles for clear visual hierarchy.
type TypographyScale struct {
	// Font sizes
	FontXS   string // 10px
	FontSM   string // 12px
	FontBase string // 14px
	FontMD   string // 16px
	FontLG   string // 18px
	FontXL   string // 24px
	FontXXL  string // 32px
	FontHuge string // 48px

	// Line heights
	LineHeightTight  float64 // 1.2
	LineHeightNormal float64 // 1.5
	LineHeightLoose  float64 // 1.8

	// Letter spacing
	LetterSpacingTight  string // -0.02em
	LetterSpacingNormal string // 0
	LetterSpacingWide   string // 0.05em
}

// DefaultTypography returns the default typography scale.
func DefaultTypography() TypographyScale {
	return TypographyScale{
		FontXS:              "10px",
		FontSM:              "12px",
		FontBase:            "14px",
		FontMD:              "16px",
		FontLG:              "18px",
		FontXL:              "24px",
		FontXXL:             "32px",
		FontHuge:            "48px",
		LineHeightTight:     1.2,
		LineHeightNormal:    1.5,
		LineHeightLoose:     1.8,
		LetterSpacingTight:  "-0.02em",
		LetterSpacingNormal: "0",
		LetterSpacingWide:   "0.05em",
	}
}

// BorderScale defines border widths and radius values for consistent shapes.
type BorderScale struct {
	// Widths
	WidthNone   int // 0
	WidthThin   int // 1px
	WidthNormal int // 2px
	WidthThick  int // 4px

	// Radius
	RadiusNone   int // 0
	RadiusXS     int // 2px
	RadiusSM     int // 4px
	RadiusMD     int // 8px
	RadiusLG     int // 12px
	RadiusXL     int // 16px
	RadiusFull   int // 9999px (fully rounded)
	RadiusCircle int // 50% (circle)
}

// DefaultBorders returns the default border scale.
func DefaultBorders() BorderScale {
	return BorderScale{
		WidthNone:   0,
		WidthThin:   1,
		WidthNormal: 2,
		WidthThick:  4,
		RadiusNone:  0,
		RadiusXS:    2,
		RadiusSM:    4,
		RadiusMD:    8,
		RadiusLG:    12,
		RadiusXL:    16,
		RadiusFull:  9999,
	}
}

// ShadowScale defines shadow values for elevation.
// Shadows create depth and visual hierarchy in the UI.
type ShadowScale struct {
	None   string // No shadow
	SM     string // Small elevation
	MD     string // Medium elevation
	LG     string // Large elevation
	XL     string // Extra large elevation
	Inner  string // Inner shadow (inset)
	Glow   string // Soft glow effect
	Accent string // Colored accent shadow
}

// DefaultShadows returns the default shadow scale.
func DefaultShadows() ShadowScale {
	return ShadowScale{
		None:   "none",
		SM:     "0 1px 2px 0 rgba(0, 0, 0, 0.05)",
		MD:     "0 4px 6px -1px rgba(0, 0, 0, 0.1)",
		LG:     "0 10px 15px -3px rgba(0, 0, 0, 0.1)",
		XL:     "0 20px 25px -5px rgba(0, 0, 0, 0.1)",
		Inner:  "inset 0 2px 4px 0 rgba(0, 0, 0, 0.06)",
		Glow:   "0 0 15px rgba(139, 233, 253, 0.5)",
		Accent: "0 0 0 3px rgba(139, 233, 253, 0.3)",
	}
}

// AnimationScale defines animation durations and easing functions.
// Consistent timing makes interactions feel smooth and predictable.
type AnimationScale struct {
	// Durations
	DurationInstant time.Duration // 50ms - Instant feedback
	DurationFast    time.Duration // 150ms - Quick transitions
	DurationNormal  time.Duration // 250ms - Standard transitions
	DurationSlow    time.Duration // 500ms - Deliberate animations
	DurationSlower  time.Duration // 750ms - Very slow animations

	// Easing (conceptual - for documentation)
	EaseLinear    string // Linear progression
	EaseIn        string // Slow start, fast end
	EaseOut       string // Fast start, slow end
	EaseInOut     string // Slow start and end
	EaseBack      string // Slight overshoot
	EaseElastic   string // Elastic bounce
	EaseBounce    string // Bounce effect
	EaseCircular  string // Circular curve
	EaseQuadratic string // Quadratic curve
}

// DefaultAnimation returns the default animation scale.
func DefaultAnimation() AnimationScale {
	return AnimationScale{
		DurationInstant: 50 * time.Millisecond,
		DurationFast:    150 * time.Millisecond,
		DurationNormal:  250 * time.Millisecond,
		DurationSlow:    500 * time.Millisecond,
		DurationSlower:  750 * time.Millisecond,
		EaseLinear:      "linear",
		EaseIn:          "ease-in",
		EaseOut:         "ease-out",
		EaseInOut:       "ease-in-out",
		EaseBack:        "cubic-bezier(0.175, 0.885, 0.32, 1.275)",
		EaseElastic:     "cubic-bezier(0.68, -0.55, 0.265, 1.55)",
		EaseBounce:      "cubic-bezier(0.175, 0.885, 0.32, 1.175)",
		EaseCircular:    "cubic-bezier(0.785, 0.135, 0.15, 0.86)",
		EaseQuadratic:   "cubic-bezier(0.455, 0.03, 0.515, 0.955)",
	}
}

// ZIndexScale defines the stacking order for layered UI elements.
// Consistent z-index values prevent layering conflicts.
type ZIndexScale struct {
	Base     int // 0 - Base layer
	Dropdown int // 1000 - Dropdowns and tooltips
	Sticky   int // 1100 - Sticky elements
	Fixed    int // 1200 - Fixed position elements
	Modal    int // 1300 - Modal dialogs
	Popover  int // 1400 - Popovers
	Tooltip  int // 1500 - Tooltips
	Toast    int // 1600 - Toast notifications
	Overlay  int // 1700 - Full-screen overlays
	Debug    int // 9999 - Debug overlays
}

// DefaultZIndex returns the default z-index scale.
func DefaultZIndex() ZIndexScale {
	return ZIndexScale{
		Base:     0,
		Dropdown: 1000,
		Sticky:   1100,
		Fixed:    1200,
		Modal:    1300,
		Popover:  1400,
		Tooltip:  1500,
		Toast:    1600,
		Overlay:  1700,
		Debug:    9999,
	}
}

// DesignTokens combines all token scales into a single system.
// This is the main entry point for accessing design tokens.
type DesignTokens struct {
	Spacing    SpacingScale
	Typography TypographyScale
	Borders    BorderScale
	Shadows    ShadowScale
	Animation  AnimationScale
	ZIndex     ZIndexScale
}

// Default returns the default design token system.
func Default() DesignTokens {
	return DesignTokens{
		Spacing:    DefaultSpacing(),
		Typography: DefaultTypography(),
		Borders:    DefaultBorders(),
		Shadows:    DefaultShadows(),
		Animation:  DefaultAnimation(),
		ZIndex:     DefaultZIndex(),
	}
}

// ApplySpacing creates a lipgloss style with the given spacing.
func ApplySpacing(spacing SpacingScale) map[string]lipgloss.Style {
	return map[string]lipgloss.Style{
		"padding-xs":  lipgloss.NewStyle().Padding(0, spacing.XS),
		"padding-sm":  lipgloss.NewStyle().Padding(0, spacing.SM),
		"padding-md":  lipgloss.NewStyle().Padding(0, spacing.MD),
		"padding-lg":  lipgloss.NewStyle().Padding(0, spacing.LG),
		"padding-xl":  lipgloss.NewStyle().Padding(0, spacing.XL),
		"padding-xxl": lipgloss.NewStyle().Padding(0, spacing.XXL),
		"margin-xs":   lipgloss.NewStyle().Margin(0, spacing.XS),
		"margin-sm":   lipgloss.NewStyle().Margin(0, spacing.SM),
		"margin-md":   lipgloss.NewStyle().Margin(0, spacing.MD),
		"margin-lg":   lipgloss.NewStyle().Margin(0, spacing.LG),
		"margin-xl":   lipgloss.NewStyle().Margin(0, spacing.XL),
		"margin-xxl":  lipgloss.NewStyle().Margin(0, spacing.XXL),
	}
}

// SemanticColors defines semantic color tokens mapped to theme colors.
// These provide meaning and context beyond raw color values.
type SemanticColors struct {
	// Interactive states
	InteractiveDefault  lipgloss.Color // Default interactive elements
	InteractiveHover    lipgloss.Color // Hover state
	InteractiveFocus    lipgloss.Color // Focus state
	InteractiveActive   lipgloss.Color // Active/pressed state
	InteractiveDisabled lipgloss.Color // Disabled state

	// Text hierarchy
	TextPrimary   lipgloss.Color // Primary text
	TextSecondary lipgloss.Color // Secondary text
	TextTertiary  lipgloss.Color // Tertiary/muted text
	TextInverse   lipgloss.Color // Inverse text (on dark bg)
	TextLink      lipgloss.Color // Hyperlinks

	// Backgrounds
	BgPrimary   lipgloss.Color // Primary background
	BgSecondary lipgloss.Color // Secondary background
	BgTertiary  lipgloss.Color // Tertiary background
	BgInverse   lipgloss.Color // Inverse background
	BgOverlay   lipgloss.Color // Overlay background

	// Borders
	BorderDefault lipgloss.Color // Default borders
	BorderFocus   lipgloss.Color // Focused borders
	BorderError   lipgloss.Color // Error borders

	// Status colors
	StatusInfo    lipgloss.Color // Informational
	StatusSuccess lipgloss.Color // Success
	StatusWarning lipgloss.Color // Warning
	StatusError   lipgloss.Color // Error
	StatusNeutral lipgloss.Color // Neutral
}

// DefaultSemanticColors returns default semantic colors.
func DefaultSemanticColors() SemanticColors {
	return SemanticColors{
		InteractiveDefault:  lipgloss.Color("62"),
		InteractiveHover:    lipgloss.Color("99"),
		InteractiveFocus:    lipgloss.Color("141"),
		InteractiveActive:   lipgloss.Color("57"),
		InteractiveDisabled: lipgloss.Color("243"),
		TextPrimary:         lipgloss.Color("255"),
		TextSecondary:       lipgloss.Color("252"),
		TextTertiary:        lipgloss.Color("243"),
		TextInverse:         lipgloss.Color("235"),
		TextLink:            lipgloss.Color("75"),
		BgPrimary:           lipgloss.Color("235"),
		BgSecondary:         lipgloss.Color("237"),
		BgTertiary:          lipgloss.Color("239"),
		BgInverse:           lipgloss.Color("255"),
		BgOverlay:           lipgloss.Color("0"),
		BorderDefault:       lipgloss.Color("238"),
		BorderFocus:         lipgloss.Color("62"),
		BorderError:         lipgloss.Color("196"),
		StatusInfo:          lipgloss.Color("75"),
		StatusSuccess:       lipgloss.Color("82"),
		StatusWarning:       lipgloss.Color("214"),
		StatusError:         lipgloss.Color("196"),
		StatusNeutral:       lipgloss.Color("243"),
	}
}

// ComponentTokens defines tokens specific to UI components.
// These are higher-level tokens that combine base tokens with semantic meaning.
type ComponentTokens struct {
	// Button
	ButtonPaddingX     int
	ButtonPaddingY     int
	ButtonBorderRadius int
	ButtonFontSize     string

	// Input
	InputPaddingX     int
	InputPaddingY     int
	InputBorderRadius int
	InputFontSize     string
	InputHeight       int

	// Card
	CardPadding      int
	CardBorderRadius int
	CardShadow       string

	// Modal
	ModalPadding      int
	ModalBorderRadius int
	ModalMaxWidth     int

	// Toast
	ToastPadding      int
	ToastBorderRadius int
	ToastMinWidth     int
	ToastMaxWidth     int

	// Tooltip
	TooltipPaddingX     int
	TooltipPaddingY     int
	TooltipBorderRadius int
	TooltipFontSize     string
}

// DefaultComponentTokens returns default component tokens.
func DefaultComponentTokens(spacing SpacingScale, typography TypographyScale, borders BorderScale, shadows ShadowScale) ComponentTokens {
	return ComponentTokens{
		// Button
		ButtonPaddingX:     spacing.MD,
		ButtonPaddingY:     spacing.SM,
		ButtonBorderRadius: borders.RadiusSM,
		ButtonFontSize:     typography.FontBase,

		// Input
		InputPaddingX:     spacing.MD,
		InputPaddingY:     spacing.SM,
		InputBorderRadius: borders.RadiusSM,
		InputFontSize:     typography.FontBase,
		InputHeight:       spacing.XL + (spacing.SM * 2),

		// Card
		CardPadding:      spacing.XL,
		CardBorderRadius: borders.RadiusMD,
		CardShadow:       shadows.MD,

		// Modal
		ModalPadding:      spacing.XXL,
		ModalBorderRadius: borders.RadiusLG,
		ModalMaxWidth:     600,

		// Toast
		ToastPadding:      spacing.MD,
		ToastBorderRadius: borders.RadiusMD,
		ToastMinWidth:     300,
		ToastMaxWidth:     500,

		// Tooltip
		TooltipPaddingX:     spacing.MD,
		TooltipPaddingY:     spacing.SM,
		TooltipBorderRadius: borders.RadiusSM,
		TooltipFontSize:     typography.FontSM,
	}
}
