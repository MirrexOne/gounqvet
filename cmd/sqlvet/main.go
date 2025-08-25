package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	internal "github.com/MirrexOne/sqlvet/internal/analyzer"
)

// main provides CLI interface for SQLVet linter
//
// This file creates a standalone version of the analyzer that can be run
// directly from the command line without golangci-lint:
//
//	sqlvet ./...
//	sqlvet -help
//	sqlvet /path/to/project
//
// Uses standard singlechecker from golang.org/x/tools, which provides:
// - Consistent CLI interface with other Go analyzers
// - Automatic support for standard flags (-json, -c=N etc.)
// - Proper output formatting for various tools
func main() {
	// NewAnalyzer() creates analyzer instance from internal package
	// singlechecker.Main() handles CLI arguments and runs analysis
	singlechecker.Main(internal.NewAnalyzer())
}
