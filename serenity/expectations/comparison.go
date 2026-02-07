package expectations

import (
	"fmt"

	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	"github.com/nchursin/serenity-go/serenity/expectations/utils"
)

// IsGreaterThanExpectation checks if a numeric value is greater than expected
type IsGreaterThanExpectation struct {
	expected interface{}
}

// NewIsGreaterThan creates a new IsGreaterThan expectation
func NewIsGreaterThan(expected interface{}) ensure.Expectation[interface{}] {
	return IsGreaterThanExpectation{expected: expected}
}

// Evaluate evaluates the greater than expectation
func (igt IsGreaterThanExpectation) Evaluate(actual interface{}) error {
	return utils.CompareValues(actual, igt.expected, ">")
}

// Description returns the expectation description
func (igt IsGreaterThanExpectation) Description() string {
	return fmt.Sprintf("is greater than %v", igt.expected)
}

// Convenience function for creating IsGreaterThan expectations
func IsGreaterThan(expected interface{}) ensure.Expectation[interface{}] {
	return NewIsGreaterThan(expected)
}

// IsLessThanExpectation checks if a numeric value is less than expected
type IsLessThanExpectation struct {
	expected interface{}
}

// NewIsLessThan creates a new IsLessThan expectation
func NewIsLessThan(expected interface{}) ensure.Expectation[interface{}] {
	return IsLessThanExpectation{expected: expected}
}

// Evaluate evaluates the less than expectation
func (ilt IsLessThanExpectation) Evaluate(actual interface{}) error {
	return utils.CompareValues(actual, ilt.expected, "<")
}

// Description returns the expectation description
func (ilt IsLessThanExpectation) Description() string {
	return fmt.Sprintf("is less than %v", ilt.expected)
}

// Convenience function for creating IsLessThan expectations
func IsLessThan(expected interface{}) ensure.Expectation[interface{}] {
	return NewIsLessThan(expected)
}
