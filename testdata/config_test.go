// Package testdata contains configuration-specific test cases
package testdata

import (
	"fmt"
	"log"
)

// ignoredFunctions - functions that should be ignored based on config
func ignoredFunctions() {
	// These should NOT trigger warnings if properly configured
	fmt.Printf("SELECT * FROM debug_table")
	fmt.Sprintf("SELECT * FROM temp_%s", "data")
	log.Printf("Executing: SELECT * FROM logs")
}

// ignoredPackages - package-level ignoring
func ignoredPackages() {
	// If this package is in ignored-packages, nothing should trigger warnings
	query := "SELECT * FROM anywhere"
	_ = query
}

// allowedPatterns - patterns that should be allowed
func allowedPatterns() {
	// These should NOT trigger warnings due to allowed patterns
	tempQuery := "SELECT * FROM temp_table_123"  // Should be allowed if pattern matches
	backupQuery := "SELECT * FROM users_backup"  // Should be allowed if pattern matches
	systemQuery := "SELECT * FROM sys.databases" // Should be allowed if pattern matches

	_ = tempQuery
	_ = backupQuery
	_ = systemQuery
}

// sqlBuilders - SQL builder configuration tests
func sqlBuilders() {
	// Mock SQL builder patterns - should trigger warning if check-sql-builders: true
	builderQuery := "SELECT * FROM users" // want `SELECT \* usage detected`
	_ = builderQuery
}

// complexScenarios - complex scenarios for configuration testing
func complexScenarios() {
	// Multiple configurations combined
	query := `
		SELECT * FROM temp_debug_table
		WHERE created_at > NOW() - INTERVAL '1 hour'
	` // want `SELECT \* usage detected`
	_ = query
}
