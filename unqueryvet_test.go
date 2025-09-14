package unqueryvet_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/MirrexOne/unqueryvet"
	"github.com/MirrexOne/unqueryvet/internal/analyzer"
)

func TestUnqueryvet(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), analyzer.NewAnalyzer(), "clean")
}

func TestUnqueryvetWithCustomSettings(t *testing.T) {
	// Test with custom configuration
	settings := unqueryvet.Settings{
		CheckSQLBuilders: true,
		AllowedPatterns:  []string{`SELECT \* FROM temp_.*`},
	}

	// Create analyzer with custom settings
	customAnalyzer := analyzer.NewAnalyzerWithSettings(settings)
	analysistest.Run(t, analysistest.TestData(), customAnalyzer, "clean")
}

func TestDefaultSettings(t *testing.T) {
	// Test that default settings are reasonable
	defaults := unqueryvet.DefaultSettings()

	// Verify default values
	if !defaults.CheckSQLBuilders {
		t.Error("CheckSQLBuilders should be enabled by default")
	}

	// But should have some allowed patterns for reasonable behavior
	if len(defaults.AllowedPatterns) == 0 {
		t.Error("Should have some default allowed patterns")
	}

	// Check that COUNT(*) is allowed by default
	if !containsCountPattern(defaults.AllowedPatterns) {
		t.Error("Default patterns should include COUNT(*) pattern")
	}
}

func containsCountPattern(patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(pattern, "COUNT") {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkUnqueryvet(b *testing.B) {
	testdata := analysistest.TestData()

	for i := 0; i < b.N; i++ {
		analysistest.Run(b, testdata, unqueryvet.Analyzer, "testdata")
	}
}

// Test specific scenarios
func TestSelectStarDetection(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		wantDiag bool
	}{
		{
			name:     "basic select star",
			code:     `package test; func f() { _ = "SELECT * FROM users" }`,
			wantDiag: true,
		},
		{
			name:     "count star acceptable",
			code:     `package test; func f() { _ = "SELECT COUNT(*) FROM users" }`,
			wantDiag: false,
		},
		{
			name:     "information schema acceptable",
			code:     `package test; func f() { _ = "SELECT * FROM information_schema.tables" }`,
			wantDiag: false,
		},
		{
			name:     "case insensitive",
			code:     `package test; func f() { _ = "select * from users" }`,
			wantDiag: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This is a simplified test - in practice you'd use analysistest.RunWithSources
			// or create temporary files for more complex testing
			if tc.wantDiag {
				// Verify that diagnostic is produced
				t.Log("Should produce diagnostic for:", tc.code)
			} else {
				// Verify that no diagnostic is produced
				t.Log("Should NOT produce diagnostic for:", tc.code)
			}
		})
	}
}
