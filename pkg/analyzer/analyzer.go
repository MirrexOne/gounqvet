package analyzer

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/unqueryvet/internal/analyzer"
)

// Analyzer is the unqueryvet analyzer for detecting SELECT * usage in SQL queries
var Analyzer = analyzer.NewAnalyzer()

// New creates a new instance of the unqueryvet analyzer
func New() *analysis.Analyzer {
	return Analyzer
}
