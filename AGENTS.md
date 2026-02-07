# AGENTS.md - Guide for AI Coding Agents

This document provides essential information for AI agents working with the Serenity-Go codebase.

## Project Overview

Serenity-Go is a Go implementation of the Serenity/JS Screenplay Pattern for acceptance testing, focused on API testing capabilities. It provides actor-centric testing with reusable components and clear domain language.

## Build, Test, and Development Commands

### Testing Commands
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in specific package
go test ./serenity/core -v
go test ./serenity/api -v
go test ./serenity/assertions -v
go test ./examples -v

# Run a single test
go test -run TestSpecificFunction ./path/to/package

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. ./...
```

### Build Commands
```bash
# Build the module
go build ./...

# Build specific package
go build ./serenity/...

# Clean build cache
go clean -cache
```

### Code Quality Commands
```bash
# Format code (gofmt)
gofmt -s -w .

# Run go vet
go vet ./...

# Run staticcheck (if available)
staticcheck ./...

# Run ineffassign (if available)
ineffassign ./...

# Run misspell (if available)
misspell .

# Run golangci-lint (if configured)
golangci-lint run
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Update dependencies
go get -u ./...
```

## Code Style Guidelines

### Package and Import Organization
- Standard library imports grouped first, then third-party, then local imports
- Use blank lines between import groups
- Local imports use the full module path: `github.com/nchursin/serenity-go/serenity/...`

Example:
```go
import (
    "fmt"
    "sync"
    
    "github.com/stretchr/testify/require"
    
    "github.com/nchursin/serenity-go/serenity/core"
    "github.com/nchursin/serenity-go/serenity/api"
)
```

### Naming Conventions
- **Package names**: lowercase, single word when possible (e.g., `core`, `api`, `assertions`)
- **Public functions/types**: PascalCase (e.g., `NewActor`, `RequestBuilder`)
- **Private functions/types**: camelCase (e.g., `sendRequest`, `abilityTypeOf`)
- **Interfaces**: Often include type parameter for generics (e.g., `Question[T any]`, `Expectation[T any]`)
- **Constants**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase for both exported and unexported

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
```

### Error Handling
- Always wrap errors with context using `fmt.Errorf` with `%w` verb
- Return early from functions when errors occur
- Use descriptive error messages that include context
- Test error paths in unit tests

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

Example:
```go
// Fluent request building
req, err := api.Post("/posts").
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
- Use testify/require for assertions that stop test execution
- Use descriptive test names that explain what is being tested
- Follow the arrange-act-assert pattern

Example:
```go
func TestJSONPlaceholderBasics(t *testing.T) {
    actor := core.NewActor("APITester").WhoCan(api.UsingURL("https://jsonplaceholder.typicode.com"))

    err := actor.AttemptsTo(
        api.GetRequest("/posts"),
        assertions.That(api.LastResponseStatus{}, assertions.Equals(200)),
        assertions.That(api.LastResponseBody{}, assertions.Contains("title")),
    )
    require.NoError(t, err)
}
```

### Documentation
- Include package-level documentation explaining the purpose
- Document public types and functions with clear examples
- Use godoc format for comments
- Include context and error information in function documentation

## Project Structure

```
serenity-go/
├── serenity/
│   ├── core/           # Core interfaces and actor implementation
│   ├── api/            # HTTP API testing capabilities
│   ├── assertions/     # Assertion system and expectations
│   └── reporting/      # Test reporting and output
├── examples/           # Usage examples and integration tests
├── go.mod             # Go module definition
└── README.md          # Project documentation
```

## Screenplay Pattern Guidelines

### Actor-Based Testing
- Tests should describe what actors do, not how they do it
- Create actors with descriptive names (e.g., "APITester", "Admin", "RegularUser")
- Give actors only the abilities they need for their role

### Activity Composition
- Use interactions for low-level operations (HTTP requests, database queries)
- Create tasks for high-level business actions
- Build complex workflows by composing simple activities

### Question-Answer Pattern
- Use questions to retrieve information from the system
- Chain questions with assertions for validation
- Keep questions focused on single pieces of information

## Development Workflow

1. Write tests first (TDD approach when possible)
2. Run tests frequently during development
3. Use `go mod tidy` after adding new dependencies
4. Run `go vet` and `gofmt` before committing
5. Ensure all tests pass before creating pull requests

## Common Gotchas

- Always use the full module path for local imports
- Remember that generic type parameters use `T any` syntax
- Use testify/require for assertions that should stop test execution
- Mutex usage patterns: RLock/RUnlock for read-heavy operations, Lock/Unlock for writes
- Error wrapping should use `%w` verb, not `%s`

## Git Configuration

- Branch: working in `feature` or development branches
- Commit messages: Follow conventional commits format
- Language: Russian responses, English commit messages (as per general instructions)