# Gounqvet Examples

This directory contains configuration examples and usage demonstrations for Gounqvet.

## Configuration Files

- `golangci.yml` - Standard configuration for golangci-lint integration
- `strict-config.yml` - Strict configuration that catches almost everything
- `permissive-config.yml` - Permissive configuration for gradual adoption

## Code Examples

See `testdata/example.go` for code examples demonstrating:

- Patterns that Gounqvet will warn about
- Recommended patterns that pass validation
- How to suppress warnings with `//nolint:gounqvet`
- SQL builder usage patterns

## Running the Examples

To see Gounqvet in action on the example file:

```bash
# Run from project root
go run ./cmd/gounqvet ./examples/testdata/example.go
```

This will show warnings for the problematic patterns in the `ExampleBadCode` function.

## Using Configuration Files

Copy any of the configuration files to your project root as `.golangci.yml` and customize as needed:

```bash
cp examples/golangci.yml .golangci.yml
# Edit .golangci.yml to fit your needs
```

Then run golangci-lint:

```bash
golangci-lint run ./...
```
