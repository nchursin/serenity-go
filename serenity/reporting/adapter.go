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
	actorName string
	startTime time.Time
}

// NewActivityTracker creates a new activity tracker (backward compatibility)
func NewActivityTracker(reporter Reporter, activity string) *ActivityTracker {
	return &ActivityTracker{
		reporter:  reporter,
		activity:  activity,
		actorName: "", // No actor name for backward compatibility
		startTime: time.Now(),
	}
}

// NewActivityTrackerWithActor creates a new activity tracker with actor name
func NewActivityTrackerWithActor(reporter Reporter, activity string, actorName string) *ActivityTracker {
	return &ActivityTracker{
		reporter:  reporter,
		activity:  activity,
		actorName: actorName,
		startTime: time.Now(),
	}
}

// getActivityDescription replaces #actor placeholder with actor name
func (at *ActivityTracker) getActivityDescription() string {
	if at.actorName == "" {
		return at.activity // No actor name, return original
	}

	// Replace #actor with actor name
	description := at.activity
	if len(description) >= 7 && description[:7] == "#actor " {
		description = at.actorName + " " + description[7:] // Replace "#actor " with actor name
	}

	return description
}

// Start starts tracking the activity
func (at *ActivityTracker) Start() {
	description := at.getActivityDescription()
	at.reporter.OnStepStart(description)
}

// Finish completes tracking the activity
func (at *ActivityTracker) Finish(err error) {
	status := StatusPassed
	var activityErr error = nil

	if err != nil {
		status = StatusFailed
		activityErr = err
	}

	description := at.getActivityDescription()
	result := &testResult{
		name:     description,
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
