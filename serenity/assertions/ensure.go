package assertions

import (
	"fmt"
	"reflect"
	"strings"

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

// EnsureBuilder helps build assertions with fluent interface (deprecated, kept for compatibility)
type EnsureBuilder[T any] struct {
	question core.Question[T]
}

// That creates a new Ensure assertion with the new API
func That[T any](question core.Question[T], expectation Expectation[T]) core.Activity {
	return &EnsureActivity[T]{
		question:    question,
		expectation: expectation,
	}
}

// Is sets the expectation for the assertion (deprecated, use That(question, expectation) instead)
func (eb *EnsureBuilder[T]) Is(expectation Expectation[T]) core.Activity {
	return That(eb.question, expectation)
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

// Basic expectation implementations

// EqualsExpectation checks if the actual value equals the expected value
type EqualsExpectation[T any] struct {
	expected T
}

// NewEquals creates a new Equals expectation
func NewEquals[T any](expected T) EqualsExpectation[T] {
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

// ContainsExpectation checks if a string contains the expected substring
type ContainsExpectation struct {
	substring string
}

// NewContains creates a new Contains expectation
func NewContains(substring string) ContainsExpectation {
	return ContainsExpectation{substring: substring}
}

// Evaluate evaluates the contains expectation
func (c ContainsExpectation) Evaluate(actual string) error {
	if !strings.Contains(actual, c.substring) {
		return fmt.Errorf("expected string to contain '%s', but got '%s'", c.substring, actual)
	}
	return nil
}

// Description returns the expectation description
func (c ContainsExpectation) Description() string {
	return fmt.Sprintf("contains '%s'", c.substring)
}

// ContainsKeyExpectation checks if a map contains the expected key
type ContainsKeyExpectation struct {
	key string
}

// NewContainsKey creates a new ContainsKey expectation
func NewContainsKey(key string) ContainsKeyExpectation {
	return ContainsKeyExpectation{key: key}
}

// Evaluate evaluates the contains key expectation
func (ck ContainsKeyExpectation) Evaluate(actual interface{}) error {
	val := reflect.ValueOf(actual)
	if val.Kind() != reflect.Map {
		return fmt.Errorf("expected a map, but got %T", actual)
	}

	// Try to convert to map[string]interface{} for string keys
	if mapStr, ok := actual.(map[string]interface{}); ok {
		if _, exists := mapStr[ck.key]; !exists {
			return fmt.Errorf("expected map to contain key '%s'", ck.key)
		}
		return nil
	}

	// Fallback to reflection for any map type
	mapKey := reflect.ValueOf(ck.key)
	if !val.MapIndex(mapKey).IsValid() {
		return fmt.Errorf("expected map to contain key '%s'", ck.key)
	}
	return nil
}

// Description returns the expectation description
func (ck ContainsKeyExpectation) Description() string {
	return fmt.Sprintf("contains key '%s'", ck.key)
}

// IsEmptyExpectation checks if a string is empty
type IsEmptyExpectation struct{}

// NewIsEmpty creates a new IsEmpty expectation
func NewIsEmpty() IsEmptyExpectation {
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

// ArrayLengthEqualsExpectation checks if an array/slice has the expected length
type ArrayLengthEqualsExpectation struct {
	expectedLength int
}

// NewArrayLengthEquals creates a new ArrayLengthEquals expectation
func NewArrayLengthEquals(expectedLength int) ArrayLengthEqualsExpectation {
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

// IsGreaterThanExpectation checks if a numeric value is greater than expected
type IsGreaterThanExpectation struct {
	expected interface{}
}

// NewIsGreaterThan creates a new IsGreaterThan expectation
func NewIsGreaterThan(expected interface{}) IsGreaterThanExpectation {
	return IsGreaterThanExpectation{expected: expected}
}

// Evaluate evaluates the greater than expectation
func (igt IsGreaterThanExpectation) Evaluate(actual interface{}) error {
	return compareValues(actual, igt.expected, ">")
}

// Description returns the expectation description
func (igt IsGreaterThanExpectation) Description() string {
	return fmt.Sprintf("is greater than %v", igt.expected)
}

// IsLessThanExpectation checks if a numeric value is less than expected
type IsLessThanExpectation struct {
	expected interface{}
}

// NewIsLessThan creates a new IsLessThan expectation
func NewIsLessThan(expected interface{}) IsLessThanExpectation {
	return IsLessThanExpectation{expected: expected}
}

// Evaluate evaluates the less than expectation
func (ilt IsLessThanExpectation) Evaluate(actual interface{}) error {
	return compareValues(actual, ilt.expected, "<")
}

// Description returns the expectation description
func (ilt IsLessThanExpectation) Description() string {
	return fmt.Sprintf("is less than %v", ilt.expected)
}

// Helper function to compare numeric values
func compareValues(actual, expected interface{}, operator string) error {
	actualFloat, err := toFloat64(actual)
	if err != nil {
		return fmt.Errorf("cannot compare actual value: %w", err)
	}

	expectedFloat, err := toFloat64(expected)
	if err != nil {
		return fmt.Errorf("cannot compare expected value: %w", err)
	}

	switch operator {
	case ">":
		if actualFloat <= expectedFloat {
			return fmt.Errorf("expected value to be greater than %v, but got %v", expected, actual)
		}
	case "<":
		if actualFloat >= expectedFloat {
			return fmt.Errorf("expected value to be less than %v, but got %v", expected, actual)
		}
	}

	return nil
}

// Helper function to convert values to float64 for comparison
func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported numeric type: %T", value)
	}
}

// Convenience functions for creating expectations
func Equals[T any](expected T) EqualsExpectation[T] {
	return NewEquals(expected)
}

func Contains(substring string) ContainsExpectation {
	return NewContains(substring)
}

func ContainsKey(key string) ContainsKeyExpectation {
	return NewContainsKey(key)
}

func IsEmpty() IsEmptyExpectation {
	return NewIsEmpty()
}

func ArrayLengthEquals(expectedLength int) ArrayLengthEqualsExpectation {
	return NewArrayLengthEquals(expectedLength)
}

func IsGreaterThan(expected interface{}) IsGreaterThanExpectation {
	return NewIsGreaterThan(expected)
}

func IsLessThan(expected interface{}) IsLessThanExpectation {
	return NewIsLessThan(expected)
}

// Ensure interaction implementation
var Ensure = ensureBuilder{}

type ensureBuilder struct{}

// StringThat creates a new Ensure assertion for string questions
func (eb ensureBuilder) StringThat(question core.Question[string]) *EnsureBuilder[string] {
	return &EnsureBuilder[string]{question: question}
}

// IntThat creates a new Ensure assertion for int questions
func (eb ensureBuilder) IntThat(question core.Question[int]) *EnsureBuilder[int] {
	return &EnsureBuilder[int]{question: question}
}

// AnyThat creates a new Ensure assertion for any type questions
func (eb ensureBuilder) AnyThat(question core.Question[interface{}]) *EnsureBuilder[interface{}] {
	return &EnsureBuilder[interface{}]{question: question}
}
