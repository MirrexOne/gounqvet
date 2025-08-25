# Contributing to SQLVet

We welcome contributions to SQLVet! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/yourusername/sqlvet.git
   cd sqlvet
   ```
3. **Install dependencies**:
   ```bash
   go mod download
   ```
4. **Run tests** to ensure everything works:
   ```bash
   go test ./...
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the guidelines below

3. **Add tests** for any new functionality

4. **Run the test suite**:
   ```bash
   go test ./...
   go test -race ./...
   ```

5. **Run the linter**:
   ```bash
   golangci-lint run
   ```

6. **Commit your changes**:
   ```bash
   git commit -am "Add your descriptive commit message"
   ```

7. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

8. **Create a Pull Request** on GitHub

### Code Style

- Follow standard Go conventions and idioms
- Use `go fmt` to format your code
- Write clear, descriptive variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and reasonably sized

### Testing

- Add tests for all new functionality
- Update existing tests when modifying behavior
- Ensure all tests pass before submitting PR
- Include both positive and negative test cases
- Add testdata examples for new detection patterns

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected bool
    }{
        {"should detect", "SELECT * FROM users", true},
        {"should not detect", "SELECT id FROM users", false},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Test logic here
        })
    }
}
```

### Adding New Detection Patterns

When adding new SQL patterns that sqlvet should detect:

1. **Add test cases** in `testdata/sqlvet.go` with `// want` comments
2. **Update the analyzer logic** in `internal/analyzer/analyzer.go`
3. **Add configuration options** if needed in `config.go`
4. **Update documentation** in README.md

### Project Structure

```
sqlvet/
├── analyzer.go              # Main analyzer export
├── config.go               # Configuration structures  
├── internal/
│   ├── analyzer/           # Core analysis logic
│   └── config/            # Configuration handling
├── testdata/              # Test cases for analysis
├── examples/              # Configuration examples
└── cmd/sqlvet/           # CLI tool
```

## Types of Contributions

### Bug Reports

When filing bug reports, please include:

- **Clear description** of the issue
- **Steps to reproduce** the problem
- **Expected vs actual behavior**
- **Go version** and operating system
- **Minimal code example** that demonstrates the issue

### Feature Requests

For feature requests, please provide:

- **Clear description** of the proposed feature
- **Use case** explaining why it would be valuable
- **Proposed implementation** approach (if you have ideas)
- **Backwards compatibility** considerations

### Documentation Improvements

Documentation contributions are very welcome:

- Fix typos or unclear explanations
- Add examples for configuration options
- Improve installation or usage instructions
- Add FAQ entries for common questions

## Code Review Process

1. **Automated checks** must pass (CI, tests, linting)
2. **Manual review** by maintainers
3. **Discussion** of any necessary changes
4. **Approval** and merge

### Review Criteria

We look for:

- **Correctness**: Does the code work as intended?
- **Testing**: Are there adequate tests?
- **Documentation**: Is new functionality documented?
- **Performance**: Does it avoid performance regressions?
- **Backwards compatibility**: Does it break existing functionality?

## Release Process

SQLVet follows semantic versioning:

- **Major version** (v1.0.0): Breaking changes
- **Minor version** (v1.1.0): New features, backwards compatible
- **Patch version** (v1.1.1): Bug fixes, backwards compatible

## Community Guidelines

- **Be respectful** and professional
- **Ask questions** if anything is unclear
- **Help others** by reviewing PRs and answering questions
- **Follow the code of conduct** (treat everyone with respect)

## Getting Help

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Code Review**: Feel free to ask questions in PR comments

## Advanced Contributing

### Working on golangci-lint Integration

If you're working on improving golangci-lint integration:

1. Test against the golangci-lint repository
2. Ensure configuration options work correctly
3. Verify nolint directive support
4. Check performance with large codebases

### Performance Considerations

- Profile your changes with large codebases
- Avoid unnecessary allocations in hot paths
- Use efficient algorithms for AST traversal
- Consider memory usage with concurrent analysis

### Integration Testing

Test sqlvet with real projects:

```bash
# Test with a real project
cd /path/to/some/go/project
go run github.com/MirrexOne/sqlvet/cmd/sqlvet ./...

# Test with golangci-lint
golangci-lint run --enable-only=sqlvet
```

## Thank You!

Your contributions help make SQLVet better for everyone. Thank you for taking the time to contribute!