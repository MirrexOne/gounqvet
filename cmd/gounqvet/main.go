// Package main provides the command-line interface for the gounqvet analyzer.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	internal "github.com/MirrexOne/gounqvet/internal/analyzer"
)

func main() {
	singlechecker.Main(internal.NewAnalyzer())
}
