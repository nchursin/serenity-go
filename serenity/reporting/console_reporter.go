package reporting

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nchursin/serenity-go/serenity/core"
)

// TODO: вныести ConsoleReporter в подпакет `console_reporter`

// Reporter handles test execution reporting
type Reporter interface {
	// OnTestStart is called when a test begins
	OnTestStart(testName string)

	// OnTestFinish is called when a test completes
	OnTestFinish(result core.TestResult)

	// OnStepStart is called when a step/activity begins
	OnStepStart(stepDescription string)

	// OnStepFinish is called when a step/activity completes
	OnStepFinish(stepResult core.TestResult)

	// SetOutput sets the output destination
	SetOutput(w io.Writer)
}

// ConsoleReporter provides console-based test reporting
type ConsoleReporter struct {
	output      io.Writer
	currentTest string
	indentLevel int
}

// NewConsoleReporter creates a new console reporter
func NewConsoleReporter() *ConsoleReporter {
	return &ConsoleReporter{
		output: os.Stdout,
	}
}

// SetOutput sets the output destination
func (cr *ConsoleReporter) SetOutput(w io.Writer) {
	cr.output = w
}

// OnTestStart is called when a test begins
func (cr *ConsoleReporter) OnTestStart(testName string) {
	cr.currentTest = testName
	cr.indentLevel = 0
	cr.writeLine("🚀 Starting: %s", testName)
}

// OnTestFinish is called when a test completes
func (cr *ConsoleReporter) OnTestFinish(result core.TestResult) {
	emoji := "✅"
	statusText := "PASSED"

	switch result.Status {
	case core.StatusFailed:
		emoji = "❌"
		statusText = "FAILED"
	case core.StatusSkipped:
		emoji = "⏭️"
		statusText = "SKIPPED"
	}

	cr.writeLine("%s %s: %s (%.2fs)", emoji, result.Name, statusText, result.Duration.Seconds())

	if result.Error != nil {
		cr.writeLine("   Error: %s", result.Error.Error())
	}

	cr.writeLine("")
}

// OnStepStart is called when a step/activity begins
func (cr *ConsoleReporter) OnStepStart(stepDescription string) {
	cr.indentLevel++
	description := cr.formatStepDescription(stepDescription)
	cr.writeLine("%s🔄 %s", cr.getIndent(), description)
}

// OnStepFinish is called when a step/activity completes
func (cr *ConsoleReporter) OnStepFinish(stepResult core.TestResult) {
	emoji := "✅"

	if stepResult.Status == core.StatusFailed {
		emoji = "❌"
	}

	description := cr.formatStepDescription(stepResult.Name)
	cr.writeLine("%s%s %s (%.2fs)", cr.getIndent(), emoji, description, stepResult.Duration.Seconds())

	if stepResult.Error != nil {
		cr.writeLine("%s   Error: %s", cr.getIndent(), stepResult.Error.Error())
	}

	cr.indentLevel--
}

// formatStepDescription formats step descriptions for better readability
func (cr *ConsoleReporter) formatStepDescription(description string) string {
	// Remove #actor prefix if present
	formatted := strings.ReplaceAll(description, "#actor ", "")

	// Capitalize first letter
	if len(formatted) > 0 {
		formatted = strings.ToUpper(formatted[:1]) + formatted[1:]
	}

	return formatted
}

// getIndent returns the current indentation string
func (cr *ConsoleReporter) getIndent() string {
	return strings.Repeat("  ", cr.indentLevel)
}

// writeLine writes a formatted line to the output
func (cr *ConsoleReporter) writeLine(format string, args ...interface{}) {
	if cr.output != nil {
		fmt.Fprintf(cr.output, format+"\n", args...)
	}
}

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
	status := core.StatusPassed
	var activityErr error = nil

	if err != nil {
		status = core.StatusFailed
		activityErr = err
	}

	result := core.TestResult{
		Name:     at.activity,
		Status:   status,
		Duration: time.Since(at.startTime),
		Error:    activityErr,
	}

	at.reporter.OnStepFinish(result)
}
