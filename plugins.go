package main

import (
	"fmt"
	"io"
	"os"

	"github.com/suprsokr/vanilladbc/pkg/dbc"
	"github.com/suprsokr/vanilladbc/pkg/dbd"
	"github.com/suprsokr/vanilladbc/pkg/plugin"
	csvplugin "github.com/suprsokr/vanilladbc-csv"
	jsonplugin "github.com/suprsokr/vanilladbc-json"
	mysqlplugin "github.com/suprsokr/vanilladbc-mysql"
)

// getPlugin returns a plugin writer instance based on the plugin name
func getPlugin(name string, outputFile string) (plugin.Writer, error) {
	var writer io.Writer = os.Stdout

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}
		writer = file
	}

	switch name {
	case "json":
		return jsonplugin.NewPretty(writer), nil
	case "csv":
		return csvplugin.New(writer), nil
	case "mysql":
		// MySQL requires configuration
		return nil, fmt.Errorf("mysql plugin requires configuration (use --mysql-* flags)")
	default:
		return nil, fmt.Errorf("unknown plugin: %s (available: json, csv, mysql)", name)
	}
}

// getReaderPlugin returns a plugin reader instance based on the plugin name
func getReaderPlugin(name string, inputFile string) (plugin.Reader, error) {
	var reader io.Reader

	if inputFile == "" || inputFile == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open input file: %w", err)
		}
		reader = file
	}

	switch name {
	case "json":
		return jsonplugin.NewReader(reader), nil
	case "csv":
		return csvplugin.NewReader(reader), nil
	case "mysql":
		// MySQL requires configuration
		return nil, fmt.Errorf("mysql plugin requires configuration (use --mysql-* flags)")
	default:
		return nil, fmt.Errorf("unknown plugin: %s (available: json, csv, mysql)", name)
	}
}

// getMySQLPlugin returns a MySQL plugin with the given configuration
func getMySQLPlugin(host string, port int, user, password, database, tableName string) (*mysqlplugin.Plugin, error) {
	config := mysqlplugin.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}
	
	return mysqlplugin.New(config, tableName)
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
