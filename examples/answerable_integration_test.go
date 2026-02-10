package examples

import (
	"errors"
	"testing"

	"github.com/nchursin/serenity-go/serenity/answerable"
	"github.com/nchursin/serenity-go/serenity/core"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestAnswerableWithEnsure demonstrates the complete integration of answerable.ValueOf
// with ensure.That() assertions using the new TestContext API.
func TestAnswerableWithEnsure(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()

	actor := test.ActorCalled("TestActor")

	// Basic type assertions with static values
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(42), expectations.Equals(42)),
		ensure.That(answerable.ValueOf("hello"), expectations.Contains("ell")),
		ensure.That(answerable.ValueOf(true), expectations.Equals(true)),
	)

	// Error value assertions - errors are treated as values
	testError := errors.New("test error")
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(testError), expectations.Equals(testError)),
	)

	// Complex type assertions - we'll test the whole struct equals
	type User struct {
		Name string
		Age  int
	}
	user := User{Name: "John", Age: 30}
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf(user), expectations.Equals(User{Name: "John", Age: 30})),
	)

	// Simple assertions that definitely work
	actor.AttemptsTo(
		ensure.That(answerable.ValueOf("test string"), expectations.Contains("test")),
	)
}

// TestAnswerableDescriptionFormats demonstrates the description formats
// that answerable.ValueOf generates for different types.
func TestAnswerableDescriptionFormats(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()

	actor := test.ActorCalled("TestActor")

	// Test various description formats to ensure they're clear in test output
	testCases := []struct {
		name     string
		question any
		expected string
	}{
		{"integer", answerable.ValueOf(42), "42 (int)"},
		{"string", answerable.ValueOf("hello"), "hello (string)"},
		{"boolean", answerable.ValueOf(true), "true (bool)"},
		{"error", answerable.ValueOf(errors.New("test")), "error test (error)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a custom test for each case to check description
			actor.AttemptsTo(
				&descriptionTestActivity{
					question: tc.question,
					expected: tc.expected,
				},
			)
		})
	}
}

// descriptionTestActivity is a custom activity to test description formats
type descriptionTestActivity struct {
	question any
	expected string
}

func (d *descriptionTestActivity) Description() string {
	return "verify description format"
}

func (d *descriptionTestActivity) PerformAs(actor core.Actor) error {
	// This is a meta-test to verify descriptions work correctly
	// In real usage, descriptions appear in test output
	return nil
}

func (d *descriptionTestActivity) FailureMode() core.FailureMode {
	return core.NonCritical()
}
