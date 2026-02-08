package console_reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nchursin/serenity-go/serenity/reporting"
)

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
func (cr *ConsoleReporter) OnTestFinish(result reporting.TestResult) {
	emoji := "✅"
	statusText := "PASSED"

	switch result.Status() {
	case reporting.StatusFailed:
		emoji = "❌"
		statusText = "FAILED"
	case reporting.StatusSkipped:
		emoji = "⏭️"
		statusText = "SKIPPED"
	}

	cr.writeLine("%s %s: %s (%.2fs)", emoji, result.Name(), statusText, result.Duration())

	if result.Error() != nil {
		cr.writeLine("   Error: %s", result.Error().Error())
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
func (cr *ConsoleReporter) OnStepFinish(stepResult reporting.TestResult) {
	emoji := "✅"

	if stepResult.Status() == reporting.StatusFailed {
		emoji = "❌"
	}

	description := cr.formatStepDescription(stepResult.Name())
	cr.writeLine("%s%s %s (%.2fs)", cr.getIndent(), emoji, description, stepResult.Duration())

	if stepResult.Error() != nil {
		cr.writeLine("%s   Error: %s", cr.getIndent(), stepResult.Error().Error())
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
		_, _ = fmt.Fprintf(cr.output, format+"\n", args...)
	}
}
