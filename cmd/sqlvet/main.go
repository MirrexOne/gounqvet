package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	internal "github.com/MirrexOne/sqlvet/internal/analyzer"
)

// main предоставляет CLI интерфейс для SQLVet линтера
//
// Этот файл создает standalone версию анализатора, которую можно запускать
// напрямую из командной строки без golangci-lint:
//
//	sqlvet ./...
//	sqlvet -help
//	sqlvet /path/to/project
//
// Использует стандартный singlechecker из golang.org/x/tools, что обеспечивает:
// - Консистентный CLI интерфейс с другими Go анализаторами
// - Автоматическую поддержку стандартных флагов (-json, -c=N и т.д.)
// - Правильное форматирование вывода для различных инструментов
func main() {
	// NewAnalyzer() создает экземпляр анализатора из internal пакета
	// singlechecker.Main() обрабатывает CLI аргументы и запускает анализ
	singlechecker.Main(internal.NewAnalyzer())
}
