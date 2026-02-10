package testing

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	"github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
)

func TestReportingIntegration(t *testing.T) {
	// Capture output
	var output bytes.Buffer
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(&output)

	test := NewSerenityTestWithReporter(t, reporter)
	defer test.Shutdown()

	// Create actor with API ability
	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Perform successful activities
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts/1"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
	)

	// Get the captured output
	capturedOutput := output.String()

	// Verify key reporting elements are present
	require.Contains(t, capturedOutput, "Starting:", "Should show test start")
	require.Contains(t, capturedOutput, "🚀", "Should show rocket start icon")
	require.Contains(t, capturedOutput, "✅", "Should show success indicators")
	require.Contains(t, capturedOutput, "APITester sends GET request", "Should show activity descriptions")
	require.Contains(t, capturedOutput, "APITester ensures that the last response", "Should show ensure descriptions")
	require.Contains(t, capturedOutput, "TestReportingIntegration", "Should show test name")
}

func TestErrorReporting(t *testing.T) {
	// Test error reporting by creating a failing test scenario
	_ = errors.New("placeholder to use errors package") // avoid unused import

	// Use a separate testing.T instance that we can control
	mockT := &testing.T{}

	// Capture output
	var output bytes.Buffer
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(&output)

	test := NewSerenityTestWithReporter(mockT, reporter)

	// Manually mark the test as failed to trigger error reporting
	mockT.Fail()

	// Shutdown to trigger error reporting
	test.Shutdown()

	// Get the captured output
	capturedOutput := output.String()

	// Verify error reporting elements are present
	require.Contains(t, capturedOutput, "🚀 Starting:", "Should show test start")
	require.Contains(t, capturedOutput, "❌", "Should show failure indicators")
	require.Contains(t, capturedOutput, "FAILED", "Should show FAILED status")
	require.Contains(t, capturedOutput, "Error:", "Should show error details")
}

func TestMultipleActorsReporting(t *testing.T) {
	// Capture output
	var output bytes.Buffer
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(&output)

	test := NewSerenityTestWithReporter(t, reporter)
	defer test.Shutdown()

	// Create multiple actors
	actor1 := test.ActorCalled("Actor1").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))
	actor2 := test.ActorCalled("Actor2").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Both actors perform activities
	actor1.AttemptsTo(api.SendGetRequest("/posts/1"))
	actor2.AttemptsTo(api.SendGetRequest("/posts/2"))

	// Get the captured output
	capturedOutput := output.String()

	// Verify both actors' activities are reported
	require.Contains(t, capturedOutput, "Actor1 sends GET request to /posts/1", "Should show first actor activity")
	require.Contains(t, capturedOutput, "Actor2 sends GET request to /posts/2", "Should show second actor activity")
	require.Contains(t, capturedOutput, "✅", "Should show success indicators")
	require.Contains(t, capturedOutput, "TestMultipleActorsReporting", "Should show test name")
}

func TestComplexWorkflowReporting(t *testing.T) {
	// Capture output
	var output bytes.Buffer
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(&output)

	test := NewSerenityTestWithReporter(t, reporter)
	defer test.Shutdown()

	actor := test.ActorCalled("WorkflowActor").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Perform a complex workflow with multiple steps
	actor.AttemptsTo(
		// Step 1: Get posts
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),

		// Step 2: Create a new post
		api.SendPostRequest("/posts").WithBody(map[string]interface{}{
			"title":  "Test Post",
			"body":   "Test body",
			"userId": 1,
		}),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(201)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("Test Post")),
	)

	// Get the captured output
	capturedOutput := output.String()

	// Verify workflow reporting
	require.Contains(t, capturedOutput, "WorkflowActor sends GET request to /posts", "Should show GET request")
	require.Contains(t, capturedOutput, "WorkflowActor sends POST request to /posts", "Should show POST request")
	require.Contains(t, capturedOutput, "WorkflowActor ensures that", "Should show assertions")
	require.Contains(t, capturedOutput, "✅", "Should show success indicators")
	require.Contains(t, capturedOutput, "TestComplexWorkflowReporting", "Should show test name")
	// Should contain multiple activity logs
	activities := countOccurrences(capturedOutput, "✅")
	require.GreaterOrEqual(t, activities, 5, "Should log multiple activities")
}

func TestConcurrentActivitiesReporting(t *testing.T) {
	// Capture output
	var output bytes.Buffer
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(&output)

	test := NewSerenityTestWithReporter(t, reporter)
	defer test.Shutdown()

	actor := test.ActorCalled("ConcurrentActor").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Perform activities concurrently
	done := make(chan bool, 2)

	go func() {
		actor.AttemptsTo(
			api.SendGetRequest("/posts/1"),
			ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		)
		done <- true
	}()

	go func() {
		actor.AttemptsTo(
			api.SendGetRequest("/posts/2"),
			ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		)
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Get the captured output
	capturedOutput := output.String()

	// Verify concurrent reporting
	require.Contains(t, capturedOutput, "ConcurrentActor sends GET request to /posts/1", "Should show first concurrent request")
	require.Contains(t, capturedOutput, "ConcurrentActor sends GET request to /posts/2", "Should show second concurrent request")
	require.Contains(t, capturedOutput, "✅", "Should show success indicators")
	require.Contains(t, capturedOutput, "TestConcurrentActivitiesReporting", "Should show test name")
}

// Helper function to count occurrences of a substring
func countOccurrences(s, substr string) int {
	count := 0
	start := 0
	for {
		idx := bytes.Index([]byte(s[start:]), []byte(substr))
		if idx == -1 {
			break
		}
		count++
		start += idx + 1
	}
	return count
}
