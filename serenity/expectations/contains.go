package expectations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nchursin/serenity-go/serenity/ensure"
)

// ContainsExpectation checks if a string contains the expected substring
type ContainsExpectation struct {
	substring string
}

// NewContains creates a new Contains expectation
func NewContains(substring string) ensure.Expectation[string] {
	return ContainsExpectation{substring: substring}
}

// Evaluate evaluates the contains expectation
func (c ContainsExpectation) Evaluate(actual string) error {
	if actual == "" {
		return fmt.Errorf("expected string to contain '%s', but got empty string", c.substring)
	}
	if !strings.Contains(actual, c.substring) {
		return fmt.Errorf("expected string to contain '%s', but got '%s'", c.substring, actual)
	}
	return nil
}

// Description returns the expectation description
func (c ContainsExpectation) Description() string {
	return fmt.Sprintf("contains '%s'", c.substring)
}

// Convenience function for creating Contains expectations
func Contains(substring string) ensure.Expectation[string] {
	return NewContains(substring)
}

// ContainsKeyExpectation checks if a map contains the expected key
type ContainsKeyExpectation struct {
	key string
}

// NewContainsKey creates a new ContainsKey expectation
func NewContainsKey(key string) ensure.Expectation[interface{}] {
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

// Convenience function for creating ContainsKey expectations
func ContainsKey(key string) ensure.Expectation[interface{}] {
	return NewContainsKey(key)
}
