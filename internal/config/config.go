package config

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config определяет структуру конфигурации для линтера sqlvet.
type Config struct {
	IgnoreDirs          []string `yaml:"ignore-dirs"`
	IgnoreFilePatterns  []string `yaml:"ignore-file-patterns"`
	IgnorePackages      []string `yaml:"ignore-packages"`
	IgnoreFunctions     []string `yaml:"ignore-functions"`
	IgnoreFunctionRegex []string `yaml:"ignore-function-regex"`
	AllowedPatterns     []string `yaml:"allowed-patterns"`
	CheckSQLBuilders    bool     `yaml:"check-sql-builders"`
	IgnoreBuilders      []string `yaml:"ignore-builders"` // Пока не используется, но зарезервировано

	// Скомпилированные регулярные выражения для производительности
	ignoreFileRegex []*regexp.Regexp
	ignoreFuncRegex []*regexp.Regexp
	allowedRegex    []*regexp.Regexp
	initOnce        sync.Once
}

// defaultConfig создает конфигурацию по умолчанию.
func defaultConfig() *Config {
	return &Config{
		IgnoreDirs:          []string{"vendor", "test", "tests"},
		IgnoreFilePatterns:  []string{"_test.go$", "\\.pb\\.go$"},
		IgnorePackages:      []string{},
		IgnoreFunctions:     []string{},
		IgnoreFunctionRegex: []string{},
		AllowedPatterns:     []string{"count\\(\\*\\)"}, // Разрешаем `count(*)` по умолчанию
		CheckSQLBuilders:    true,
		IgnoreBuilders:      []string{},
	}
}

// NewFromGolangciLint разбирает конфигурацию, предоставленную golangci-lint.
func NewFromGolangciLint(settings map[string]any) (*Config, error) {
	cfg := defaultConfig()

	if settings == nil {
		cfg.compileRegex()
		return cfg, nil
	}

	// golangci-lint передает настройки как map[string]any.
	// Мы преобразуем их в YAML, а затем в нашу структуру для простоты.
	yamlBytes, err := yaml.Marshal(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal linter settings to YAML: %w", err)
	}

	if err := yaml.Unmarshal(yamlBytes, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal linter settings: %w", err)
	}

	cfg.compileRegex()
	return cfg, nil
}

// compileRegex компилирует все строковые паттерны в регулярные выражения.
func (c *Config) compileRegex() {
	c.initOnce.Do(func() {
		c.ignoreFileRegex = compileRegexList(c.IgnoreFilePatterns)
		c.ignoreFuncRegex = compileRegexList(c.IgnoreFunctionRegex)
		c.allowedRegex = compileRegexList(c.AllowedPatterns)
	})
}

func compileRegexList(patterns []string) []*regexp.Regexp {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err == nil {
			compiled = append(compiled, re)
		}
		// Ошибки компиляции игнорируются, неверные паттерны просто не будут работать.
	}
	return compiled
}

// ShouldIgnoreFile проверяет, следует ли игнорировать файл.
func (c *Config) ShouldIgnoreFile(filePath string) bool {
	for _, dir := range c.IgnoreDirs {
		if strings.Contains(filePath, dir+"/") {
			return true
		}
	}
	for _, re := range c.ignoreFileRegex {
		if re.MatchString(filePath) {
			return true
		}
	}
	return false
}

// ShouldIgnoreFunction проверяет, следует ли игнорировать функцию.
func (c *Config) ShouldIgnoreFunction(funcName string) bool {
	for _, name := range c.IgnoreFunctions {
		if name == funcName {
			return true
		}
	}
	for _, re := range c.ignoreFuncRegex {
		if re.MatchString(funcName) {
			return true
		}
	}
	return false
}

// IsAllowed проверяет, соответствует ли строка разрешенному паттерну.
func (c *Config) IsAllowed(value string) bool {
	for _, re := range c.allowedRegex {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}
