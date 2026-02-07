package expectations

import (
	"fmt"
	"reflect"

	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
)

// IsEmptyExpectation checks if a string is empty
type IsEmptyExpectation struct{}

// NewIsEmpty creates a new IsEmpty expectation
func NewIsEmpty() ensure.Expectation[interface{}] {
	return IsEmptyExpectation{}
}

// Evaluate evaluates the is empty expectation
func (ie IsEmptyExpectation) Evaluate(actual interface{}) error {
	val := reflect.ValueOf(actual)

	switch val.Kind() {
	case reflect.String:
		if val.String() != "" {
			return fmt.Errorf("expected string to be empty, but got '%s'", val.String())
		}
	case reflect.Slice, reflect.Array:
		if val.Len() != 0 {
			return fmt.Errorf("expected slice/array to be empty, but got %d elements", val.Len())
		}
	case reflect.Map:
		if val.Len() != 0 {
			return fmt.Errorf("expected map to be empty, but got %d elements", val.Len())
		}
	default:
		return fmt.Errorf("IsEmpty expectation only works with strings, slices, arrays, and maps, but got %T", actual)
	}

	return nil
}

// Description returns the expectation description
func (ie IsEmptyExpectation) Description() string {
	return "is empty"
}

// Convenience function for creating IsEmpty expectations
func IsEmpty() ensure.Expectation[interface{}] {
	return NewIsEmpty()
}

// ArrayLengthEqualsExpectation checks if an array/slice has the expected length
type ArrayLengthEqualsExpectation struct {
	expectedLength int
}

// NewArrayLengthEquals creates a new ArrayLengthEquals expectation
func NewArrayLengthEquals(expectedLength int) ensure.Expectation[interface{}] {
	return ArrayLengthEqualsExpectation{expectedLength: expectedLength}
}

// Evaluate evaluates the array length expectation
func (ale ArrayLengthEqualsExpectation) Evaluate(actual interface{}) error {
	val := reflect.ValueOf(actual)

	var length int
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		length = val.Len()
	case reflect.String:
		length = len(val.String())
	default:
		return fmt.Errorf("ArrayLengthEquals expectation only works with arrays, slices, and strings, but got %T", actual)
	}

	if length != ale.expectedLength {
		return fmt.Errorf("expected length to be %d, but got %d", ale.expectedLength, length)
	}
	return nil
}

// Description returns the expectation description
func (ale ArrayLengthEqualsExpectation) Description() string {
	return fmt.Sprintf("has length %d", ale.expectedLength)
}

// Convenience function for creating ArrayLengthEquals expectations
func ArrayLengthEquals(expectedLength int) ensure.Expectation[interface{}] {
	return NewArrayLengthEquals(expectedLength)
}
