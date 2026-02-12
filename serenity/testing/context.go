// Package testing provides the TestContext API for simplified Serenity/JS testing in Go.
//
// The TestContext API eliminates the need for manual error handling in tests by
// automatically managing test failures through the testing.TB interface.
//
// Key Features:
//
//   - Automatic error handling through TestContext
//   - Actor lifecycle management with defer.Shutdown()
//   - Integrated reporting capabilities
//   - Support for multiple actors in single test
//   - Thread-safe actor management
//
// Basic Usage:
//
//	test := serenity.NewSerenityTest(t)
//	defer test.Shutdown()
//
//	actor := test.ActorCalled("APITester").WhoCan(
//		api.CallAnApiAt("https://api.example.com"),
//	)
//
//	actor.AttemptsTo(
//		api.SendGetRequest("/users"),
//		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
//	)
//
// Multiple Actors:
//
//	test := serenity.NewSerenityTest(t)
//	defer test.Shutdown()
//
//	admin := test.ActorCalled("Admin").WhoCan(api.CallAnApiAt(apiURL))
//	user := test.ActorCalled("User").WhoCan(api.CallAnApiAt(apiURL))
//
//	admin.AttemptsTo(api.SendPostRequest("/users", userData))
//	user.AttemptsTo(api.SendGetRequest("/users/1"))
//
// Custom Reporting:
//
//	reporter := custom.NewJSONReporter()
//	test := serenity.NewSerenityTestWithReporter(t, reporter)
//	defer test.Shutdown()
//
// Error Handling:
//
//	Unlike the legacy API where errors need to be manually handled:
//
//	// Legacy approach
//	err := actor.AttemptsTo(activity)
//	require.NoError(t, err)
//
//	// TestContext API - automatic error handling
//	actor.AttemptsTo(activity) // Errors automatically fail the test
//
// Thread Safety:
//
//	All actor operations are thread-safe. Multiple goroutines can safely
//	use actors created from the same SerenityTest instance.
package testing

//go:generate mockgen -source=context.go -destination=mocks/mock_test_context.go -package=mocks

// TestContext provides an interface for test operations and logging
// that wraps the standard testing.TB interface
type TestContext interface {
	// Name returns the name of the test
	Name() string

	// Logf logs a formatted message
	Logf(format string, args ...interface{})

	// Errorf logs a formatted error message and marks the test as failed
	Errorf(format string, args ...interface{})

	// FailNow marks the test as failed and stops execution
	FailNow()

	// Failed returns true if the test has already failed
	Failed() bool
}
