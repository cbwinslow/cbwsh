// Package main demonstrates the cbwsh component registry.
//
// This example shows how to browse, search, and use the component registry
// to discover available UI components.
package main

import (
	"fmt"
	"strings"

	"github.com/cbwinslow/cbwsh/pkg/ui/components"
	"github.com/cbwinslow/cbwsh/pkg/ui/registry"
	"github.com/cbwinslow/cbwsh/pkg/ui/tokens"
)

func main() {
	// Initialize the registry
	reg := registry.DefaultRegistry()
	tok := tokens.Default()
	colors := tokens.DefaultSemanticColors()

	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║           cbwsh Component Registry - Browser                  ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Example 1: List all components
	fmt.Println("📦 All Available Components:")
	fmt.Println()

	allComponents := reg.List()
	for _, component := range allComponents {
		badge := components.NewBadge(component.Category)
		badge.Tokens = tok
		badge.SemanticColors = colors

		switch component.Category {
		case "interactive":
			badge.Variant = components.BadgeInfo
		case "layout":
			badge.Variant = components.BadgeSuccess
		case "display":
			badge.Variant = components.BadgeWarning
		}

		fmt.Printf("  • %s %s\n", component.DisplayName, badge.Render())
		fmt.Printf("    %s\n", component.Description)
		fmt.Println()
	}

	// Example 2: Browse by category
	divider := components.NewDivider(64)
	divider.Tokens = tok
	divider.SemanticColors = colors
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("📂 Components by Category:")
	fmt.Println()

	categories := reg.Categories()
	for _, category := range categories {
		categoryBadge := components.NewBadge(strings.ToUpper(category))
		categoryBadge.Tokens = tok
		categoryBadge.SemanticColors = colors

		switch category {
		case "interactive":
			categoryBadge.Variant = components.BadgeInfo
		case "layout":
			categoryBadge.Variant = components.BadgeSuccess
		case "display":
			categoryBadge.Variant = components.BadgeWarning
		}

		fmt.Printf("  %s\n", categoryBadge.Render())
		fmt.Println()

		categoryComponents := reg.ListByCategory(category)
		for _, component := range categoryComponents {
			fmt.Printf("    • %s - %s\n", component.DisplayName, component.Description)
		}
		fmt.Println()
	}

	// Example 3: Search components
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("🔍 Search Examples:")
	fmt.Println()

	searches := []string{"button", "status", "layout", "container"}
	for _, query := range searches {
		results := reg.Search(query)
		fmt.Printf("  Search: %q → %d result(s)\n", query, len(results))
		for _, component := range results {
			fmt.Printf("    • %s\n", component.DisplayName)
		}
		fmt.Println()
	}

	// Example 4: Component details
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("📝 Component Details Example:")
	fmt.Println()

	buttonComponent, _ := reg.Get("button")
	card := components.NewCard(
		"🔘 "+buttonComponent.DisplayName,
		fmt.Sprintf("%s\n\nVersion: %s\nCategory: %s\nTags: %s",
			buttonComponent.Description,
			buttonComponent.Version,
			buttonComponent.Category,
			strings.Join(buttonComponent.Tags, ", "),
		),
	)
	card.Footer = "Files: " + strings.Join(buttonComponent.Files, ", ")
	card.Tokens = tok
	card.SemanticColors = colors
	fmt.Println(card.Render())
	fmt.Println()

	// Example 5: Show example code
	fmt.Println("💡 Example Usage:")
	fmt.Println()
	fmt.Println("```go")
	fmt.Println(buttonComponent.Example)
	fmt.Println("```")
	fmt.Println()

	// Example 6: Component composition showcase
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("🎨 Component Composition Showcase:")
	fmt.Println()

	// Create a composed UI using registry components
	showcaseCard := components.NewCard("Component Registry", "Discover and use terminal UI components")

	// Create a stack of status indicators
	statusStack := components.NewStack(components.StackVertical)
	statusStack.Spacing = tok.Spacing.SM
	statusStack.Tokens = tok

	componentCount := len(allComponents)
	categoryCount := len(categories)

	status1 := components.NewStatusIndicator(
		fmt.Sprintf("%d components available", componentCount),
		components.StatusTypeSuccess,
	)
	status1.ShowIcon = true
	status1.Tokens = tok
	status1.SemanticColors = colors

	status2 := components.NewStatusIndicator(
		fmt.Sprintf("%d categories", categoryCount),
		components.StatusTypeInfo,
	)
	status2.ShowIcon = true
	status2.Tokens = tok
	status2.SemanticColors = colors

	statusStack.Add(status1.Render())
	statusStack.Add(status2.Render())

	showcaseCard.Content = statusStack.Render()

	// Create badge stack
	badgeStack := components.NewStack(components.StackHorizontal)
	badgeStack.Spacing = tok.Spacing.MD
	badgeStack.Tokens = tok

	for _, category := range categories {
		catBadge := components.NewBadge(category)
		catBadge.Tokens = tok
		catBadge.SemanticColors = colors

		switch category {
		case "interactive":
			catBadge.Variant = components.BadgeInfo
		case "layout":
			catBadge.Variant = components.BadgeSuccess
		case "display":
			catBadge.Variant = components.BadgeWarning
		}

		badgeStack.Add(catBadge.Render())
	}

	showcaseCard.Footer = badgeStack.Render()
	showcaseCard.Tokens = tok
	showcaseCard.SemanticColors = colors

	fmt.Println(showcaseCard.Render())
	fmt.Println()

	// Example 7: Registry stats
	fmt.Println(divider.Render())
	fmt.Println()

	fmt.Println("📊 Registry Statistics:")
	fmt.Println()

	statsCard := components.NewCard("Registry Info", "")

	var stats []string
	stats = append(stats, fmt.Sprintf("Total Components: %d", componentCount))
	stats = append(stats, fmt.Sprintf("Categories: %d", categoryCount))

	for _, category := range categories {
		count := len(reg.ListByCategory(category))
		stats = append(stats, fmt.Sprintf("  • %s: %d", category, count))
	}

	statsCard.Content = strings.Join(stats, "\n")
	statsCard.Tokens = tok
	statsCard.SemanticColors = colors

	fmt.Println(statsCard.Render())
	fmt.Println()

	fmt.Println("✨ Component Registry Browser Complete!")
	fmt.Println()
	fmt.Println("For more information:")
	fmt.Println("  • See registry/README.md for registry documentation")
	fmt.Println("  • See DESIGN_SYSTEM.md for design system documentation")
	fmt.Println("  • Run examples/design_system.go to see component examples")
}
