// Package testdata contains test cases for sqlvet analyzer
package testdata

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Basic SELECT * detection in string literals
func basicSelectStar() {
	query := "SELECT * FROM users"           // want `SELECT \* usage detected`
	anotherQuery := "select * from products" // want `SELECT \* usage detected`
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
	rows, _ := db.Query("SELECT * FROM orders") // want `SELECT \* usage detected`
	_ = rows

	// This should also trigger warning
	stmt, _ := db.Prepare("SELECT * FROM customers") // want `SELECT \* usage detected`
	_ = stmt
}

// Multiline SQL queries
func multilineSQL() {
	query := `
		SELECT * 
		FROM users 
		WHERE active = true` // want `SELECT \* usage detected`
	_ = query

	// With comments - should still detect
	queryWithComments := `
		-- Get all user data
		SELECT * FROM users  -- want SELECT \* usage detected
		WHERE id > 100
	`
	_ = queryWithComments
}

// SQL Builders (Squirrel examples)
func sqlBuilders() {
	// Direct SELECT * in SQL builder - should trigger warning
	query1 := squirrel.Select("*").From("users") // want `SELECT \* usage detected`
	_ = query1

	// Empty Select() followed by explicit columns - should be OK
	query2 := squirrel.Select().Columns("name", "email").From("users")
	_ = query2

	// Chained Select().Columns("*") - should trigger warning
	query3 := squirrel.Select().Columns("*").From("users") // want `SELECT \* usage detected`
	_ = query3
}

// Edge cases and complex scenarios
func edgeCases() {
	// String interpolation with SELECT *
	table := "users"
	query := fmt.Sprintf("SELECT * FROM %s", table) // want `SELECT \* usage detected`
	_ = query

	// SELECT * in stored procedure calls
	procCall := "CALL get_user_data('SELECT * FROM users')" // want `SELECT \* usage detected`
	_ = procCall

	// Nested queries
	nestedQuery := `
		SELECT u.name 
		FROM users u 
		WHERE u.id IN (
			SELECT * FROM user_permissions  -- want SELECT \* usage detected
			WHERE role = 'admin'
		)
	`
	_ = nestedQuery
}

// Ignored patterns with nolint
func ignoredPatterns() {
	// This should be ignored due to nolint comment
	query1 := "SELECT * FROM users" //nolint:sqlvet
	_ = query1

	// This should also be ignored
	//nolint:sqlvet
	query2 := "SELECT * FROM products"
	_ = query2

	// General nolint should also work
	query3 := "SELECT * FROM orders" //nolint
	_ = query3
}

// Function arguments and variables
func functionArguments() {
	executeQuery("SELECT * FROM users") // want `SELECT \* usage detected`
	
	// Variable assignment
	var userQuery = "SELECT * FROM users" // want `SELECT \* usage detected`
	_ = userQuery
}

func executeQuery(query string) {
	fmt.Println("Executing:", query)
}

// SQL builders with variables
func sqlBuildersWithVariables() {
	// Variable holding empty Select() - should trigger warning if no columns added
	query := squirrel.Select() // This will be flagged if no Columns() call follows
	_ = query                  // want `SELECT \* usage detected`

	// Variable with proper columns - should be OK
	goodQuery := squirrel.Select()
	goodQuery = goodQuery.Columns("id", "name").From("users")
	_ = goodQuery
}

// Testing ignored functions
func ignoredFunctions() {
	// fmt.Printf should be ignored by default
	fmt.Printf("Query: SELECT * FROM debug_table")
	
	// Custom debugging function (if configured to ignore)
	debugQuery("SELECT * FROM temp_data")
}

func debugQuery(query string) {
	fmt.Println("DEBUG:", query)
}