package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MirrexOne/sqlvet/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NewAnalyzer creates the SQLVet analyzer with enhanced logic for production use
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "sqlvet",
		Doc:      "detects SELECT * in SQL queries and SQL builders, preventing performance issues and encouraging explicit column selection",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// NewAnalyzerWithSettings creates analyzer with provided settings for golangci-lint integration
func NewAnalyzerWithSettings(s config.SQLVetSettings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "sqlvet",
		Doc:      "detects SELECT * in SQL queries and SQL builders, preventing performance issues and encouraging explicit column selection",
		Run:      func(pass *analysis.Pass) (any, error) {
			return RunWithConfig(pass, &s)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// RunWithConfig performs analysis with provided configuration
// This is the main entry point for configured analysis
func RunWithConfig(pass *analysis.Pass, cfg *config.SQLVetSettings) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Use provided configuration or default if nil
	if cfg == nil {
		defaultSettings := config.DefaultSettings()
		cfg = &defaultSettings
	}

	// Define AST node types we're interested in
	nodeFilter := []ast.Node{
		(*ast.BasicLit)(nil), // String literals
		(*ast.CallExpr)(nil), // Function/method calls
		(*ast.File)(nil),     // Files (for SQL builder analysis)
	}

	// Walk through all AST nodes and analyze them
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Check if we should skip this file
		if shouldSkipFile(pass, n, cfg) {
			return
		}

		// Check for nolint comments before analyzing the node
		if hasNolintComment(pass, n) {
			return
		}

		switch node := n.(type) {
		case *ast.File:
			// Analyze SQL builders only if enabled in configuration
			if cfg.CheckSQLBuilders {
				analyzeSQLBuilders(pass, node, cfg)
			}
		case *ast.BasicLit:
			// Check string literals for SELECT * usage
			checkBasicLit(pass, node, cfg)
		case *ast.CallExpr:
			// Analyze function calls for SQL with SELECT * usage
			checkCallExpr(pass, node, cfg)
		}
	})

	return nil, nil
}

// run performs the main analysis of Go code files for SELECT * usage
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Define AST node types we're interested in
	nodeFilter := []ast.Node{
		(*ast.BasicLit)(nil), // String literals
		(*ast.CallExpr)(nil), // Function/method calls
		(*ast.File)(nil),     // Files (for SQL builder analysis)
	}

	// Always use default settings since passing settings through ResultOf doesn't work reliably
	defaultSettings := config.DefaultSettings()
	cfg := &defaultSettings

	// Walk through all AST nodes and analyze them
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Check if we should skip this file
		if shouldSkipFile(pass, n, cfg) {
			return
		}

		// Check for nolint comments before analyzing the node
		if hasNolintComment(pass, n) {
			return
		}

		switch node := n.(type) {
		case *ast.File:
			// Analyze SQL builders only if enabled in configuration
			if cfg.CheckSQLBuilders {
				analyzeSQLBuilders(pass, node, cfg)
			}
		case *ast.BasicLit:
			// Check string literals for SELECT * usage
			checkBasicLit(pass, node, cfg)
		case *ast.CallExpr:
			// Analyze function calls for SQL with SELECT * usage
			checkCallExpr(pass, node, cfg)
		}
	})

	return nil, nil
}

// hasNolintComment checks for //nolint:sqlvet comment before the node
// This provides standard nolint directive support for golangci-lint
func hasNolintComment(pass *analysis.Pass, node ast.Node) bool {
	pos := pass.Fset.Position(node.Pos())

	for _, file := range pass.Files {
		if pass.Fset.Position(file.Pos()).Filename != pos.Filename {
			continue
		}

		// Check all comments in the file
		for _, commentGroup := range file.Comments {
			// Comment should be before the node and on the previous or same line
			commentEnd := pass.Fset.Position(commentGroup.End())
			if commentEnd.Filename == pos.Filename &&
				commentGroup.End() < node.Pos() &&
				pos.Line-commentEnd.Line <= 1 {

				for _, comment := range commentGroup.List {
					text := comment.Text
					// Support various nolint comment variants
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

// checkBasicLit checks string literals for SELECT * usage
// Based on the normalizeSQLQuery logic
func checkBasicLit(pass *analysis.Pass, lit *ast.BasicLit, cfg *config.SQLVetSettings) {
	if lit.Kind != token.STRING {
		return
	}

	// Normalize SQL query using advanced logic
	content := normalizeSQLQuery(lit.Value)
	if isSelectStarQuery(content, cfg) {
		pass.Report(analysis.Diagnostic{
			Pos:     lit.Pos(),
			Message: getWarningMessage(),
		})
	}
}

// checkCallExpr analyzes function calls for SQL with SELECT * usage
// Includes checking arguments and SQL builders
func checkCallExpr(pass *analysis.Pass, call *ast.CallExpr, cfg *config.SQLVetSettings) {
	// Skip ignored functions and packages
	if isIgnoredFunctionOrPackage(call, cfg) {
		return
	}

	// Check SQL builders for SELECT * in arguments
	if cfg.CheckSQLBuilders && isSQLBuilderSelectStar(call, cfg) {
		pass.Report(analysis.Diagnostic{
			Pos:     call.Pos(),
			Message: getDetailedWarningMessage("sql_builder"),
		})
		return
	}

	// Check function call arguments for strings with SELECT *
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			content := normalizeSQLQuery(lit.Value)
			if isSelectStarQuery(content, cfg) {
				pass.Report(analysis.Diagnostic{
					Pos:     lit.Pos(),
					Message: getWarningMessage(),
				})
			}
		}
	}
}

// isIgnoredFunctionOrPackage checks if function call should be ignored
// Supports both direct functions and package methods
func isIgnoredFunctionOrPackage(call *ast.CallExpr, cfg *config.SQLVetSettings) bool {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		// Direct function call (e.g., myFunc())
		for _, fn := range cfg.IgnoredFunctions {
			if fun.Name == fn {
				return true
			}
		}

	case *ast.SelectorExpr:
		// Package method call (e.g., pkg.Method())
		if ident, ok := fun.X.(*ast.Ident); ok {
			// Check ignored packages
			for _, pkg := range cfg.IgnoredPackages {
				if ident.Name == pkg {
					return true
				}
			}

			// Check full function name (pkg.Method)
			fullName := fmt.Sprintf("%s.%s", ident.Name, fun.Sel.Name)
			for _, fn := range cfg.IgnoredFunctions {
				if fullName == fn {
					return true
				}
			}
		}
	}
	return false
}


// shouldSkipFile determines if file should be skipped based on configuration
func shouldSkipFile(pass *analysis.Pass, node ast.Node, cfg *config.SQLVetSettings) bool {
	pos := pass.Fset.Position(node.Pos())
	filename := pos.Filename

	// For analysis tests, don't skip files in testdata directories
	// This is important for golangci-lint integration testing
	if strings.Contains(filename, "analysistest") || strings.Contains(filename, "testdata") {
		return false
	}

	// Check ignored file patterns
	for _, pattern := range cfg.IgnoredFilePatterns {
		// Check both base filename and full path
		matched, err := filepath.Match(pattern, filepath.Base(filename))
		if err == nil && matched {
			return true
		}

		matched, err = filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}

	// Check ignored directories
	for _, dir := range cfg.IgnoredDirectories {
		if isFileInDirectory(filename, dir) {
			return true
		}
	}

	return false
}

// isFileInDirectory checks if file is in the specified directory
func isFileInDirectory(path, dir string) bool {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if strings.EqualFold(segment, dir) && i < len(segments)-1 {
			return true
		}
	}
	return false
}

// normalizeSQLQuery normalizes SQL query for analysis
// Advanced implementation with escape sequence handling
// NormalizeSQLQuery is exported for testing
func NormalizeSQLQuery(query string) string {
	return normalizeSQLQuery(query)
}

func normalizeSQLQuery(query string) string {
	if len(query) < 2 {
		return query
	}

	first, last := query[0], query[len(query)-1]

	// 1. Handle different quote types with escape sequence processing
	if first == '"' && last == '"' {
		// For regular strings check for escape sequences
		if !strings.Contains(query, "\\") {
			query = trimQuotes(query)
		} else if unquoted, err := strconv.Unquote(query); err == nil {
			// Use standard Go unquoting for proper escape sequence handling
			query = unquoted
		} else {
			// Fallback: simple quote removal
			query = trimQuotes(query)
		}
	} else if first == '`' && last == '`' {
		// Raw strings - simply remove backticks
		query = trimQuotes(query)
	}

	// 2. Process comments line by line before normalization
	lines := strings.Split(query, "\n")
	var processedParts []string

	for _, line := range lines {
		// Remove comments from current line
		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}

		// Add non-empty lines
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			processedParts = append(processedParts, trimmed)
		}
	}

	// 3. Reassemble query and normalize
	query = strings.Join(processedParts, " ")
	query = strings.ToUpper(query)
	query = strings.ReplaceAll(query, "\t", " ")
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	return strings.TrimSpace(query)
}

// trimQuotes removes first and last character (quotes)
func trimQuotes(query string) string {
	return query[1 : len(query)-1]
}

// isSelectStarQuery determines if query contains SELECT *
// Enhanced version with allowed patterns support
// IsSelectStarQuery is exported for testing
func IsSelectStarQuery(query string, cfg *config.SQLVetSettings) bool {
	return isSelectStarQuery(query, cfg)
}

func isSelectStarQuery(query string, cfg *config.SQLVetSettings) bool {
	// Check allowed patterns first - if query matches an allowed pattern, ignore it
	for _, pattern := range cfg.AllowedPatterns {
		if matched, _ := regexp.MatchString(pattern, query); matched {
			return false
		}
	}

	// Check for SELECT * in query (case-insensitive)
	upperQuery := strings.ToUpper(query)
	if strings.Contains(upperQuery, "SELECT *") {
		// Ensure this is actually an SQL query by checking for SQL keywords
		sqlKeywords := []string{"FROM", "WHERE", "JOIN", "GROUP", "ORDER", "HAVING", "UNION", "LIMIT"}
		for _, keyword := range sqlKeywords {
			if strings.Contains(upperQuery, keyword) {
				return true
			}
		}
		
		// Also check if it's just "SELECT *" without other keywords (still problematic)
		trimmed := strings.TrimSpace(upperQuery)
		if trimmed == "SELECT *" {
			return true
		}
	}
	return false
}

// getWarningMessage returns informative warning message
func getWarningMessage() string {
	return "avoid SELECT * - explicitly specify needed columns for better performance, maintainability and stability"
}

// getDetailedWarningMessage returns context-specific warning message
func getDetailedWarningMessage(context string) string {
	switch context {
	case "sql_builder":
		return "avoid SELECT * in SQL builder - explicitly specify columns to prevent unnecessary data transfer and schema change issues"
	case "nested":
		return "avoid SELECT * in subquery - can cause performance issues and unexpected results when schema changes"
	case "empty_select":
		return "SQL builder Select() without columns defaults to SELECT * - add specific columns with .Columns() method"
	default:
		return "avoid SELECT * - explicitly specify needed columns for better performance, maintainability and stability"
	}
}

// isSQLBuilderSelectStar checks SQL builder method calls for SELECT * usage
func isSQLBuilderSelectStar(call *ast.CallExpr, cfg *config.SQLVetSettings) bool {
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Check that this is a Select method call
	if fun.Sel == nil || fun.Sel.Name != "Select" {
		return false
	}

	if len(call.Args) == 0 {
		return false
	}

	// Check Select method arguments for "*" or empty strings
	for _, arg := range call.Args {
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			value := strings.Trim(lit.Value, "`\"")
			// Consider both "*" and empty strings in Select() as problematic
			if value == "*" || value == "" {
				return true
			}
		}
	}

	return false
}

// analyzeSQLBuilders performs advanced SQL builder analysis
// Key logic for handling edge-cases like Select().Columns("*")
func analyzeSQLBuilders(pass *analysis.Pass, file *ast.File, cfg *config.SQLVetSettings) {
	// Track SQL builder variables and their state
	builderVars := make(map[string]*ast.CallExpr) // Variables with empty Select() calls
	hasColumns := make(map[string]bool)           // Flag: were columns added for variable

	// First pass: find variables created with empty Select() calls
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.AssignStmt:
			// Analyze assignments like: query := builder.Select()
			for i, expr := range node.Rhs {
				if call, ok := expr.(*ast.CallExpr); ok {
					if isEmptySelectCall(call) {
						// Found empty Select() call, remember the variable
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

	// Second pass: check usage of Columns/Column methods
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.CallExpr:
			if sel, ok := node.Fun.(*ast.SelectorExpr); ok {
				// Check calls to Columns() or Column() methods
				if sel.Sel != nil && (sel.Sel.Name == "Columns" || sel.Sel.Name == "Column") {
					// Check for "*" in arguments
					if hasStarInColumns(node) {
						pass.Report(analysis.Diagnostic{
							Pos:     node.Pos(),
							Message: getDetailedWarningMessage("sql_builder"),
						})
					}

					// Update variable state - columns were added
					if ident, ok := sel.X.(*ast.Ident); ok {
						if _, exists := builderVars[ident.Name]; exists {
							if !hasStarInColumns(node) {
								hasColumns[ident.Name] = true
							}
						}
					}
				}
			}

			// Check call chains like builder.Select().Columns("*")
			if isSelectWithColumns(node) {
				if hasStarInColumns(node) {
					if sel, ok := node.Fun.(*ast.SelectorExpr); ok && sel.Sel != nil {
						pass.Report(analysis.Diagnostic{
							Pos:     node.Pos(),
							Message: getDetailedWarningMessage("sql_builder"),
						})
					}
				}
				return true
			}
		}
		return true
	})

	// Final check: warn about builders with empty Select() without subsequent columns
	for varName, call := range builderVars {
		if !hasColumns[varName] {
			pass.Report(analysis.Diagnostic{
				Pos:     call.Pos(),
				Message: getDetailedWarningMessage("empty_select"),
			})
		}
	}
}

// isEmptySelectCall checks if call is an empty Select()
func isEmptySelectCall(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if sel.Sel != nil && sel.Sel.Name == "Select" && len(call.Args) == 0 {
			return true
		}
	}
	return false
}

// isSelectWithColumns checks call chains like Select().Columns()
func isSelectWithColumns(call *ast.CallExpr) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if sel.Sel != nil && (sel.Sel.Name == "Columns" || sel.Sel.Name == "Column") {
			// Check that previous call in chain is Select()
			if innerCall, ok := sel.X.(*ast.CallExpr); ok {
				return isEmptySelectCall(innerCall)
			}
		}
	}
	return false
}

// hasStarInColumns checks if call arguments contain "*" symbol
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
