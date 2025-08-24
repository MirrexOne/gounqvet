package main

import (
	internal "github.com/MirrexOne/sqlvet/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

// AnalyzerPlugin предоставляет интеграцию SQLVet с golangci-lint
// Это единственная публичная точка входа для внешних инструментов
//
// Функция должна называться именно "AnalyzerPlugin" и возвращать слайс анализаторов
// для совместимости с plugin системой golangci-lint
func AnalyzerPlugin() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		internal.NewAnalyzer(),
	}
}

// Этот файл компилируется как Go plugin для golangci-lint
// При подключении в .golangci.yaml, golangci-lint будет:
// 1. Компилировать этот файл как plugin
// 2. Загружать plugin динамически
// 3. Вызывать функцию AnalyzerPlugin()
// 4. Интегрировать возвращенные анализаторы в свой пайплайн
//
// Преимущества такого подхода:
// - Нет необходимости модифицировать golangci-lint
// - Автоматическая поддержка всех флагов и опций golangci-lint
// - Консистентный вывод с другими линтерами
// - Поддержка //nolint комментариев работает автоматически
