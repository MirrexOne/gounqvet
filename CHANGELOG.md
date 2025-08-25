# Changelog

All notable changes to SQLVet will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Integration with golangci-lint v2.x
- Improved error messages with detailed context
- Support for detecting empty Select() calls in SQL builders

### Fixed
- Configuration compatibility with golangci-lint v2
- Removed unused parameter in analyzeSQLBuilders function
- Cleaned up duplicate testdata files

## [1.0.0] - 2025-08-25

### Added
- Initial release of SQLVet linter
- Detection of `SELECT *` in SQL string literals
- Support for SQL builders (Squirrel, GORM, SQLBoiler)
- Configurable ignored functions and packages
- Configurable allowed patterns for legitimate `SELECT *` usage
- Support for `//nolint:sqlvet` directive
- Comprehensive test suite with >80% coverage
- Integration test framework
- Benchmark tests for performance validation
- Example configurations (strict, permissive, standard)
- golangci-lint integration support
- Detailed documentation and usage examples

### Features
- **Core Detection**: Identifies `SELECT *` in string literals and SQL builders
- **Smart Filtering**: Automatically allows COUNT(*), MAX(*), MIN(*) and system queries
- **Configuration**: Extensive configuration options for different use cases
- **Performance**: Efficient AST traversal with minimal overhead
- **Integration**: Seamless integration with golangci-lint
- **Suppression**: Standard `//nolint:sqlvet` comment support

### Configuration Defaults
- Ignored file patterns: `*_test.go`, `*.pb.go`, `*_gen.go`, `*.gen.go`, `*_generated.go`
- Ignored directories: `vendor`, `testdata`, `migrations`, `generated`, `.git`, `node_modules`
- Allowed patterns: COUNT(*), MAX(*), MIN(*), information_schema, pg_catalog, sys schema queries

## [0.1.0] - 2025-08-20 (Pre-release)

### Added
- Basic SELECT * detection in string literals
- Initial analyzer implementation
- Basic test coverage
- Initial documentation

---

## Version History

- **v1.0.0** - Production-ready release with full golangci-lint integration
- **v0.1.0** - Initial pre-release for testing

## Upgrade Guide

### From v0.x to v1.0

1. Update your import paths if using as a library
2. Review new configuration options in `.golangci.yml`
3. Test with your codebase to ensure compatibility

### golangci-lint Integration

To use with golangci-lint v2.x:

```yaml
version: "2"

linters:
  enable:
    - sqlvet

linters-settings:
  sqlvet:
    check-sql-builders: true
```

## Future Plans

- [ ] Support for more SQL builders
- [ ] Detection of other SQL anti-patterns
- [ ] Performance optimizations
- [ ] IDE integration plugins
- [ ] Configuration presets for popular frameworks

## Links

- [GitHub Repository](https://github.com/MirrexOne/sqlvet)
- [Issue Tracker](https://github.com/MirrexOne/sqlvet/issues)
- [Documentation](https://godoc.org/github.com/MirrexOne/sqlvet)
