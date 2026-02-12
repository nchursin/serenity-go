package testing

import (
	"fmt"
	"sync"
	"time"

	"github.com/nchursin/serenity-go/serenity/abilities"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/reporting"
	"github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
)

// ReporterProvider provides access to reporter adapter
type ReporterProvider interface {
	// GetReporterAdapter returns the test runner adapter for reporting
	GetReporterAdapter() *reporting.TestRunnerAdapter
}

// SerenityTest manages the lifecycle of test actors and provides the TestContext API.
// This interface serves as the main entry point for using the simplified testing approach.
//
// Lifecycle Management:
//  1. Create test instance with NewSerenityTest() or NewSerenityTestWithReporter()
//  2. Create actors using ActorCalled()
//  3. Execute test activities
//  4. Call Shutdown() to clean up resources (typically via defer)
//
// Thread Safety:
//
//	All SerenityTest methods are thread-safe. Multiple goroutines can safely
//	create and use actors from the same test instance.
type SerenityTest interface {
	// TestContext returns the embedded testing.TB interface.
	// This method provides access to the underlying testing framework.
	TestContext() TestContext

	// ActorCalled creates a new test-aware actor with the specified name.
	// The actor is automatically configured with TestContext error handling.
	//
	// Parameters:
	//	name - Human-readable name for the actor (used in reporting)
	//
	// Returns:
	//	An Actor instance configured for automatic error handling
	ActorCalled(name string) core.Actor

	// Shutdown cleans up resources and finalizes the test.
	// This method should be called via defer after creating the test instance.
	// Failure to call Shutdown() may result in resource leaks.
	//
	// Example:
	//	test := serenity.NewSerenityTest(t)
	//	defer test.Shutdown() // Ensure cleanup
	//
	// Side effects:
	//	- Flushes any pending reports
	//	- Cleans up actor resources
	//	- Finalizes test metrics
	Shutdown()

	// GetReporterAdapter returns the test runner adapter for reporting
	GetReporterAdapter() *reporting.TestRunnerAdapter
}

// Test Lifecycle Examples:
//
// Basic Test Structure:
//
//	func TestAPIEndpoints(t *testing.T) {
//		test := serenity.NewSerenityTest(t)
//		defer test.Shutdown()
//
//		actor := test.ActorCalled("APITester").WhoCan(
//			api.CallAnApiAt("https://jsonplaceholder.typicode.com"),
//		)
//
//		actor.AttemptsTo(
//			api.SendGetRequest("/posts"),
//			ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
//			ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
//		)
//	}
//
// Test with Custom Reporter:
//
//	func TestWithCustomReporting(t *testing.T) {
//		reporter := custom.NewJSONReporter()
//		test := serenity.NewSerenityTestWithReporter(t, reporter)
//		defer test.Shutdown()
//
//		actor := test.ActorCalled("ReportedUser").WhoCan(api.CallAnApiAt(apiURL))
//		actor.AttemptsTo(api.SendGetRequest("/health"))
//	}

// testResult implements the TestResult interface
type testResult struct {
	name     string
	status   reporting.Status
	duration time.Duration
	err      error
}

// Name returns the test name
func (tr *testResult) Name() string {
	return tr.name
}

// Status returns the test status
func (tr *testResult) Status() reporting.Status {
	return tr.status
}

// Duration returns the test duration in seconds
func (tr *testResult) Duration() float64 {
	return tr.duration.Seconds()
}

// Error returns the test error, if any
func (tr *testResult) Error() error {
	return tr.err
}

// serenityTest implements SerenityTest
type serenityTest struct {
	ctx       TestContext
	actors    map[string]core.Actor
	mutex     sync.RWMutex
	adapter   *reporting.TestRunnerAdapter
	startTime time.Time
	testName  string
}

// NewSerenityTest creates a new SerenityTest instance
func NewSerenityTest(t TestContext) SerenityTest {
	return NewSerenityTestWithReporter(t, console_reporter.NewConsoleReporter())
}

// NewSerenityTestWithReporter creates a new SerenityTest instance with a reporter
func NewSerenityTestWithReporter(t TestContext, reporter reporting.Reporter) SerenityTest {
	var adapter *reporting.TestRunnerAdapter
	if reporter != nil {
		adapter = reporting.NewTestRunnerAdapter(reporter)
	}

	testName := t.Name()

	// Notify reporter that test is starting
	if reporter != nil {
		reporter.OnTestStart(testName)
	}

	return &serenityTest{
		ctx:       t,
		actors:    make(map[string]core.Actor),
		adapter:   adapter,
		startTime: time.Now(),
		testName:  testName,
	}
}

// ActorCalled returns an actor with the given name
func (st *serenityTest) ActorCalled(name string) core.Actor {
	st.mutex.RLock()
	actor, exists := st.actors[name]
	st.mutex.RUnlock()

	if exists {
		return actor
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Double-check after acquiring write lock
	if actor, exists := st.actors[name]; exists {
		return actor
	}

	// Create new actor with test context and reporter
	actor = &testActor{
		name:        name,
		abilities:   make([]abilities.Ability, 0),
		testContext: st.ctx,
		reporter:    st.adapter,
	}

	st.actors[name] = actor
	return actor
}

// TestContext returns the embedded testing.TB interface.
// This method provides access to the underlying testing framework.
func (st *serenityTest) TestContext() TestContext {
	return st.ctx
}

// GetReporterAdapter returns the test runner adapter for reporting
func (st *serenityTest) GetReporterAdapter() *reporting.TestRunnerAdapter {
	return st.adapter
}

// Shutdown cleans up resources
func (st *serenityTest) Shutdown() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Create test result
	duration := time.Since(st.startTime)
	status := reporting.StatusPassed
	var testErr error

	if st.ctx.Failed() {
		status = reporting.StatusFailed
		testErr = fmt.Errorf("test failed")
	}

	result := &testResult{
		name:     st.testName,
		status:   status,
		duration: duration,
		err:      testErr,
	}

	// Notify reporter that test is finished
	if st.adapter != nil && st.adapter.GetReporter() != nil {
		st.adapter.GetReporter().OnTestFinish(result)
	}

	// Clear actors map
	st.actors = make(map[string]core.Actor)
}
