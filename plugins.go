package main

import (
	"fmt"
	"io"

	"github.com/suprsokr/vanilladbc/pkg/dbc"
	"github.com/suprsokr/vanilladbc/pkg/dbd"
	"github.com/suprsokr/vanilladbc/pkg/plugin"
)

// getPlugin returns a plugin instance based on the plugin name
func getPlugin(name string, outputFile string) (plugin.Writer, error) {
	// var writer io.Writer = os.Stdout
	//
	// if outputFile != "" {
	// 	file, err := os.Create(outputFile)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to create output file: %w", err)
	// 	}
	// 	writer = file
	// }

	switch name {
	case "json":
		// JSON plugin would be imported here if available
		// For now, return error indicating plugin not available
		return nil, fmt.Errorf("json plugin not available - install vanilladbc-json package")
	default:
		return nil, fmt.Errorf("unknown plugin: %s", name)
	}
}

// streamToPlugin reads DBC records and streams them to a plugin
func streamToPlugin(r io.Reader, dbdDef *dbd.DBDefinition, build dbd.Build, plugin plugin.Writer, recordCount *int) error {
	iter, err := dbc.NewIterator(r, dbdDef, build)
	if err != nil {
		return err
	}

	for iter.Next() {
		if err := plugin.WriteRecord(iter.Record()); err != nil {
			return fmt.Errorf("plugin failed at record %d: %w", iter.Index(), err)
		}
		*recordCount++
	}

	return iter.Err()
}
