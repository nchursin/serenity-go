# Serenity-Go: Screenplay Pattern Testing Framework for Go

![CI](https://github.com/nchursin/serenity-go/workflows/CI/badge.svg) ![codecov](https://codecov.io/gh/nchursin/serenity-go/graph/badge.svg) ![Version](https://img.shields.io/github/v/release/nchursin/serenity-go)

A Go implementation of the Serenity/JS Screenplay Pattern for acceptance testing, focused on API testing capabilities.

## Overview

Serenity-Go brings the power of the Screenplay Pattern to Go testing, providing:

- **Actor-centric testing** - Tests describe what actors do, not how they do it
- **Reusable components** - Build a library of reusable tasks and interactions
- **Clear domain language** - Tests that read like business requirements
- **Modular design** - Use only what you need for your testing scenarios
- **Framework agnostic** - Works with any Go test runner

## Quick Start

### Installation

```bash
go get github.com/nchursin/serenity-go
```

### Basic Example

```go
package main

import (
    "testing"
    "github.com/stretchr/testify/require"

    "github.com/nchursin/serenity-go/serenity/abilities/api"
    "github.com/nchursin/serenity-go/serenity/core"
    "github.com/nchursin/serenity-go/serenity/expectations"
    "github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

func TestAPI(t *testing.T) {
    // Create an actor with API ability
    actor := core.NewActor("APITester").WhoCan(
        api.CallAnApiAt("https://jsonplaceholder.typicode.com"),
    )

    // Define test data
    newPost := map[string]interface{}{
        "title":  "Test Post",
        "body":   "This is a test post",
        "userId": 1,
    }

    // Test the API
    err := actor.AttemptsTo(
        api.SendPostRequest("/posts").
            WithBody(newPost),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("Test Post")),
    )

    require.NoError(t, err)
}
```

## Core Concepts

### Actors

Actors represent people or systems interacting with your application:

```go
// Create an actor
actor := core.NewActor("John Doe")

// Give the actor abilities to interact with your system
actor = actor.WhoCan(api.CallAnApiAt("https://api.example.com"))
```

### Abilities

Abilities enable actors to interact with different interfaces:

```go
// HTTP API ability
apiAbility := api.CallAnApiAt("https://api.example.com")

// Actor with multiple abilities
actor := core.NewActor("TestUser").WhoCan(
    apiAbility,
    // ... other abilities
)
```

### Activities

Activities represent actions that actors perform:

#### Interactions (low-level actions)
```go
// Send HTTP requests with fluent interface
api.SendGetRequest("/users")
api.SendPostRequest("/posts").WithBody(postData)
api.SendPutRequest("/users/1").WithBody(updatedUser)
api.SendDeleteRequest("/posts/123")
```

#### Tasks (high-level business actions)
```go
// Define reusable task
createUserTask := core.Where(
    "creates a new user",
    core.Do("creates a new user", func(a core.Actor) error {
        req, err := api.Post("/users").
            With(userData).
            Build()
        if err != nil {
            return err
        }
        return api.SendRequest(req).PerformAs(a)
    }),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
)

// Use the task
actor.AttemptsTo(createUserTask)
```

### Questions

Questions retrieve information from the system:

```go
// Built-in questions
ensure.That(api.LastResponseStatus{}, expectations.Equals(200))
ensure.That(api.LastResponseBody{}, expectations.Contains("success"))
ensure.That(api.NewResponseHeader("content-type"), expectations.Contains("json"))
```

### Assertions

Verify that expectations are met:

```go
ensure.That(question, expectations.Equals(expected))
ensure.That(question, expectations.Contains(substring))
ensure.That(question, expectations.IsEmpty())
ensure.That(question, expectations.ArrayLengthEquals(5))
ensure.That(question, expectations.IsGreaterThan(10))
ensure.That(question, expectations.ContainsKey("id"))

// Custom validation with Satisfies
ensure.That(answerable.ValueOf(value), expectations.Satisfies("custom description", func(actual T) error {
    // Your validation logic here
    return nil // or error with description
}))
```

## API Testing

### HTTP Requests

```go
// GET request
err := actor.AttemptsTo(
    api.SendGetRequest("/posts"),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
)

// POST request with JSON data
newPost := map[string]interface{}{
    "title":  "New Post",
    "body":   "Post content",
    "userId": 1,
}

err = actor.AttemptsTo(
    api.SendPostRequest("/posts").
        WithBody(newPost),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
)

// PUT request with headers
err = actor.AttemptsTo(
    api.SendPutRequest("/posts/1").
        WithHeader("Authorization", "Bearer token").
        WithBody(updatedData),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
)

// DELETE request
err = actor.AttemptsTo(
    api.SendDeleteRequest("/posts/1"),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
)
```

### Request Building

```go
// Fluent request building
err = actor.AttemptsTo(
    api.SendPostRequest("/posts").
        WithHeader("Content-Type", "application/json").
        WithHeader("Authorization", "Bearer token").
        WithBody(postData),
)
```

### Response Validation

```go
err := actor.AttemptsTo(
    api.SendGetRequest("/posts/1").
        WithHeader("Accept", "application/json"),
    ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
    ensure.That(api.NewResponseHeader("content-type"), expectations.Contains("json")),
)
```

## Working Examples

The `examples/` directory contains working examples with real APIs:

- `basic_test.go` - JSONPlaceholder API testing examples including basic operations, error scenarios, and multiple actors
- `jsonplaceholder_test.go` - Additional JSONPlaceholder API examples with CRUD operations
- `satisfies_demo_test.go` - Comprehensive examples of custom `Satisfies` expectations including go-cmp integration

Run examples:

```bash
go test ./examples -v
```

For detailed `Satisfies` examples, see [docs/SATISFIES_EXAMPLES.md](docs/SATISFIES_EXAMPLES.md).

## Architecture

### Core Components

- **serenity/core/** - Screenplay Pattern interfaces (Actor, Activity, Question, Task)
- **serenity/abilities/api/** - HTTP API testing capabilities
- **serenity/expectations/** - Assertion system and expectations
- **serenity/expectations/ensure/** - Ensure-style assertions
- **serenity/reporting/** - Test reporting and output

### Design Principles

1. **Composable** - Build complex behaviors from simple components
2. **Reusable** - Create libraries of tasks and interactions
3. **Readable** - Tests that read like business specifications
4. **Extensible** - Add new abilities and integrations
5. **Type-safe** - Leverage Go's type system for safety

## Advanced Usage

### Custom Interactions

```go
customInteraction := core.Do("performs custom action", func(actor core.Actor) error {
    // Your custom logic here
    return nil
})

actor.AttemptsTo(customInteraction)
```

### Custom Questions

```go
customQuestion := core.Of[int]("asks for custom value", func(actor core.Actor) (int, error) {
    // Your custom logic here
    return 42, nil
})

ensure.That(customQuestion, expectations.Equals(42))
```

### Custom Expectations with Satisfies

Create custom expectations using the `Satisfies` function for complex validation logic:

```go
// Simple custom validation
actor.AttemptsTo(
    ensure.That(answerable.ValueOf(age), expectations.Satisfies("is positive number", func(actual int) error {
        if actual <= 0 {
            return fmt.Errorf("expected positive value, but got %d", actual)
        }
        return nil
    })),
)

// Advanced struct comparison with go-cmp
actor.AttemptsTo(
    ensure.That(answerable.ValueOf(actualUser), expectations.Satisfies("matches expected user", func(actual User) error {
        if diff := cmp.Diff(expectedUser, actual); diff != "" {
            return fmt.Errorf("user mismatch:\n%s", diff)
        }
        return nil
    })),
)

// Complex business logic validation
actor.AttemptsTo(
    ensure.That(answerable.ValueOf(order), expectations.Satisfies("has valid order data", func(actual Order) error {
        if !strings.HasPrefix(actual.ID, "ORD-") {
            return fmt.Errorf("order ID must start with ORD-, got %s", actual.ID)
        }
        if actual.Amount <= 0 {
            return fmt.Errorf("order amount must be positive, got %f", actual.Amount)
        }
        // Add more validation logic...
        return nil
    })),
)
```

The `Satisfies` function takes:
- A description string that appears in test failure messages
- A validation function that returns `nil` for success or an error for failure

This enables powerful, type-safe custom validations while maintaining the Screenplay Pattern's readable test structure.

### Task Composition

```go
// Build complex workflows from simple tasks
setupTask := core.Where("setup test data", setupDataAction)
testTask := core.Where("run test scenario", testAction)
cleanupTask := core.Where("cleanup test data", cleanupAction)

actor.AttemptsTo(
    setupTask,
    testTask,
    cleanupTask,
)
```

### Multiple Actors

```go
admin := core.NewActor("Admin").WhoCan(api.UsingURL(baseURL))
user := core.NewActor("RegularUser").WhoCan(api.UsingURL(baseURL))

// Admin creates resources
err := admin.AttemptsTo(createResourceTask)

// User interacts with resources
err = user.AttemptsTo(accessResourceTask)
```

## Comparison with Serenity/JS

This Go implementation follows the same design principles as Serenity/JS:

| Serenity/JS | Serenity-Go |
|--------------|-------------|
| `actorCalled('John')` | `core.NewActor("John")` |
| `WhoCan(CallAnAPI.using(...))` | `WhoCan(api.CallAnApiAt(...))` |
| `attemptsTo(Send.a(...))` | `AttemptsTo(api.SendGetRequest(...))` |
| `Ensure.that(LastResponse.status(), equals(200))` | `ensure.That(api.LastResponseStatus{}, expectations.Equals(200))` |

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License
Apache 2.0 - see LICENSE file for details.

