package testing

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nchursin/serenity-go/serenity/abilities"
	"github.com/nchursin/serenity-go/serenity/core"
)

// SerenityTest provides a test context and manages actors for a test
type SerenityTest interface {
	// ActorCalled returns an actor with the given name, creating it if necessary
	ActorCalled(name string) core.Actor

	// Shutdown cleans up resources and should be called with defer
	Shutdown()
}

// serenityTest implements SerenityTest
type serenityTest struct {
	ctx    TestContext
	actors map[string]core.Actor
	mutex  sync.RWMutex
}

// NewSerenityTest creates a new SerenityTest instance
func NewSerenityTest(t testing.TB) SerenityTest {
	return &serenityTest{
		ctx:    NewTestContext(t),
		actors: make(map[string]core.Actor),
	}
}

// ActorCalled returns an actor with the given name
func (st *serenityTest) ActorCalled(name string) core.Actor {
	st.mutex.RLock()
	actor, exists := st.actors[name]
	st.mutex.RUnlock()

	if exists {
		return actor
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Double-check after acquiring write lock
	if actor, exists := st.actors[name]; exists {
		return actor
	}

	// Create new actor with test context
	actor = &testActor{
		name:        name,
		abilities:   make([]abilities.Ability, 0),
		testContext: st.ctx,
	}

	st.actors[name] = actor
	return actor
}

// Shutdown cleans up resources
func (st *serenityTest) Shutdown() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	// Clear actors map
	st.actors = make(map[string]core.Actor)
}

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

// Helper function to check if ability matches type
func abilityMatchesType(ability, abilityType abilities.Ability) bool {
	return fmt.Sprintf("%T", ability) == fmt.Sprintf("%T", abilityType)
}
