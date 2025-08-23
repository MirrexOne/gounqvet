package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/mirrexone/sqlvet/internal/analyzer"
)

// This is a plugin entry point for golangci-lint.
// It's a map of analyzer names to their constructors.
// golangci-lint will look for a variable named `New` in the plugin package.
var New = analyzer.NewFromGolangciLint

// main - это точка входа для запуска линтера как отдельного инструмента.
func main() {
	multichecker.Main(
		analyzer.NewDefaultAnalyzer(),
	)
}
