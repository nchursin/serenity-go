package core

// This file defines failure handling modes for activities in Serenity-Go.
// Failure modes determine how the test execution should proceed when
// activities encounter errors during execution.
//
// Failure Modes:
//
//	FailFast          - Stop execution immediately on first error (default)
//	ErrorButContinue  - Log error but continue with remaining activities
//	Ignore            - Completely ignore failures and continue
//
// Usage Examples:
//
//	// Default behavior (FailFast)
//	actor.AttemptsTo(
//		core.Do("critical setup", setupSystem),
//		core.Do("main operation", mainOperation),     // Won't execute if setup fails
//		core.Do("cleanup", cleanup),                  // Won't execute if above fails
//	)
//
//	// Non-critical operations
//	actor.AttemptsTo(
//		core.Do("main operation", mainOperation),
//		core.Do("log metrics", logMetrics).WithFailureMode(core.NonCritical()),
//		core.Do("send notification", sendNotification).WithFailureMode(core.NonCritical()),
//		core.Do("cleanup", cleanup),                   // Always executes
//	)
//
//	// Optional operations
//	actor.AttemptsTo(
//		core.Do("main operation", mainOperation),
//		core.Do("optional cleanup", cleanup).WithFailureMode(core.Optional()),
//		core.Do("final verification", verification),    // Always executes
//	)
//
// Choosing the Right Mode:
//
//	- Use FailFast for critical operations where failure invalidates the test
//	- Use ErrorButContinue for non-critical operations where failure is notable
//	- Use Ignore for truly optional operations where failure is expected/acceptable
//

// FailureMode defines how activities should handle failures.
// This enum determines whether test execution continues or stops
// when an activity encounters an error.
//
// Default Behavior:
//
//	All activities use FailFast mode unless explicitly overridden.
type FailureMode int

const (
	// FailFast stops execution immediately when an activity fails.
	// This is the default and most commonly used failure mode.
	//
	// Use Cases:
	//	- Critical setup operations
	//	- Main test scenarios
	//	- Dependencies required for subsequent steps
	//	- Validation steps that invalidate the test if failed
	//
	// Example:
	//
	//	actor.AttemptsTo(
	//		core.Do("establishes database connection", connectToDB), // Critical
	//		core.Do("creates test data", createTestData),               // Won't run if above fails
	//		core.Do("runs main test", mainTest),                       // Won't run if above fails
	//	)
	//
	// Behavior:
	//	- Stops execution immediately on error
	//	- Returns the first error encountered
	//	- Subsequent activities are not executed
	FailFast FailureMode = iota

	// ErrorButContinue logs the error but continues with remaining activities.
	// Use this for non-critical operations where you want to know about
	// failures but don't want them to stop the entire test.
	//
	// Use Cases:
	//	- Monitoring and metrics collection
	//	- Non-essential notifications
	//	- Cache operations
	//	- Audit logging
	//	- Performance measurements
	//
	// Example:
	//
	//	actor.AttemptsTo(
	//		core.Do("executes main business logic", businessLogic),
	//		core.Do("collects performance metrics", collectMetrics).WithFailureMode(core.NonCritical()),
	//		core.Do("sends usage statistics", sendStats).WithFailureMode(core.NonCritical()),
	//		core.Do("verifies result", verifyResult), // Always runs
	//	)
	//
	// Behavior:
	//	- Logs the error but continues execution
	//	- Does not return the error unless all activities fail
	//	- Later activities continue to execute
	ErrorButContinue

	// Ignore completely ignores the failure and continues.
	// Use this for truly optional operations where failure is acceptable
	// or expected in certain scenarios.
	//
	// Use Cases:
	//	- Optional cleanup operations
	//	- Best-effort notifications
	//	- Resource cleanup where failure doesn't matter
	//	- Precondition checks that may legitimately fail
	//	- Debugging or diagnostic operations
	//
	// Example:
	//
	//	actor.AttemptsTo(
	//		core.Do("executes main test", mainTest),
	//		core.Do("cleans up temporary files", cleanupTempFiles).WithFailureMode(core.Optional()),
	//		core.Do("attempts to send debug info", sendDebugInfo).WithFailureMode(core.Optional()),
	//		core.Do("final verification", finalVerification), // Always runs
	//	)
	//
	// Behavior:
	//	- Completely ignores any errors
	//	- Never returns errors from this activity
	//	- Execution continues regardless of success or failure
	Ignore
)

// Critical returns a failure mode that stops execution on failure.
// This is a semantic function that returns FailFast mode.
// Use this when you want to explicitly indicate critical operations.
//
// Returns:
//   - FailureMode: Always returns FailFast
//
// Usage:
//
//	// Explicit critical operation
//	actor.AttemptsTo(
//		core.Do("establishes database connection", connectDB).WithFailureMode(core.Critical()),
//		core.Do("creates user account", createUser),
//	)
//
//	// Equivalent to default behavior
//	actor.AttemptsTo(
//		core.Do("establishes database connection", connectDB), // Already critical by default
//		core.Do("creates user account", createUser),
//	)
//
// When to Use:
//   - For clarity in complex test scenarios
//   - When explicitly documenting critical dependencies
//   - In reusable test components to ensure critical behavior
func Critical() FailureMode { return FailFast }

// NonCritical returns a failure mode that logs errors but continues.
// This is a semantic function that returns ErrorButContinue mode.
// Use this for operations that should be noted but not stop execution.
//
// Returns:
//   - FailureMode: Always returns ErrorButContinue
//
// Usage:
//
//	actor.AttemptsTo(
//		core.Do("main business operation", businessLogic),
//		core.Do("collects usage metrics", collectMetrics).WithFailureMode(core.NonCritical()),
//		core.Do("sends user notification", sendNotification).WithFailureMode(core.NonCritical()),
//		core.Do("verifies result", verifyResult),
//	)
//
// When to Use:
//   - Monitoring and telemetry operations
//   - Non-essential user communications
//   - Performance measurements
//   - Audit and logging activities
func NonCritical() FailureMode { return ErrorButContinue }

// Optional returns a failure mode that ignores errors completely.
// This is a semantic function that returns Ignore mode.
// Use this for operations where failure is acceptable or expected.
//
// Returns:
//   - FailureMode: Always returns Ignore
//
// Usage:
//
//	actor.AttemptsTo(
//		core.Do("main test execution", mainTest),
//		core.Do("attempts to cleanup resources", cleanup).WithFailureMode(core.Optional()),
//		core.Do("tries to send debug report", sendDebugReport).WithFailureMode(core.Optional()),
//		core.Do("final assertion", finalAssertion),
//	)
//
// When to Use:
//   - Optional cleanup operations
//   - Best-effort notifications
//   - Resource cleanup where failure is acceptable
//   - Precondition checks that may legitimately fail
//   - Debug or diagnostic operations
func Optional() FailureMode { return Ignore }
