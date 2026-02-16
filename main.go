package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "info":
		if err := cmdInfo(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "convert":
		if err := cmdConvert(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("vanilladbc - Vanilla WoW DBC file converter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  vanilladbc info <dbd_file> <build>")
	fmt.Println("  vanilladbc convert <dbc_file> <dbd_file> <build> --plugin <plugin> [--output <file>]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  info     - Show DBC table schema information")
	fmt.Println("  convert  - Convert DBC file using a plugin")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Show schema")
	fmt.Println("  vanilladbc info Spell.dbd 1.12.1.5875")
	fmt.Println()
	fmt.Println("  # Convert to JSON")
	fmt.Println("  vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin json --output spell.json")
	fmt.Println()
	fmt.Println("Available Plugins:")
	fmt.Println("  json     - Convert to/from JSON format")
	fmt.Println()
	fmt.Println("For more information about plugins, see:")
	fmt.Println("  https://github.com/suprsokr/vanilladbc-json")
}

func cmdInfo() error {
	if len(os.Args) < 4 {
		return fmt.Errorf("usage: vanilladbc info <dbd_file> <build>")
	}

	dbdFile := os.Args[2]
	buildStr := os.Args[3]

	return runInfo(dbdFile, buildStr)
}

func cmdConvert() error {
	if len(os.Args) < 6 {
		return fmt.Errorf("usage: vanilladbc convert <dbc_file> <dbd_file> <build> --plugin <plugin> [--output <file>]")
	}

	dbcFile := os.Args[2]
	dbdFile := os.Args[3]
	buildStr := os.Args[4]

	// Parse flags
	var pluginName string
	var outputFile string

	for i := 5; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--plugin":
			if i+1 >= len(os.Args) {
				return fmt.Errorf("--plugin requires a value")
			}
			pluginName = os.Args[i+1]
			i++
		case "--output":
			if i+1 >= len(os.Args) {
				return fmt.Errorf("--output requires a value")
			}
			outputFile = os.Args[i+1]
			i++
		}
	}

	if pluginName == "" {
		return fmt.Errorf("--plugin is required")
	}

	return runConvert(dbcFile, dbdFile, buildStr, pluginName, outputFile)
}
