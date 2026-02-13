// Package answerable provides utilities for converting static values and functions into core.Question[T] objects.
//
// The primary use case is to enable the use of static values and dynamic calculations in ensure.That() assertions:
//
//	// Static values
//	ensure.That(answerable.ValueOf(4), expectations.Equals(4))
//	ensure.That(answerable.ValueOf("hello"), expectations.Contains("ell"))
//	ensure.That(answerable.ValueOf(user), expectations.HasField("Name", "John"))
//
//	// Dynamic functions
//	ensure.That(answerable.ResultOf("user count", func(actor core.Actor) (int, error) {
//		db := actor.AbilityTo(DatabaseAbility{}).(DatabaseAbility)
//		return db.CountUsers(), nil
//	}), expectations.GreaterThan(0))
//
// ValueOf creates a Question[T] that returns the static value when asked by any actor.
// This is particularly useful when you want to assert against static values rather than
// dynamic system state.
//
// ResultOf creates a Question[T] from a function with a custom description. This allows
// for dynamic calculations and complex operations directly in assertions while maintaining
// readable test reports.
//
// Static value examples:
//
//	// Basic types
//	ensure.That(answerable.ValueOf(42), expectations.Equals(42))
//	ensure.That(answerable.ValueOf("test"), expectations.Contains("es"))
//	ensure.That(answerable.ValueOf(true), expectations.Equals(true))
//
//	// Complex types
//	user := User{Name: "John", Age: 30}
//	ensure.That(answerable.ValueOf(user), expectations.HasField("Name", "John"))
//
//	// Error values (errors are treated as values, not as failures)
//	err := fmt.Errorf("something went wrong")
//	ensure.That(answerable.ValueOf(err), expectations.Equals(err))
//
//	// Pointers and nil handling
//	var ptr *string
//	ensure.That(answerable.ValueOf(ptr), expectations.IsNil())
//
// Dynamic function examples:
//
//	// Simple calculations
//	ensure.That(answerable.ResultOf("calculated age", func(actor core.Actor) (int, error) {
//		return 25, nil
//	}), expectations.Equals(25))
//
//	// Using actor properties
//	ensure.That(answerable.ResultOf("actor greeting", func(actor core.Actor) (string, error) {
//		return "Hello, " + actor.Name(), nil
//	}), expectations.Contains("Hello"))
//
//	// Complex operations with error handling
//	ensure.That(answerable.ResultOf("user from database", func(actor core.Actor) (*User, error) {
//		db := actor.AbilityTo(DatabaseAbility{}).(DatabaseAbility)
//		return db.GetUser("123")
//	}), expectations.NotNil())
//
// The created Question[T] from ValueOf is independent of any actor context - it will always
// return the same static value regardless of which actor asks the question.
//
// The created Question[T] from ResultOf executes the provided function each time it is asked,
// allowing for dynamic behavior and actor-dependent calculations.
package answerable

import (
	"context"

	"github.com/nchursin/serenity-go/serenity/core"
)

// ValueOf creates a core.Question[T] that returns the provided static value
// when answered by any actor.
//
// The value is treated as-is, even if it's an error type. This means that
// error values are passed through as the answer rather than being treated
// as failure conditions.
//
// Parameters:
//   - value: The static value to be wrapped as a Question
//
// Returns:
//   - core.Question[T]: A question that always returns the provided value
//
// Example:
//
//	q := answerable.ValueOf(42)
//	result, err := q.AnsweredBy(actor) // result = 42, err = nil
func ValueOf[T any](value T) core.Question[T] {
	return &valueQuestion[T]{value: value}
}

// ResultOf creates a core.Question[T] from a function with a custom description.
//
// The function is executed when the question is answered by an actor.
// This allows for dynamic calculations and complex operations in assertions.
//
// Parameters:
//   - description: Human-readable description for test reports
//   - fn: Function that takes an actor and context, returns (value, error)
//
// Returns:
//   - core.Question[T]: A question that executes the function when answered
//
// Example:
//
//	ensure.That(
//		answerable.ResultOf("user count", func(actor core.Actor, ctx context.Context) (int, error) {
//			db := actor.AbilityTo(DatabaseAbility{}).(DatabaseAbility)
//			return db.CountUsers(), nil
//		}),
//		expectations.GreaterThan(0),
//	)
func ResultOf[T any](description string, fn func(core.Actor, context.Context) (T, error)) core.Question[T] {
	if fn == nil {
		panic("ResultOf: function parameter cannot be nil")
	}
	return &functionQuestion[T]{
		description: description,
		function:    fn,
	}
}
