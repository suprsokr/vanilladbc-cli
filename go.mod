module github.com/suprsokr/vanilladbc-cli

go 1.25.5

replace github.com/suprsokr/vanilladbc => ../vanilladbc

require github.com/suprsokr/vanilladbc v0.0.0-00010101000000-000000000000

require (
	github.com/suprsokr/vanilladbc-csv v0.0.0-00010101000000-000000000000
	github.com/suprsokr/vanilladbc-json v0.0.0-00010101000000-000000000000
	github.com/suprsokr/vanilladbc-mysql v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
)

replace github.com/suprsokr/vanilladbc-json => ../vanilladbc-json

replace github.com/suprsokr/vanilladbc-csv => ../vanilladbc-csv

replace github.com/suprsokr/vanilladbc-mysql => ../vanilladbc-mysql
