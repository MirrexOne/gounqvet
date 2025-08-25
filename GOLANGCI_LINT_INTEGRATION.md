# Integration Guide for golangci-lint

This guide explains how to integrate Gounqvet into golangci-lint.

## Prerequisites

1. Fork golangci-lint repository (you already have it)
2. Gounqvet v1.0.0 is published and tagged
3. Local golangci-lint development environment

## Step-by-Step Integration

### 1. Navigate to your golangci-lint fork

```bash
cd ~/path/to/golangci-lint
git checkout master
git pull upstream master
git checkout -b add-gounqvet-linter
```

### 2. Add Gounqvet dependency

```bash
go get github.com/MirrexOne/gounqvet@v1.0.0
```

### 3. Create the linter wrapper

Create file `pkg/golinters/gounqvet/gounqvet.go`:

```go
package gounqvet

import (
	"github.com/MirrexOne/gounqvet"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"golang.org/x/tools/go/analysis"
)

func New(settings *config.GounqvetSettings) *goanalysis.Linter {
	cfg := gounqvet.DefaultSettings()
	
	if settings != nil {
		if settings.CheckSQLBuilders != nil {
			cfg.CheckSQLBuilders = *settings.CheckSQLBuilders
		}
		if settings.IgnoredFunctions != nil {
			cfg.IgnoredFunctions = settings.IgnoredFunctions
		}
		if settings.IgnoredPackages != nil {
			cfg.IgnoredPackages = settings.IgnoredPackages
		}
		if settings.AllowedPatterns != nil {
			cfg.AllowedPatterns = settings.AllowedPatterns
		}
		if settings.IgnoredFilePatterns != nil {
			cfg.IgnoredFilePatterns = settings.IgnoredFilePatterns
		}
		if settings.IgnoredDirectories != nil {
			cfg.IgnoredDirectories = settings.IgnoredDirectories
		}
	}

	analyzer := gounqvet.NewWithConfig(&cfg)
	
	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
```

### 4. Add configuration structure

Add to `pkg/config/linters_settings.go`:

```go
type GounqvetSettings struct {
	CheckSQLBuilders    *bool    `mapstructure:"check-sql-builders"`
	IgnoredFunctions    []string `mapstructure:"ignored-functions"`
	IgnoredPackages     []string `mapstructure:"ignored-packages"`
	AllowedPatterns     []string `mapstructure:"allowed-patterns"`
	IgnoredFilePatterns []string `mapstructure:"ignored-file-patterns"`
	IgnoredDirectories  []string `mapstructure:"ignored-directories"`
}
```

Add to the `LintersSettings` struct:

```go
type LintersSettings struct {
	// ... other linters
	Gounqvet GounqvetSettings `mapstructure:"gounqvet"`
	// ... other linters
}
```

### 5. Register the linter

Add to `pkg/lint/lintersdb/builder_linter.go`:

```go
import "github.com/golangci/golangci-lint/pkg/golinters/gounqvet"

// In the getLinterConfigs function, add:
{
	lc.WithSince("v1.62.0").
		WithPresets(linter.PresetSQL, linter.PresetStyle).
		WithURL("https://github.com/MirrexOne/gounqvet").
		WithLoadForGoAnalysis(),
	
	linter.NewConfig(golinters.NewGounqvet(&m.cfg.LintersSettings.Gounqvet)).
		WithEnabledByDefault(false),
},
```

### 6. Add to linter list

Add to `pkg/lint/lintersdb/linters_list.go`:

```go
const (
	// ... other linters
	GounqvetName = "gounqvet"
)
```

### 7. Update configuration schema

Add to `.golangci.reference.yml`:

```yaml
linters-settings:
  gounqvet:
    # Check SQL builders in addition to string literals
    # Default: true
    check-sql-builders: true
    
    # Functions to ignore during analysis
    # Default: []
    ignored-functions:
      - fmt.Printf
      - log.Printf
    
    # Packages to ignore completely
    # Default: []
    ignored-packages:
      - testing
    
    # Regex patterns for allowed SELECT * usage
    # Default includes: COUNT(*), MAX(*), MIN(*), information_schema, pg_catalog, sys
    allowed-patterns:
      - "SELECT \\* FROM temp_.*"
    
    # File patterns to ignore
    # Default: ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"]
    ignored-file-patterns:
      - "*_mock.go"
    
    # Directories to ignore
    # Default: ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"]
    ignored-directories:
      - "examples"
```

### 8. Add tests

Create `test/testdata/gounqvet.go`:

```go
//golangcitest:args -Egounqvet
package testdata

import "database/sql"

func testGounqvet() {
	var db *sql.DB
	
	// Should trigger gounqvet
	query := "SELECT * FROM users" // want "avoid SELECT \\* - explicitly specify needed columns"
	db.Query(query)
	
	// Should not trigger (COUNT(*) is allowed)
	countQuery := "SELECT COUNT(*) FROM users"
	db.Query(countQuery)
	
	// Should not trigger with nolint
	debugQuery := "SELECT * FROM debug" //nolint:gounqvet
	db.Query(debugQuery)
}
```

### 9. Update documentation

Add to `docs/linters.md`:

```markdown
## gounqvet

Gounqvet detects `SELECT *` usage in SQL queries and SQL builders, encouraging explicit column selection for better performance, maintainability, and API stability.

### Why avoid SELECT *?

- **Performance**: Selecting unnecessary columns wastes network bandwidth and memory
- **Maintainability**: Schema changes can break your application unexpectedly
- **Security**: May expose sensitive data that shouldn't be returned
- **API Stability**: Adding new columns can break clients that depend on column order

### Configuration

```yaml
linters-settings:
  gounqvet:
    check-sql-builders: true
    ignored-functions:
      - fmt.Printf
```
```

### 10. Run tests

```bash
# Build golangci-lint
make build

# Run tests
make test

# Test specifically with gounqvet
./golangci-lint run --enable-only gounqvet ./test/testdata/gounqvet.go
```

## Creating the Pull Request

### PR Title
```
feat: add gounqvet linter for detecting SELECT * in SQL queries
```

### PR Description Template
```markdown
## Description

This PR adds `gounqvet` - a linter that detects `SELECT *` usage in SQL queries and SQL builders, encouraging explicit column selection.

## What does gounqvet do?

- Detects `SELECT *` in SQL string literals
- Finds `SELECT *` usage in SQL builders (Squirrel, GORM, etc.)
- Supports configuration for allowed patterns and ignored functions
- Provides context-specific error messages

## Why is this useful?

Using `SELECT *` in production code can lead to:
- Performance issues (fetching unnecessary columns)
- Maintenance problems (schema changes breaking code)
- Security risks (exposing sensitive data)
- API instability (column order dependencies)

## Configuration

```yaml
linters-settings:
  gounqvet:
    check-sql-builders: true
    ignored-functions:
      - fmt.Printf
    allowed-patterns:
      - "SELECT \\* FROM information_schema\\..*"
```

## Testing

- [x] Added unit tests
- [x] Added integration tests
- [x] Tested with real projects
- [x] Documentation updated

## Related

- Linter repository: https://github.com/MirrexOne/gounqvet
- Issue: #[issue_number] (if applicable)

## Checklist

- [x] I have added tests for my changes
- [x] I have updated the documentation
- [x] I have updated `.golangci.reference.yml`
- [x] I have added the linter to the builders
- [x] All tests pass locally
```

## Common Issues and Solutions

### Issue: Import cycle
Solution: Use the internal package pattern

### Issue: Configuration not loading
Solution: Check mapstructure tags match YAML keys

### Issue: Linter not found
Solution: Ensure it's registered in builder_linter.go

## Final Checklist

- [ ] Gounqvet v1.0.0 is published on GitHub
- [ ] go.mod in golangci-lint includes gounqvet dependency
- [ ] Linter wrapper created
- [ ] Configuration structure added
- [ ] Linter registered in builder
- [ ] Reference configuration updated
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] PR created with detailed description
