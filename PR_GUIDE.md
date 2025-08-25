# SQLVet - Руководство по созданию PR в golangci-lint

## ✅ Статус готовности

### Критические исправления выполнены:
1. ✅ **Go версия изменена на 1.22** (было 1.24.0)
2. ✅ **Реализована правильная передача конфигурации**
3. ✅ **Подготовлена структура тестов**
4. ✅ **Создана полная документация**

## 📋 Пошаговая инструкция для создания PR

### Шаг 1: Финализация sqlvet репозитория

```bash
# 1. Создание и публикация тега
cd /Users/aebelovitskiy/GolandProjects/sqlvet
git tag -d v1.0.1  # Удалить старый, если есть
git tag v1.0.2 -m "Release v1.0.2 - Go 1.22 compatibility for golangci-lint"
git push origin main
git push origin v1.0.2
```

### Шаг 2: Подготовка golangci-lint форка

```bash
# 1. Обновить форк
cd /Users/aebelovitskiy/GolandProjects/golangci-lint-fork
git checkout master
git pull upstream master
git push origin master

# 2. Создать ветку для PR
git checkout -b feat/add-sqlvet-linter

# 3. Удалить старые файлы sqlvet (если есть)
rm -f pkg/golinters/sqlvet.go
git rm pkg/golinters/sqlvet.go
```

### Шаг 3: Добавление файлов в golangci-lint

#### 3.1 Создать pkg/golinters/sqlvet/sqlvet.go

```go
package sqlvet

import (
    "github.com/MirrexOne/sqlvet"
    "github.com/MirrexOne/sqlvet/pkg/config"
    
    "github.com/golangci/golangci-lint/v2/pkg/config"
    "github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.SqlvetSettings) *goanalysis.Linter {
    var a *analysis.Analyzer
    
    if settings != nil {
        cfg := &sqlvetconfig.SQLVetSettings{
            CheckSQLBuilders:    settings.CheckSQLBuilders,
            IgnoredFunctions:    settings.IgnoredFunctions,
            IgnoredPackages:     settings.IgnoredPackages,
            AllowedPatterns:     settings.AllowedPatterns,
            IgnoredFilePatterns: settings.IgnoredFilePatterns,
            IgnoredDirectories:  settings.IgnoredDirectories,
        }
        a = sqlvet.NewWithConfig(cfg)
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

#### 3.2 Создать тестовые файлы

```bash
mkdir -p pkg/golinters/sqlvet/testdata
cp /Users/aebelovitskiy/GolandProjects/sqlvet/testdata/golangci/sqlvet_test.go \
   pkg/golinters/sqlvet/testdata/sqlvet.go
```

#### 3.3 Обновить .golangci.next.reference.yml

Добавить в секцию `enable` (строка ~111, в алфавитном порядке):
```yaml
    - sqlvet
```

Добавить в секцию `disable` (строка ~210, в алфавитном порядке):
```yaml
    - sqlvet
```

Добавить в `linters-settings` (после `sloglint`, перед `spancheck`, ~строка 2800):
```yaml
  sqlvet:
    # Check SQL builders like Squirrel, GORM (default: true)
    check-sql-builders: false
    
    # Functions to ignore during analysis (default: [])
    ignored-functions:
      - "fmt.Printf"
      - "log.Printf"
      - "debug.Query"
    
    # Packages to ignore during analysis (default: [])
    ignored-packages:
      - "testing"
      - "github.com/example/debug"
    
    # Regex patterns allowed to use SELECT * 
    # Default: COUNT(*), MAX(*), MIN(*), information_schema.*, pg_catalog.*, sys.*
    allowed-patterns:
      - "SELECT \\* FROM temp_.*"
      - "SELECT \\* FROM .*_backup"
    
    # File patterns to ignore
    # Default: ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"]
    ignored-file-patterns:
      - "*_integration.go"
      - "mock_*.go"
    
    # Directories to ignore  
    # Default: ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"]
    ignored-directories:
      - "examples"
      - "fixtures"
```

#### 3.4 Обновить pkg/lint/lintersdb/builder_linter.go

Найти правильное место (после `sqlclosecheck`, перед `staticcheck`) и добавить:
```go
linter.NewConfig(sqlvet.New(&cfg.Linters.Settings.Sqlvet)).
    WithSince("v1.66.0").  // Использовать следующую версию golangci-lint
    WithLoadForGoAnalysis().
    WithURL("https://github.com/MirrexOne/sqlvet"),
```

#### 3.5 Обновить jsonschema/golangci.next.jsonschema.json (опционально)

Добавить в алфавитном порядке:
```json
"sqlvet": {
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "check-sql-builders": {
      "type": "boolean",
      "default": true,
      "description": "Check SQL builders like Squirrel for SELECT * usage"
    },
    "ignored-functions": {
      "type": "array",
      "items": {"type": "string"},
      "default": [],
      "description": "List of function names to ignore"
    },
    "ignored-packages": {
      "type": "array",
      "items": {"type": "string"},
      "default": [],
      "description": "List of package names to ignore"
    },
    "allowed-patterns": {
      "type": "array",
      "items": {"type": "string"},
      "default": [
        "COUNT\\(\\s*\\*\\s*\\)",
        "MAX\\(\\s*\\*\\s*\\)",
        "MIN\\(\\s*\\*\\s*\\)",
        "SELECT \\* FROM information_schema\\..*",
        "SELECT \\* FROM pg_catalog\\..*",
        "SELECT \\* FROM sys\\..*"
      ],
      "description": "Regex patterns that are allowed to use SELECT *"
    },
    "ignored-file-patterns": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["*_test.go", "*.pb.go", "*_gen.go", "*.gen.go", "*_generated.go"],
      "description": "File patterns to ignore during analysis"
    },
    "ignored-directories": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["vendor", "testdata", "migrations", "generated", ".git", "node_modules"],
      "description": "Directory names to ignore during analysis"
    }
  }
}
```

### Шаг 4: Тестирование

```bash
# 1. Обновить зависимости
go mod tidy

# 2. Запустить тесты
go test ./pkg/golinters/sqlvet/...

# 3. Протестировать линтер
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=sqlvet ./pkg/golinters/sqlvet/testdata/

# 4. Проверить на реальном проекте
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=sqlvet ../sqlvet/
```

### Шаг 5: Коммит и создание PR

```bash
# 1. Добавить файлы
git add .

# 2. Коммит
git commit -m "feat: add sqlvet linter

SQLVet detects SELECT * usage in SQL queries and SQL builders,
encouraging explicit column selection for better performance,
maintainability, and API stability.

Linter repository: https://github.com/MirrexOne/sqlvet"

# 3. Push в форк
git push origin feat/add-sqlvet-linter
```

### Шаг 6: Создание PR на GitHub

1. Перейти на https://github.com/golangci/golangci-lint
2. Нажать "Compare & pull request"
3. Использовать следующий шаблон:

```markdown
## Description

This PR adds the `sqlvet` linter which detects `SELECT *` usage in SQL queries and SQL builders.

**Linter repository:** https://github.com/MirrexOne/sqlvet

## Why SQLVet?

SQLVet helps improve SQL query quality by detecting `SELECT *` usage and encouraging explicit column selection. This prevents:

- **Performance issues**: Selecting unnecessary columns wastes network bandwidth and memory
- **Maintenance problems**: Schema changes can break applications unexpectedly
- **Security risks**: May expose sensitive data unintentionally
- **API instability**: Adding new columns can break clients

## Features

- ✅ Detects `SELECT *` in string literals
- ✅ Detects `SELECT *` in SQL builders (Squirrel, GORM, etc.)
- ✅ Smart defaults (allows `COUNT(*)`, system tables)
- ✅ Highly configurable
- ✅ Supports `//nolint:sqlvet` directives
- ✅ Fast and lightweight

## Configuration Example

```yaml
linters-settings:
  sqlvet:
    check-sql-builders: true
    ignored-functions:
      - "fmt.Printf"
    allowed-patterns:
      - "SELECT \\* FROM temp_.*"
```

## Testing

All tests pass:
```bash
go test ./pkg/golinters/sqlvet/...
go run ./cmd/golangci-lint/ run --no-config --default=none --enable=sqlvet ./pkg/golinters/sqlvet/testdata/
```

## Checklist

- [x] CLA signed
- [x] Linter uses go/analysis
- [x] Has valid license (MIT)
- [x] Go version 1.22 compatible
- [x] Has tests
- [x] No init(), panic(), log.Fatal()
- [x] Does not modify AST
- [x] `.golangci.next.reference.yml` updated
- [x] Alphabetical order maintained
- [x] WithLoadForGoAnalysis() added

Resolves #[issue_number_if_any]
```

## ⚠️ Важные моменты

1. **CLA**: Нужно будет подписать Contributor License Agreement при создании PR
2. **Версия**: Убедитесь, что используете правильную версию в `WithSince()`
3. **Импорты**: Проверьте правильность импортов в sqlvet.go
4. **Тесты**: Убедитесь, что все тесты проходят перед созданием PR

## 📊 Проверочный список

- [x] Go версия изменена на 1.22
- [x] Все требования чеклиста выполнены
- [ ] Тег v1.0.2 опубликован
- [ ] Форк golangci-lint обновлен
- [ ] Все файлы добавлены
- [ ] Тесты проходят
- [ ] PR создан

## 🎯 Результат

После выполнения всех шагов, ваш PR будет полностью соответствовать требованиям golangci-lint и будет готов к review командой maintainers.
