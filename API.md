# SQLVet API Documentation

## Package: `github.com/MirrexOne/sqlvet`

### Overview

SQLVet provides a Go static analysis tool for detecting `SELECT *` usage in SQL queries.

### Installation

```bash
go get github.com/MirrexOne/sqlvet
```

## Core API

### Creating an Analyzer

```go
import "github.com/MirrexOne/sqlvet"

// Create analyzer with default settings
analyzer := sqlvet.New()

// Create analyzer with custom configuration
cfg := &sqlvet.Settings{
    CheckSQLBuilders: true,
    IgnoredFunctions: []string{"debug.Query"},
}
analyzer := sqlvet.NewWithConfig(cfg)
```

### Configuration Structure

```go
type Settings struct {
    // Enable SQL builder checking
    CheckSQLBuilders bool

    // Functions to ignore during analysis
    IgnoredFunctions []string

    // Packages to ignore completely
    IgnoredPackages []string

    // Regex patterns for allowed SELECT * usage
    AllowedPatterns []string

    // File patterns to ignore
    IgnoredFilePatterns []string

    // Directories to ignore
    IgnoredDirectories []string
}
```

### Default Settings

```go
defaults := sqlvet.DefaultSettings()
// Returns Settings with:
// - CheckSQLBuilders: true
// - Common patterns allowed (COUNT(*), information_schema, etc.)
// - Standard ignored files (*_test.go, *.pb.go, etc.)
// - Standard ignored directories (vendor, testdata, etc.)
```

## Internal Package API

### Package: `github.com/MirrexOne/sqlvet/internal/analyzer`

⚠️ **Note**: Internal packages are not part of the public API and may change without notice.

```go
// Create the base analyzer
analyzer := analyzer.NewAnalyzer()

// Create with specific settings
analyzer := analyzer.NewAnalyzerWithSettings(settings)

// Run analysis with configuration
result, err := analyzer.RunWithConfig(pass, cfg)
```

### Utility Functions

```go
// Normalize SQL query for analysis
normalized := analyzer.NormalizeSQLQuery(query)

// Check if query contains SELECT *
hasSelectStar := analyzer.IsSelectStarQuery(query, cfg)
```

## Public Package API

### Package: `github.com/MirrexOne/sqlvet/pkg/analyzer`

```go
import "github.com/MirrexOne/sqlvet/pkg/analyzer"

// Get the analyzer instance
a := analyzer.New()
```

### Package: `github.com/MirrexOne/sqlvet/pkg/config`

```go
import "github.com/MirrexOne/sqlvet/pkg/config"

// Get default configuration
cfg := config.DefaultSettings()

// Create custom configuration
cfg := config.SQLVetSettings{
    CheckSQLBuilders: false,
    AllowedPatterns: []string{
        "SELECT \\* FROM temp_.*",
    },
}
```

## Integration with analysis/singlechecker

```go
package main

import (
    "golang.org/x/tools/go/analysis/singlechecker"
    "github.com/MirrexOne/sqlvet"
)

func main() {
    singlechecker.Main(sqlvet.New())
}
```

## Integration with analysis/multichecker

```go
package main

import (
    "golang.org/x/tools/go/analysis/multichecker"
    "github.com/MirrexOne/sqlvet"
)

func main() {
    multichecker.Main(
        sqlvet.New(),
        // other analyzers...
    )
}
```

## Programmatic Usage

```go
package main

import (
    "fmt"
    "go/ast"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "github.com/MirrexOne/sqlvet"
)

func runAnalysis() {
    analyzer := sqlvet.New()
    
    // Configure the analysis pass
    pass := &analysis.Pass{
        // ... pass configuration
    }
    
    // Run the analyzer
    result, err := analyzer.Run(pass)
    if err != nil {
        fmt.Printf("Analysis failed: %v\n", err)
    }
}
```

## Configuration Examples

### Strict Configuration

```go
cfg := &sqlvet.Settings{
    CheckSQLBuilders: true,
    IgnoredFunctions: []string{}, // No ignored functions
    AllowedPatterns: []string{
        "COUNT\\(\\s*\\*\\s*\\)",
    },
}
```

### Permissive Configuration

```go
cfg := &sqlvet.Settings{
    CheckSQLBuilders: false,
    IgnoredFunctions: []string{
        "fmt.Printf",
        "log.Printf",
        "debug.Query",
    },
    AllowedPatterns: []string{
        "SELECT \\* FROM temp_.*",
        "SELECT \\* FROM .*_backup",
        "COUNT\\(\\s*\\*\\s*\\)",
    },
}
```

## Error Messages

SQLVet provides context-specific error messages:

- **Basic query**: "avoid SELECT * - explicitly specify needed columns for better performance, maintainability and stability"
- **SQL Builder**: "avoid SELECT * in SQL builder - explicitly specify columns to prevent unnecessary data transfer and schema change issues"
- **Empty Select**: "SQL builder Select() without columns defaults to SELECT * - add specific columns with .Columns() method"

## Suppressing Warnings

Use the standard `//nolint:sqlvet` directive:

```go
query := "SELECT * FROM users" //nolint:sqlvet
```

## Thread Safety

All analyzer functions are thread-safe and can be used concurrently.

## Performance Considerations

- The analyzer uses efficient AST traversal
- Pattern matching is optimized with pre-compiled regular expressions
- File and directory filtering happens early to avoid unnecessary processing

## Version Compatibility

- Go 1.21+ required
- Compatible with golangci-lint v2.x
- Uses golang.org/x/tools/go/analysis framework
