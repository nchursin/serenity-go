package examples

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

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

	test := serenity.NewSerenityTestWithReporter(context.Background(), t, reporter)

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
	// Create file for output in current directory
	file, err := os.Create("test_report.txt")
	require.NoError(t, err)

	// Create reporter that writes to file
	reporter := console_reporter.NewConsoleReporter()
	reporter.SetOutput(file)

	test := serenity.NewSerenityTestWithReporter(context.Background(), t, reporter)

	apiTester := test.ActorCalled("FileReporter").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	apiTester.AttemptsTo(
		api.SendGetRequest("/users"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)

	// Ensure file is closed and flushed
	file.Close()

	// Verify file was created and contains content
	content, err := os.ReadFile("test_report.txt")
	require.NoError(t, err)
	require.Contains(t, string(content), "Starting: TestReportingToFile")
	require.Contains(t, string(content), "FileReporter sends GET request")
}
