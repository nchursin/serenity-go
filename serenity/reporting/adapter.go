package reporting

import "time"

// TestRunnerAdapter provides integration with test runners
type TestRunnerAdapter struct {
	reporter Reporter
}

// NewTestRunnerAdapter creates a new test runner adapter
func NewTestRunnerAdapter(reporter Reporter) *TestRunnerAdapter {
	return &TestRunnerAdapter{
		reporter: reporter,
	}
}

// GetReporter returns the underlying reporter
func (tra *TestRunnerAdapter) GetReporter() Reporter {
	return tra.reporter
}

// ActivityTracker tracks activity execution for reporting
type ActivityTracker struct {
	reporter  Reporter
	activity  string
	startTime time.Time
}

// NewActivityTracker creates a new activity tracker
func NewActivityTracker(reporter Reporter, activity string) *ActivityTracker {
	return &ActivityTracker{
		reporter:  reporter,
		activity:  activity,
		startTime: time.Now(),
	}
}

// Start starts tracking the activity
func (at *ActivityTracker) Start() {
	at.reporter.OnStepStart(at.activity)
}

// Finish completes tracking the activity
func (at *ActivityTracker) Finish(err error) {
	status := StatusPassed
	var activityErr error = nil

	if err != nil {
		status = StatusFailed
		activityErr = err
	}

	result := &testResult{
		name:     at.activity,
		status:   status,
		duration: time.Since(at.startTime).Seconds(),
		error:    activityErr,
	}

	at.reporter.OnStepFinish(result)
}

// testResult implements TestResult interface
type testResult struct {
	name     string
	status   Status
	duration float64
	error    error
}

func (tr *testResult) Name() string      { return tr.name }
func (tr *testResult) Status() Status    { return tr.status }
func (tr *testResult) Duration() float64 { return tr.duration }
func (tr *testResult) Error() error      { return tr.error }
