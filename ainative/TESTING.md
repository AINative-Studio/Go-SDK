# Go SDK Testing Guide

This guide explains how to run tests for the AINative Go SDK, specifically the embeddings operations.

---

## Quick Start

### Run All Tests
```bash
cd /Users/aideveloper/core/developer-tools/sdks/go
go test ./ainative -v
```

### Run Specific Test Suites
```bash
# Run only Generate tests
go test ./ainative -v -run TestGenerate

# Run only EmbedAndStore tests
go test ./ainative -v -run TestEmbedAndStore

# Run only SemanticSearch tests
go test ./ainative -v -run TestSemanticSearch

# Run only error handling tests
go test ./ainative -v -run TestHTTPError
```

---

## Coverage Analysis

### Generate Coverage Report
```bash
go test ./ainative -cover
```

### Detailed Coverage Report
```bash
go test ./ainative -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

This will open an HTML report in your browser showing line-by-line coverage.

### Coverage by Package
```bash
go test ./ainative -coverprofile=coverage.out -covermode=count
go tool cover -func=coverage.out
```

---

## Integration Tests

Integration tests require a live API key and project ID.

### Setup
```bash
export AINATIVE_API_KEY="your-api-key-here"
export AINATIVE_PROJECT_ID="your-project-id-here"
```

### Run Integration Tests
```bash
go test -tags=integration ./ainative -v
```

### Run Specific Integration Test
```bash
go test -tags=integration ./ainative -v -run TestIntegration_Generate
go test -tags=integration ./ainative -v -run TestIntegration_FullWorkflow
```

---

## Race Condition Detection

Check for race conditions:
```bash
go test ./ainative -race
```

For integration tests:
```bash
go test -tags=integration ./ainative -race
```

---

## Benchmarking

Run benchmarks (if any):
```bash
go test ./ainative -bench=. -benchmem
```

---

## Test Organization

### Unit Tests (`embeddings_test.go`)
- **27 tests** covering all embeddings operations
- Uses `httptest.NewServer` for HTTP mocking
- No external dependencies required
- Fast execution (< 1 second)

**Test Categories**:
- Generate operation (6 tests)
- EmbedAndStore operation (4 tests)
- SemanticSearch operation (7 tests)
- Other operations (3 tests)
- Error handling (7 tests)

### Integration Tests (`embeddings_integration_test.go`)
- **7 tests** for end-to-end validation
- Requires live API credentials
- Tests against actual API
- Build tag: `integration`

**Test Categories**:
- Individual operation tests (6 tests)
- Full workflow test (1 test)

---

## Expected Output

### Successful Test Run
```
=== RUN   TestGenerate_Success
--- PASS: TestGenerate_Success (0.00s)
=== RUN   TestGenerate_EmptyTexts
--- PASS: TestGenerate_EmptyTexts (0.00s)
=== RUN   TestGenerate_DefaultModel
--- PASS: TestGenerate_DefaultModel (0.00s)
...
PASS
coverage: 87.5% of statements
ok      ainative        0.234s
```

### Coverage Report Example
```
ainative/embeddings.go:108:     Generate                87.5%
ainative/embeddings.go:206:     EmbedAndStore           90.0%
ainative/embeddings.go:287:     SemanticSearch          92.3%
ainative/embeddings.go:362:     ListModels              100.0%
ainative/embeddings.go:393:     HealthCheck             100.0%
ainative/embeddings.go:425:     GetUsage                100.0%
total:                          (statements)            87.8%
```

---

## Continuous Integration

### GitHub Actions Example
```yaml
name: Go Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test ./ainative -v -cover

      - name: Run race detector
        run: go test ./ainative -race
```

---

## Troubleshooting

### Tests Fail to Compile
```bash
# Ensure dependencies are installed
go mod tidy
go mod download
```

### Coverage Too Low
```bash
# Generate detailed coverage report to identify gaps
go test ./ainative -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Integration Tests Fail
```bash
# Verify environment variables are set
echo $AINATIVE_API_KEY
echo $AINATIVE_PROJECT_ID

# Check API connectivity
curl -H "Authorization: Bearer $AINATIVE_API_KEY" \
  https://api.ainative.studio/health
```

---

## Test Quality Checklist

✅ Tests are isolated (no shared state)
✅ Tests are deterministic (consistent results)
✅ Tests are fast (< 1 second for unit tests)
✅ Tests have clear names
✅ Tests follow AAA pattern (Arrange, Act, Assert)
✅ Error cases are tested
✅ Edge cases are covered
✅ Integration tests are separated with build tags

---

## Contributing

When adding new features to the embeddings service:

1. **Add unit tests** in `embeddings_test.go`
2. **Add integration test** in `embeddings_integration_test.go` (if applicable)
3. **Run all tests** before committing
4. **Check coverage** meets minimum 80%
5. **Run race detector** to ensure thread safety

---

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Library](https://github.com/stretchr/testify)
- [Go Test Coverage](https://go.dev/blog/cover)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

---

**Last Updated**: 2025-12-21
