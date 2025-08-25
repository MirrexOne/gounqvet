// Package config provides configuration for the SQLVet analyzer
package config

// SQLVetSettings represents configuration for the sqlvet linter
// This structure is designed to be compatible with golangci-lint integration
type SQLVetSettings struct {
	// CheckSQLBuilders enables checking SQL builders like Squirrel for SELECT * usage
	CheckSQLBuilders bool `mapstructure:"check-sql-builders" json:"check-sql-builders" yaml:"check-sql-builders"`

	// IgnoredFunctions is a list of function names to ignore during analysis
	// Example: ["fmt.Printf", "log.Debug"]
	IgnoredFunctions []string `mapstructure:"ignored-functions" json:"ignored-functions" yaml:"ignored-functions"`

	// IgnoredPackages is a list of package names to ignore during analysis
	// Example: ["testing", "debug"]
	IgnoredPackages []string `mapstructure:"ignored-packages" json:"ignored-packages" yaml:"ignored-packages"`

	// AllowedPatterns is a list of regex patterns that are allowed to use SELECT *
	// Example: ["SELECT \\* FROM temp_.*", "SELECT \\* FROM .*_backup"]
	AllowedPatterns []string `mapstructure:"allowed-patterns" json:"allowed-patterns" yaml:"allowed-patterns"`

	// IgnoredFilePatterns is a list of file patterns to ignore during analysis
	// Example: ["*_test.go", "*.gen.go"]
	IgnoredFilePatterns []string `mapstructure:"ignored-file-patterns" json:"ignored-file-patterns" yaml:"ignored-file-patterns"`

	// IgnoredDirectories is a list of directory names to ignore during analysis
	// Example: ["vendor", "testdata", "migrations"]
	IgnoredDirectories []string `mapstructure:"ignored-directories" json:"ignored-directories" yaml:"ignored-directories"`
}

// DefaultSettings returns the default configuration for sqlvet
func DefaultSettings() SQLVetSettings {
	return SQLVetSettings{
		CheckSQLBuilders: true,
		IgnoredFunctions: []string{},
		IgnoredPackages:  []string{},
		AllowedPatterns: []string{
			`(?i)COUNT\(\s*\*\s*\)`,
			`(?i)MAX\(\s*\*\s*\)`,
			`(?i)MIN\(\s*\*\s*\)`,
			`(?i)SELECT \* FROM information_schema\..*`,
			`(?i)SELECT \* FROM pg_catalog\..*`,
			`(?i)SELECT \* FROM sys\..*`,
		},
		IgnoredFilePatterns: []string{
			"*_test.go",
			"*.pb.go",
			"*_gen.go",
			"*.gen.go",
			"*_generated.go",
		},
		IgnoredDirectories: []string{
			"vendor",
			"testdata",
			"migrations",
			"generated",
			".git",
			"node_modules",
		},
	}
}