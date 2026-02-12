# TestContext API Migration Guide

This guide helps migrate from the legacy Serenity API to the new TestContext API.

## Quick Comparison

### Legacy API
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

### TestContext API
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

## Migration Steps

1. **Replace `require.NoError(t, err)` calls with TestContext API**
   - Remove manual error handling from test logic
   - Let TestContext automatically handle failures

2. **Remove manual error handling from test logic**
   - Delete `require.NoError(t, err)` calls
   - Delete `if err != nil { return }` patterns
   - Trust TestContext to handle errors appropriately

3. **Use `defer test.Shutdown()` immediately after test creation**
   - This pattern remains the same in both APIs
   - Ensures proper resource cleanup

4. **Keep existing actor and ability definitions unchanged**
   - `ActorCalled()` method works the same way
   - `WhoCan()` chaining remains identical
   - Ability definitions are unchanged

## Key Differences

| Feature | Legacy API | TestContext API |
|---------|-------------|-----------------|
| Error Handling | Manual (`require.NoError`) | Automatic (built-in) |
| Code Verbosity | More verbose | Cleaner, focused |
| Error Messages | Basic require errors | Rich context with actor names |
| Test Failures | May require custom handling | Immediate with stack traces |
| Reporting | Optional | Integrated |

## Before and After Examples

### Complex Workflow - Before
```go
func TestComplexWorkflowLegacy(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    actor := test.ActorCalled("WorkflowActor").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Multiple manual error checks
    err := actor.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
    require.NoError(t, err)

    err = actor.AttemptsTo(
        api.SendPostRequest("/posts").WithBody(map[string]interface{}{
            "title":  "Test Post",
            "body":   "Test body",
            "userId": 1,
        }),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
    )
    require.NoError(t, err)

    err = actor.AttemptsTo(
        ensure.That(api.LastResponseBody{}, expectations.Contains("Test Post")),
    )
    require.NoError(t, err)
}
```

### Complex Workflow - After
```go
func TestComplexWorkflowNew(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    actor := test.ActorCalled("WorkflowActor").WhoCan(
        api.CallAnApiAt("https://api.example.com"),
    )

    // Clean, focused on test logic
    actor.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),

        api.SendPostRequest("/posts").WithBody(map[string]interface{}{
            "title":  "Test Post",
            "body":   "Test body",
            "userId": 1,
        }),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("Test Post")),
    )
    // No manual error handling needed!
}
```

## Benefits of TestContext API

### **Cleaner Tests**
- No boilerplate error handling code
- Focus on actual test logic
- Reduced cognitive load

### **Better Errors**
- Automatic stack traces with context
- Actor names included in error messages  
- Rich failure information

### **Thread Safety**
- Built-in concurrent testing support
- Safe sharing of actors across goroutines
- Consistent behavior in parallel scenarios

### **Consistency**
- Uniform error handling across all tests
- Standardized failure reporting
- Predictable test behavior

## Common Migration Patterns

### Error Handling Replacement
```go
// Before (Legacy)
if err := actor.AttemptsTo(activity); err != nil {
    t.Errorf("Activity failed: %v", err)
    return
}

// After (TestContext)  
actor.AttemptsTo(activity) // Errors automatically handled
```

### Multiple Activities
```go
// Before (Legacy)
err1 := actor.AttemptsTo(activity1)
require.NoError(t, err1)

err2 := actor.AttemptsTo(activity2)  
require.NoError(t, err2)

// After (TestContext)
actor.AttemptsTo(
    activity1,
    activity2,
) // All errors handled automatically
```

### Concurrent Testing
```go
// Before (Legacy) - Complex error coordination
var wg sync.WaitGroup
var errors []error
var mu sync.Mutex

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        if err := actor.AttemptsTo(activity); err != nil {
            mu.Lock()
            errors = append(errors, err)
            mu.Unlock()
        }
    }(i)
}
wg.Wait()
// Manual error aggregation and checking...

// After (TestContext) - Automatic error coordination
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        actor.AttemptsTo(activity) // Errors automatically fail test
    }(i)
}
wg.Wait()
```

## Troubleshooting

### Test Still Passes When It Should Fail
- Ensure you're using `test.ActorCalled()` not a different actor creation method
- Verify `defer test.Shutdown()` is called immediately after test creation
- Check that you're not catching and swallowing errors

### Missing Error Context
- Make sure actor names are descriptive for better error messages
- Use `ensure.That()` for assertions instead of manual checks
- Let TestContext handle all activity errors

### Performance Issues
- TestContext API has minimal overhead
- Profile if you suspect performance regressions
- Most issues come from test setup, not error handling

## Rollback Strategy

If you need to rollback to legacy API during migration:

1. Keep the same test structure
2. Add back `require.NoError(t, err)` after each `AttemptsTo` call  
3. Handle errors manually as needed
4. Remove automatic error handling expectations

The TestContext and legacy APIs can coexist in the same codebase, allowing gradual migration.