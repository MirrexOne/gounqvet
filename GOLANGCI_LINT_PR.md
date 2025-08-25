# SQLVet - PR Checklist for golangci-lint Integration

## PR Description

**Linter Repository**: https://github.com/MirrexOne/sqlvet

**Description**: SQLVet is a Go static analysis tool that detects `SELECT *` usage in SQL queries and SQL builders, encouraging explicit column selection for better performance, maintainability, and API stability.

## Completed Requirements

### Linter Requirements
- ✅ **Not a duplicate** - Unique linter for SQL query analysis focusing on SELECT * detection
- ✅ **Valid license** - MIT License (see LICENSE file)
- ✅ **Go version** - Using Go 1.24.0 (compatible with requirement)
- ✅ **CI and tests** - Has GitHub Actions CI and comprehensive tests
- ✅ **Uses go/analysis** - Built on golang.org/x/tools/go/analysis API
- ✅ **Valid tag** - v1.0.1
- ✅ **No init()** - Verified, no init functions
- ✅ **No panic()** - Verified, no panic calls
- ✅ **No log.Fatal()** or os.Exit()** - Verified, none present
- ✅ **Does not modify AST** - Only analyzes, doesn't modify
- ✅ **Has tests** - Comprehensive test suite included

### Linter Tests Inside golangci-lint
- ✅ **Has std lib imports** - Uses fmt, database/sql in tests
- ✅ **Integration tests without configuration** - Default configuration tests included
- ✅ **Integration tests with configuration** - Tests with various settings included

### .golangci.next.reference.yml Updates Needed
```yaml
# Add to linters list (alphabetical order)
enable:
  - sqlvet

disable:
  - sqlvet

# Add configuration section (alphabetical order)
linters-settings:
  sqlvet:
    # Enable checking SQL builders like Squirrel (default: true)
    check-sql-builders: true
    
    # Functions to ignore during analysis (default: [])
    ignored-functions:
      - "fmt.Printf"
      - "log.Printf"
    
    # Packages to ignore during analysis (default: [])
    ignored-packages:
      - "testing"
    
    # Regex patterns that are allowed to use SELECT * (default patterns included)
    allowed-patterns:
      - "COUNT\\(\\s*\\*\\s*\\)"
      - "SELECT \\* FROM information_schema\\..*"
    
    # File patterns to ignore (default: ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"])
    ignored-file-patterns:
      - "*_test.go"
    
    # Directories to ignore (default: ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"])
    ignored-directories:
      - "vendor"
```

### Other Requirements
- ✅ **Files have same name as linter** - sqlvet.go, sqlvet tests
- ✅ **Not added to .golangci.yml** - Only .golangci.next.reference.yml
- ✅ **Alphabetical order** - Will be maintained in PR
- ✅ **Load mode** - Uses WithLoadForGoAnalysis() for TypesInfo
- ✅ **Version** - Set to v1.65.0 (already in builder_linter.go)
- ✅ **URL** - https://github.com/MirrexOne/sqlvet

## Key Features

1. **SELECT * Detection**: Identifies problematic `SELECT *` usage in SQL queries
2. **SQL Builder Support**: Works with Squirrel, GORM, and other SQL builders
3. **Configurable**: Extensive configuration options for customization
4. **Smart Defaults**: Allows COUNT(*), system tables, etc. by default
5. **nolint Support**: Standard //nolint:sqlvet suppression
6. **Fast & Lightweight**: Efficient AST traversal with minimal overhead

## Integration Status

The linter is already partially integrated in golangci-lint v1.65.0 but needs the following updates:

1. ✅ Fixed configuration passing to analyzer
2. ✅ Implemented NewAnalyzerWithSettings properly
3. ✅ Added comprehensive tests
4. ⏳ Need to update .golangci.next.reference.yml
5. ⏳ Need to update pkg/golinters/sqlvet.go to use correct imports

## Testing

Run tests:
```bash
# In sqlvet directory
go test ./...

# In golangci-lint with sqlvet
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=sqlvet ./pkg/golinters/sqlvet/testdata/
```

## Benefits

- **Performance**: Prevents selecting unnecessary columns
- **Maintainability**: Makes schema changes more predictable
- **Security**: Avoids exposing sensitive data unintentionally  
- **API Stability**: Prevents breaking changes from new columns

## Notes

This linter focuses specifically on SQL query quality and complements existing linters by addressing a common performance and maintainability issue in Go applications that work with databases.
