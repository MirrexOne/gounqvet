package analyzer

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/gounqvet/internal/analyzer"
)

// Analyzer is the gounqvet analyzer for detecting SELECT * usage in SQL queries
var Analyzer = analyzer.NewAnalyzer()

// New creates a new instance of the gounqvet analyzer
func New() *analysis.Analyzer {
	return Analyzer
}
