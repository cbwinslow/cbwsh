# cbwsh Component Registry

A component registry system for cbwsh terminal UI components, inspired by [shadcn/ui](https://ui.shadcn.com/).

## Overview

The cbwsh component registry provides a centralized catalog of reusable terminal UI components built with the Bubble Tea ecosystem. Components follow Figma-inspired design principles and use design tokens for consistency.

## Features

- 🎨 **Curated Components**: Hand-picked, high-quality terminal UI components
- 🔍 **Searchable**: Find components by name, description, category, or tags
- 📦 **Categorized**: Components organized by type (interactive, layout, display)
- 📝 **Well-Documented**: Each component includes description, example, and metadata
- 🎯 **Design System**: All components use consistent design tokens
- 💻 **Copy-Paste Ready**: Example code ready to use in your projects

## Available Components

### Interactive

- **Button** - A versatile button component with multiple variants and sizes

### Layout

- **Card** - Container component for grouping related content
- **Stack** - Organize items vertically or horizontally
- **Divider** - Visual separator between content sections

### Display

- **Badge** - Small label for status indicators and counts
- **Status Indicator** - Display status with icon and label

## Component Categories

### Interactive Components
Components that users interact with (buttons, inputs, etc.)

### Layout Components
Components for structuring and organizing content (cards, stacks, grids, etc.)

### Display Components
Components for displaying information (badges, status indicators, alerts, etc.)

## Usage

### Using the Registry API

```go
package main

import (
    "fmt"
    "github.com/cbwinslow/cbwsh/pkg/ui/registry"
)

func main() {
    // Load the default registry
    reg := registry.DefaultRegistry()
    
    // List all components
    components := reg.List()
    for _, c := range components {
        fmt.Printf("%s - %s\n", c.DisplayName, c.Description)
    }
    
    // Get a specific component
    button, err := reg.Get("button")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Example:\n%s\n", button.Example)
    
    // Search for components
    results := reg.Search("status")
    for _, c := range results {
        fmt.Printf("Found: %s\n", c.DisplayName)
    }
    
    // List components by category
    layoutComponents := reg.ListByCategory("layout")
    for _, c := range layoutComponents {
        fmt.Printf("Layout: %s\n", c.DisplayName)
    }
}
```

### Using Components in Your Code

All components are available in the `pkg/ui/components` package:

```go
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
    
    // Create a button
    btn := components.NewButton("Save Changes")
    btn.Variant = components.ButtonPrimary
    btn.Size = components.ButtonSizeMedium
    btn.Icon = "💾"
    btn.Tokens = tok
    btn.SemanticColors = colors
    fmt.Println(btn.Render())
    
    // Create a card
    card := components.NewCard("User Profile", "John Doe\njohn@example.com")
    card.Footer = "Last updated: 2 minutes ago"
    card.Elevated = true
    card.Tokens = tok
    card.SemanticColors = colors
    fmt.Println(card.Render())
    
    // Create a stack with multiple items
    stack := components.NewStack(components.StackVertical)
    stack.Spacing = tok.Spacing.MD
    
    badge1 := components.NewBadge("NEW")
    badge1.Variant = components.BadgeInfo
    badge1.Tokens = tok
    badge1.SemanticColors = colors
    
    badge2 := components.NewBadge("ACTIVE")
    badge2.Variant = components.BadgeSuccess
    badge2.Tokens = tok
    badge2.SemanticColors = colors
    
    stack.Add(badge1.Render())
    stack.Add(badge2.Render())
    fmt.Println(stack.Render())
}
```

## Component Metadata

Each component in the registry includes:

- **Name**: Unique identifier (e.g., "button")
- **Display Name**: Human-readable name (e.g., "Button")
- **Description**: What the component does
- **Category**: Component type (interactive, layout, display)
- **Version**: Semantic version
- **Files**: Source files containing the component
- **Example**: Ready-to-use example code
- **Tags**: Searchable keywords
- **Dependencies**: Other components this depends on

## Design System Integration

All components use the cbwsh design system:

- **Design Tokens**: Consistent spacing, typography, colors, borders, animations
- **Semantic Colors**: Meaningful color tokens (interactive states, text hierarchy, status)
- **Component Composition**: Components can be nested and composed together

See [DESIGN_SYSTEM.md](../DESIGN_SYSTEM.md) for complete design system documentation.

## Adding Components to the Registry

To add a new component to the registry:

1. Create the component in `pkg/ui/components/`
2. Add component metadata to `registry/components.json`
3. Include example usage code
4. Add appropriate tags for searchability
5. Update this README

Example metadata:

```json
{
  "name": "my-component",
  "display_name": "My Component",
  "description": "A brief description of what this component does",
  "category": "interactive",
  "files": [
    "pkg/ui/components/my-component.go"
  ],
  "example": "component := components.NewMyComponent()\noutput := component.Render()",
  "version": "1.0.0",
  "tags": [
    "interactive",
    "custom"
  ]
}
```

## Component Guidelines

When creating components for the registry:

1. **Use Design Tokens**: Always use design tokens instead of hard-coded values
2. **Follow Conventions**: Follow existing component structure and naming
3. **Provide Examples**: Include clear, working example code
4. **Document Variants**: Clearly document all available variants and options
5. **Compose When Possible**: Build on existing components rather than starting from scratch
6. **Test Thoroughly**: Include comprehensive tests for all functionality
7. **Maintain Accessibility**: Consider keyboard navigation and screen readers

## Registry File Format

Components are stored in JSON format at `registry/components.json`:

```json
[
  {
    "name": "component-name",
    "display_name": "Component Display Name",
    "description": "Component description",
    "category": "interactive|layout|display",
    "files": ["pkg/ui/components/file.go"],
    "example": "Example code",
    "version": "1.0.0",
    "tags": ["tag1", "tag2"],
    "dependencies": ["other-component"]
  }
]
```

## Examples

See the [examples directory](../examples/) for complete working examples:

```bash
# Run the design system example to see all components
cd examples
go run design_system.go
```

## Contributing

We welcome contributions! To add a component to the registry:

1. Create the component following our design system guidelines
2. Add comprehensive tests
3. Update the registry JSON file
4. Add example usage
5. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## Future Plans

- CLI tool for browsing and searching components (`cbwsh registry`)
- Interactive component browser
- Component templates and generators
- Version management and updates
- Community component submissions
- Component playground/previewer

## Inspiration

This registry system is inspired by:

- [shadcn/ui](https://ui.shadcn.com/) - Component registry for React
- [Figma](https://www.figma.com/) - Design system principles
- [Storybook](https://storybook.js.org/) - Component documentation

## License

MIT License - See [LICENSE](../LICENSE) for details

---

*Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and the Charm ecosystem*
