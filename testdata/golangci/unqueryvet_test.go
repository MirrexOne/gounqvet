package testdata

import (
	"database/sql"
	"fmt"
)

// TestBasicSelectStar tests basic SELECT * detection
func TestBasicSelectStar() {
	// These should trigger warnings
	query1 := "SELECT * FROM users"    // want `SELECT star usage detected`
	query2 := "select * from products" // want `SELECT star usage detected`
	query3 := `SELECT * FROM orders`   // want `SELECT star usage detected`

	// These should NOT trigger warnings (allowed patterns)
	countQuery := "SELECT COUNT(*) FROM users"
	maxQuery := "SELECT MAX(*) FROM items"
	minQuery := "SELECT MIN(*) FROM scores"
	sysQuery := "SELECT * FROM information_schema.tables"
	pgQuery := "SELECT * FROM pg_catalog.pg_tables"

	// Using queries to avoid unused variable warnings
	fmt.Println(query1, query2, query3, countQuery, maxQuery, minQuery, sysQuery, pgQuery)
}

// TestSQLInFunctionCalls tests detection in function calls
func TestSQLInFunctionCalls() {
	db, _ := sql.Open("postgres", "")

	// Should trigger warnings
	rows, _ := db.Query("SELECT * FROM orders WHERE status = ?", "active") // want `SELECT star usage detected`
	stmt, _ := db.Prepare("SELECT * FROM customers")                       // want `SELECT star usage detected`

	// Should NOT trigger warnings
	goodRows, _ := db.Query("SELECT id, name FROM users")
	goodStmt, _ := db.Prepare("SELECT COUNT(*) FROM orders")

	defer rows.Close()
	defer stmt.Close()
	defer goodRows.Close()
	defer goodStmt.Close()
}

// TestMultilineSQL tests multiline SQL queries
func TestMultilineSQL() {
	// Should trigger warning
	query := `
		SELECT * 
		FROM users 
		WHERE active = true` // want `SELECT star usage detected`

	// Should NOT trigger warning (explicit columns)
	goodQuery := `
		SELECT id, name, email
		FROM users
		WHERE active = true`

	fmt.Println(query, goodQuery)
}

// TestNolintDirective tests nolint suppression
func TestNolintDirective() {
	// These should be ignored due to nolint comments
	query1 := "SELECT * FROM users" //nolint:unqueryvet

	//nolint:unqueryvet
	query2 := "SELECT * FROM products"

	// General nolint should also work
	query3 := "SELECT * FROM orders" //nolint

	fmt.Println(query1, query2, query3)
}

// TestNestedQueries tests nested SQL queries
func TestNestedQueries() {
	query := `
		SELECT u.name 
		FROM users u 
		WHERE u.id IN (
			SELECT * FROM user_permissions  -- want SELECT star usage detected
			WHERE role = 'admin'
		)`

	fmt.Println(query)
}

// TestStringInterpolation tests SELECT * in string formatting
func TestStringInterpolation() {
	table := "users"
	query := fmt.Sprintf("SELECT * FROM %s", table) // want `SELECT star usage detected`

	// Good practice - explicit columns
	goodQuery := fmt.Sprintf("SELECT id, name FROM %s", table)

	fmt.Println(query, goodQuery)
}
