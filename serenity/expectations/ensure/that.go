package ensure

import (
	"fmt"

	"github.com/nchursin/serenity-go/serenity/core"
)

// Expectation represents an expectation that can be evaluated against actual values
type Expectation[T any] interface {
	// Evaluate evaluates the expectation against the actual value
	Evaluate(actual T) error

	// Description returns a human-readable description of the expectation
	Description() string
}

// EnsureActivity represents an assertion that a question's answer meets an expectation
type EnsureActivity[T any] struct {
	question    core.Question[T]
	expectation Expectation[T]
}

// That creates a new Ensure assertion with the new API
func That[T any](question core.Question[T], expectation Expectation[T]) core.Activity {
	return &EnsureActivity[T]{
		question:    question,
		expectation: expectation,
	}
}

// Description returns the activity description
func (e *EnsureActivity[T]) Description() string {
	return fmt.Sprintf("#actor ensures that %s %s", e.question.Description(), e.expectation.Description())
}

// PerformAs executes the ensure activity
func (e *EnsureActivity[T]) PerformAs(actor core.Actor) error {
	actual, err := e.question.AnsweredBy(actor)
	if err != nil {
		return fmt.Errorf("failed to answer question '%s': %w", e.question.Description(), err)
	}

	if evaluateErr := e.expectation.Evaluate(actual); evaluateErr != nil {
		return fmt.Errorf("assertion failed for '%s': %w", e.question.Description(), evaluateErr)
	}

	return nil
}

// FailureMode returns the failure mode for ensure activities (default: FailFast)
func (e *EnsureActivity[T]) FailureMode() core.FailureMode {
	return core.NonCritical()
}
