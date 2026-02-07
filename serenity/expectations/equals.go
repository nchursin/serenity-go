package expectations

import (
	"fmt"
	"reflect"

	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

// EqualsExpectation checks if the actual value equals the expected value
type EqualsExpectation[T any] struct {
	expected T
}

// NewEquals creates a new Equals expectation
func NewEquals[T any](expected T) ensure.Expectation[T] {
	return EqualsExpectation[T]{expected: expected}
}

// Evaluate evaluates the equals expectation
func (eq EqualsExpectation[T]) Evaluate(actual T) error {
	if !reflect.DeepEqual(actual, eq.expected) {
		return fmt.Errorf("expected %v, but got %v", eq.expected, actual)
	}
	return nil
}

// Description returns the expectation description
func (eq EqualsExpectation[T]) Description() string {
	return fmt.Sprintf("equals %v", eq.expected)
}

// Convenience function for creating Equals expectations
func Equals[T any](expected T) ensure.Expectation[T] {
	return NewEquals(expected)
}
