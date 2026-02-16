# vanilladbc-cli

Command-line tool for converting World of Warcraft Vanilla (1.0.0 - 1.12.3) DBC files to various formats.

## Features

- **Info Command** - Display DBC table schema information
- **Convert Command** - Convert DBC files using plugins
- **Plugin System** - Extensible plugin architecture for format conversion

## Installation

```bash
go install github.com/suprsokr/vanilladbc-cli@latest
```

Or build from source:

```bash
git clone https://github.com/suprsokr/vanilladbc-cli.git
cd vanilladbc-cli
go build
```

## Usage

### Show Table Schema

```bash
vanilladbc info <dbd_file> <build>

# Example
vanilladbc info Spell.dbd 1.12.1.5875
```

Output:
```
Table: Spell
Build: 1.12.1.5875

Total Columns Defined: 145
Fields in Build: 86

Field Definitions:
------------------
  1. ID                             int<32> (ID)
  2. School                         int<32>
  3. Category                       int<32> -> SpellCategory::ID
  ...
```

### Convert DBC Files

```bash
vanilladbc convert <dbc_file> <dbd_file> <build> --plugin <plugin> [--output <file>]

# Example: Convert to JSON
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin json --output spell.json

# Example: Output to stdout
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin json
```

## Available Plugins

Plugins are separate packages that implement format conversion:

### JSON Plugin
- **Package**: [vanilladbc-json](https://github.com/suprsokr/vanilladbc-json)
- **Install**: `go get github.com/suprsokr/vanilladbc-json`
- **Usage**: `--plugin json`

### Future Plugins
- **MySQL** - Convert to MySQL INSERT statements
- **CSV** - Convert to CSV format
- **XML** - Convert to XML format
- **SQL** - Convert to generic SQL

## Creating a Plugin

Plugins implement the `plugin.Writer` interface from [vanilladbc](https://github.com/suprsokr/vanilladbc):

```go
package myplugin

import (
    "github.com/suprsokr/vanilladbc/pkg/plugin"
    "github.com/suprsokr/vanilladbc/pkg/dbc"
    "github.com/suprsokr/vanilladbc/pkg/dbd"
)

type MyPlugin struct {
    // your fields
}

func (p *MyPlugin) WriteHeader(versionDef *dbd.VersionDefinition, 
                                columns map[string]dbd.ColumnDefinition) error {
    // Setup your output format
    return nil
}

func (p *MyPlugin) WriteRecord(record dbc.Record) error {
    // Convert and write each record
    return nil
}

func (p *MyPlugin) WriteFooter() error {
    // Finalize output
    return nil
}
```

Then register your plugin in `vanilladbc-cli/plugins.go` or use Go's plugin system.

## Dependencies

- [vanilladbc](https://github.com/suprsokr/vanilladbc) - Core DBC/DBD library
- [VanillaDBDefs](https://github.com/suprsokr/VanillaDBDefs) - Database definitions

## Related Projects

- [vanilladbc](https://github.com/suprsokr/vanilladbc) - Core library
- [vanilladbc-json](https://github.com/suprsokr/vanilladbc-json) - JSON plugin
- [VanillaDBDefs](https://github.com/suprsokr/VanillaDBDefs) - Vanilla database definitions

## License

MIT License - See LICENSE file for details
