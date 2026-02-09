package testing

// TestContext provides an interface for test operations and logging
// that wraps the standard testing.TB interface
type TestContext interface {
	// Logf logs a formatted message
	Logf(format string, args ...interface{})

	// Errorf logs a formatted error message and marks the test as failed
	Errorf(format string, args ...interface{})

	// FailNow marks the test as failed and stops execution
	FailNow()

	// Failed returns true if the test has already failed
	Failed() bool
}
