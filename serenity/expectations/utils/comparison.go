package utils

import (
	"fmt"
)

// compareValues compares numeric values using the specified operator
func CompareValues(actual, expected interface{}, operator string) error {
	actualFloat, err := ToFloat64(actual)
	if err != nil {
		return fmt.Errorf("cannot compare actual value: %w", err)
	}

	expectedFloat, err := ToFloat64(expected)
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

// ToFloat64 converts values to float64 for comparison
func ToFloat64(value interface{}) (float64, error) {
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
