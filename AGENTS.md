# AGENTS.md - Guide for AI Coding Agents

This document provides essential information for AI agents working with the Serenity-Go codebase.

## Project Overview

Serenity-Go is a Go implementation of the Serenity/JS Screenplay Pattern for acceptance testing, focused on API testing capabilities. It provides actor-centric testing with reusable components and clear domain language.

## Build, Test, and Development Commands

### Primary Commands (Use Makefile)
```bash
# Full development cycle
make all              # clean deps mocks fmt lint test

# Testing commands
make test             # go test ./...
make test-v           # go test -v ./...
make test-coverage    # go test -cover ./...
make test-bench       # go test -bench=. ./...

# Code quality
make fmt              # gofmt -s -w .
make fmt-check        # check formatting without modifying
make lint             # golangci-lint run
make vet              # go vet ./...
make check            # fmt-check lint test
make ci               # fmt lint test (for CI)

# Dependencies and Build
make deps             # go mod download && go mod tidy
make build            # go build ./...
make clean            # go clean -cache

# Mocks
make mocks            # go generate ./...
make mocks-clean      # find . -name "mock_*.go" -delete
```

### Direct Go Commands (Fallback)
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in specific package
go test ./serenity/core -v
go test ./serenity/abilities/api -v
go test ./serenity/expectations -v
go test ./serenity/testing -v
go test ./examples -v

# Run a single test
go test -run TestSpecificFunction ./path/to/package

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...

# Build the module
go build ./...

# Clean build cache
go clean -cache

# Generate mocks
go generate ./...

# Dependency management
go mod download
go mod tidy
go mod verify
go get -u ./...
```

## Code Style Guidelines

### Package and Import Organization
- Standard library imports grouped first, then third-party, then local imports
- Use blank lines between import groups
- Local imports use the full module path: `github.com/nchursin/serenity-go/serenity/...`
- golangci-lint enforces goimports formatting automatically

Example:
```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"

    "github.com/nchursin/serenity-go/serenity/core"
    "github.com/nchursin/serenity-go/serenity/abilities/api"
    "github.com/nchursin/serenity-go/serenity/testing"
)
```

### Naming Conventions
- **Package names**: lowercase, single word when possible (e.g., `core`, `api`, `expectations`, `testing`)
- **Public functions/types**: PascalCase (e.g., `NewActor`, `RequestBuilder`, `NewSerenityTest`)
- **Private functions/types**: camelCase (e.g., `sendRequest`, `abilityTypeOf`, `callAnAPI`)
- **Interfaces**: Often include type parameter for generics (e.g., `Question[T any]`, `Expectation[T any]`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase for both exported and unexported
- **Mock files**: prefixed with `mock_` and generated in `mocks/` subdirectories

### Type and Interface Design
- Use generics for type-safe interfaces with `T any` syntax
- Interface methods should have clear, descriptive names
- Separate interfaces for different concerns (Actor, Activity, Question, etc.)
- Use composition over inheritance where possible

Example:
```go
type Question[T any] interface {
    AnsweredBy(actor Actor) (T, error)
    Description() string
}

type Actor interface {
    Name() string
    WhoCan(abilities ...abilities.Ability) Actor
    AbilityTo(ability abilities.Ability) (abilities.Ability, error)
    AttemptsTo(activities ...Activity)
    AnswersTo(question Question[any]) (any, bool)
}
```

### Error Handling
- Always wrap errors with context using `fmt.Errorf` with `%w` verb
- Return early from functions when errors occur
- Use descriptive error messages that include context
- Test error paths in unit tests
- New TestContext API automatically handles test failures without manual `require.NoError`

Example:
```go
if err := someOperation(); err != nil {
    return fmt.Errorf("failed to perform operation: %w", err)
}
```

### Function and Method Organization
- Keep functions focused and single-purpose
- Use builder patterns for complex object construction
- Chain method calls where it improves readability
- Use descriptive function names that explain what they do
- Split large files into focused files by responsibility

Example:
```go
// Fluent request building
req, err := api.NewRequestBuilder("POST", "/posts").
    WithHeader("Content-Type", "application/json").
    WithHeader("Authorization", "Bearer token").
    With(postData).
    Build()
```

### Struct Organization
- Fields should be ordered logically (public first, then private)
- Use embedded types only when it provides clear value
- Include JSON tags for structs that are serialized
- Use pointer types for optional fields

Example:
```go
type TestResult struct {
    Name     string        `json:"name"`
    Status   Status        `json:"status"`
    Duration time.Duration `json:"duration"`
    Error    error         `json:"error,omitempty"`
}
```

### Concurrency
- Use mutexes for protecting shared state (RWMutex for read-heavy patterns)
- Keep critical sections as small as possible
- Use defer statements for unlock operations
- Consider using channels for communication between goroutines

Example:
```go
func (a *actor) WhoCan(abilities ...Ability) Actor {
    a.mutex.Lock()
    defer a.mutex.Unlock()

    a.abilities = append(a.abilities, abilities...)
    return a
}
```

### Testing Patterns
- Write table-driven tests when testing multiple scenarios
- Use testify/require for assertions that stop test execution (legacy approach)
- New TestContext API eliminates need for manual `require.NoError` calls
- Use descriptive test names that explain what is being tested
- Follow the arrange-act-assert pattern
- Use gomock for interface mocking when needed

Example (New TestContext API):
```go
func TestJSONPlaceholderBasics(t *testing.T) {
    test := serenity.NewSerenityTest(t)

    apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

    apiTester.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
    )
}
```

Example (Legacy with require):
```go
func TestJSONPlaceholderBasics(t *testing.T) {
    test := serenity.NewSerenityTest(t)

    actor := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

    err := actor.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
    require.NoError(t, err)  // Still needed for legacy mode
}
```

### Linting Configuration (.golangci.yml)
- **Line length**: 120 characters
- **Enabled linters**: errcheck, gosec, govet, ineffassign, misspell, staticcheck, unconvert, unused
- **Exclusions**: _test.go files, examples/ directory, generated files
- **Formatters**: gofmt, goimports with local prefix `github.com/nchursin/serenity-go`
- **Mock files**: Generated with `go generate` using gomock

## Project Structure

```
serenity-go/
├── serenity/
│   ├── core/              # Core interfaces and actor implementation
│   ├── abilities/         # Actor abilities (API, etc.)
│   │   ├── ability.go     # Base ability interface
│   │   └── api/           # HTTP API testing capabilities
│   ├── expectations/      # Assertion system and expectations
│   │   ├── ensure/        # Ensure-style assertions
│   │   ├── equals.go      # Equality expectations
│   │   ├── contains.go    # Contains expectations
│   │   ├── comparison.go  # Comparison expectations
│   │   └── collection.go   # Collection expectations
│   ├── testing/           # TestContext and testing utilities
│   │   ├── context.go     # New TestContext API
│   │   ├── actor.go       # Test-aware actor implementation
│   │   └── mocks/         # Generated mocks
│   └── reporting/         # Test reporting and output
├── examples/              # Usage examples and integration tests
│   ├── demo_new_api_test.go
│   ├── basic_test_new_api_test.go
│   └── ability/
├── go.mod                 # Go module definition
├── Makefile              # Build and development commands
├── .golangci.yml         # Linting configuration
├── .gitignore            # Git ignore patterns
└── README.md             # Project documentation
```

## Development Workflow

1. **Setup**: Run `make deps` to ensure dependencies are current
2. **Development**: Use `make fmt` and `make lint` frequently during coding. **Always** use TDD approach.
3. **Testing**: Run `make test-v` for verbose output during development
4. **Pre-commit**: Run `make check` (fmt-check, lint, test) before committing
5. **CI**: Use `make ci` for automated pipeline (fmt, lint, test)
6. **Mocks**: Regenerate mocks with `make mocks` when interfaces change

## Common Gotchas

- Always use the full module path for local imports
- Remember that generic type parameters use `T any` syntax
- Use TestContext API for new tests to avoid manual error handling
- Use testify/require for assertions that should stop test execution (legacy)
- Mutex usage patterns: RLock/RUnlock for read-heavy operations, Lock/Unlock for writes
- Error wrapping should use `%w` verb, not `%s`
- golangci-lint will automatically format imports with goimports
- Line length is enforced at 120 characters
- Examples directory is excluded from most linting rules
- Mock files should be regenerated when interface definitions change
- Use `go generate` or `make mocks` to regenerate all mocks at once
- Use convetional commit messages when asked to commit changes

## Dependencies

- **Go**: Version 1.23.4
- **Testify**: v1.11.1 for assertions and test utilities
- **Gomock**: v0.6.0 for interface mocking
- **Golangci-lint**: For code quality and formatting enforcement

Ensure all dependencies are up to date with `make deps` before development.
