package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	cfg "github.com/MirrexOne/sqlvet/internal/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer создает анализатор SQLVet с улучшенной логикой
// На основе вашей реализации, но адаптированной для производственного использования
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "sqlvet",
		Doc:      "detects SELECT * in SQL queries and SQL builders, preventing performance issues and encouraging explicit column selection",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// run выполняет основной анализ файлов Go кода на наличие SELECT *
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Определяем типы узлов AST, которые нас интересуют
	nodeFilter := []ast.Node{
		(*ast.BasicLit)(nil), // Строковые литералы
		(*ast.CallExpr)(nil), // Вызовы функций/методов
		(*ast.File)(nil),     // Файлы (для анализа SQL билдеров)
	}

	config := cfg.NewConfig()

	// Проходим по всем узлам AST и анализируем их
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Проверяем, нужно ли пропустить этот файл
		if shouldSkipFile(pass, n, config) {
			return
		}

		// Проверяем nolint комментарии перед анализом узла
		if hasNolintComment(pass, n) {
			return
		}

		switch node := n.(type) {
		case *ast.File:
			// Анализируем SQL билдеры только если это включено в конфигурации
			if config.CheckSQLBuilders {
				analyzeSQLBuilders(pass, node, config)
			}
		case *ast.BasicLit:
			// Проверяем строковые литералы на наличие SELECT *
			checkBasicLit(pass, node, config)
		case *ast.CallExpr:
			// Анализируем вызовы функций на наличие SQL с SELECT *
			checkCallExpr(pass, node, config)
		}
	})

	return nil, nil
}

// hasNolintComment проверяет наличие //nolint:sqlvet комментария перед узлом
// Это обеспечивает стандартную поддержку nolint директив golangci-lint
func hasNolintComment(pass *analysis.Pass, node ast.Node) bool {
	pos := pass.Fset.Position(node.Pos())

	for _, file := range pass.Files {
		if pass.Fset.Position(file.Pos()).Filename != pos.Filename {
			continue
		}

		// Проверяем все комментарии в файле
		for _, commentGroup := range file.Comments {
			// Комментарий должен быть перед узлом и на предыдущей или той же строке
			commentEnd := pass.Fset.Position(commentGroup.End())
			if commentEnd.Filename == pos.Filename &&
				commentGroup.End() < node.Pos() &&
				pos.Line-commentEnd.Line <= 1 {

				for _, comment := range commentGroup.List {
					text := comment.Text
					// Поддерживаем различные варианты nolint комментариев
					if strings.Contains(text, "nolint:sqlvet") ||
						strings.Contains(text, "nolint") {
						return true
					}
				}
			}
		}
	}

	return false
}

// checkBasicLit проверяет строковые литералы на наличие SELECT *
// Базируется на вашей логике normalizeSQLQuery
func checkBasicLit(pass *analysis.Pass, lit *ast.BasicLit, config *cfg.Config) {
	if lit.Kind != token.STRING {
		return
	}

	// Нормализуем SQL запрос используя вашу продвинутую логику
	content := normalizeSQLQuery(lit.Value)
	if isSelectStarQuery(content, config) {
		message := getWarningMessage()
		pass.Reportf(lit.Pos(), "%s", message)
	}
}

// checkCallExpr анализирует вызовы функций на наличие SQL с SELECT *
// Включает проверку аргументов и SQL билдеров
func checkCallExpr(pass *analysis.Pass, call *ast.CallExpr, config *cfg.Config) {
	// Пропускаем игнорируемые функции и пакеты
	if isIgnoredFunctionOrPackage(call, config) {
		return
	}

	// Проверяем SQL билдеры на SELECT * в аргументах
	if config.CheckSQLBuilders && isSQLBuilderSelectStar(call, config) {
		message := getWarningMessage()
		pass.Reportf(call.Pos(), "%s", message)
		return
	}

	// Проверяем аргументы вызова функции на наличие строк с SELECT *
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			content := normalizeSQLQuery(lit.Value)
			if isSelectStarQuery(content, config) {
				message := getWarningMessage()
				pass.Reportf(lit.Pos(), "%s", message)
			}
		}
	}
}

// isIgnoredFunctionOrPackage проверяет, нужно ли игнорировать вызов функции
// Поддерживает как прямые функции, так и методы пакетов
func isIgnoredFunctionOrPackage(call *ast.CallExpr, config *cfg.Config) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Прямой вызов функции (например, myFunc())
		for _, fn := range config.IgnoredFunctions {
			if fun.Name == fn {
				return true
			}
		}

	case *ast.SelectorExpr:
		// Вызов метода пакета (например, pkg.Method())
		if ident, ok := fun.X.(*ast.Ident); ok {
			// Проверяем игнорируемые пакеты
			for _, pkg := range config.IgnoredPackages {
				if strings.EqualFold(ident.Name, pkg) {
					return true
				}
			}

			// Проверяем полное имя функции (pkg.Method)
			fullName := fmt.Sprintf("%s.%s", ident.Name, fun.Sel.Name)
			for _, fn := range config.IgnoredFunctions {
				if strings.EqualFold(fullName, fn) {
					return true
				}
			}
		}
	}
	return false
}

// shouldSkipFile определяет, нужно ли пропустить файл на основе конфигурации
func shouldSkipFile(pass *analysis.Pass, node ast.Node, config *cfg.Config) bool {
	pos := pass.Fset.Position(node.Pos())
	filename := pos.Filename

	// Проверяем паттерны игнорируемых файлов
	for _, pattern := range config.IgnoredFilePatterns {
		// Проверяем как базовое имя файла, так и полный путь
		matched, err := filepath.Match(pattern, filepath.Base(filename))
		if err == nil && matched {
			return true
		}

		matched, err = filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}

	// Проверяем игнорируемые директории
	for _, dir := range config.IgnoredDirectories {
		if isFileInDirectory(filename, dir) {
			return true
		}
	}

	return false
}

// isFileInDirectory проверяет, находится ли файл в указанной директории
func isFileInDirectory(path, dir string) bool {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if strings.EqualFold(segment, dir) && i < len(segments)-1 {
			return true
		}
	}
	return false
}

// normalizeSQLQuery нормализует SQL запрос для анализа
// Это ваша продвинутая реализация с обработкой escape-последовательностей
func normalizeSQLQuery(query string) string {
	if len(query) < 2 {
		return query
	}

	first, last := query[0], query[len(query)-1]

	// 1. Обработка различных типов кавычек с учетом escape-последовательностей
	if first == '"' && last == '"' {
		// Для обычных строк проверяем наличие escape-последовательностей
		if !strings.Contains(query, "\\") {
			query = trimQuotes(query)
		} else if unquoted, err := strconv.Unquote(query); err == nil {
			// Используем стандартный Go unquoting для правильной обработки escape-последовательностей
			query = unquoted
		} else {
			// Fallback: простое удаление кавычек
			query = trimQuotes(query)
		}
	} else if first == '`' && last == '`' {
		// Raw strings - просто удаляем backticks
		query = trimQuotes(query)
	}

	// 2. Обработка комментариев построчно до нормализации
	lines := strings.Split(query, "\n")
	var processedParts []string

	for _, line := range lines {
		// Удаляем комментарии из текущей строки
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}

		// Добавляем непустые строки
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			processedParts = append(processedParts, trimmed)
		}
	}

	// 3. Собираем запрос обратно и нормализуем
	query = strings.Join(processedParts, " ")
	query = strings.ToUpper(query)
	query = strings.ReplaceAll(query, "\t", " ")
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	return strings.TrimSpace(query)
}

// trimQuotes удаляет первый и последний символ (кавычки)
func trimQuotes(query string) string {
	return query[1 : len(query)-1]
}

// isSelectStarQuery определяет, содержит ли запрос SELECT *
// Улучшенная версия с поддержкой разрешенных паттернов
func isSelectStarQuery(query string, config *cfg.Config) bool {
	const sqlKeyword = "SELECT *"

	// Проверяем разрешенные паттерны - если запрос соответствует разрешенному паттерну, игнорируем
	for _, pattern := range config.AllowedPatterns {
		if matched, _ := regexp.MatchString(pattern, query); matched {
			return false
		}
	}

	// Проверяем наличие SELECT * в запросе
	if strings.Contains(query, sqlKeyword) {
		// Убеждаемся, что это действительно SQL запрос, проверяя наличие SQL ключевых слов
		sqlKeywords := []string{"FROM", "WHERE", "JOIN", "GROUP", "ORDER", "HAVING"}
		for _, keyword := range sqlKeywords {
			if strings.Contains(query, keyword) {
				return true
			}
		}
	}
	return false
}

// getWarningMessage возвращает стандартное сообщение предупреждения
func getWarningMessage() string {
	return "SELECT * usage detected"
}

// isSQLBuilderSelectStar проверяет вызовы методов SQL билдеров на наличие SELECT *
func isSQLBuilderSelectStar(call *ast.CallExpr, config *cfg.Config) bool {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Проверяем, что это вызов метода Select
	if fun.Sel == nil || fun.Sel.Name != "Select" {
		return false
	}

	if len(call.Args) == 0 {
		return false
	}

	// Проверяем аргументы метода Select на наличие "*" или пустых строк
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			value := strings.Trim(lit.Value, "`\"")
			// Считаем проблематичными как "*", так и пустые строки в Select()
			if value == "*" || value == "" {
				return true
			}
		}
	}

	return false
}

// analyzeSQLBuilders выполняет продвинутый анализ SQL билдеров
// Это ваша ключевая логика для обработки edge-cases как Select().Columns("*")
func analyzeSQLBuilders(pass *analysis.Pass, file *ast.File, config *cfg.Config) {
	// Отслеживаем переменные SQL билдеров и их состояние
	builderVars := make(map[string]*ast.CallExpr) // Переменные с пустыми Select() вызовами
	hasColumns := make(map[string]bool)           // Флаг: были ли добавлены колонки для переменной

	// Первый проход: находим переменные, созданные с пустыми Select() вызовами
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			// Анализируем присваивания вида: query := builder.Select()
			for i, expr := range node.Rhs {
				if call, ok := expr.(*ast.CallExpr); ok {
					if isEmptySelectCall(call) {
						// Нашли пустой Select() вызов, запоминаем переменную
						if i < len(node.Lhs) {
							if ident, ok := node.Lhs[i].(*ast.Ident); ok {
								builderVars[ident.Name] = call
								hasColumns[ident.Name] = false
							}
						}
					}
				}
			}
		}
		return true
	})

	// Второй проход: проверяем использование методов Columns/Column
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				// Проверяем вызовы методов Columns() или Column()
				if sel.Sel != nil && (sel.Sel.Name == "Columns" || sel.Sel.Name == "Column") {
					// Проверяем на наличие "*" в аргументах
					if hasStarInColumns(node) {
						pass.Reportf(node.Pos(), "SELECT * usage detected")
					}

					// Обновляем состояние переменной - колонки были добавлены
					if ident, ok := sel.X.(*ast.Ident); ok {
						if _, exists := builderVars[ident.Name]; exists {
							if !hasStarInColumns(node) {
								hasColumns[ident.Name] = true
							}
						}
					}
				}
			}

			// Проверяем цепочки вызовов вида builder.Select().Columns("*")
			if isSelectWithColumns(node) {
				if hasStarInColumns(node) {
					if sel, ok := node.Fun.(*ast.SelectorExpr); ok && sel.Sel != nil {
						pass.Reportf(node.Pos(), "SELECT * usage detected")
					}
				}
				return true
			}
		}
		return true
	})

	// Финальная проверка: предупреждаем о билдерах с пустым Select() без последующих колонок
	for varName, call := range builderVars {
		if !hasColumns[varName] {
			pass.Reportf(call.Pos(), "SELECT * usage detected")
		}
	}
}

// isEmptySelectCall проверяет, является ли вызов пустым Select()
func isEmptySelectCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if sel.Sel != nil && sel.Sel.Name == "Select" && len(call.Args) == 0 {
			return true
		}
	}
	return false
}

// isSelectWithColumns проверяет цепочки вызовов вида Select().Columns()
func isSelectWithColumns(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if sel.Sel != nil && (sel.Sel.Name == "Columns" || sel.Sel.Name == "Column") {
			// Проверяем, что предыдущий вызов в цепочке - это Select()
			if innerCall, ok := sel.X.(*ast.CallExpr); ok {
				return isEmptySelectCall(innerCall)
			}
		}
	}
	return false
}

// hasStarInColumns проверяет, содержат ли аргументы вызова символ "*"
func hasStarInColumns(call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			value := strings.Trim(lit.Value, "`\"")
			if value == "*" {
				return true
			}
		}
	}
	return false
}
