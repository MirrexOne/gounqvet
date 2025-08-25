package examples

import (
	"database/sql"
)

// This file demonstrates Gounqvet warnings and how to fix them
// NOTE: These are examples only - they don't actually execute

// ExampleBadCode shows patterns that Gounqvet will warn about
func ExampleBadCode() {
	var db *sql.DB

	// BAD: Direct SELECT * in string literal
	// gounqvet will warn about this
	query1 := "SELECT * FROM users WHERE active = true"
	_ = query1
	_ = db // db would be used in real code

	// BAD: SELECT * in function call
	// gounqvet will warn about this
	// In real code: rows, err := db.Query("SELECT * FROM products")
	_ = "SELECT * FROM products"

	// BAD: Multiline SELECT *
	// gounqvet will warn about this
	multilineQuery := `
		SELECT *
		FROM orders
		WHERE status = 'pending'
	`
	_ = multilineQuery
}

// ExampleGoodCode shows patterns that Gounqvet approves
func ExampleGoodCode() {
	var db *sql.DB
	_ = db // db would be used in real code

	// GOOD: Explicit column selection
	query1 := "SELECT id, name, email FROM users WHERE active = true"
	_ = query1

	// GOOD: Specific columns in function call
	// In real code: rows, err := db.Query("SELECT id, name, price FROM products")
	_ = "SELECT id, name, price FROM products"

	// GOOD: Multiline with explicit columns
	multilineQuery := `
		SELECT id, customer_id, total, status
		FROM orders
		WHERE status = 'pending'
	`
	_ = multilineQuery

	// GOOD: COUNT(*) is allowed by default
	countQuery := "SELECT COUNT(*) FROM users"
	_ = countQuery

	// GOOD: Information schema queries are allowed
	schemaQuery := "SELECT * FROM information_schema.tables WHERE table_name = 'users'"
	_ = schemaQuery
}

// ExampleSuppression shows how to suppress Gounqvet warnings
func ExampleSuppression() {
	var db *sql.DB
	_ = db // db would be used in real code

	// Using nolint directive to suppress warning
	debugQuery := "SELECT * FROM debug_logs" //nolint:gounqvet
	_ = debugQuery

	// Alternative: use for temporary debugging
	// Remove before committing
	tempQuery := "SELECT * FROM temp_table" //nolint:gounqvet // temporary for debugging
	_ = tempQuery
}

// ExampleSQLBuilders shows SQL builder patterns
func ExampleSQLBuilders() {
	// Pseudo-code examples (requires actual SQL builder libraries)

	// BAD: SELECT * in SQL builder
	// query := squirrel.Select("*").From("users")

	// BAD: Empty Select defaults to SELECT *
	// query := squirrel.Select().From("users")

	// GOOD: Explicit columns in SQL builder
	// query := squirrel.Select("id", "name", "email").From("users")

	// GOOD: Using Columns method
	// query := squirrel.Select().Columns("id", "name").From("users")
}
