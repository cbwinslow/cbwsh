// Package main demonstrates the cbwsh design system usage.
//
// This example shows how to use the design tokens and components
// to build a cohesive terminal UI following Figma design principles.
package main

import (
	"fmt"

	"github.com/cbwinslow/cbwsh/pkg/ui/components"
	"github.com/cbwinslow/cbwsh/pkg/ui/tokens"
)

func main() {
	// Initialize design tokens
	tok := tokens.Default()
	colors := tokens.DefaultSemanticColors()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║         cbwsh Design System - Component Examples              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Example 1: Buttons
	fmt.Println("📌 Button Components:")
	fmt.Println()

	primaryBtn := components.NewButton("Save Changes")
	primaryBtn.Variant = components.ButtonPrimary
	primaryBtn.Size = components.ButtonSizeMedium
	primaryBtn.Icon = "💾"
	primaryBtn.Tokens = tok
	primaryBtn.SemanticColors = colors
	fmt.Println("  " + primaryBtn.Render())

	secondaryBtn := components.NewButton("Cancel")
	secondaryBtn.Variant = components.ButtonSecondary
	secondaryBtn.Size = components.ButtonSizeMedium
	secondaryBtn.Tokens = tok
	secondaryBtn.SemanticColors = colors
	fmt.Println("  " + secondaryBtn.Render())

	dangerBtn := components.NewButton("Delete Account")
	dangerBtn.Variant = components.ButtonDanger
	dangerBtn.Size = components.ButtonSizeLarge
	dangerBtn.Icon = "⚠️"
	dangerBtn.Tokens = tok
	dangerBtn.SemanticColors = colors
	fmt.Println("  " + dangerBtn.Render())

	fmt.Println()

	// Example 2: Cards
	fmt.Println("📌 Card Components:")
	fmt.Println()

	card1 := components.NewCard("User Profile", "Name: John Doe\nEmail: john@example.com\nRole: Developer")
	card1.Footer = "Last updated: 2 minutes ago"
	card1.Tokens = tok
	card1.SemanticColors = colors
	fmt.Println(card1.Render())
	fmt.Println()

	// Example 3: Badges
	fmt.Println("📌 Badge Components:")
	fmt.Println()

	badges := []struct {
		label   string
		variant components.BadgeVariant
	}{
		{"NEW", components.BadgeInfo},
		{"ACTIVE", components.BadgeSuccess},
		{"WARNING", components.BadgeWarning},
		{"ERROR", components.BadgeError},
		{"DRAFT", components.BadgeDefault},
	}

	badgeStack := components.NewStack(components.StackHorizontal)
	badgeStack.Spacing = tok.Spacing.MD
	badgeStack.Tokens = tok

	for _, b := range badges {
		badge := components.NewBadge(b.label)
		badge.Variant = b.variant
		badge.Tokens = tok
		badge.SemanticColors = colors
		badgeStack.Add(badge.Render())
	}

	fmt.Println("  " + badgeStack.Render())
	fmt.Println()

	// Example 4: Status Indicators
	fmt.Println("📌 Status Indicators:")
	fmt.Println()

	statuses := []struct {
		label  string
		status components.StatusType
	}{
		{"Server Running", components.StatusTypeSuccess},
		{"Deployment Pending", components.StatusTypeWarning},
		{"Build Failed", components.StatusTypeError},
		{"New Updates Available", components.StatusTypeInfo},
		{"Idle", components.StatusTypeNeutral},
	}

	statusStack := components.NewStack(components.StackVertical)
	statusStack.Spacing = tok.Spacing.SM
	statusStack.Tokens = tok

	for _, s := range statuses {
		status := components.NewStatusIndicator(s.label, s.status)
		status.ShowIcon = true
		status.Tokens = tok
		status.SemanticColors = colors
		statusStack.Add("  " + status.Render())
	}

	fmt.Println(statusStack.Render())
	fmt.Println()

	// Example 5: Complex Composition
	fmt.Println("📌 Complex Composition Example:")
	fmt.Println()

	// Create a card with multiple elements
	composedCard := components.NewCard(
		"🚀 Deployment Status",
		"",
	)

	// Create content with status indicators
	content := components.NewStack(components.StackVertical)
	content.Spacing = tok.Spacing.SM
	content.Tokens = tok

	success := components.NewStatusIndicator("Backend deployed", components.StatusTypeSuccess)
	success.Tokens = tok
	success.SemanticColors = colors

	pending := components.NewStatusIndicator("Frontend deploying", components.StatusTypeWarning)
	pending.Tokens = tok
	pending.SemanticColors = colors

	content.Add(success.Render())
	content.Add(pending.Render())

	composedCard.Content = content.Render()

	// Create footer with badge
	badge := components.NewBadge("v2.0.1")
	badge.Variant = components.BadgeInfo
	badge.Tokens = tok
	badge.SemanticColors = colors

	composedCard.Footer = "Version: " + badge.Render()
	composedCard.Tokens = tok
	composedCard.SemanticColors = colors

	fmt.Println(composedCard.Render())
	fmt.Println()

	// Example 6: Design Token Usage
	fmt.Println("📌 Design Tokens Reference:")
	fmt.Println()
	fmt.Println("  Spacing Scale (4px grid):")
	fmt.Printf("    XXS: %dpx, XS: %dpx, SM: %dpx, MD: %dpx\n",
		tok.Spacing.XXS, tok.Spacing.XS, tok.Spacing.SM, tok.Spacing.MD)
	fmt.Printf("    LG: %dpx, XL: %dpx, XXL: %dpx, XXXL: %dpx\n",
		tok.Spacing.LG, tok.Spacing.XL, tok.Spacing.XXL, tok.Spacing.XXXL)
	fmt.Println()

	fmt.Println("  Typography Scale:")
	fmt.Printf("    XS: %s, SM: %s, Base: %s, MD: %s\n",
		tok.Typography.FontXS, tok.Typography.FontSM,
		tok.Typography.FontBase, tok.Typography.FontMD)
	fmt.Printf("    LG: %s, XL: %s, XXL: %s, Huge: %s\n",
		tok.Typography.FontLG, tok.Typography.FontXL,
		tok.Typography.FontXXL, tok.Typography.FontHuge)
	fmt.Println()

	fmt.Println("  Animation Durations:")
	fmt.Printf("    Instant: %v, Fast: %v, Normal: %v\n",
		tok.Animation.DurationInstant,
		tok.Animation.DurationFast,
		tok.Animation.DurationNormal)
	fmt.Printf("    Slow: %v, Slower: %v\n",
		tok.Animation.DurationSlow,
		tok.Animation.DurationSlower)
	fmt.Println()

	fmt.Println("  Z-Index Scale:")
	fmt.Printf("    Base: %d, Dropdown: %d, Modal: %d\n",
		tok.ZIndex.Base, tok.ZIndex.Dropdown, tok.ZIndex.Modal)
	fmt.Printf("    Toast: %d, Overlay: %d, Debug: %d\n",
		tok.ZIndex.Toast, tok.ZIndex.Overlay, tok.ZIndex.Debug)
	fmt.Println()

	// Example 7: Layout with Stacks
	fmt.Println("📌 Layout Example (Vertical Stack):")
	fmt.Println()

	layoutStack := components.NewStack(components.StackVertical)
	layoutStack.Spacing = tok.Spacing.LG
	layoutStack.Tokens = tok

	header := components.NewCard("Header", "This is the header section")
	header.Tokens = tok
	header.SemanticColors = colors
	layoutStack.Add(header.Render())

	body := components.NewCard("Main Content", "This is the main content area\nwith multiple lines\nof text.")
	body.Tokens = tok
	body.SemanticColors = colors
	layoutStack.Add(body.Render())

	footer := components.NewCard("Footer", "This is the footer section")
	footer.Tokens = tok
	footer.SemanticColors = colors
	layoutStack.Add(footer.Render())

	fmt.Println(layoutStack.Render())
	fmt.Println()

	// Example 8: Dividers
	fmt.Println("📌 Divider Component:")
	fmt.Println()

	divider := components.NewDivider(60)
	divider.Tokens = tok
	divider.SemanticColors = colors
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("✨ Design System Examples Complete!")
	fmt.Println()
	fmt.Println("For more information, see DESIGN_SYSTEM.md")
}
