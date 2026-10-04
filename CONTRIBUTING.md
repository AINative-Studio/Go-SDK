# Contributing to the AINative Go SDK

Thanks for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git

### Setting Up Your Development Environment

```bash
git clone https://github.com/yourusername/Go-SDK.git
cd Go-SDK
go mod download
```

## Development Workflow

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Run integration tests (requires a real API key)
AINATIVE_API_KEY=your-key go test -tags=integration ./...
```

All new features should include tests.

### Code Style

Run `gofmt` and `go vet` before committing:
```bash
gofmt -w .
go vet ./...
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes and add tests
4. Run `go test ./...`
5. Push to your fork and open a pull request against this repo

### Pull Request Checklist

- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`gofmt`)
- [ ] `go vet` passes
- [ ] Documentation is updated (README, doc comments)
- [ ] No decrease in test coverage

## Reporting Bugs

Please include:
1. A clear description of the issue
2. Your environment (Go version, OS, module version)
3. Minimal reproduction code
4. Expected vs. actual behavior

## Feature Requests

Please include the use case, a proposed solution, alternatives considered, and example usage.

## Questions?

- Check the [Documentation](https://docs.ainative.studio/sdk/go)
- Open a [GitHub Issue](https://github.com/AINative-Studio/Go-SDK/issues)

## License

By contributing, you agree that your contributions will be licensed under the MIT License (see [LICENSE](LICENSE)).
