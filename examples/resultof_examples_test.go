package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/nchursin/serenity-go/serenity/answerable"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestResultOf demonstrates usage of answerable.ResultOf for dynamic assertions
func TestResultOf(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("ResultOfTester")

	// Basic usage with different data types
	actor.AttemptsTo(
		// Simple calculation
		ensure.That(
			answerable.ResultOf("calculated value", func(actor core.Actor, ctx context.Context) (int, error) {
				return 42, nil
			}),
			expectations.Equals(42),
		),

		// String manipulation
		ensure.That(
			answerable.ResultOf("greeting message", func(actor core.Actor, ctx context.Context) (string, error) {
				return "Hello, " + actor.Name(), nil
			}),
			expectations.Contains("Hello, ResultOfTester"),
		),

		// Boolean logic
		ensure.That(
			answerable.ResultOf("validation result", func(actor core.Actor, ctx context.Context) (bool, error) {
				name := actor.Name()
				return len(name) > 5, nil
			}),
			expectations.Equals(true),
		),

		// Float calculation
		ensure.That(
			answerable.ResultOf("division result", func(actor core.Actor, ctx context.Context) (float64, error) {
				return 10.0 / 2.0, nil
			}),
			expectations.Equals(5.0),
		),
	)
}

// TestResultOfCalculations demonstrates using ResultOf for calculations
func TestResultOfCalculations(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	dataProcessor := test.ActorCalled("DataProcessor")

	// Example: Processing numeric data
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	dataProcessor.AttemptsTo(
		ensure.That(
			answerable.ResultOf("sum of even numbers", func(actor core.Actor, ctx context.Context) (int, error) {
				sum := 0
				for _, num := range numbers {
					if num%2 == 0 {
						sum += num
					}
				}
				return sum, nil
			}),
			expectations.Equals(30), // 2+4+6+8+10 = 30
		),

		ensure.That(
			answerable.ResultOf("average of numbers", func(actor core.Actor, ctx context.Context) (float64, error) {
				if len(numbers) == 0 {
					return 0, fmt.Errorf("no numbers to calculate average")
				}
				sum := 0
				for _, num := range numbers {
					sum += num
				}
				return float64(sum) / float64(len(numbers)), nil
			}),
			expectations.Equals(5.5),
		),

		ensure.That(
			answerable.ResultOf("maximum value", func(actor core.Actor, ctx context.Context) (int, error) {
				if len(numbers) == 0 {
					return 0, fmt.Errorf("no numbers to find maximum")
				}
				max := numbers[0]
				for _, num := range numbers {
					if num > max {
						max = num
					}
				}
				return max, nil
			}),
			expectations.Equals(10),
		),
	)
}

// TestResultWithErrorHandling demonstrates error handling in ResultOf functions
func TestResultWithErrorHandling(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("ErrorTestActor")

	// Example: Functions that can return errors
	actor.AttemptsTo(
		// This will succeed
		ensure.That(
			answerable.ResultOf("successful operation", func(actor core.Actor, ctx context.Context) (string, error) {
				return "success", nil
			}),
			expectations.Equals("success"),
		),

		// Example: Conditional logic with error handling
		ensure.That(
			answerable.ResultOf("safe division", func(actor core.Actor, ctx context.Context) (float64, error) {
				numerator, denominator := 10, 2
				if denominator == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return float64(numerator) / float64(denominator), nil
			}),
			expectations.Equals(5.0),
		),

		// Example: String manipulation with validation
		ensure.That(
			answerable.ResultOf("validated string", func(actor core.Actor, ctx context.Context) (string, error) {
				input := "test"
				if input == "" {
					return "", fmt.Errorf("input cannot be empty")
				}
				return fmt.Sprintf("processed: %s", input), nil
			}),
			expectations.Equals("processed: test"),
		),
	)
}

// TestResultOfWithActor demonstrates using actor properties in ResultOf
func TestResultOfWithActor(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor1 := test.ActorCalled("Actor1")
	actor2 := test.ActorCalled("Actor2")

	// Both actors use the same ResultOf question but get different results
	greetingQuestion := answerable.ResultOf("personalized greeting", func(actor core.Actor, ctx context.Context) (string, error) {
		return fmt.Sprintf("Hello, %s!", actor.Name()), nil
	})

	actor1.AttemptsTo(
		ensure.That(greetingQuestion, expectations.Equals("Hello, Actor1!")),
	)

	actor2.AttemptsTo(
		ensure.That(greetingQuestion, expectations.Equals("Hello, Actor2!")),
	)

	// Different actors can have different calculations
	nameLengthQuestion := answerable.ResultOf("name length", func(actor core.Actor, ctx context.Context) (int, error) {
		return len(actor.Name()), nil
	})

	actor1.AttemptsTo(
		ensure.That(nameLengthQuestion, expectations.Equals(6)), // "Actor1"
	)

	actor2.AttemptsTo(
		ensure.That(nameLengthQuestion, expectations.Equals(6)), // "Actor2"
	)
}

// TestResultOfComplexTypes demonstrates using ResultOf with complex data structures
func TestResultOfComplexTypes(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("ComplexTypeTester")

	type User struct {
		Name  string
		Age   int
		Email string
	}

	actor.AttemptsTo(
		// Create and validate a struct
		ensure.That(
			answerable.ResultOf("user creation", func(actor core.Actor, ctx context.Context) (User, error) {
				return User{
					Name:  "John Doe",
					Age:   30,
					Email: "john@example.com",
				}, nil
			}),
			expectations.Equals(User{Name: "John Doe", Age: 30, Email: "john@example.com"}),
		),

		// Work with slices
		ensure.That(
			answerable.ResultOf("user names list", func(actor core.Actor, ctx context.Context) ([]string, error) {
				users := []User{
					{Name: "Alice", Age: 25, Email: "alice@example.com"},
					{Name: "Bob", Age: 30, Email: "bob@example.com"},
					{Name: "Charlie", Age: 35, Email: "charlie@example.com"},
				}

				names := make([]string, len(users))
				for i, user := range users {
					names[i] = user.Name
				}
				return names, nil
			}),
			expectations.Satisfies("contains Bob", func(names []string) error {
				for _, name := range names {
					if name == "Bob" {
						return nil
					}
				}
				return fmt.Errorf("expected to find Bob in names, but got %v", names)
			}),
		),

		// Work with maps
		ensure.That(
			answerable.ResultOf("user data map", func(actor core.Actor, ctx context.Context) (map[string]interface{}, error) {
				return map[string]interface{}{
					"name":    "Jane Doe",
					"age":     28,
					"active":  true,
					"balance": 1000.50,
				}, nil
			}),
			expectations.Equals(map[string]interface{}{
				"name":    "Jane Doe",
				"age":     28,
				"active":  true,
				"balance": 1000.50,
			}),
		),

		// Pointer operations
		ensure.That(
			answerable.ResultOf("pointer to string", func(actor core.Actor, ctx context.Context) (*string, error) {
				message := "Hello from pointer"
				return &message, nil
			}),
			expectations.Satisfies("is not nil", func(p *string) error {
				if p == nil {
					return fmt.Errorf("expected pointer to not be nil")
				}
				return nil
			}),
		),
	)
}

// TestResultOfMixedWithStatic demonstrates mixing ResultOf with ValueOf
func TestResultOfMixedWithStatic(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("MixedTester")

	// Mix of static values and ResultOf functions
	actor.AttemptsTo(
		// Static value
		ensure.That(answerable.ValueOf(42), expectations.Equals(42)),

		// Dynamic calculation using static values
		ensure.That(
			answerable.ResultOf("double calculation", func(actor core.Actor, ctx context.Context) (int, error) {
				staticValue := 42
				return staticValue * 2, nil
			}),
			expectations.Equals(84),
		),

		// Complex calculation
		ensure.That(
			answerable.ResultOf("complex math", func(actor core.Actor, ctx context.Context) (int, error) {
				a := 10
				b := 20
				c := a + b  // 30
				d := c * 2  // 60
				e := d - 10 // 50
				return e, nil
			}),
			expectations.Equals(50),
		),

		// String operations
		ensure.That(
			answerable.ResultOf("string building", func(actor core.Actor, ctx context.Context) (string, error) {
				parts := []string{"Hello", "from", "ResultOf"}
				result := ""
				for i, part := range parts {
					if i > 0 {
						result += " "
					}
					result += part
				}
				return result, nil
			}),
			expectations.Equals("Hello from ResultOf"),
		),
	)
}

// TestResultOfEdgeCases demonstrates edge cases and special scenarios
func TestResultOfEdgeCases(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	actor := test.ActorCalled("EdgeCaseTester")

	actor.AttemptsTo(
		// Empty slice
		ensure.That(
			answerable.ResultOf("empty slice", func(actor core.Actor, ctx context.Context) ([]int, error) {
				return []int{}, nil
			}),
			expectations.Equals([]int{}),
		),

		// Nil pointer
		ensure.That(
			answerable.ResultOf("nil pointer", func(actor core.Actor, ctx context.Context) (*int, error) {
				return nil, nil
			}),
			expectations.Equals((*int)(nil)),
		),

		// Empty string
		ensure.That(
			answerable.ResultOf("empty string", func(actor core.Actor, ctx context.Context) (string, error) {
				return "", nil
			}),
			expectations.Equals(""),
		),

		// Zero values
		ensure.That(
			answerable.ResultOf("zero int", func(actor core.Actor, ctx context.Context) (int, error) {
				return 0, nil
			}),
			expectations.Equals(0),
		),

		ensure.That(
			answerable.ResultOf("zero float", func(actor core.Actor, ctx context.Context) (float64, error) {
				return 0.0, nil
			}),
			expectations.Equals(0.0),
		),

		ensure.That(
			answerable.ResultOf("zero bool", func(actor core.Actor, ctx context.Context) (bool, error) {
				return false, nil
			}),
			expectations.Equals(false),
		),
	)
}
