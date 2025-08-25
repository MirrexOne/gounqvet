// Package clean contains minimal test cases for sqlvet analyzer  
package clean

// Simple test case
func test() {
	query := "SELECT * FROM users" // want "SELECT star usage detected"
	_ = query
}

// Acceptable pattern
func testAcceptable() {
	query := "SELECT COUNT(*) FROM users"
	_ = query
}