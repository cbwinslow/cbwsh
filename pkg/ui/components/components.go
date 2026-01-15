// Package components provides reusable UI components following Figma design principles.
//
// Components in this package follow these design principles:
//   - Composition: Components are built from smaller, reusable pieces
//   - Consistency: All components use the same design tokens
//   - Flexibility: Components support variants and customization
//   - Accessibility: Components are keyboard-navigable and screen-reader friendly
//
// Following Figma's component model, each component:
//   - Has well-defined variants (e.g., Button sizes: sm, md, lg)
//   - Uses design tokens for spacing, colors, typography
//   - Supports composition (components within components)
//   - Has clear, descriptive prop types
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cbwinslow/cbwsh/pkg/ui/tokens"
)

// ButtonVariant defines different button styles.
type ButtonVariant string

const (
	ButtonPrimary   ButtonVariant = "primary"
	ButtonSecondary ButtonVariant = "secondary"
	ButtonDanger    ButtonVariant = "danger"
	ButtonGhost     ButtonVariant = "ghost"
	ButtonLink      ButtonVariant = "link"
)

// ButtonSize defines button sizes.
type ButtonSize string

const (
	ButtonSizeSmall  ButtonSize = "sm"
	ButtonSizeMedium ButtonSize = "md"
	ButtonSizeLarge  ButtonSize = "lg"
)

// Button represents a button component with Figma-inspired variants.
type Button struct {
	// Content is the button text
	Content string
	// Variant is the button style
	Variant ButtonVariant
	// Size is the button size
	Size ButtonSize
	// Disabled indicates if the button is disabled
	Disabled bool
	// FullWidth makes the button take full width
	FullWidth bool
	// Icon is an optional icon prefix
	Icon string
	// Tokens are the design tokens to use
	Tokens tokens.DesignTokens
	// SemanticColors are the semantic color tokens
	SemanticColors tokens.SemanticColors
}

// NewButton creates a new button with default values.
func NewButton(content string) *Button {
	return &Button{
		Content:        content,
		Variant:        ButtonPrimary,
		Size:           ButtonSizeMedium,
		Disabled:       false,
		FullWidth:      false,
		Tokens:         tokens.Default(),
		SemanticColors: tokens.DefaultSemanticColors(),
	}
}

// Render renders the button to a string.
func (b *Button) Render() string {
	style := b.buildStyle()
	content := b.Content
	if b.Icon != "" {
		content = b.Icon + " " + content
	}
	return style.Render(content)
}

func (b *Button) buildStyle() lipgloss.Style {
	spacing := b.Tokens.Spacing

	// Base style
	style := lipgloss.NewStyle()

	// Size
	switch b.Size {
	case ButtonSizeSmall:
		style = style.Padding(0, spacing.SM)
	case ButtonSizeMedium:
		style = style.Padding(0, spacing.MD)
	case ButtonSizeLarge:
		style = style.Padding(0, spacing.LG)
	}

	// Border radius
	// Note: Lipgloss doesn't support border-radius directly, so we use borders
	style = style.Border(lipgloss.RoundedBorder())

	// Variant
	if b.Disabled {
		style = style.
			Foreground(b.SemanticColors.InteractiveDisabled).
			BorderForeground(b.SemanticColors.InteractiveDisabled)
	} else {
		switch b.Variant {
		case ButtonPrimary:
			style = style.
				Foreground(b.SemanticColors.TextInverse).
				Background(b.SemanticColors.InteractiveDefault).
				BorderForeground(b.SemanticColors.InteractiveDefault).
				Bold(true)
		case ButtonSecondary:
			style = style.
				Foreground(b.SemanticColors.TextPrimary).
				BorderForeground(b.SemanticColors.BorderDefault)
		case ButtonDanger:
			style = style.
				Foreground(b.SemanticColors.TextInverse).
				Background(b.SemanticColors.StatusError).
				BorderForeground(b.SemanticColors.StatusError).
				Bold(true)
		case ButtonGhost:
			style = style.
				Foreground(b.SemanticColors.TextPrimary).
				BorderForeground(lipgloss.Color("transparent"))
		case ButtonLink:
			style = style.
				Foreground(b.SemanticColors.TextLink).
				BorderForeground(lipgloss.Color("transparent")).
				Underline(true)
		}
	}

	// Full width
	if b.FullWidth {
		style = style.Width(100)
	}

	return style
}

// Card represents a card component with consistent styling.
type Card struct {
	// Title is the card title
	Title string
	// Content is the card content
	Content string
	// Footer is optional footer content
	Footer string
	// Elevated adds shadow for depth
	Elevated bool
	// Bordered adds a border
	Bordered bool
	// Padding controls internal spacing
	Padding int
	// Tokens are the design tokens
	Tokens tokens.DesignTokens
	// SemanticColors are the semantic color tokens
	SemanticColors tokens.SemanticColors
}

// NewCard creates a new card.
func NewCard(title, content string) *Card {
	return &Card{
		Title:          title,
		Content:        content,
		Elevated:       true,
		Bordered:       true,
		Padding:        tokens.Default().Spacing.XL,
		Tokens:         tokens.Default(),
		SemanticColors: tokens.DefaultSemanticColors(),
	}
}

// Render renders the card to a string.
func (c *Card) Render() string {
	var sections []string

	titleStyle := lipgloss.NewStyle().
		Foreground(c.SemanticColors.TextPrimary).
		Bold(true).
		Padding(0, 0, 1, 0)

	contentStyle := lipgloss.NewStyle().
		Foreground(c.SemanticColors.TextSecondary)

	footerStyle := lipgloss.NewStyle().
		Foreground(c.SemanticColors.TextTertiary).
		Padding(1, 0, 0, 0)

	if c.Title != "" {
		sections = append(sections, titleStyle.Render(c.Title))
	}
	if c.Content != "" {
		sections = append(sections, contentStyle.Render(c.Content))
	}
	if c.Footer != "" {
		sections = append(sections, footerStyle.Render(c.Footer))
	}

	inner := lipgloss.JoinVertical(lipgloss.Left, sections...)

	cardStyle := lipgloss.NewStyle().
		Padding(0, c.Padding/2).
		Background(c.SemanticColors.BgSecondary)

	if c.Bordered {
		cardStyle = cardStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(c.SemanticColors.BorderDefault)
	}

	return cardStyle.Render(inner)
}

// Badge represents a small labeled component.
type Badge struct {
	// Label is the badge text
	Label string
	// Variant is the badge style
	Variant BadgeVariant
	// Tokens are the design tokens
	Tokens tokens.DesignTokens
	// SemanticColors are the semantic color tokens
	SemanticColors tokens.SemanticColors
}

// BadgeVariant defines badge styles.
type BadgeVariant string

const (
	BadgeDefault BadgeVariant = "default"
	BadgeSuccess BadgeVariant = "success"
	BadgeWarning BadgeVariant = "warning"
	BadgeError   BadgeVariant = "error"
	BadgeInfo    BadgeVariant = "info"
)

// NewBadge creates a new badge.
func NewBadge(label string) *Badge {
	return &Badge{
		Label:          label,
		Variant:        BadgeDefault,
		Tokens:         tokens.Default(),
		SemanticColors: tokens.DefaultSemanticColors(),
	}
}

// Render renders the badge to a string.
func (b *Badge) Render() string {
	style := lipgloss.NewStyle().
		Padding(0, b.Tokens.Spacing.SM).
		Bold(true)

	switch b.Variant {
	case BadgeSuccess:
		style = style.
			Foreground(b.SemanticColors.TextInverse).
			Background(b.SemanticColors.StatusSuccess)
	case BadgeWarning:
		style = style.
			Foreground(b.SemanticColors.TextInverse).
			Background(b.SemanticColors.StatusWarning)
	case BadgeError:
		style = style.
			Foreground(b.SemanticColors.TextInverse).
			Background(b.SemanticColors.StatusError)
	case BadgeInfo:
		style = style.
			Foreground(b.SemanticColors.TextInverse).
			Background(b.SemanticColors.StatusInfo)
	default:
		style = style.
			Foreground(b.SemanticColors.TextPrimary).
			Background(b.SemanticColors.BgTertiary)
	}

	return style.Render(b.Label)
}

// Stack represents a vertical or horizontal layout container.
type Stack struct {
	// Direction is the stack direction
	Direction StackDirection
	// Spacing is the gap between items
	Spacing int
	// Items are the stack items
	Items []string
	// Alignment controls item alignment
	Alignment StackAlignment
	// Tokens are the design tokens
	Tokens tokens.DesignTokens
}

// StackDirection defines stack direction.
type StackDirection string

const (
	StackVertical   StackDirection = "vertical"
	StackHorizontal StackDirection = "horizontal"
)

// StackAlignment defines item alignment.
type StackAlignment string

const (
	StackAlignStart  StackAlignment = "start"
	StackAlignCenter StackAlignment = "center"
	StackAlignEnd    StackAlignment = "end"
)

// NewStack creates a new stack.
func NewStack(direction StackDirection) *Stack {
	return &Stack{
		Direction: direction,
		Spacing:   tokens.Default().Spacing.MD,
		Items:     make([]string, 0),
		Alignment: StackAlignStart,
		Tokens:    tokens.Default(),
	}
}

// Add adds an item to the stack.
func (s *Stack) Add(item string) *Stack {
	s.Items = append(s.Items, item)
	return s
}

// Render renders the stack to a string.
func (s *Stack) Render() string {
	if len(s.Items) == 0 {
		return ""
	}

	var spacer string
	if s.Direction == StackVertical {
		spacer = strings.Repeat("\n", s.Spacing/4)
	} else {
		spacer = strings.Repeat(" ", s.Spacing)
	}

	if s.Direction == StackVertical {
		alignment := lipgloss.Left
		switch s.Alignment {
		case StackAlignCenter:
			alignment = lipgloss.Center
		case StackAlignEnd:
			alignment = lipgloss.Right
		}
		return lipgloss.JoinVertical(alignment, s.Items...)
	}

	return strings.Join(s.Items, spacer)
}

// Divider represents a visual separator.
type Divider struct {
	// Length is the divider length
	Length int
	// Char is the divider character
	Char string
	// Tokens are the design tokens
	Tokens tokens.DesignTokens
	// SemanticColors are the semantic color tokens
	SemanticColors tokens.SemanticColors
}

// NewDivider creates a new divider.
func NewDivider(length int) *Divider {
	return &Divider{
		Length:         length,
		Char:           "─",
		Tokens:         tokens.Default(),
		SemanticColors: tokens.DefaultSemanticColors(),
	}
}

// Render renders the divider to a string.
func (d *Divider) Render() string {
	style := lipgloss.NewStyle().
		Foreground(d.SemanticColors.BorderDefault)

	divider := strings.Repeat(d.Char, d.Length)
	return style.Render(divider)
}

// StatusIndicator represents a status indicator with icon and label.
type StatusIndicator struct {
	// Label is the status text
	Label string
	// Status is the status type
	Status StatusType
	// ShowIcon indicates whether to show an icon
	ShowIcon bool
	// Tokens are the design tokens
	Tokens tokens.DesignTokens
	// SemanticColors are the semantic color tokens
	SemanticColors tokens.SemanticColors
}

// StatusType defines status types.
type StatusType string

const (
	StatusTypeSuccess StatusType = "success"
	StatusTypeWarning StatusType = "warning"
	StatusTypeError   StatusType = "error"
	StatusTypeInfo    StatusType = "info"
	StatusTypeNeutral StatusType = "neutral"
)

// NewStatusIndicator creates a new status indicator.
func NewStatusIndicator(label string, status StatusType) *StatusIndicator {
	return &StatusIndicator{
		Label:          label,
		Status:         status,
		ShowIcon:       true,
		Tokens:         tokens.Default(),
		SemanticColors: tokens.DefaultSemanticColors(),
	}
}

// Render renders the status indicator to a string.
func (s *StatusIndicator) Render() string {
	var color lipgloss.Color
	var icon string

	switch s.Status {
	case StatusTypeSuccess:
		color = s.SemanticColors.StatusSuccess
		icon = "✓"
	case StatusTypeWarning:
		color = s.SemanticColors.StatusWarning
		icon = "⚠"
	case StatusTypeError:
		color = s.SemanticColors.StatusError
		icon = "✗"
	case StatusTypeInfo:
		color = s.SemanticColors.StatusInfo
		icon = "ℹ"
	default:
		color = s.SemanticColors.StatusNeutral
		icon = "●"
	}

	style := lipgloss.NewStyle().Foreground(color)

	content := s.Label
	if s.ShowIcon {
		content = icon + " " + content
	}

	return style.Render(content)
}
