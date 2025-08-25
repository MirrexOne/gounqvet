// Package analyzer provides the SQLVet analyzer for detecting SELECT * usage
package analyzer

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/sqlvet/internal/analyzer"
)

// Analyzer is the sqlvet analyzer for detecting SELECT * usage in SQL queries
var Analyzer = analyzer.NewAnalyzer()

// New creates a new instance of the sqlvet analyzer
func New() *analysis.Analyzer {
	return Analyzer
}