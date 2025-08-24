package internal

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config определяет минимальную конфигурацию SQLVet с разумными умолчаниями
type Config struct {
	// Основные настройки для фильтрации файлов и директорий
	IgnoredDirectories  []string `yaml:"ignore-dirs,omitempty"`
	IgnoredFilePatterns []string `yaml:"ignore-file-patterns,omitempty"`

	// Фильтрация по пакетам и функциям
	IgnoredPackages  []string `yaml:"ignore-packages,omitempty"`
	IgnoredFunctions []string `yaml:"ignore-functions,omitempty"`

	// Паттерны для более гибкого игнорирования
	IgnoredFunctionPatterns []string `yaml:"ignore-function-patterns,omitempty"`

	// Разрешенные SQL паттерны (whitelist для исключений)
	AllowedPatterns []string `yaml:"allowed-patterns,omitempty"`

	// Настройки анализа SQL-билдеров
	CheckSQLBuilders bool     `yaml:"check-builders,omitempty"`
	IgnoredBuilders  []string `yaml:"ignored-builders,omitempty"`
}

// GolangCIConfig представляет структуру .golangci.yaml для интеграции
type GolangCIConfig struct {
	LintersSettings struct {
		Sqlvet Config `yaml:"sqlvet,omitempty"`
	} `yaml:"linters-settings,omitempty"`
}

// NewConfig создает конфигурацию с минимальными разумными умолчаниями
// Философия: все должно работать из коробки без настройки
func NewConfig() *Config {
	cfg := &Config{
		// Стандартные директории, которые обычно нужно игнорировать
		IgnoredDirectories: []string{
			"vendor",       // Зависимости
			".git",         // Git репозиторий
			"node_modules", // Node.js зависимости
			"testdata",     // Тестовые данные
		},

		// Стандартные паттерны файлов для игнорирования
		IgnoredFilePatterns: []string{
			"*_test.go",      // Тестовые файлы
			"*.pb.go",        // Protobuf сгенерированные файлы
			"*_gen.go",       // Сгенерированные файлы
			"*_generated.go", // Альтернативный паттерн для сгенерированных файлов
			"mock_*.go",      // Mock файлы
			"*_mock.go",      // Альтернативный паттерн для mock файлов
		},

		// Пакеты, которые обычно содержат отладочный или служебный код
		IgnoredPackages: []string{
			"main", // Основной пакет часто содержит служебный код
			"test", // Тестовые пакеты
			"mock", // Mock пакеты
		},

		// Функции, которые часто используются для отладки или логирования
		IgnoredFunctions: []string{
			"fmt.Printf",
			"fmt.Sprintf",
			"log.Printf",
			"testing.T.Log",
			"testing.T.Logf",
		},

		// Паттерны функций для игнорирования (используют filepath.Match синтаксис)
		IgnoredFunctionPatterns: []string{
			"Test*",      // Тестовые функции
			"Benchmark*", // Benchmark функции
			"Example*",   // Example функции
			"*_test",     // Функции, заканчивающиеся на _test
			"Debug*",     // Отладочные функции
		},

		// Разрешенные SQL паттерны - исключения из правил
		// Эти паттерны НЕ будут вызывать предупреждений даже при наличии SELECT *
		AllowedPatterns: []string{
			`SELECT COUNT\(\*\)`, // COUNT(*) обычно приемлем
			`SELECT MAX\(\*\)`,   // MAX(*) тоже может быть нормальным
			`SELECT MIN\(\*\)`,   // MIN(*) аналогично
			// Системные запросы к information_schema часто используют SELECT *
			`information_schema`,
			`INFORMATION_SCHEMA`,
		},

		// По умолчанию проверяем SQL-билдеры - это основная функциональность
		CheckSQLBuilders: true,

		// Игнорируемые билдеры - здесь можно указать ORM или специфические билдеры
		// По умолчанию пустой список означает, что проверяем все
		IgnoredBuilders: []string{
			// Можно добавить: "gorm", "xorm", "ent" если нужно игнорировать ORM
		},
	}

	// Попытка загрузить дополнительную конфигурацию из .golangci.yaml
	if golangciPath := findGolangCIConfig(); golangciPath != "" {
		if loaded := loadConfigFromGolangCI(golangciPath); loaded != nil {
			mergeConfigs(cfg, loaded)
		}
	}

	return cfg
}

// findGolangCIConfig ищет файл конфигурации golangci-lint в текущей директории и выше
func findGolangCIConfig() string {
	possibleNames := []string{
		".golangci.yaml",
		".golangci.yml",
		"golangci.yaml",
		"golangci.yml",
	}

	dir, _ := os.Getwd()

	// Поднимаемся по директориям вверх до корня
	for dir != "/" && dir != "." && dir != "" {
		for _, name := range possibleNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Достигли корня
		}
		dir = parent
	}

	return ""
}

// loadConfigFromGolangCI загружает конфигурацию из файла .golangci.yaml
func loadConfigFromGolangCI(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var golangCI GolangCIConfig
	if err := yaml.Unmarshal(data, &golangCI); err != nil {
		return nil
	}

	return &golangCI.LintersSettings.Sqlvet
}

// mergeConfigs объединяет базовую конфигурацию с загруженной из файла
// Приоритет отдается настройкам из файла
func mergeConfigs(base, loaded *Config) {
	if loaded == nil {
		return
	}

	// Объединяем списки, сохраняя уникальность элементов
	if len(loaded.IgnoredDirectories) > 0 {
		base.IgnoredDirectories = mergeLists(base.IgnoredDirectories, loaded.IgnoredDirectories)
	}

	if len(loaded.IgnoredFilePatterns) > 0 {
		base.IgnoredFilePatterns = mergeLists(base.IgnoredFilePatterns, loaded.IgnoredFilePatterns)
	}

	if len(loaded.IgnoredPackages) > 0 {
		base.IgnoredPackages = mergeLists(base.IgnoredPackages, loaded.IgnoredPackages)
	}

	if len(loaded.IgnoredFunctions) > 0 {
		base.IgnoredFunctions = mergeLists(base.IgnoredFunctions, loaded.IgnoredFunctions)
	}

	if len(loaded.IgnoredFunctionPatterns) > 0 {
		base.IgnoredFunctionPatterns = mergeLists(base.IgnoredFunctionPatterns, loaded.IgnoredFunctionPatterns)
	}

	if len(loaded.AllowedPatterns) > 0 {
		base.AllowedPatterns = mergeLists(base.AllowedPatterns, loaded.AllowedPatterns)
	}

	if len(loaded.IgnoredBuilders) > 0 {
		base.IgnoredBuilders = mergeLists(base.IgnoredBuilders, loaded.IgnoredBuilders)
	}

	// Булевые значения просто перезаписываем
	base.CheckSQLBuilders = loaded.CheckSQLBuilders
}

// mergeLists объединяет два списка строк, удаляя дубликаты
func mergeLists(base, additional []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(base)+len(additional))

	// Добавляем элементы из базового списка
	for _, item := range base {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	// Добавляем элементы из дополнительного списка
	for _, item := range additional {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Вспомогательные методы для проверки правил игнорирования

// IsPackageIgnored проверяет, игнорируется ли указанный пакет
func (c *Config) IsPackageIgnored(pkgName string) bool {
	for _, ignored := range c.IgnoredPackages {
		if strings.EqualFold(pkgName, ignored) {
			return true
		}
	}
	return false
}

// IsFunctionIgnored проверяет, игнорируется ли указанная функция
// Поддерживает как точные совпадения, так и паттерны
func (c *Config) IsFunctionIgnored(funcName string) bool {
	// Проверяем точные совпадения
	for _, ignored := range c.IgnoredFunctions {
		if funcName == ignored {
			return true
		}
	}

	// Проверяем паттерны
	for _, pattern := range c.IgnoredFunctionPatterns {
		if matched, _ := filepath.Match(pattern, funcName); matched {
			return true
		}
	}

	return false
}

// IsBuilderIgnored проверяет, игнорируется ли указанный SQL-билдер
func (c *Config) IsBuilderIgnored(builderName string) bool {
	if builderName == "" {
		return false
	}

	lowerBuilder := strings.ToLower(builderName)

	for _, ignored := range c.IgnoredBuilders {
		if strings.Contains(lowerBuilder, strings.ToLower(ignored)) {
			return true
		}
	}
	return false
}

// IsPatternAllowed проверяет, соответствует ли SQL запрос разрешенному паттерну
func (c *Config) IsPatternAllowed(query string) bool {
	for _, pattern := range c.AllowedPatterns {
		if strings.Contains(strings.ToUpper(query), strings.ToUpper(pattern)) {
			return true
		}
	}
	return false
}
