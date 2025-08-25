# Contributing to SQLVet

Thank you for your interest in contributing to SQLVet! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Check existing issues to avoid duplicates
2. Use the issue templates when available
3. Provide clear descriptions and reproducible examples

### Suggesting Features

Feature requests are welcome! Please:
1. Explain the problem your feature would solve
2. Provide use cases and examples
3. Consider if it aligns with the project's goals

### Submitting Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Follow the coding standards** (see below)
3. **Write tests** for new functionality
4. **Update documentation** as needed
5. **Ensure all tests pass** before submitting
6. **Follow commit message guidelines** (see COMMIT_STYLE.md)

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/sqlvet.git
cd sqlvet

# Add upstream remote
git remote add upstream https://github.com/MirrexOne/sqlvet.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Run linter
golangci-lint run ./...
```

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Keep functions focused and small
- Add comments for exported functions and types
- Handle errors appropriately

### Package Structure

```
sqlvet/
├── cmd/           # Command-line interface
├── internal/      # Internal packages
├── pkg/           # Public packages
├── testdata/      # Test fixtures
├── examples/      # Usage examples
└── docs/          # Documentation
```

### Testing Guidelines

- Write unit tests for new functions
- Use table-driven tests where appropriate
- Include edge cases and error scenarios
- Maintain test coverage above 80%

Example test structure:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"test case 1", "input", "expected"},
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Commit Message Guidelines

See [COMMIT_STYLE.md](COMMIT_STYLE.md) for detailed commit message format.

Quick reference:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions or changes
- `refactor:` - Code refactoring
- `chore:` - Maintenance tasks

## Pull Request Process

1. **Update your fork**:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes** and commit them

4. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

5. **Create a Pull Request** on GitHub

### PR Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Linter passes (`golangci-lint run ./...`)
- [ ] Documentation updated
- [ ] Commit messages follow guidelines
- [ ] Branch is up-to-date with main

## Adding New Linter Rules

When adding new detection patterns:

1. **Update the analyzer** in `internal/analyzer/analyzer.go`
2. **Add test cases** in `internal/analyzer/analyzer_test.go`
3. **Add integration tests** in `testdata/`
4. **Update documentation** in README.md
5. **Add configuration options** if needed

## Release Process

Releases follow semantic versioning (MAJOR.MINOR.PATCH):

- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

To create a release:
1. Update CHANGELOG.md
2. Create a git tag: `git tag vX.Y.Z`
3. Push the tag: `git push origin vX.Y.Z`

## Getting Help

- Check the [documentation](README.md)
- Review [existing issues](https://github.com/MirrexOne/sqlvet/issues)
- Ask questions in [discussions](https://github.com/MirrexOne/sqlvet/discussions)

## Recognition

Contributors will be recognized in:
- The project's contributor list
- Release notes for significant contributions
- CHANGELOG.md for feature additions

Thank you for contributing to SQLVet!
