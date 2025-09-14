package a

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Basic SELECT * detection in string literals
func basicSelectStar() {
	query := "SELECT * FROM users"           // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	anotherQuery := "select * from products" // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = query
	_ = anotherQuery
}

// Acceptable patterns that should NOT trigger warnings
func acceptablePatterns() {
	// COUNT(*) is acceptable
	countQuery := "SELECT COUNT(*) FROM users"
	_ = countQuery

	// Information schema queries are acceptable (system tables)
	sysQuery := "SELECT * FROM information_schema.tables"
	_ = sysQuery

	// Catalog queries are acceptable
	pgQuery := "SELECT * FROM pg_catalog.pg_tables"
	_ = pgQuery
}

// SQL in function calls
func sqlInFunctionCalls() {
	db, _ := sql.Open("postgres", "")

	// This should trigger warning
	rows, _ := db.Query("SELECT * FROM orders") // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = rows

	// This should also trigger warning
	stmt, _ := db.Prepare("SELECT * FROM customers") // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = stmt
}

// Multiline SQL queries
func multilineSQL() {
	query := ` // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
		SELECT * 
		FROM users
		WHERE active = true
	`
	_ = query
}

// Functions that will trigger warnings with default config
func defaultBehavior() {
	// These will trigger warnings with minimal default config
	fmt.Printf("SELECT * FROM debug_table")      // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	fmt.Sprintf("SELECT * FROM temp_%s", "data") // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	log.Printf("Executing: SELECT * FROM logs")  // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"

	// File operations
	_, _ = os.Open("file.sql")
}

// Nolint directives
func nolintDirectives() {
	// This should be ignored due to nolint comment
	query := "SELECT * FROM temp_table" //nolint:unqueryvet
	_ = query

	// This should be ignored due to general nolint
	anotherQuery := "SELECT * FROM backup" //nolint
	_ = anotherQuery
}
