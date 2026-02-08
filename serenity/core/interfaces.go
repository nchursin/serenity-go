package core

import (
	"time"

	"github.com/nchursin/serenity-go/serenity/abilities"
)

// Actor represents a person or external system interacting with the system under test.
// Actors have abilities that enable them to perform activities.
type Actor interface {
	// Name returns the actor's name
	Name() string

	// WhoCan gives the actor additional abilities to interact with the system
	WhoCan(abilities ...abilities.Ability) Actor

	// AbilityTo retrieves a specific ability from the actor
	AbilityTo(ability abilities.Ability) (abilities.Ability, error)

	// AttemptsTo performs one or more activities
	AttemptsTo(activities ...Activity) error

	// AnswersTo answers a question about the system state
	AnswersTo(question Question[any]) (any, error)
}

// Activity represents an action that an actor can perform
type Activity interface {
	// PerformAs executes the activity as the given actor
	PerformAs(actor Actor) error

	// Description returns a human-readable description of the activity
	Description() string
}

// Interaction represents a low-level activity (atomic operation)
type Interaction interface {
	Activity
}

// Task represents a high-level business-focused activity composed of interactions
type Task interface {
	Activity
}

// Question enables actors to retrieve information from the system
type Question[T any] interface {
	// AnsweredBy returns the answer when asked by the given actor
	AnsweredBy(actor Actor) (T, error)

	// Description returns a human-readable description of what the question asks
	Description() string
}

// TestResult represents the outcome of a test
type TestResult struct {
	Name     string        `json:"name"`
	Status   Status        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    error         `json:"error,omitempty"`
}

// Status represents the test execution status
type Status int

const (
	StatusPending Status = iota
	StatusRunning
	StatusPassed
	StatusFailed
	StatusSkipped
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusRunning:
		return "running"
	case StatusPassed:
		return "passed"
	case StatusFailed:
		return "failed"
	case StatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}
