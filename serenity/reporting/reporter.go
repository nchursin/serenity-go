package reporting

import "io"

// Reporter handles test execution reporting
type Reporter interface {
	// OnTestStart is called when a test begins
	OnTestStart(testName string)

	// OnTestFinish is called when a test completes
	OnTestFinish(result TestResult)

	// OnStepStart is called when a step/activity begins
	OnStepStart(stepDescription string)

	// OnStepFinish is called when a step/activity completes
	OnStepFinish(stepResult TestResult)

	// SetOutput sets the output destination
	SetOutput(w io.Writer)
}

// TestResult represents the result of a test or step execution
type TestResult interface {
	Name() string
	Status() Status
	Duration() float64
	Error() error
}

// Status represents the status of a test or step
type Status int

const (
	StatusPassed Status = iota
	StatusFailed
	StatusSkipped
)
