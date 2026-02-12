# TestContext API Best Practices

This document demonstrates recommended patterns for using the TestContext API in Serenity/JS testing.

## Basic Usage Pattern

```go
func TestAPIExample(t *testing.T) {
    // ALWAYS use defer Shutdown() immediately after test creation
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    // Use descriptive actor names for better reporting
    apiUser := test.ActorCalled("APIUser").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Chain activities for logical grouping
    apiUser.AttemptsTo(
        api.SendPostRequest("/users").WithBody(map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        }),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("id")),
    )
}
```

## Concurrent Testing Pattern

```go
func TestConcurrentOperations(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    // Actors are thread-safe and can be shared across goroutines
    actor := test.ActorCalled("ConcurrentUser").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Use sync.WaitGroup for coordinated concurrent operations
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(requestID int) {
            defer wg.Done()
            actor.AttemptsTo(
                api.SendPostRequest("/jobs").WithBody(map[string]interface{}{
                    "id": requestID,
                }),
                ensure.That(api.LastResponseStatus{}, expectations.Equals(202)),
            )
        }(i)
    }
    wg.Wait()
}
```

## Error Handling Pattern

```go
func TestErrorScenarios(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    actor := test.ActorCalled("ErrorProneUser").WhoCan(
        api.CallAnApiAt("https://invalid-domain-that-does-not-exist.com"),
    )

    // No need for require.NoError - errors automatically fail the test
    actor.AttemptsTo(api.SendGetRequest("/endpoint"))
    // Test will fail here with automatic error reporting
}
```

## Multiple Actors Pattern

```go
func TestMultipleRoles(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    // Create specialized actors for different roles
    admin := test.ActorCalled("Admin").WhoCan(api.CallAnApiAt("https://api.example.com"))
    user := test.ActorCalled("User").WhoCan(api.CallAnApiAt("https://api.example.com"))

    // Admin creates resources
    admin.AttemptsTo(
        api.SendPostRequest("/users").WithBody(map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
        }),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
    )

    // User accesses resources
    user.AttemptsTo(
        api.SendGetRequest("/users/1"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
}
```

## Custom Reporting Pattern

```go
func TestWithCustomReporting(t *testing.T) {
    // Create custom reporter for advanced reporting needs
    reporter := &customReporter{
        outputFile:       "test-results.json",
        includeStackTrace: true,
    }

    test := serenity.NewSerenityTestWithReporter(t, reporter)
    defer test.Shutdown()

    actor := test.ActorCalled("ReportedUser").WhoCan(api.CallAnApiAt("https://api.example.com"))
    actor.AttemptsTo(api.SendGetRequest("/health"))
}

// Custom reporter implementation example
type customReporter struct {
    outputFile       string
    includeStackTrace bool
}

func (r *customReporter) OnTestStart(testName string) {
    // Custom test start logic
}

func (r *customReporter) OnTestFinish(result interface{}) {
    // Custom test finish logic - write to file, etc.
}

func (r *customReporter) OnActivityStart(activityName string, actorName string) {
    // Custom activity start logic
}

func (r *customReporter) OnActivityFinish(activityName string, actorName string, err error) {
    // Custom activity finish logic
}
```

## Key Principles

1. **Always call `defer test.Shutdown()` immediately after test creation**
2. **Use descriptive actor names** for better reporting
3. **Chain related activities** together for logical grouping
4. **Leverage automatic error handling** - no need for manual `require.NoError`
5. **Actors are thread-safe** - can be shared across goroutines
6. **Use `ensure.That`** for assertions in TestContext API

## Migration from Legacy API

### Before (Legacy)
```go
func TestLegacyApproach(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    actor := test.ActorCalled("APIUser").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Manual error handling required
    err := actor.AttemptsTo(
        api.SendGetRequest("/users"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
    require.NoError(t, err) // Manual error checking
}
```

### After (TestContext API)
```go
func TestNewApproach(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    actor := test.ActorCalled("APIUser").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Automatic error handling - no require.NoError needed
    actor.AttemptsTo(
        api.SendGetRequest("/users"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
}
```

## Benefits of TestContext API

- **Cleaner Tests**: No boilerplate error handling
- **Better Errors**: Automatic stack traces and context
- **Thread Safety**: Built-in concurrent testing support
- **Consistency**: Uniform error handling across all tests