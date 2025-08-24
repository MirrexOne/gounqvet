package sqlvet

import (
	"github.com/MirrexOne/sqlvet/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

// Analyzer экспортированная переменная для golangci-lint
var Analyzer *analysis.Analyzer = analyzer.NewAnalyzer()