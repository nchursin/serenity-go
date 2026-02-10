# ConsoleReporter Integration Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Integrate ConsoleReporter into the main Serenity-Go test flow to provide detailed console output of test execution with step-by-step activity tracking.

**Architecture:** Add optional reporter support to SerenityTest that wraps activities with ActivityTracker for automatic console reporting during test execution.

**Tech Stack:** Go 1.23.4, Serenity-Go testing framework, ConsoleReporter, ActivityTracker, existing TestContext API.

---

### Task 1: Extend SerenityTest interface to support Reporter

**Files:**
- Modify: `serenity/testing/serenity_test_manager.go:10-17`
- Test: `serenity/testing/serenity_test_manager_test.go`

**Step 1: Write the failing test**

```go
func TestSerenityTestWithConsoleReporter(t *testing.T) {
    // Create a SerenityTest with console reporter
    test := NewSerenityTestWithReporter(t, nil)
    defer test.Shutdown()
    
    actor := test.ActorCalled("TestActor")
    require.NotNil(t, actor)
    
    // Verify that reporter is configured
    adapter := test.GetReporterAdapter()
    require.NotNil(t, adapter)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./serenity/testing -v -run TestSerenityTestWithConsoleReporter`
Expected: FAIL with "NewSerenityTestWithReporter not defined"

**Step 3: Write minimal implementation**

```go
// ReporterProvider provides access to the reporter adapter
type ReporterProvider interface {
    GetReporterAdapter() *reporting.TestRunnerAdapter
}

// Extend SerenityTest interface
type SerenityTest interface {
    ActorCalled(name string) core.Actor
    Shutdown()
    GetReporterAdapter() *reporting.TestRunnerAdapter // New method
}

// NewSerenityTestWithReporter creates a SerenityTest with optional reporter
func NewSerenityTestWithReporter(t TestContext, reporter reporting.Reporter) SerenityTest {
    var adapter *reporting.TestRunnerAdapter
    if reporter != nil {
        adapter = reporting.NewTestRunnerAdapter(reporter)
    }
    
    return &serenityTest{
        ctx:     t,
        actors:  make(map[string]core.Actor),
        adapter: adapter, // Add adapter field
    }
}
```

**Step 4: Update serenityTest struct**

```go
type serenityTest struct {
    ctx     TestContext
    actors  map[string]core.Actor
    mutex   sync.RWMutex
    adapter *reporting.TestRunnerAdapter // New field
}
```

**Step 5: Add GetReporterAdapter method**

```go
func (st *serenityTest) GetReporterAdapter() *reporting.TestRunnerAdapter {
    return st.adapter
}
```

**Step 6: Run test to verify it passes**

Run: `go test ./serenity/testing -v -run TestSerenityTestWithConsoleReporter`
Expected: PASS

**Step 7: Commit**

```bash
git add serenity/testing/serenity_test_manager.go serenity/testing/serenity_test_manager_test.go
git commit -m "feat: extend SerenityTest to support optional reporter"
```

---

### Task 2: Update NewSerenityTest to use ConsoleReporter by default

**Files:**
- Modify: `serenity/testing/serenity_test_manager.go:26-32`
- Test: `serenity/testing/serenity_test_manager_test.go`

**Step 1: Write the failing test**

```go
func TestNewSerenityTestUsesConsoleReporter(t *testing.T) {
    test := NewSerenityTest(t)
    defer test.Shutdown()
    
    adapter := test.GetReporterAdapter()
    require.NotNil(t, adapter)
    
    // Verify it's a ConsoleReporter
    reporter := adapter.GetReporter()
    _, isConsole := reporter.(*console_reporter.ConsoleReporter)
    require.True(t, isConsole, "Expected ConsoleReporter")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./serenity/testing -v -run TestNewSerenityTestUsesConsoleReporter`
Expected: FAIL - adapter is nil

**Step 3: Update NewSerenityTest to use ConsoleReporter**

```go
import (
    "github.com/nchursin/serenity-go/serenity/reporting"
    "github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
)

// NewSerenityTest creates a new SerenityTest instance with ConsoleReporter
func NewSerenityTest(t TestContext) SerenityTest {
    reporter := console_reporter.NewConsoleReporter()
    return NewSerenityTestWithReporter(t, reporter)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./serenity/testing -v -run TestNewSerenityTestUsesConsoleReporter`
Expected: PASS

**Step 5: Commit**

```bash
git add serenity/testing/serenity_test_manager.go serenity/testing/serenity_test_manager_test.go
git commit -m "feat: make SerenityTest use ConsoleReporter by default"
```

---

### Task 3: Integrate ActivityTracker into testActor.AttemptsTo

**Files:**
- Modify: `serenity/testing/actor.go:47-64`
- Test: `serenity/testing/actor_test.go`

**Step 1: Write the failing test**

```go
func TestTestActorAttemptsToWithReporting(t *testing.T) {
    mockReporter := mocks.NewMockReporter(gomock.NewController(t))
    mockTestContext := mocks.NewMockTestContext(gomock.NewController(t))
    
    // Expect OnStepStart for first activity
    mockReporter.EXPECT().OnStepStart("Send GET request to /posts")
    mockReporter.EXPECT().OnStepFinish(gomock.Any()).Times(1)
    
    // Expect no error from test context
    mockTestContext.EXPECT().Failed().Return(false).AnyTimes()
    
    // Create test actor with mock reporter
    adapter := reporting.NewTestRunnerAdapter(mockReporter)
    test := &serenityTest{
        ctx:     mockTestContext,
        actors:  make(map[string]core.Actor),
        adapter: adapter,
    }
    
    actor := test.ActorCalled("TestActor")
    
    // Mock activity
    mockActivity := mocks.NewMockActivity(gomock.NewController(t))
    mockActivity.EXPECT().PerformAs(gomock.Any()).Return(nil)
    mockActivity.EXPECT().Description().Return("Send GET request to /posts")
    mockActivity.EXPECT().FailureMode().Return(core.FailFast)
    
    actor.AttemptsTo(mockActivity)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./serenity/testing -v -run TestTestActorAttemptsToWithReporting`
Expected: FAIL - no reporter integration

**Step 3: Update testActor to support reporting**

```go
// testActor is a wrapper around core.Actor that includes test context
type testActor struct {
    name        string
    abilities   []abilities.Ability
    testContext TestContext
    mutex       sync.RWMutex
    reporter    *reporting.TestRunnerAdapter // New field
}
```

**Step 4: Update serenityTest.ActorCalled to pass reporter**

```go
// Create new actor with test context and reporter
actor := &testActor{
    name:        name,
    abilities:   make([]abilities.Ability, 0),
    testContext: st.ctx,
    reporter:    st.adapter, // Pass reporter
}
```

**Step 5: Update AttemptsTo to use ActivityTracker**

```go
// AttemptsTo executes activities with error handling through test context and reporting
func (ta *testActor) AttemptsTo(activities ...core.Activity) {
    for _, activity := range activities {
        var tracker *reporting.ActivityTracker
        
        // Create activity tracker if reporter is available
        if ta.reporter != nil {
            tracker = reporting.NewActivityTracker(ta.reporter.GetReporter(), activity.Description())
            tracker.Start()
        }
        
        err := activity.PerformAs(ta)
        
        // Finish tracking with result
        if tracker != nil {
            tracker.Finish(err)
        }
        
        if err != nil {
            switch activity.FailureMode() {
            case core.FailFast:
                ta.testContext.Errorf("Critical activity error '%s' failed: %v", activity.Description(), err)
                ta.testContext.FailNow()
                return
            case core.ErrorButContinue:
                ta.testContext.Errorf("Non-critical activity error '%s' failed: %v", activity.Description(), err)
            case core.Ignore:
                ta.testContext.Logf("Ignore activity error '%s' failed: %v", activity.Description(), err)
            }
        }
    }
}
```

**Step 6: Run test to verify it passes**

Run: `go test ./serenity/testing -v -run TestTestActorAttemptsToWithReporting`
Expected: PASS

**Step 7: Commit**

```bash
git add serenity/testing/actor.go serenity/testing/actor_test.go
git commit -m "feat: integrate ActivityTracker into testActor.AttemptsTo"
```

---

### Task 4: Add test lifecycle reporting to SerenityTest

**Files:**
- Modify: `serenity/testing/serenity_test_manager.go:26-32`, `serenity/testing/serenity_test_manager.go:63-70`
- Test: `serenity/testing/serenity_test_manager_test.go`

**Step 1: Write the failing test**

```go
func TestSerenityTestLifecycleReporting(t *testing.T) {
    mockReporter := mocks.NewMockReporter(gomock.NewController(t))
    mockTestContext := mocks.NewMockTestContext(gomock.NewController(t))
    
    // Expect test lifecycle events
    mockReporter.EXPECT().OnTestStart("TestExample")
    mockReporter.EXPECT().OnTestFinish(gomock.Any()).Do(func(result reporting.TestResult) {
        require.Equal(t, "TestExample", result.Name())
        require.Equal(t, reporting.StatusPassed, result.Status())
    })
    
    mockTestContext.EXPECT().Failed().Return(false)
    
    test := NewSerenityTestWithReporter(mockTestContext, mockReporter)
    
    // Simulate test end
    test.Shutdown()
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./serenity/testing -v -run TestSerenityTestLifecycleReporting`
Expected: FAIL - no lifecycle reporting

**Step 3: Update NewSerenityTestWithReporter to start test reporting**

```go
// NewSerenityTestWithReporter creates a SerenityTest with optional reporter
func NewSerenityTestWithReporter(t TestContext, reporter reporting.Reporter) SerenityTest {
    var adapter *reporting.TestRunnerAdapter
    var testName string
    
    // Extract test name from TestContext if possible
    if tb, ok := t.(interface{ Name() string }); ok {
        testName = tb.Name()
    } else {
        testName = "Unknown"
    }
    
    if reporter != nil {
        adapter = reporting.NewTestRunnerAdapter(reporter)
        reporter.OnTestStart(testName)
    }
    
    return &serenityTest{
        ctx:       t,
        actors:    make(map[string]core.Actor),
        adapter:   adapter,
        startTime: time.Now(), // Add start time tracking
        testName:  testName,   // Add test name field
    }
}
```

**Step 4: Update serenityTest struct**

```go
type serenityTest struct {
    ctx       TestContext
    actors    map[string]core.Actor
    mutex     sync.RWMutex
    adapter   *reporting.TestRunnerAdapter
    startTime time.Time // New field
    testName  string    // New field
}
```

**Step 5: Update Shutdown to report test completion**

```go
import (
    "time"
)

// Shutdown cleans up resources and reports test completion
func (st *serenityTest) Shutdown() {
    st.mutex.Lock()
    defer st.mutex.Unlock()
    
    // Report test completion if reporter is available
    if st.adapter != nil {
        status := reporting.StatusPassed
        if st.ctx.Failed() {
            status = reporting.StatusFailed
        }
        
        result := &reporting.TestResult{
            Name:     st.testName,
            Status:   status,
            Duration: time.Since(st.startTime).Seconds(),
            Error:    nil, // TestContext doesn't expose error directly
        }
        
        st.adapter.GetReporter().OnTestFinish(result)
    }
    
    // Clear actors map
    st.actors = make(map[string]core.Actor)
}
```

**Step 6: Fix TestResult implementation issue**

Create a helper struct in serenity_test_manager.go:

```go
// testResult implements TestResult interface for reporting
type testResult struct {
    name     string
    status   reporting.Status
    duration float64
    error    error
}

func (tr *testResult) Name() string      { return tr.name }
func (tr *testResult) Status() reporting.Status { return tr.status }
func (tr *testResult) Duration() float64 { return tr.duration }
func (tr *testResult) Error() error      { return tr.error }
```

**Step 7: Update Shutdown to use local testResult**

```go
result := &testResult{
    Name:     st.testName,
    Status:   status,
    Duration: time.Since(st.startTime).Seconds(),
    Error:    nil,
}
```

**Step 8: Run test to verify it passes**

Run: `go test ./serenity/testing -v -run TestSerenityTestLifecycleReporting`
Expected: PASS

**Step 9: Commit**

```bash
git add serenity/testing/serenity_test_manager.go serenity/testing/serenity_test_manager_test.go
git commit -m "feat: add test lifecycle reporting to SerenityTest"
```

---

### Task 5: Update existing examples to demonstrate reporting

**Files:**
- Modify: `examples/basic_test_new_api_test.go:12-32`
- Test: Run the updated examples

**Step 1: Update basic example to show reporting in action**

```go
// TestJSONPlaceholderBasicsNewAPI demonstrates basic API testing with JSONPlaceholder using new TestContext API
func TestJSONPlaceholderBasicsNewAPI(t *testing.T) {
    test := serenity.NewSerenityTest(t)
    defer test.Shutdown()

    apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

    // Test GET posts - should return existing posts
    apiTester.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
    )

    // Test GET users - should return existing users
    apiTester.AttemptsTo(
        api.SendGetRequest("/users"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("email")),
    )
    
    // The console output will now show detailed step-by-step execution
}
```

**Step 2: Create a new example demonstrating reporting features**

Create: `examples/reporting_demo_test.go`

```go
package examples

import (
    "testing"
    "os"
    
    "github.com/nchursin/serenity-go/serenity/abilities/api"
    "github.com/nchursin/serenity-go/serenity/expectations"
    "github.com/nchursin/serenity-go/serenity/expectations/ensure"
    "github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
    serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestConsoleReportingDemo demonstrates console reporting features
func TestConsoleReportingDemo(t *testing.T) {
    // Create custom console reporter with different output
    reporter := console_reporter.NewConsoleReporter()
    
    test := serenity.NewSerenityTestWithReporter(t, reporter)
    defer test.Shutdown()

    apiTester := test.ActorCalled("DemoAPITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

    // This will show detailed console output with emojis and timing
    apiTester.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
    )
    
    // Another sequence to demonstrate step tracking
    apiTester.AttemptsTo(
        api.SendGetRequest("/posts/1"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
        ensure.That(api.LastResponseBody{}, expectations.Contains("userId")),
    )
}

// TestReportingToFile demonstrates outputting report to file
func TestReportingToFile(t *testing.T) {
    // Create file for output
    file, err := os.Create("test_report.txt")
    require.NoError(t, err)
    defer file.Close()
    
    // Create reporter that writes to file
    reporter := console_reporter.NewConsoleReporter()
    reporter.SetOutput(file)
    
    test := serenity.NewSerenityTestWithReporter(t, reporter)
    defer test.Shutdown()

    apiTester := test.ActorCalled("FileReporter").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))
    
    apiTester.AttemptsTo(
        api.SendGetRequest("/users"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
}
```

**Step 3: Run tests to see console output**

Run: `go test ./examples -v -run TestJSONPlaceholderBasicsNewAPI`
Expected: PASS with detailed console output showing steps

Run: `go test ./examples -v -run TestConsoleReportingDemo`
Expected: PASS with console reporting demonstration

**Step 4: Check file output**

Run: `go test ./examples -v -run TestReportingToFile && cat test_report.txt`
Expected: File contains formatted test report

**Step 5: Commit**

```bash
git add examples/basic_test_new_api_test.go examples/reporting_demo_test.go test_report.txt
git commit -m "feat: update examples to demonstrate console reporting"
```

---

### Task 6: Integration testing and validation

**Files:**
- Create: `serenity/testing/integration_test.go`
- Test: `make test`

**Step 1: Write comprehensive integration test**

```go
package testing

import (
    "bytes"
    "strings"
    "testing"
    
    "github.com/nchursin/serenity-go/serenity/abilities/api"
    "github.com/nchursin/serenity-go/serenity/expectations"
    "github.com/nchursin/serenity-go/serenity/expectations/ensure"
    "github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
    "github.com/stretchr/testify/require"
)

// TestReportingIntegration tests the complete reporting flow
func TestReportingIntegration(t *testing.T) {
    // Capture output
    var buf bytes.Buffer
    
    // Create reporter with captured output
    reporter := console_reporter.NewConsoleReporter()
    reporter.SetOutput(&buf)
    
    test := NewSerenityTestWithReporter(t, reporter)
    defer test.Shutdown()

    actor := test.ActorCalled("IntegrationTester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))
    
    // Perform activities
    actor.AttemptsTo(
        api.SendGetRequest("/posts"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
    
    // Check captured output
    output := buf.String()
    require.Contains(t, output, "Starting:")
    require.Contains(t, output, "IntegrationTester")
    require.Contains(t, output, "Send GET request")
    require.Contains(t, output, "PASSED")
}

// TestErrorReporting tests error reporting in activities
func TestErrorReporting(t *testing.T) {
    var buf bytes.Buffer
    
    reporter := console_reporter.NewConsoleReporter()
    reporter.SetOutput(&buf)
    
    test := NewSerenityTestWithReporter(t, reporter)
    defer test.Shutdown()

    actor := test.ActorCalled("ErrorTester")
    
    // Simulate failing activity
    mockActivity := &mockFailingActivity{}
    actor.AttemptsTo(mockActivity)
    
    output := buf.String()
    require.Contains(t, output, "FAILED")
    require.Contains(t, output, "Mock activity failed")
}

// mockFailingActivity simulates a failing activity
type mockFailingActivity struct{}

func (m *mockFailingActivity) PerformAs(actor core.Actor) error {
    return fmt.Errorf("Mock activity failed")
}

func (m *mockFailingActivity) Description() string {
    return "Mock failing activity"
}

func (m *mockFailingActivity) FailureMode() core.FailureMode {
    return core.FailFast
}
```

**Step 2: Run integration tests**

Run: `go test ./serenity/testing -v -run TestReportingIntegration`
Expected: PASS

Run: `go test ./serenity/testing -v -run TestErrorReporting`
Expected: PASS

**Step 3: Run full test suite**

Run: `make test`
Expected: All tests pass

**Step 4: Run linting**

Run: `make lint`
Expected: No linting errors

**Step 5: Run with coverage**

Run: `make test-coverage`
Expected: Good coverage on new reporting code

**Step 6: Commit**

```bash
git add serenity/testing/integration_test.go
git commit -m "feat: add comprehensive integration tests for reporting"
```

---

### Task 7: Documentation updates

**Files:**
- Create: `docs/reporting.md`
- Modify: `README.md`

**Step 1: Create reporting documentation**

Create: `docs/reporting.md`

```markdown
# Console Reporting in Serenity-Go

## Overview

Serenity-Go provides built-in console reporting that shows detailed execution of tests and activities with emojis, timing, and structured output.

## Usage

### Basic Usage

```go
func TestMyAPI(t *testing.T) {
    test := serenity.NewSerenityTest(t)  // Uses ConsoleReporter by default
    defer test.Shutdown()

    actor := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://api.example.com"))
    
    actor.AttemptsTo(
        api.SendGetRequest("/users"),
        ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
    )
}
```

### Custom Reporter Configuration

```go
// Create custom reporter
reporter := console_reporter.NewConsoleReporter()
test := serenity.NewSerenityTestWithReporter(t, reporter)
defer test.Shutdown()
```

### File Output

```go
file, _ := os.Create("test-report.txt")
reporter := console_reporter.NewConsoleReporter()
reporter.SetOutput(file)
```

## Output Format

The console reporter provides:
- 🚀 Test start indicators
- 🔄 Activity start indicators  
- ✅/❌ Activity completion status
- ⏱️ Execution timing
- 📝 Error details when failures occur

## Integration

The reporting system automatically integrates with:
- TestContext API
- Activity execution
- Error handling
- Test lifecycle management
```

**Step 2: Update README**

Add section to README.md about reporting features.

**Step 3: Commit**

```bash
git add docs/reporting.md README.md
git commit -m "docs: add documentation for console reporting features"
```

---

## Verification Commands

After completing all tasks:

```bash
# Run all tests
make test

# Run with verbose output to see reporting
go test ./examples -v -run TestConsoleReportingDemo

# Check code quality
make lint

# Generate coverage
make test-coverage

# Verify integration
go test ./serenity/testing -v -run TestReportingIntegration
```

Expected: All tests pass, console output shows detailed step-by-step execution with emojis and timing information.