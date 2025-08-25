// Package gounqvet provides a Go static analysis tool that detects SELECT * usage
package gounqvet

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/gounqvet/internal/analyzer"
	"github.com/MirrexOne/gounqvet/pkg/config"
)

// Analyzer is the main gounqvet analyzer instance
// This is the primary export that golangci-lint will use
var Analyzer = analyzer.NewAnalyzer()

// New creates a new instance of the gounqvet analyzer
func New() *analysis.Analyzer {
	return Analyzer
}

// NewWithConfig creates a new analyzer instance with custom configuration
// This is the recommended way to use gounqvet with custom settings
func NewWithConfig(cfg *config.GounqvetSettings) *analysis.Analyzer {
	if cfg == nil {
		return Analyzer
	}
	return analyzer.NewAnalyzerWithSettings(*cfg)
}
