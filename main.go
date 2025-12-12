// cbwsh - Custom Bubble Tea Shell
//
// A modern, modular terminal shell built with the Bubble Tea ecosystem.
// cbwsh provides rich TUI components, AI integration, multi-pane support,
// and advanced shell features in an elegant terminal interface.
//
// Usage:
//
//	cbwsh [flags]
//
// The shell starts in interactive mode by default. Use keyboard shortcuts
// to access various features:
//   - Ctrl+Q: Quit
//   - Ctrl+?: Help
//   - Ctrl+A: AI assist mode
//   - Ctrl+M: Toggle AI monitor
//
// For more information, visit: https://github.com/cbwinslow/cbwsh
package main

import (
	"fmt"
	"os"

	"github.com/cbwinslow/cbwsh/internal/app"
)

// main is the entry point for the cbwsh application.
// It initializes and runs the shell, handling any errors gracefully.
func main() {
	// Run the application and handle errors
	if err := app.Run(); err != nil {
		// Print error to stderr with a clear prefix
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		
		// Exit with non-zero status code to indicate failure
		os.Exit(1)
	}
	
	// Successful execution - exit with status 0 (implicit)
}
