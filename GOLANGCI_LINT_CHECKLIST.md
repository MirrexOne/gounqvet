# SQLVet - golangci-lint Integration Checklist

## ✅ Requirements Status

### General Requirements
- [ ] **The CLA must be signed** - Will be signed when creating PR

### Pull Request Description
- ✅ **Link to linter repository**: https://github.com/MirrexOne/sqlvet
- ✅ **Short description**: SQLVet detects `SELECT *` usage in SQL queries and SQL builders, encouraging explicit column selection for better performance, maintainability, and API stability.

### Linter Requirements
- ✅ **Not a duplicate**: Unique linter for SQL query analysis focusing on SELECT * detection
- ✅ **Valid license**: MIT License (see LICENSE file) with proper author and year
- ✅ **Go version <= 1.22**: Now using `go 1.22` ✅ FIXED
- ✅ **Has CI and tests**: GitHub Actions CI configured, comprehensive test suite
- ✅ **Uses go/analysis**: Built on `golang.org/x/tools/go/analysis`
- ✅ **Valid tag**: `v1.0.1` (will push after final review)
- ✅ **No init()**: Verified - no init functions found
- ✅ **No panic()**: Verified - no panic calls
- ✅ **No log.Fatal(), os.Exit()**: Verified - none present
- ✅ **Does not modify AST**: Only analyzes, doesn't modify
- ⏳ **No false positives/negatives**: Team will help verify
- ⏳ **Has tests inside golangci-lint**: Will be added in PR

### Linter Tests Inside golangci-lint
- ✅ **Has std lib imports**: Uses `fmt`, `database/sql` in tests
- ✅ **Integration tests without configuration**: Prepared in `testdata/golangci/sqlvet_test.go`
- ✅ **Integration tests with configuration**: Prepared with various settings

### .golangci.next.reference.yml Updates
- ⏳ **File must be updated**: Will be done in PR
- ✅ **File .golangci.reference.yml NOT edited**: Will not touch this file
- ⏳ **Added to enable/disable lists**: Will add in alphabetical order
- ⏳ **Configuration section added**: Will add with non-default values and descriptions

### Other Requirements
- ✅ **Files have same name as linter**: `sqlvet.go`, `sqlvet` tests
- ✅ **.golangci.yml not edited**: Will not modify
- ⏳ **Alphabetical order**: Will maintain in `lintersdb/builder_linter.go` and `.golangci.next.reference.yml`
- ✅ **Load mode**: Uses `LoadModeTypesInfo`, requires `WithLoadForGoAnalysis()`
- ⏳ **Version in WithSince**: Will use next minor version (e.g., v1.66.0)
- ✅ **WithURL()**: Already set to https://github.com/MirrexOne/sqlvet

### Recommendations
- ⏳ **jsonschema/golangci.next.jsonschema.json updated**: Will update in PR
- ✅ **jsonschema/golangci.jsonschema.json NOT edited**: Will not touch
- ✅ **Has readme and linting**: Comprehensive README.md exists
- ⏳ **Published as binary**: Can add GitHub releases with binaries
- ✅ **Has .gitignore**: Proper .gitignore file exists
- ✅ **Tag not recreated**: Will create new tags only
- ✅ **Uses main as default branch**: Confirmed

## Files to be added/modified in golangci-lint PR:

### 1. pkg/golinters/sqlvet/sqlvet.go
```go
package sqlvet

import (
    "github.com/MirrexOne/sqlvet"
    "github.com/golangci/golangci-lint/v2/pkg/config"
    "github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.SqlvetSettings) *goanalysis.Linter {
    var a *analysis.Analyzer
    if settings != nil {
        a = sqlvet.NewWithConfig(&sqlvet.Config{
            CheckSQLBuilders:    settings.CheckSQLBuilders,
            IgnoredFunctions:    settings.IgnoredFunctions,
            IgnoredPackages:     settings.IgnoredPackages,
            AllowedPatterns:     settings.AllowedPatterns,
            IgnoredFilePatterns: settings.IgnoredFilePatterns,
            IgnoredDirectories:  settings.IgnoredDirectories,
        })
    } else {
        a = sqlvet.Analyzer
    }
    
    return goanalysis.NewLinter(
        a.Name,
        a.Doc,
        []*analysis.Analyzer{a},
        nil,
    ).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
```

### 2. pkg/golinters/sqlvet/testdata/sqlvet.go
```go
// Test file content (copy from testdata/golangci/sqlvet_test.go)
```

### 3. pkg/golinters/sqlvet/sqlvet_test.go
```go
package sqlvet

import (
    "testing"
    "github.com/golangci/golangci-lint/v2/test/testshared"
)

func TestSQLVet(t *testing.T) {
    testshared.NewLintRunner(t).Run(testshared.Args{
        Dir: "testdata",
        Control: []string{"sqlvet"},
    })
}
```

### 4. .golangci.next.reference.yml additions

In `enable` section (alphabetical order):
```yaml
- sqlvet
```

In `disable` section (alphabetical order):
```yaml
- sqlvet
```

In `linters-settings` section (alphabetical order after `spancheck`):
```yaml
sqlvet:
  # Check SQL builders like Squirrel (default: true)
  check-sql-builders: false
  
  # Functions to ignore (default: [])
  ignored-functions:
    - "fmt.Printf"
    - "log.Printf"
  
  # Packages to ignore (default: [])
  ignored-packages:
    - "testing"
  
  # Patterns allowed to use SELECT * (default: COUNT(*), MAX(*), MIN(*), information_schema, pg_catalog, sys)
  allowed-patterns:
    - "SELECT \\* FROM temp_.*"
  
  # File patterns to ignore (default: ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"])
  ignored-file-patterns:
    - "*_integration.go"
  
  # Directories to ignore (default: ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"])
  ignored-directories:
    - "examples"
```

### 5. pkg/lint/lintersdb/builder_linter.go modification

Add in alphabetical order (after spancheck, before staticcheck):
```go
linter.NewConfig(sqlvet.New(&cfg.Linters.Settings.Sqlvet)).
    WithSince("v1.66.0").
    WithLoadForGoAnalysis().
    WithURL("https://github.com/MirrexOne/sqlvet"),
```

### 6. jsonschema/golangci.next.jsonschema.json

Add in alphabetical order:
```json
"sqlvet": {
  "type": "object",
  "properties": {
    "check-sql-builders": {
      "type": "boolean",
      "default": true
    },
    "ignored-functions": {
      "type": "array",
      "items": {"type": "string"},
      "default": []
    },
    "ignored-packages": {
      "type": "array", 
      "items": {"type": "string"},
      "default": []
    },
    "allowed-patterns": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["COUNT\\(\\s*\\*\\s*\\)", "SELECT \\* FROM information_schema\\..*", "SELECT \\* FROM pg_catalog\\..*", "SELECT \\* FROM sys\\..*"]
    },
    "ignored-file-patterns": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"]
    },
    "ignored-directories": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"]
    }
  }
}
```

## Pre-PR Checklist

- [x] Fix Go version to 1.22
- [ ] Commit and push changes
- [ ] Create and push tag v1.0.1
- [ ] Fork golangci-lint if not already done
- [ ] Create feature branch in fork
- [ ] Add all required files
- [ ] Run tests locally
- [ ] Sign CLA
- [ ] Create PR with proper description

## PR Template

```markdown
## Description

This PR adds the `sqlvet` linter which detects `SELECT *` usage in SQL queries and SQL builders.

**Linter repository:** https://github.com/MirrexOne/sqlvet

**Short description:** SQLVet encourages explicit column selection in SQL queries by detecting `SELECT *` usage, improving performance, maintainability, and API stability.

## Checklist

- [ ] CLA signed
- [ ] Linter repository link provided
- [ ] Short description provided
- [ ] Tests added
- [ ] `.golangci.next.reference.yml` updated
- [ ] Files follow naming convention
- [ ] Alphabetical order maintained

## Testing

```bash
go test ./pkg/golinters/sqlvet/...
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=sqlvet ./pkg/golinters/sqlvet/testdata/
```
```
