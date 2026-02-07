# Serenity-Go: Screenplay Pattern Testing Framework for Go

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

    "github.com/nchursin/serenity-go/serenity/api"
    "github.com/nchursin/serenity-go/serenity/assertions"
    "github.com/nchursin/serenity-go/serenity/core"
)

func TestAPI(t *testing.T) {
    // Create an actor with API ability
    actor := core.NewActor("APITester").WhoCan(
        api.UsingURL("https://jsonplaceholder.typicode.com"),
    )

    // Define test data
    newPost := map[string]interface{}{
        "title":  "Test Post",
        "body":   "This is a test post",
        "userId": 1,
    }

    // Test the API
    err := actor.AttemptsTo(
        api.PostRequest("/posts").With(newPost),
        assertions.That(api.LastResponseStatus{}, assertions.Equals(201)),
        assertions.That(api.LastResponseBody{}, assertions.Contains("Test Post")),
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
actor = actor.WhoCan(api.UsingURL("https://api.example.com"))
```

### Abilities

Abilities enable actors to interact with different interfaces:

```go
// HTTP API ability
apiAbility := api.UsingURL("https://api.example.com")

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
// Send HTTP request
api.GetRequest("/users")
api.PostRequest("/posts").With(postData)
api.PutRequest("/users/1").With(updatedUser)
api.DeleteRequest("/posts/123")
```

#### Tasks (high-level business actions)
```go
// Define reusable task
createUserTask := core.Where(
    "creates a new user",
    api.PostRequest("/users").With(userData),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(201)),
)

// Use the task
actor.AttemptsTo(createUserTask)
```

### Questions

Questions retrieve information from the system:

```go
// Built-in questions
assertions.That(api.LastResponseStatus{}, assertions.Equals(200))
assertions.That(api.LastResponseBody{}, assertions.Contains("success"))
assertions.That(api.NewResponseHeader("content-type"), assertions.Contains("json"))
```

### Assertions

Verify that expectations are met:

```go
assertions.That(question, assertions.Equals(expected))
assertions.That(question, assertions.Contains(substring))
assertions.That(question, assertions.IsEmpty())
assertions.That(question, assertions.ArrayLengthEquals(5))
assertions.That(question, assertions.IsGreaterThan(10))
assertions.That(question, assertions.ContainsKey("id"))
```

## API Testing

### HTTP Requests

```go
// GET request
err := actor.AttemptsTo(
    api.GetRequest("/posts"),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(200)),
)

// POST request with JSON data
newPost := map[string]interface{}{
    "title":  "New Post",
    "body":   "Post content",
    "userId": 1,
}

err = actor.AttemptsTo(
    api.PostRequest("/posts").With(newPost),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(201)),
)

// PUT request with headers
err = actor.AttemptsTo(
    core.NewInteraction("updates a post", func(a core.Actor) error {
        req, err := api.Put("/posts/1").
            WithHeader("Authorization", "Bearer token").
            With(updatedData).
            Build()
        if err != nil {
            return err
        }
        return api.SendRequest(req).PerformAs(a)
    }),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(200)),
)

// DELETE request
err = actor.AttemptsTo(
    api.DeleteRequest("/posts/1"),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(200)),
)
```

### Request Building

```go
// Fluent request building
req, err := api.Post("/posts").
    WithHeader("Content-Type", "application/json").
    WithHeader("Authorization", "Bearer token").
    With(postData).
    Build()

if err != nil {
    return err
}

err = actor.AttemptsTo(api.SendRequest(req))
```

### Response Validation

```go
err := actor.AttemptsTo(
    api.GetRequest("/posts/1"),
    assertions.That(api.LastResponseStatus{}, assertions.Equals(200)),
    assertions.That(api.LastResponseBody{}, assertions.Contains("title")),
    assertions.That(api.NewResponseHeader("content-type"), assertions.Contains("json")),
)
```

## Working Examples

The `examples/` directory contains working examples with real APIs:

- `basic_test.go` - Core functionality demonstrations
- `jsonplaceholder_test.go` - CRUD operations with JSONPlaceholder API

Run examples:

```bash
go test ./examples -v
```

## Architecture

### Core Components

- **serenity/core/** - Screenplay Pattern interfaces (Actor, Activity, Question, Task)
- **serenity/api/** - HTTP API testing capabilities
- **serenity/assertions/** - Assertion system and expectations
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
customInteraction := core.NewInteraction("performs custom action", func(actor core.Actor) error {
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

assertions.That(customQuestion, assertions.Equals(42))
```

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
| `WhoCan(CallAnAPI.using(...))` | `WhoCan(api.UsingURL(...))` |
| `attemptsTo(Send.a(...))` | `AttemptsTo(api.PostRequest(...))` |
| `Ensure.that(LastResponse.status(), equals(200))` | `assertions.That(api.LastResponseStatus{}, assertions.Equals(200))` |

## Development Status

This is an MVP implementation focused on API testing. Future enhancements may include:

- [ ] Database testing abilities
- [ ] gRPC testing support
- [ ] Advanced reporting capabilities
- [ ] Integration with popular Go test frameworks
- [ ] Web UI testing capabilities

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License
Apache 2.0 - see LICENSE file for details.
