// Package sqlvet provides static analysis for SQL queries to detect SELECT * usage
package sqlvet

import (
	"golang.org/x/tools/go/analysis"

	"github.com/MirrexOne/sqlvet/pkg/analyzer"
)

// Analyzer is the main sqlvet analyzer for golangci-lint integration
var Analyzer *analysis.Analyzer = analyzer.New()
