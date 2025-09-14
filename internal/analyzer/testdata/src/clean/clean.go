// Package clean contains minimal test cases for unqueryvet analyzer
package clean

// Simple test case
func test() {
	query := "SELECT * FROM users" // want "avoid SELECT \\* - explicitly specify needed columns for better performance, maintainability and stability"
	_ = query
}

// Acceptable pattern
func testAcceptable() {
	query := "SELECT COUNT(*) FROM users"
	_ = query
}
