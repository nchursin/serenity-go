package testing

import (
	"sync"

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
func NewSerenityTest(t TestContext) SerenityTest {
	return &serenityTest{
		ctx:    t,
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
