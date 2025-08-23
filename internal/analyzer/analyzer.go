package analyzer

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/MirrexOne/sqlvet"
)

const (
	linterName = "sqlvet"
	doc        = `checks for "SELECT *" in string literals and SQL builders`
	reportMsg  = `avoid "SELECT *", specify columns explicitly`
)

// newAnalyzer создает новый анализатор с переданной конфигурацией.
func newAnalyzer(cfg *config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     linterName,
		Doc:      doc,
		Run:      run(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// NewDefaultAnalyzer создает анализатор с конфигурацией по умолчанию.
func NewDefaultAnalyzer() *analysis.Analyzer {
	cfg, _ := config.NewFromGolangciLint(nil)
	return newAnalyzer(cfg)
}

// NewFromGolangciLint предоставляет точку входа для golangci-lint.
// Он получает конфигурацию из настроек golangci-lint.
func NewFromGolangciLint(settings map[string]any) (*analysis.Analyzer, error) {
	cfg, err := config.NewFromGolangciLint(settings)
	if err != nil {
		return nil, err
	}
	return newAnalyzer(cfg), nil
}

func run(cfg *config.Config) func(pass *analysis.Pass) (any, error) {
	return func(pass *analysis.Pass) (any, error) {
		// Проверяем, нужно ли игнорировать файл
		if cfg.ShouldIgnoreFile(pass.Fset.File(pass.Files[0].Pos()).Name()) {
			return nil, nil
		}

		inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		nodeFilter := []ast.Node{
			(*ast.BasicLit)(nil), // Для поиска в строковых литералах: "SELECT *"
			(*ast.CallExpr)(nil), // Для поиска в вызовах функций: .Columns("*")
		}

		inspector.Preorder(nodeFilter, func(n ast.Node) {
			// Проверяем, нужно ли игнорировать узел на основе комментариев `nolint`
			// golangci-lint делает это автоматически, но для standalone запуска это может быть полезно.

			switch node := n.(type) {
			case *ast.BasicLit:
				// Анализируем только строковые литералы
				if node.Kind == token.STRING {
					checkStringLiteral(pass, cfg, node, node.Value)
				}
			case *ast.CallExpr:
				// Анализируем вызовы SQL-билдеров, если включено
				if cfg.CheckSQLBuilders {
					checkSQLBuilderCall(pass, cfg, node)
				}
			}
		})

		return nil, nil
	}
}

// checkStringLiteral проверяет строковый литерал на наличие "SELECT *".
func checkStringLiteral(pass *analysis.Pass, cfg *config.Config, node ast.Node, value string) {
	// Удаляем кавычки и приводим к нижнему регистру для надежного поиска
	cleanValue := strings.ToLower(strings.Trim(value, "`\""))
	// Используем регулярное выражение для точного совпадения
	re := regexp.MustCompile(`\bselect\s+\*`)

	if re.MatchString(cleanValue) {
		// Проверяем, не разрешен ли этот паттерн
		if cfg.IsAllowed(cleanValue) {
			return
		}
		pass.Reportf(node.Pos(), reportMsg)
	}
}

// checkSQLBuilderCall проверяет вызов функции, чтобы определить, является ли он вызовом SQL-билдера со "*".
func checkSQLBuilderCall(pass *analysis.Pass, cfg *config.Config, call *ast.CallExpr) {
	// Мы ищем вызовы методов, например, `builder.Columns("*")`
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	funcName := selector.Sel.Name

	// Проверяем, не игнорируется ли эта функция или ее паттерн
	if cfg.ShouldIgnoreFunction(funcName) {
		return
	}

	// Паттерны методов, которые могут содержать выбор всех столбцов.
	// `Select` может быть как с аргументами, так и без (например, squirrel.Select()).
	targetMethods := map[string]bool{
		"Select":  true,
		"Columns": true,
		"Column":  true,
	}

	if !targetMethods[funcName] {
		return
	}

	// Сценарий 1: Пустой вызов `Select()`.
	// Некоторые билдеры (например, squirrel) используют `Select()` без аргументов,
	// а столбцы указываются позже через `.Columns()`.
	// Обнаружение этого требует более сложного анализа потока данных,
	// который мы здесь намеренно опускаем для независимости от реализации.
	// Вместо этого мы сосредоточимся на явных нарушениях.

	// Сценарий 2: Методы с аргументом "*".
	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING && strings.Trim(lit.Value, "`\"") == "*" {
			// Проверяем, не является ли родительский вызов частью игнорируемого билдера.
			// Эта логика может быть усложнена для анализа типа `selector.X`.
			// Для простоты, мы полагаемся на игнорирование по имени функции.
			pass.Reportf(arg.Pos(), reportMsg)
			return // Достаточно одного нарушения на вызов
		}
	}
}
