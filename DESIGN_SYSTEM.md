# cbwsh Design System

## Overview

The cbwsh design system is inspired by Figma's design principles and MCP server best practices. It provides a consistent, composable, and scalable foundation for building terminal user interfaces.

## Design Principles

### 1. **Consistency**
All components use the same design tokens for spacing, colors, typography, and other visual properties. This ensures a cohesive look and feel across the entire application.

### 2. **Composition**
Components are built from smaller, reusable pieces. This allows for flexible layouts and easy customization without duplicating code.

### 3. **Accessibility**
Components are designed with keyboard navigation and screen reader support in mind, following best practices for terminal UI accessibility.

### 4. **Flexibility**
Each component supports variants and customization options, allowing developers to adapt them to specific use cases while maintaining visual consistency.

## Design Tokens

Design tokens are the visual design atoms of the system. They define fundamental values that components use consistently.

### Spacing Scale

Based on a **4px grid system** for predictable, harmonious spacing:

```go
spacing := tokens.DefaultSpacing()

spacing.XXS   // 1px  - Hairline spacing
spacing.XS    // 2px  - Extra small
spacing.SM    // 4px  - Small (base unit)
spacing.MD    // 8px  - Medium
spacing.LG    // 12px - Large
spacing.XL    // 16px - Extra large
spacing.XXL   // 24px - Double extra large
spacing.XXXL  // 32px - Triple extra large
spacing.Huge  // 48px - Huge spacing
```

**Usage Example:**
```go
style := lipgloss.NewStyle().Padding(0, spacing.MD)
```

### Typography Scale

Type scale for clear visual hierarchy:

```go
typography := tokens.DefaultTypography()

typography.FontXS    // 10px - Extra small text
typography.FontSM    // 12px - Small text
typography.FontBase  // 14px - Body text (default)
typography.FontMD    // 16px - Medium text
typography.FontLG    // 18px - Large text
typography.FontXL    // 24px - Extra large headings
typography.FontXXL   // 32px - Double extra large headings
typography.FontHuge  // 48px - Display text
```

**Line Heights:**
- `LineHeightTight`: 1.2 - Compact text
- `LineHeightNormal`: 1.5 - Standard readability
- `LineHeightLoose`: 1.8 - Extra breathing room

### Border Scale

Consistent border widths and radii:

```go
borders := tokens.DefaultBorders()

// Widths
borders.WidthThin   // 1px
borders.WidthNormal // 2px
borders.WidthThick  // 4px

// Radius
borders.RadiusXS  // 2px  - Subtle rounding
borders.RadiusSM  // 4px  - Small rounding
borders.RadiusMD  // 8px  - Medium rounding
borders.RadiusLG  // 12px - Large rounding
borders.RadiusXL  // 16px - Extra large rounding
borders.RadiusFull // 9999px - Fully rounded (pill shape)
```

### Animation Tokens

Consistent timing for smooth interactions:

```go
animation := tokens.DefaultAnimation()

animation.DurationInstant // 50ms  - Instant feedback
animation.DurationFast    // 150ms - Quick transitions
animation.DurationNormal  // 250ms - Standard transitions
animation.DurationSlow    // 500ms - Deliberate animations
animation.DurationSlower  // 750ms - Very slow animations
```

**Easing functions** (conceptual):
- `EaseLinear`: Linear progression
- `EaseInOut`: Smooth acceleration and deceleration
- `EaseBack`: Slight overshoot for emphasis
- `EaseElastic`: Bouncy, playful motion

### Z-Index Scale

Layering system to prevent z-index conflicts:

```go
zIndex := tokens.DefaultZIndex()

zIndex.Base     // 0    - Base layer
zIndex.Dropdown // 1000 - Dropdowns and tooltips
zIndex.Modal    // 1300 - Modal dialogs
zIndex.Toast    // 1600 - Toast notifications
zIndex.Overlay  // 1700 - Full-screen overlays
zIndex.Debug    // 9999 - Debug overlays (always on top)
```

### Semantic Colors

Colors with meaning, not just values:

```go
colors := tokens.DefaultSemanticColors()

// Interactive states
colors.InteractiveDefault  // Default state
colors.InteractiveHover    // Hover state
colors.InteractiveFocus    // Focus state
colors.InteractiveDisabled // Disabled state

// Text hierarchy
colors.TextPrimary   // Primary text
colors.TextSecondary // Secondary text
colors.TextTertiary  // Muted text

// Status colors
colors.StatusSuccess // Green - Success messages
colors.StatusWarning // Yellow - Warnings
colors.StatusError   // Red - Errors
colors.StatusInfo    // Blue - Information
```

## Components

### Button

A versatile button component with multiple variants and sizes.

**Variants:**
- `ButtonPrimary`: Primary action (bold, colored background)
- `ButtonSecondary`: Secondary action (outlined)
- `ButtonDanger`: Destructive action (red)
- `ButtonGhost`: Subtle action (no border)
- `ButtonLink`: Text link style

**Sizes:**
- `ButtonSizeSmall`: Compact
- `ButtonSizeMedium`: Standard
- `ButtonSizeLarge`: Prominent

**Example:**
```go
import "github.com/cbwinslow/cbwsh/pkg/ui/components"

// Create a primary button
btn := components.NewButton("Save Changes")
btn.Variant = components.ButtonPrimary
btn.Size = components.ButtonSizeMedium
btn.Icon = "💾"

// Render the button
output := btn.Render()
```

### Card

A container component for grouping related content.

**Features:**
- Title, content, and optional footer sections
- Elevation (shadow) for depth
- Border styling
- Configurable padding

**Example:**
```go
card := components.NewCard("User Profile", "John Doe\njohn@example.com")
card.Footer = "Last updated: 2 minutes ago"
card.Elevated = true
card.Bordered = true

output := card.Render()
```

### Badge

A small label component for status indicators and counts.

**Variants:**
- `BadgeDefault`: Neutral
- `BadgeSuccess`: Green - Success/active
- `BadgeWarning`: Yellow - Warning/pending
- `BadgeError`: Red - Error/critical
- `BadgeInfo`: Blue - Information

**Example:**
```go
badge := components.NewBadge("NEW")
badge.Variant = components.BadgeInfo

output := badge.Render()
```

### Stack

A layout component for organizing items vertically or horizontally.

**Directions:**
- `StackVertical`: Items stacked top to bottom
- `StackHorizontal`: Items placed left to right

**Alignment:**
- `StackAlignStart`: Align to start
- `StackAlignCenter`: Center align
- `StackAlignEnd`: Align to end

**Example:**
```go
stack := components.NewStack(components.StackVertical)
stack.Spacing = 8 // Use tokens.Spacing.MD
stack.Add("Item 1")
stack.Add("Item 2")
stack.Add("Item 3")

output := stack.Render()
```

### Status Indicator

A component for displaying status with icon and label.

**Status Types:**
- `StatusTypeSuccess`: ✓ Green
- `StatusTypeWarning`: ⚠ Yellow
- `StatusTypeError`: ✗ Red
- `StatusTypeInfo`: ℹ Blue
- `StatusTypeNeutral`: ● Gray

**Example:**
```go
status := components.NewStatusIndicator("Server Running", components.StatusTypeSuccess)
status.ShowIcon = true

output := status.Render()
```

### Divider

A visual separator between content sections.

**Example:**
```go
divider := components.NewDivider(50)
divider.Char = "─"

output := divider.Render()
```

## Component Composition

Components can be composed together to create complex UIs:

```go
// Create a card with buttons
card := components.NewCard("Confirm Action", "Are you sure you want to continue?")

// Create button stack
btnStack := components.NewStack(components.StackHorizontal)
btnStack.Spacing = 8

cancelBtn := components.NewButton("Cancel")
cancelBtn.Variant = components.ButtonSecondary

confirmBtn := components.NewButton("Confirm")
confirmBtn.Variant = components.ButtonPrimary

btnStack.Add(cancelBtn.Render())
btnStack.Add(confirmBtn.Render())

card.Footer = btnStack.Render()
output := card.Render()
```

## Best Practices

### 1. Use Design Tokens

Always use design tokens instead of hard-coded values:

```go
// ✅ Good
style := lipgloss.NewStyle().Padding(0, tokens.Default().Spacing.MD)

// ❌ Bad
style := lipgloss.NewStyle().Padding(0, 8)
```

### 2. Compose Components

Build complex UIs by composing smaller components:

```go
// Create a status card
badge := components.NewBadge("ACTIVE")
status := components.NewStatusIndicator("Running", components.StatusTypeSuccess)
card := components.NewCard("Service Status", badge.Render() + " " + status.Render())
```

### 3. Use Semantic Colors

Use semantic color tokens for meaning:

```go
// ✅ Good
style := lipgloss.NewStyle().Foreground(colors.StatusError)

// ❌ Bad
style := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
```

### 4. Maintain Spacing Consistency

Use the spacing scale for all padding and margins:

```go
// Stack items with consistent spacing
stack := components.NewStack(components.StackVertical)
stack.Spacing = tokens.Default().Spacing.MD // Use token
```

### 5. Component Variants Over Custom Styles

Use component variants instead of custom styling:

```go
// ✅ Good
btn := components.NewButton("Delete")
btn.Variant = components.ButtonDanger

// ❌ Bad
btn := components.NewButton("Delete")
// Custom red styling bypassing variants
```

## Theme Integration

The design system integrates with cbwsh's existing theme system. Themes provide color values that design tokens reference.

**Example Theme Integration:**
```go
import (
    "github.com/cbwinslow/cbwsh/pkg/ui/themes"
    "github.com/cbwinslow/cbwsh/pkg/ui/tokens"
)

// Load a theme
themeManager := themes.NewManager("~/.cbwsh/themes")
theme, _ := themeManager.Load("dracula")

// Create semantic colors from theme
semanticColors := tokens.DefaultSemanticColors()
// Map theme colors to semantic tokens as needed

// Use with components
btn := components.NewButton("Click me")
btn.SemanticColors = semanticColors
```

## Migration Guide

### Migrating Existing Components

To migrate existing UI code to use the design system:

1. **Replace hard-coded values with tokens:**
   ```go
   // Before
   style := lipgloss.NewStyle().Padding(0, 8)
   
   // After
   tokens := tokens.Default()
   style := lipgloss.NewStyle().Padding(0, tokens.Spacing.MD)
   ```

2. **Use semantic colors:**
   ```go
   // Before
   style := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
   
   // After
   colors := tokens.DefaultSemanticColors()
   style := lipgloss.NewStyle().Foreground(colors.StatusError)
   ```

3. **Adopt component composition:**
   ```go
   // Before: Custom styled elements
   
   // After: Use components
   btn := components.NewButton("Save")
   btn.Variant = components.ButtonPrimary
   ```

## Contributing

When adding new components or extending the design system:

1. Follow the existing token structure
2. Document all variants and options
3. Add comprehensive tests
4. Update this design system documentation
5. Ensure accessibility (keyboard navigation, screen readers)
6. Maintain visual consistency with other components

## Resources

- **Figma Design Principles**: https://www.figma.com/best-practices/
- **MCP Design Systems**: https://designsystem.university/articles/an-introduction-to-mcp-for-design-systems
- **Charm UI Components**: https://github.com/charmbracelet

---

*Last updated: January 2025*
