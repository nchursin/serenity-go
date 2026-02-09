package testing

import (
	"fmt"
	"sync"

	"github.com/nchursin/serenity-go/serenity/abilities"
	"github.com/nchursin/serenity-go/serenity/core"
)

// testActor is a wrapper around core.Actor that includes test context
type testActor struct {
	name        string
	abilities   []abilities.Ability
	testContext TestContext
	mutex       sync.RWMutex
}

// Name returns the actor's name
func (ta *testActor) Name() string {
	return ta.name
}

// WhoCan adds abilities to the actor
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

// AttemptsTo executes activities with error handling through test context
func (ta *testActor) AttemptsTo(activities ...core.Activity) {
	for _, activity := range activities {
		if err := activity.PerformAs(ta); err != nil {
			switch activity.FailureMode() {
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
