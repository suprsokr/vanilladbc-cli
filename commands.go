package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/suprsokr/vanilladbc/pkg/dbd"
)

func runInfo(dbdFile string, buildStr string) error {
	// Parse DBD file
	fmt.Printf("Parsing DBD file: %s\n", dbdFile)
	dbdDef, err := dbd.ParseFile(dbdFile)
	if err != nil {
		return fmt.Errorf("failed to parse DBD file: %w", err)
	}

	// Get version definition
	build, err := dbd.NewBuild(buildStr)
	if err != nil {
		return fmt.Errorf("invalid build: %w", err)
	}

	versionDef, err := dbdDef.GetVersionDefinition(*build)
	if err != nil {
		return fmt.Errorf("failed to get version definition: %w", err)
	}

	// Print info
	tableName := filepath.Base(dbdFile)
	tableName = tableName[:len(tableName)-4] // Remove .dbd extension

	fmt.Println()
	fmt.Printf("Table: %s\n", tableName)
	fmt.Printf("Build: %s\n", buildStr)
	fmt.Println()
	fmt.Printf("Total Columns Defined: %d\n", len(dbdDef.Columns))
	fmt.Printf("Fields in Build: %d\n", len(versionDef.Definitions))
	fmt.Println()
	fmt.Println("Field Definitions:")
	fmt.Println("------------------")

	for i, def := range versionDef.Definitions {
		colDef := dbdDef.Columns[def.Column]
		typeStr := string(colDef.Type)

		if colDef.Type == dbd.TypeInt || colDef.Type == dbd.TypeUInt {
			signStr := ""
			if def.IsUnsigned {
				signStr = "unsigned "
			}
			typeStr = fmt.Sprintf("%s%s<%d>", signStr, typeStr, def.Size)
		}

		arrayStr := ""
		if def.ArraySize > 0 {
			arrayStr = fmt.Sprintf("[%d]", def.ArraySize)
		}

		idStr := ""
		if def.IsID {
			idStr = " (ID)"
		}

		fkStr := ""
		if colDef.ForeignTable != "" {
			fkStr = fmt.Sprintf(" -> %s::%s", colDef.ForeignTable, colDef.ForeignColumn)
		}

		fmt.Printf("%3d. %-30s %s%s%s%s\n", i+1, def.Column, typeStr, arrayStr, idStr, fkStr)
	}

	return nil
}

func runConvert(dbcFile, dbdFile, buildStr, pluginName, outputFile string) error {
	// Parse DBD file
	fmt.Printf("Parsing DBD file: %s\n", dbdFile)
	dbdDef, err := dbd.ParseFile(dbdFile)
	if err != nil {
		return fmt.Errorf("failed to parse DBD file: %w", err)
	}

	// Get plugin
	plugin, err := getPlugin(pluginName, outputFile)
	if err != nil {
		return fmt.Errorf("failed to get plugin: %w", err)
	}

	// Get version definition
	build, err := dbd.NewBuild(buildStr)
	if err != nil {
		return fmt.Errorf("invalid build: %w", err)
	}

	versionDef, err := dbdDef.GetVersionDefinition(*build)
	if err != nil {
		return fmt.Errorf("failed to get version definition: %w", err)
	}

	// Open DBC file
	file, err := os.Open(dbcFile)
	if err != nil {
		return fmt.Errorf("failed to open DBC file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Printf("Converting %s using %s plugin...\n", dbcFile, pluginName)
	if err := plugin.WriteHeader(versionDef, dbdDef.Columns); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Stream records through plugin
	recordCount := 0
	err = streamToPlugin(file, dbdDef, *build, plugin, &recordCount)
	if err != nil {
		return fmt.Errorf("failed to convert records: %w", err)
	}

	// Write footer
	if err := plugin.WriteFooter(); err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	fmt.Printf("Successfully converted %d records\n", recordCount)
	if outputFile != "" {
		fmt.Printf("Output written to: %s\n", outputFile)
	}

	return nil
}

func runImport(inputFile, dbdFile, buildStr, pluginName, outputFile string) error {
	// Parse DBD file
	fmt.Printf("Parsing DBD file: %s\n", dbdFile)
	dbdDef, err := dbd.ParseFile(dbdFile)
	if err != nil {
		return fmt.Errorf("failed to parse DBD file: %w", err)
	}

	// Get version definition
	build, err := dbd.NewBuild(buildStr)
	if err != nil {
		return fmt.Errorf("invalid build: %w", err)
	}

	versionDef, err := dbdDef.GetVersionDefinition(*build)
	if err != nil {
		return fmt.Errorf("failed to get version definition: %w", err)
	}

	// Get reader plugin
	readerPlugin, err := getReaderPlugin(pluginName, inputFile)
	if err != nil {
		return fmt.Errorf("failed to get reader plugin: %w", err)
	}
	defer readerPlugin.Close()

	// Set schema for plugins that need it (JSON, CSV)
	if setter, ok := readerPlugin.(interface {
		SetSchema(*dbd.VersionDefinition, map[string]dbd.ColumnDefinition)
	}); ok {
		setter.SetSchema(versionDef, dbdDef.Columns)
	}

	// Read header
	fmt.Printf("Reading from %s using %s plugin...\n", inputFile, pluginName)
	_, _, err = readerPlugin.ReadHeader()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Read all records
	var records []interface{}
	recordCount := 0
	for {
		record, err := readerPlugin.ReadRecord()
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}
		if record == nil {
			break // No more records
		}
		records = append(records, record)
		recordCount++
	}

	fmt.Printf("Read %d records\n", recordCount)

	// TODO: Write records to DBC file using the DBC writer from vanilladbc library
	fmt.Printf("Warning: DBC writing not yet fully implemented\n")
	fmt.Printf("Records would be written to: %s\n", outputFile)

	return nil
}
