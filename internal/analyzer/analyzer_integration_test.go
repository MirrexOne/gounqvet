package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/MirrexOne/unqueryvet/internal/analyzer"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(), "a")
}

func TestAnalyzerWithSettings(t *testing.T) {
	testdata := analysistest.TestData()

	// Test with default settings
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(), "clean")

	// Test with SQL builders detection
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(), "integration")
}
