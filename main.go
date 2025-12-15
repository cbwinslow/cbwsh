// Package main is the entry point for cbwsh - Custom Bubble Tea Shell
// cbwsh is a modern, modular terminal shell built with the Bubble Tea ecosystem,
// featuring rich TUI components, AI integration, and extensive customization.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/cbwinslow/cbwsh/internal/app"
)

var (
	// Version information (set via ldflags during build)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	// Command-line flags
	showVersion = flag.Bool("version", false, "Show version information and exit")
	showHelp    = flag.Bool("help", false, "Show help information and exit")
	configFile  = flag.String("config", "", "Path to configuration file (default: ~/.cbwsh/config.yaml)")
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Handle version flag
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// Set up panic recovery for graceful error handling
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %v\n", r)
			fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", debug.Stack())
			os.Exit(1)
		}
	}()

	// Run the application with error handling
	if err := app.Run(*configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// printVersion displays version information
func printVersion() {
	fmt.Printf("cbwsh version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built:  %s\n", date)
}

// printHelp displays usage information
func printHelp() {
	fmt.Println("cbwsh - Custom Bubble Tea Shell")
	fmt.Println()
	fmt.Println("A modern, modular terminal shell with rich TUI components,")
	fmt.Println("AI integration, and extensive customization options.")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  cbwsh [flags]")
	fmt.Println()
	fmt.Println("FLAGS:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("KEY BINDINGS:")
	fmt.Println("  Ctrl+Q          Quit")
	fmt.Println("  Ctrl+? or F1    Show help")
	fmt.Println("  Ctrl+A          AI assist mode")
	fmt.Println("  Ctrl+M          Toggle AI monitor")
	fmt.Println("  Ctrl+L          Clear screen")
	fmt.Println("  Enter           Execute command")
	fmt.Println("  Tab             Autocomplete")
	fmt.Println()
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Default config: ~/.cbwsh/config.yaml")
	fmt.Println()
	fmt.Println("For more information, visit:")
	fmt.Println("  https://github.com/cbwinslow/cbwsh")
}
