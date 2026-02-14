package examples

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/nchursin/serenity-go/serenity/answerable"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestSatisfiesBasic demonstrates basic usage of Satisfies expectation
func TestSatisfiesBasic(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("BasicTester")

	// Test positive number validation
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(42), expectations.Satisfies("is positive number", func(actual int) error {
			if actual <= 0 {
				return fmt.Errorf("expected positive value, but got %d", actual)
			}
			return nil
		})),
	)

	// Test string validation
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf("hello@world.com"), expectations.Satisfies("contains valid email", func(actual string) error {
			if !strings.Contains(actual, "@") {
				return fmt.Errorf("missing @ symbol in email")
			}
			if !strings.Contains(actual, ".") {
				return fmt.Errorf("missing domain in email")
			}
			return nil
		})),
	)
}

// TestSatisfiesWithStructs demonstrates custom struct validation
func TestSatisfiesWithStructs(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("StructTester")

	type User struct {
		Name string
		Age  int
	}

	validUser := User{Name: "Alice", Age: 25}
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(validUser), expectations.Satisfies("has valid user data", func(actual User) error {
			if actual.Name == "" {
				return fmt.Errorf("name is empty")
			}
			if actual.Age < 18 {
				return fmt.Errorf("age %d is too young (minimum 18)", actual.Age)
			}
			if actual.Age > 100 {
				return fmt.Errorf("age %d is unrealistic (maximum 100)", actual.Age)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithCmpStructComparison demonstrates struct comparison using go-cmp
func TestSatisfiesWithCmpStructComparison(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("CmpStructTester")

	type User struct {
		Name string
		Age  int
		Tags []string
	}

	expected := User{Name: "Alice", Age: 25, Tags: []string{"admin", "user"}}
	actual := User{Name: "Alice", Age: 25, Tags: []string{"admin", "user"}}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(actual), expectations.Satisfies("matches expected user structure", func(actual User) error {
			if diff := cmp.Diff(expected, actual); diff != "" {
				return fmt.Errorf("user struct mismatch (-expected +actual):\n%s", diff)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithCmpWithOptions demonstrates advanced cmp usage with options
func TestSatisfiesWithCmpWithOptions(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("CmpOptionsTester")

	type TimestampedUser struct {
		ID        int
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	expected := TimestampedUser{
		ID:        1,
		Name:      "Bob",
		CreatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	actual := TimestampedUser{
		ID:        1,
		Name:      "Bob",
		CreatedAt: time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC),  // Different timestamp
		UpdatedAt: time.Date(2023, 1, 1, 12, 30, 0, 0, time.UTC), // Different timestamp
	}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(actual), expectations.Satisfies("matches user ignoring timestamps", func(actual TimestampedUser) error {
			if diff := cmp.Diff(expected, actual,
				cmpopts.IgnoreFields(TimestampedUser{}, "CreatedAt", "UpdatedAt"),
				cmpopts.EquateEmpty()); diff != "" {
				return fmt.Errorf("user struct mismatch (-expected +actual):\n%s", diff)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithCmpSliceComparison demonstrates slice comparison with sorting
func TestSatisfiesWithCmpSliceComparison(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("CmpSliceTester")

	type Item struct {
		ID   int
		Name string
	}

	expected := []Item{
		{ID: 2, Name: "item2"},
		{ID: 1, Name: "item1"},
	}

	actual := []Item{
		{ID: 1, Name: "item1"},
		{ID: 2, Name: "item2"},
	}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(actual), expectations.Satisfies("matches items ignoring order", func(actual []Item) error {
			if diff := cmp.Diff(expected, actual,
				cmpopts.SortSlices(func(a, b Item) bool { return a.ID < b.ID })); diff != "" {
				return fmt.Errorf("items slice mismatch (-expected +actual):\n%s", diff)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithCmpTransform demonstrates transformer usage
func TestSatisfiesWithCmpTransform(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("CmpTransformTester")

	type APIResponse struct {
		Data struct {
			User struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"user"`
		} `json:"data"`
	}

	type User struct {
		FirstName string
		LastName  string
	}

	expected := User{FirstName: "John", LastName: "Doe"}
	actual := APIResponse{
		Data: struct {
			User struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"user"`
		}{
			User: struct {
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			}{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
	}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(actual), expectations.Satisfies("matches expected user data", func(actual APIResponse) error {
			if diff := cmp.Diff(expected, User{
				FirstName: actual.Data.User.FirstName,
				LastName:  actual.Data.User.LastName,
			}); diff != "" {
				return fmt.Errorf("user data mismatch (-expected +actual):\n%s", diff)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithComplexValidation demonstrates complex business logic validation
func TestSatisfiesWithComplexValidation(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("ComplexValidatorTester")

	type OrderItem struct {
		ProductID string
		Quantity  int
		Price     float64
	}

	type Order struct {
		ID        string
		Amount    float64
		Currency  string
		Status    string
		CreatedAt time.Time
		Items     []OrderItem
	}

	order := Order{
		ID:        "ORD-123",
		Amount:    150.0,
		Currency:  "USD",
		Status:    "confirmed",
		CreatedAt: time.Now(),
		Items: []OrderItem{
			{ProductID: "PROD-1", Quantity: 2, Price: 50.0},
			{ProductID: "PROD-2", Quantity: 1, Price: 50.0},
		},
	}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(order), expectations.Satisfies("has valid order data", func(actual Order) error {
			// Validate ID format
			if !strings.HasPrefix(actual.ID, "ORD-") {
				return fmt.Errorf("order ID must start with ORD-, got %s", actual.ID)
			}

			// Validate amount
			if actual.Amount <= 0 {
				return fmt.Errorf("order amount must be positive, got %f", actual.Amount)
			}

			// Validate currency
			validCurrencies := []string{"USD", "EUR", "GBP"}
			currencyValid := false
			for _, currency := range validCurrencies {
				if actual.Currency == currency {
					currencyValid = true
					break
				}
			}
			if !currencyValid {
				return fmt.Errorf("invalid currency %s, valid options: %v", actual.Currency, validCurrencies)
			}

			// Validate status
			validStatuses := []string{"pending", "confirmed", "shipped", "delivered"}
			statusValid := false
			for _, status := range validStatuses {
				if actual.Status == status {
					statusValid = true
					break
				}
			}
			if !statusValid {
				return fmt.Errorf("invalid status %s, valid options: %v", actual.Status, validStatuses)
			}

			// Validate items
			if len(actual.Items) == 0 {
				return fmt.Errorf("order must have at least one item")
			}

			// Calculate total and validate
			var calculatedTotal float64
			for _, item := range actual.Items {
				if item.Quantity <= 0 {
					return fmt.Errorf("item quantity must be positive, got %d for product %s", item.Quantity, item.ProductID)
				}
				if item.Price <= 0 {
					return fmt.Errorf("item price must be positive, got %f for product %s", item.Price, item.ProductID)
				}
				calculatedTotal += float64(item.Quantity) * item.Price
			}

			// Allow small floating point differences
			if diff := actual.Amount - calculatedTotal; diff > 0.01 || diff < -0.01 {
				return fmt.Errorf("order amount %f doesn't match calculated total %f", actual.Amount, calculatedTotal)
			}

			return nil
		})),
	)
}

// TestSatisfiesErrorMessages demonstrates how error messages appear in test output
func TestSatisfiesErrorMessages(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("ErrorMessagesTester")

	// This will fail and show the custom error message
	// Uncomment to see error message format:
	/*
		actor.AttemptsTo(
			ensure.That(answerable.ValueOf(-5), expectations.Satisfies("is positive number", func(actual int) error {
				if actual <= 0 {
					return fmt.Errorf("expected positive value, but got %d", actual)
				}
				return nil
			})),
		)
	*/

	// Simple test to ensure actor is used
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(5), expectations.Satisfies("is positive number", func(actual int) error {
			if actual <= 0 {
				return fmt.Errorf("expected positive value, but got %d", actual)
			}
			return nil
		})),
	)
}

// TestSatisfiesWithMaps demonstrates map validation
func TestSatisfiesWithMaps(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("MapValidatorTester")

	config := map[string]interface{}{
		"database": map[string]interface{}{
			"host":     "localhost",
			"port":     5432,
			"ssl_mode": "require",
		},
		"logging": map[string]interface{}{
			"level":  "info",
			"format": "json",
		},
		"features": []interface{}{"auth", "api", "web"},
	}

	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(config), expectations.Satisfies("has valid configuration", func(actual map[string]interface{}) error {
			// Check required sections
			requiredSections := []string{"database", "logging", "features"}
			for _, section := range requiredSections {
				if _, exists := actual[section]; !exists {
					return fmt.Errorf("missing required config section: %s", section)
				}
			}

			// Validate database config
			dbConfig, ok := actual["database"].(map[string]interface{})
			if !ok {
				return fmt.Errorf("database config must be a map")
			}

			requiredDBFields := []string{"host", "port", "ssl_mode"}
			for _, field := range requiredDBFields {
				if _, exists := dbConfig[field]; !exists {
					return fmt.Errorf("missing required database field: %s", field)
				}
			}

			// Validate port
			port, ok := dbConfig["port"].(int)
			if !ok {
				return fmt.Errorf("database port must be an integer")
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("database port must be between 1 and 65535, got %d", port)
			}

			// Validate features
			features, ok := actual["features"].([]interface{})
			if !ok {
				return fmt.Errorf("features must be an array")
			}
			if len(features) == 0 {
				return fmt.Errorf("features array cannot be empty")
			}

			return nil
		})),
	)
}
