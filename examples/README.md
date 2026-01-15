# cbwsh Examples

This directory contains example programs demonstrating various features of cbwsh.

## Available Examples

### Design System Example (`design_system.go`)

The `design_system.go` example showcases the Figma-inspired design system, including:

- Design tokens (spacing, typography, borders, animations, z-index)
- Component library (buttons, cards, badges, status indicators, etc.)
- Component composition patterns
- Layout management with stacks
- Semantic color usage

### Component Registry Browser (`registry_browser.go`)

The `registry_browser.go` example demonstrates the component registry system inspired by shadcn/ui:

- Browse all available components
- Search components by name, description, or tags
- Filter components by category
- View component details and metadata
- See example usage code
- Component composition showcase

### Running the Examples

```bash
cd examples
go run design_system.go      # Design system demo
go run registry_browser.go   # Component registry browser
```

### What You'll See

The example demonstrates:

1. **Button Components** - Multiple variants (primary, secondary, danger) and sizes
2. **Card Components** - Structured containers with title, content, and footer
3. **Badge Components** - Status indicators with color-coded variants
4. **Status Indicators** - Icon + label combinations for status display
5. **Complex Composition** - Building sophisticated UIs by composing simple components
6. **Design Tokens Reference** - The complete token system values
7. **Layout Examples** - Using stacks for vertical and horizontal layouts
8. **Dividers** - Visual separators between content sections

## Learning Resources

- **[Design System Documentation](../DESIGN_SYSTEM.md)** - Complete guide to the design system
- **[Usage Guide](../USAGE.md)** - General cbwsh usage instructions
- **[Roadmap](../ROADMAP.md)** - Future features and plans

## Creating Your Own Examples

To create a new example:

1. Create a new `.go` file in this directory
2. Import the necessary cbwsh packages
3. Add comprehensive comments explaining what the example demonstrates
4. Add the example to this README

Example template:

```go
// Package main demonstrates [feature name].
package main

import (
    "fmt"
    "github.com/cbwinslow/cbwsh/pkg/ui/components"
    "github.com/cbwinslow/cbwsh/pkg/ui/tokens"
)

func main() {
    // Your example code here
    fmt.Println("Example output")
}
```

## Contributing

We welcome contributions of example programs! Good examples:

- Focus on a single feature or concept
- Include clear explanations and comments
- Demonstrate best practices
- Are self-contained and easy to run
- Show practical, real-world usage

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more information.
