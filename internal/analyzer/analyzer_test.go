package internal

import (
	"testing"

	cfg "github.com/MirrexOne/sqlvet/internal/config"
)

func TestNormalizeSQLQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple quoted string",
			input:    `"SELECT * FROM users"`,
			expected: "SELECT * FROM USERS",
		},
		{
			name:     "backtick string",
			input:    "`SELECT * FROM users`",
			expected: "SELECT * FROM USERS",
		},
		{
			name:     "string with escape sequences",
			input:    `"SELECT * FROM \"users\""`,
			expected: "SELECT * FROM \"USERS\"",
		},
		{
			name:     "multiline string with tabs and newlines",
			input:    `"SELECT *\n\tFROM users\n\tWHERE id = 1"`,
			expected: "SELECT * FROM USERS WHERE ID = 1",
		},
		{
			name:     "string with SQL comment",
			input:    `"SELECT * FROM users -- this is a comment"`,
			expected: "SELECT * FROM USERS",
		},
		{
			name:     "string with multiple spaces",
			input:    `"SELECT   *   FROM   users"`,
			expected: "SELECT * FROM USERS",
		},
		{
			name:     "complex string with all features",
			input:    `"SELECT *\n\tFROM \"users\"\n\t-- comment\n\tWHERE id = 1"`,
			expected: "SELECT * FROM \"USERS\" WHERE ID = 1",
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: "",
		},
		{
			name:     "string too short",
			input:    `"a"`,
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSQLQuery(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeSQLQuery(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsSelectStarQuery(t *testing.T) {
	config := cfg.NewConfig()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "simple SELECT * with FROM",
			input:    "SELECT * FROM users",
			expected: true,
		},
		{
			name:     "SELECT * with WHERE clause",
			input:    "SELECT * FROM users WHERE active = 1",
			expected: true,
		},
		{
			name:     "SELECT * with JOIN",
			input:    "SELECT * FROM users JOIN orders ON users.id = orders.user_id",
			expected: true,
		},
		{
			name:     "SELECT with explicit columns",
			input:    "SELECT id, name FROM users",
			expected: false,
		},
		{
			name:     "SELECT COUNT(*) - should be allowed by default",
			input:    "SELECT COUNT(*) FROM users",
			expected: false,
		},
		{
			name:     "SELECT * without SQL keywords",
			input:    "SELECT *",
			expected: false,
		},
		{
			name:     "INSERT statement",
			input:    "INSERT INTO users VALUES (1, 'John')",
			expected: false,
		},
		{
			name:     "UPDATE statement",
			input:    "UPDATE users SET name = 'Jane' WHERE id = 1",
			expected: false,
		},
		{
			name:     "complex SELECT * query",
			input:    "SELECT * FROM users WHERE active = 1 ORDER BY created_at DESC",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSelectStarQuery(tt.input, config)
			if result != tt.expected {
				t.Errorf("isSelectStarQuery(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConfigLoading(t *testing.T) {
	config := cfg.NewConfig()

	expectedDirs := []string{"vendor", ".git", "node_modules", "testdata"}
	for _, expectedDir := range expectedDirs {
		found := false
		for _, dir := range config.IgnoredDirectories {
			if dir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s to be in default ignored directories", expectedDir)
		}
	}

	expectedPatterns := []string{"*_test.go", "*.pb.go", "*_gen.go"}
	for _, expectedPattern := range expectedPatterns {
		found := false
		for _, pattern := range config.IgnoredFilePatterns {
			if pattern == expectedPattern {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s to be in default ignored file patterns", expectedPattern)
		}
	}
}

func TestAllowedPatterns(t *testing.T) {
	config := cfg.NewConfig()

	countQuery := "SELECT COUNT(*) FROM users"
	if isSelectStarQuery(countQuery, config) {
		t.Error("COUNT(*) should be allowed by default allowed patterns")
	}

	schemaQuery := "SELECT * FROM information_schema.tables"
	if isSelectStarQuery(schemaQuery, config) {
		t.Error("information_schema queries should be allowed by default")
	}

	normalQuery := "SELECT * FROM users WHERE active = 1"
	if !isSelectStarQuery(normalQuery, config) {
		t.Error("Normal SELECT * queries should not be allowed")
	}
}

func TestIsFileInDirectory(t *testing.T) {
	tests := []struct {
		name      string
		filepath  string
		directory string
		expected  bool
	}{
		{
			name:      "file in vendor directory",
			filepath:  "/project/vendor/github.com/pkg/module.go",
			directory: "vendor",
			expected:  true,
		},
		{
			name:      "file in nested vendor",
			filepath:  "/project/subproject/vendor/module.go",
			directory: "vendor",
			expected:  true,
		},
		{
			name:      "file not in vendor",
			filepath:  "/project/src/main.go",
			directory: "vendor",
			expected:  false,
		},
		{
			name:      "file in git directory",
			filepath:  "/project/.git/config",
			directory: ".git",
			expected:  true,
		},
		{
			name:      "vendor in filename but not directory",
			filepath:  "/project/src/vendor_utils.go",
			directory: "vendor",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFileInDirectory(tt.filepath, tt.directory)
			if result != tt.expected {
				t.Errorf("isFileInDirectory(%q, %q) = %v, want %v",
					tt.filepath, tt.directory, result, tt.expected)
			}
		})
	}
}
