package answerable

import (
	"fmt"
	"reflect"

	"github.com/nchursin/serenity-go/serenity/core"
)

// valueQuestion[T] implements core.Question[T] for static values.
// It always returns the same value regardless of which actor asks the question.
type valueQuestion[T any] struct {
	value T
}

// AnsweredBy returns the static value when asked by any actor.
// For error types, the error is returned as the value rather than as a failure condition.
func (v *valueQuestion[T]) AnsweredBy(actor core.Actor) (T, error) {
	return v.value, nil
}

// Description returns a human-readable description of the value.
// Format: "value (type)" for normal values, "error message (error)" for error types.
func (v *valueQuestion[T]) Description() string {
	// Special handling for error types using reflection
	if isError(v.value) {
		// Convert to any for type assertion on generic types
		if err, ok := any(v.value).(error); ok {
			return fmt.Sprintf("error %v (error)", err)
		}
	}

	return fmt.Sprintf("%v (%T)", v.value, v.value)
}

// isError checks if the provided value is of error type using reflection.
// This allows us to provide special formatting for error values.
func isError(value any) bool {
	if value == nil {
		return false
	}

	// Use reflection to check if the value implements the error interface
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	return reflect.TypeOf(value).Implements(errorType)
}
