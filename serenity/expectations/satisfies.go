package expectations

import (
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

// SatisfiesExpectation represents a custom expectation that evaluates using a provided function
type SatisfiesExpectation[T any] struct {
	description string
	fn          func(T) error
}

// Satisfies creates a new custom expectation with a description and validation function
//
// The function should return nil if the expectation is met, or an error describing the failure.
//
// Example:
//
//	actor.AttemptsTo(
//		ensure.That(answerable.ValueOf(value), expectations.Satisfies("is positive number", func(actual int) error {
//			if actual <= 0 {
//				return fmt.Errorf("expected positive value, but got %v", actual)
//			}
//			return nil
//		})),
//	)
func Satisfies[T any](description string, fn func(T) error) ensure.Expectation[T] {
	return &SatisfiesExpectation[T]{
		description: description,
		fn:          fn,
	}
}

// Evaluate evaluates the expectation by calling the provided function
func (s *SatisfiesExpectation[T]) Evaluate(actual T) error {
	return s.fn(actual)
}

// Description returns the custom description of the expectation
func (s *SatisfiesExpectation[T]) Description() string {
	return s.description
}
