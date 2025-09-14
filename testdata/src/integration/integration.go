package integration

import (
	"database/sql"
	"fmt"
)

// Test cases for golangci-lint integration tests

// Basic SELECT * detection
func basicSelectStar() {
	query := "SELECT * FROM users" // want "SELECT \\* usage detected"
	_ = query
}

// SQL with FROM keyword
func selectStarWithFrom() {
	db, _ := sql.Open("", "")
	rows, _ := db.Query("SELECT * FROM orders WHERE status = ?", "active") // want "SELECT \\* usage detected"
	_ = rows
}

// Acceptable COUNT(*) usage
func acceptableCount() {
	query := "SELECT COUNT(*) FROM users" // OK - should not trigger warning
	_ = query
}

// Information schema queries (should be allowed)
func informationSchemaQuery() {
	query := "SELECT * FROM information_schema.tables" // OK - should not trigger warning
	_ = query
}

// PostgreSQL catalog queries (should be allowed)
func pgCatalogQuery() {
	query := "SELECT * FROM pg_catalog.pg_tables" // OK - should not trigger warning
	_ = query
}

// Case insensitive detection
func caseInsensitive() {
	query := "select * from users where active = 1" // want "SELECT \\* usage detected"
	_ = query
}

// Multi-line query
func multiLineQuery() {
	query := `SELECT *
		FROM users 
		WHERE active = 1` // want "SELECT \\* usage detected"
	_ = query
}

// Ignored function
func ignoredFunction() {
	fmt.Printf("SELECT * FROM debug_table") // OK - Printf is ignored by default
}

// With nolint comment
func withNolintComment() {
	query := "SELECT * FROM users" //nolint:unqueryvet
	_ = query
}

// Complex query with JOIN
func complexQuery() {
	query := "SELECT * FROM users JOIN orders ON users.id = orders.user_id" // want "SELECT \\* usage detected"
	_ = query
}

// NOT a SQL query (no SQL keywords)
func notSQLQuery() {
	text := "SELECT * some random text" // OK - no SQL keywords like FROM, WHERE etc.
	_ = text
}

// SQL Builder tests
func sqlBuilderSelectStar() {
	// This would require actual SQL builder imports for real testing
	// but demonstrates the pattern expected
	// builder.Select("*").From("users") // would want "SELECT \\* usage detected"
}

// Raw string literals
func rawStringLiterals() {
	query := `SELECT * FROM products WHERE active = 1` // want "SELECT \\* usage detected"
	_ = query
}

// String concatenation (complex case)
func stringConcatenation() {
	base := "SELECT * FROM"
	table := "users"
	query := base + " " + table // This is harder to detect but should ideally trigger
	_ = query
}
