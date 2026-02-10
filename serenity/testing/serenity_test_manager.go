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

// SerenityTest provides a test context and manages actors for a test
type SerenityTest interface {
	// ActorCalled returns an actor with the given name, creating it if necessary
	ActorCalled(name string) core.Actor

	// Shutdown cleans up resources and should be called with defer
	Shutdown()

	// GetReporterAdapter returns the test runner adapter for reporting
	GetReporterAdapter() *reporting.TestRunnerAdapter
}

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
