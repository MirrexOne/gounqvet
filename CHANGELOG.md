# Changelog

All notable changes to SQLVet will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- golangci-lint native integration support
- Comprehensive configuration system with mapstructure tags
- Support for `//nolint:sqlvet` and general `//nolint` directives
- Extensive test suite with functional tests
- CI/CD pipeline with GitHub Actions
- Multiple configuration examples (strict, permissive, standard)
- Detailed documentation and contribution guidelines

### Changed
- Improved analyzer performance with better AST traversal
- Enhanced SQL pattern detection with edge case handling
- Better error messages with context information
- Restructured project for golangci-lint integration

### Fixed
- Fixed false positives with COUNT(*) and other aggregate functions
- Improved handling of multiline SQL queries
- Better detection of SQL builders with method chaining
- Fixed nolint directive processing

## [1.0.0] - 2024-08-24

### Added
- Initial release of SQLVet
- Basic SELECT * detection in string literals
- SQL builder support for popular libraries
- Configuration system for customization
- Command-line interface
- Core analysis engine using go/analysis framework

### Features
- Detects `SELECT *` in Go string literals
- Supports SQL builders like Squirrel
- Configurable ignore patterns and functions
- Fast parallel analysis
- Integration with existing Go tooling

[Unreleased]: https://github.com/MirrexOne/sqlvet/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/MirrexOne/sqlvet/releases/tag/v1.0.0