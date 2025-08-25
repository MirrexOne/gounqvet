// Package sqlvet provides a Go static analysis tool that detects SELECT * usage
package sqlvet

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/sqlvet/internal/analyzer"
	"github.com/MirrexOne/sqlvet/pkg/config"
)

// Analyzer is the main sqlvet analyzer instance
// This is the primary export that golangci-lint will use
var Analyzer = analyzer.NewAnalyzer()

// New creates a new instance of the sqlvet analyzer
func New() *analysis.Analyzer {
	return Analyzer
}

// NewWithConfig creates a new analyzer instance with custom configuration
// This is the recommended way to use sqlvet with custom settings
func NewWithConfig(cfg *config.SQLVetSettings) *analysis.Analyzer {
	if cfg == nil {
		return Analyzer
	}
	return analyzer.NewAnalyzerWithSettings(*cfg)
}
