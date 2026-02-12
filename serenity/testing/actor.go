package testing

import (
	"fmt"
	"sync"

	"github.com/nchursin/serenity-go/serenity/abilities"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/reporting"
)

// testActor implements the Actor interface with TestContext integration.
// This actor automatically handles errors through the embedded TestContext,
// providing a seamless testing experience without manual error checking.
//
// Key Features:
//   - Automatic error propagation to test framework
//   - Thread-safe operations
//   - Integrated reporting capabilities
//   - Support for all standard Actor methods
type testActor struct {
	name        string                       // Actor name for reporting
	abilities   []abilities.Ability          // Actor abilities
	testContext TestContext                  // Embedded test context for error handling
	reporter    *reporting.TestRunnerAdapter // Integrated reporter for activity tracking
	mutex       sync.RWMutex                 // Mutex for thread-safe operations
}

// Name returns the actor's name
func (ta *testActor) Name() string {
	return ta.name
}

// WhoCan adds abilities to the actor and returns the same actor instance for chaining.
// This method is thread-safe and can be called multiple times.
//
// Example:
//
//	actor := test.ActorCalled("APIUser").
//		WhoCan(api.CallAnApiAt("https://api.example.com")).
//		WhoCan(db.ConnectToDatabase("postgres://localhost/test"))
//
// Parameters:
//
//	abilities - List of abilities to add to the actor
//
// Returns:
//
//	The same actor instance with added abilities for method chaining
func (ta *testActor) WhoCan(abilities ...abilities.Ability) core.Actor {
	ta.mutex.Lock()
	defer ta.mutex.Unlock()

	ta.abilities = append(ta.abilities, abilities...)
	return ta
}

// AbilityTo returns the specified ability
func (ta *testActor) AbilityTo(abilityType abilities.Ability) (abilities.Ability, error) {
	ta.mutex.RLock()
	defer ta.mutex.RUnlock()

	for _, ability := range ta.abilities {
		if abilityMatchesType(ability, abilityType) {
			return ability, nil
		}
	}

	return nil, fmt.Errorf("actor '%s' does not have the required ability", ta.name)
}

// AttemptsTo executes activities and automatically handles any errors through TestContext.
// Unlike the legacy API, no manual error checking is required - failures automatically
// fail the test with descriptive error messages.
//
// Example:
//
//	// TestContext API - automatic error handling
//	actor.AttemptsTo(
//		api.SendGetRequest("/users"),
//		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
//	)
//
//	// Legacy API comparison
//	err := legacyActor.AttemptsTo(
//		api.SendGetRequest("/users"),
//		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
//	)
//	require.NoError(t, err) // Manual error handling required
//
// Parameters:
//
//	activities - List of activities to execute in order
//
// This method automatically handles different failure modes:
//   - FailFast: Stops test execution immediately on error
//   - ErrorButContinue: Logs error but continues with remaining activities
//   - Ignore: Silently ignores the error and continues
func (ta *testActor) AttemptsTo(activities ...core.Activity) {
	for _, activity := range activities {
		var tracker *reporting.ActivityTracker
		if ta.reporter != nil {
			tracker = reporting.NewActivityTrackerWithActor(ta.reporter.GetReporter(), activity.Description(), ta.name)
			tracker.Start()
		}

		err := activity.PerformAs(ta)

		if tracker != nil {
			tracker.Finish(err)
		}

		if err != nil {
			failureMode := activity.FailureMode()
			switch failureMode {
			case core.FailFast:
				ta.testContext.Errorf("Critical activity error '%s' failed: %v", activity.Description(), err)
				ta.testContext.FailNow()
				return
			case core.ErrorButContinue:
				ta.testContext.Errorf("Non-critical activity error '%s' failed: %v", activity.Description(), err)
			case core.Ignore:
				ta.testContext.Logf("Ignore activity error '%s' failed: %v", activity.Description(), err)
				// Do nothing
			}
		}
	}
}

// AnswersTo answers questions with boolean success flag
func (ta *testActor) AnswersTo(question core.Question[any]) (any, bool) {
	result, err := question.AnsweredBy(ta)
	if err != nil {
		ta.testContext.Errorf("Failed to answer question '%s': %v", question.Description(), err)
		return nil, false
	}
	return result, true
}
