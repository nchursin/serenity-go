package core

import (
	"fmt"
	"sync"
)

// actor implements the Actor interface
type actor struct {
	name      string
	abilities []Ability
	mutex     sync.RWMutex
}

// NewActor creates a new actor with the given name
func NewActor(name string) Actor {
	return &actor{
		name:      name,
		abilities: make([]Ability, 0),
	}
}

// Name returns the actor's name
func (a *actor) Name() string {
	return a.name
}

// WhoCan adds abilities to the actor
func (a *actor) WhoCan(abilities ...Ability) Actor {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.abilities = append(a.abilities, abilities...)
	return a
}

// AbilityTo retrieves a specific ability from the actor
func (a *actor) AbilityTo(targetAbility Ability) (Ability, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	// Check if we have the target ability type
	for _, ability := range a.abilities {
		if abilityTypeOf(targetAbility) == abilityTypeOf(ability) {
			return ability, nil
		}
	}

	return nil, fmt.Errorf("actor %s does not have ability %T", a.name, targetAbility)
}

// AttemptsTo performs one or more activities
func (a *actor) AttemptsTo(activities ...Activity) error {
	for _, activity := range activities {
		if err := activity.PerformAs(a); err != nil {
			return fmt.Errorf("failed to perform activity '%s': %w", activity.Description(), err)
		}
	}
	return nil
}

// AnswersTo answers a question about the system state
func (a *actor) AnswersTo(question Question[any]) (any, error) {
	return question.AnsweredBy(a)
}

// abilityTypeOf returns the type of an ability for comparison
func abilityTypeOf(ability Ability) string {
	return fmt.Sprintf("%T", ability)
}

// NewActor creates a new actor with the given name (exported function)
func New(name string) Actor {
	return NewActor(name)
}

// ActorCalled creates or retrieves an actor by name (similar to Serenity/JS actorCalled)
var actors = make(map[string]Actor)
var actorsMutex sync.RWMutex

func ActorCalled(name string) Actor {
	actorsMutex.RLock()
	if existing, exists := actors[name]; exists {
		actorsMutex.RUnlock()
		return existing
	}
	actorsMutex.RUnlock()

	actorsMutex.Lock()
	defer actorsMutex.Unlock()

	// Double-check after acquiring write lock
	if existing, exists := actors[name]; exists {
		return existing
	}

	newActor := NewActor(name)
	actors[name] = newActor
	return newActor
}

// ForgetActor removes an actor from the registry
func ForgetActor(name string) {
	actorsMutex.Lock()
	defer actorsMutex.Unlock()
	delete(actors, name)
}

// ForgetAllActors clears the actor registry
func ForgetAllActors() {
	actorsMutex.Lock()
	defer actorsMutex.Unlock()
	actors = make(map[string]Actor)
}
