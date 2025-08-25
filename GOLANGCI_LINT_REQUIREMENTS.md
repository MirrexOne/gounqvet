# Полные требования для успешного PR в golangci-lint

## 🏷️ Tag vs Release - Что нужно для golangci-lint

### Для golangci-lint ОБЯЗАТЕЛЬНО нужен:
1. **Git Tag** - обязательно ✅
2. **GitHub Release** - рекомендуется, но не обязательно ⚠️
3. **Опубликованный модуль на pkg.go.dev** - обязательно ✅

### Почему Tag обязателен:
- golangci-lint использует `go get` для загрузки линтера
- Go modules требуют семантические версии (v1.0.0, v1.1.0 и т.д.)
- Без тега модуль будет загружаться как `v0.0.0-{timestamp}-{commit}`

## ✅ Полный чеклист для успешного PR

### 1. Репозиторий sqlvet - ОБЯЗАТЕЛЬНЫЕ требования

#### A. Код и структура ✅
```bash
✅ Использует go/analysis API
✅ Go версия 1.22 (не выше!)
✅ Нет panic(), log.Fatal(), os.Exit()
✅ Нет init() функций
✅ MIT лицензия с корректным годом и автором
✅ Есть README.md с описанием
✅ Есть .gitignore
```

#### B. Git Tag - КРИТИЧЕСКИ ВАЖНО ⚠️
```bash
# Проверить текущие теги
git tag -l

# Создать тег (если еще не создан)
git tag v1.1.0 -m "Release v1.1.0 - golangci-lint ready"

# Отправить тег на GitHub - ОБЯЗАТЕЛЬНО!
git push origin v1.1.0
```

#### C. Публикация модуля Go ⚠️
После push тега, Go автоматически индексирует модуль, НО можно ускорить:
```bash
# Форсировать индексацию на proxy.golang.org
curl https://proxy.golang.org/github.com/MirrexOne/sqlvet/@v/v1.1.0.info

# Проверить доступность
go get github.com/MirrexOne/sqlvet@v1.1.0
```

#### D. GitHub Release - РЕКОМЕНДУЕТСЯ (но не обязательно)
```bash
# Через GitHub CLI
gh release create v1.1.0 \
  --title "v1.1.0 - Improved error messages" \
  --notes "## Changes
- Informative error messages with context
- Full golangci-lint compatibility
- Go 1.22 support
- Enhanced SQL builder detection"

# Или через веб-интерфейс GitHub:
# 1. Перейти на https://github.com/MirrexOne/sqlvet/releases
# 2. Нажать "Draft a new release"
# 3. Выбрать тег v1.1.0
# 4. Добавить описание
```

### 2. PR в golangci-lint - ОБЯЗАТЕЛЬНЫЕ файлы

#### A. Проверка доступности линтера ⚠️
```bash
# КРИТИЧЕСКИ ВАЖНО - golangci-lint должен мочь загрузить линтер!
go get -u github.com/MirrexOne/sqlvet@v1.1.0

# Если не работает, значит тег не опубликован
```

#### B. Файлы для изменения в golangci-lint:

##### 1. `pkg/golinters/sqlvet/sqlvet.go` ✅
```go
package sqlvet

import (
    "github.com/MirrexOne/sqlvet"
    sqlvetconfig "github.com/MirrexOne/sqlvet/pkg/config"
    
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

##### 2. `pkg/lint/lintersdb/builder_linter.go` ✅
Добавить в алфавитном порядке (после sqlclosecheck):
```go
linter.NewConfig(sqlvet.New(&cfg.Linters.Settings.Sqlvet)).
    WithSince("v1.66.0"). // Следующая версия golangci-lint!
    WithLoadForGoAnalysis().
    WithURL("https://github.com/MirrexOne/sqlvet"),
```

##### 3. `.golangci.next.reference.yml` ✅
```yaml
# В enable (строка ~110)
- sqlvet

# В disable (строка ~210)  
- sqlvet

# В linters-settings (после sloglint, ~2800)
sqlvet:
  # Check SQL builders (default: true)
  check-sql-builders: false
  # ... остальные настройки
```

##### 4. `pkg/golinters/sqlvet/testdata/sqlvet.go` ✅
Минимум один тестовый файл с примерами

##### 5. `go.mod` в golangci-lint ✅
Автоматически обновится при `go mod tidy`

### 3. Процесс создания PR

#### Шаг 1: Убедиться, что линтер доступен
```bash
# В любой директории
go get github.com/MirrexOne/sqlvet@v1.1.0
# Должно скачаться без ошибок!
```

#### Шаг 2: Подготовка форка golangci-lint
```bash
cd /Users/aebelovitskiy/GolandProjects/golangci-lint-fork
git checkout master
git pull upstream master
git checkout -b feat/add-sqlvet-linter
```

#### Шаг 3: Добавление файлов
```bash
# Создать директорию
mkdir -p pkg/golinters/sqlvet/testdata

# Добавить файлы (см. выше)
# ...

# Обновить зависимости
go get github.com/MirrexOne/sqlvet@v1.1.0
go mod tidy
```

#### Шаг 4: Тестирование
```bash
# Собрать golangci-lint
go build ./cmd/golangci-lint/

# Протестировать
./golangci-lint run --no-config --default=none --enable=sqlvet ./...
```

#### Шаг 5: Коммит и PR
```bash
git add .
git commit -m "feat: add sqlvet linter"
git push origin feat/add-sqlvet-linter
```

## ⚠️ Частые ошибки и их решения

### Ошибка 1: "cannot find module"
**Причина:** Тег не опубликован на GitHub
**Решение:** `git push origin v1.1.0`

### Ошибка 2: "unknown revision v1.1.0"
**Причина:** Go proxy еще не синхронизировался
**Решение:** Подождать 5-10 минут или форсировать через curl

### Ошибка 3: "invalid version"
**Причина:** Неправильный формат тега (должен быть vX.Y.Z)
**Решение:** Создать правильный тег v1.1.0

### Ошибка 4: Tests fail in golangci-lint
**Причина:** Тестовые сообщения не совпадают
**Решение:** Обновить тесты с правильными сообщениями об ошибках

## 📋 Финальный чеклист перед PR

- [ ] **sqlvet репозиторий:**
  - [ ] Код на main ветке закоммичен
  - [ ] Тег v1.1.0 создан локально
  - [ ] Тег отправлен на GitHub (`git push origin v1.1.0`)
  - [ ] Модуль доступен через `go get`
  - [ ] (Опционально) GitHub Release создан

- [ ] **golangci-lint PR:**
  - [ ] Форк обновлен с upstream
  - [ ] Новая ветка создана
  - [ ] Все файлы добавлены
  - [ ] `go mod tidy` выполнен
  - [ ] Тесты проходят
  - [ ] PR создан с правильным описанием

## 🎯 Команды для быстрого выполнения

```bash
# 1. В sqlvet - публикация
cd /Users/aebelovitskiy/GolandProjects/sqlvet
git push origin main
git push origin v1.1.0

# 2. Проверка доступности (подождать 2-3 минуты)
go get github.com/MirrexOne/sqlvet@v1.1.0

# 3. Создание GitHub Release (опционально)
gh release create v1.1.0 --title "v1.1.0" --notes "golangci-lint ready"

# 4. В golangci-lint - создание PR
cd ../golangci-lint-fork
git checkout -b feat/add-sqlvet-linter
# ... добавить файлы ...
go get github.com/MirrexOne/sqlvet@v1.1.0
go mod tidy
git add .
git commit -m "feat: add sqlvet linter"
git push origin feat/add-sqlvet-linter
```

## ✅ Готовность: 95%

Осталось только:
1. Опубликовать тег на GitHub
2. Дождаться индексации Go proxy (2-5 минут)
3. Создать PR в golangci-lint
