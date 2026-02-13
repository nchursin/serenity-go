package core

import (
	"context"
	"fmt"
)

// This file provides concrete implementations of the Question interface
// defined in interfaces.go. These implementations enable type-safe
// queries about system state using Go generics.
//
// Key Implementation:
//
//	question[T] - Generic question implementation with custom ask function
//
// Factory Functions:
//
//	NewQuestion() - Creates a new question with explicit function
//	Of()          - Convenience function for creating questions
//
// Usage Examples:
//
//	// Create a question using NewQuestion
//	userCount := core.NewQuestion[int]("number of users", func(actor core.Actor) (int, error) {
//		db := actor.AbilityTo(&database.DatabaseAbility{}).(database.DatabaseAbility)
//		return db.QueryRow("SELECT COUNT(*) FROM users").Int()
//	})
//
//	// Create a question using Of (convenience)
//	userName := core.Of("current user name", func(actor core.Actor) (string, error) {
//		session := actor.AbilityTo(&auth.SessionAbility{}).(auth.SessionAbility)
//		return session.GetCurrentUser().Name
//	})
//
// Type Safety:
//
//	// Compile-time type checking ensures correct return types
//	var count int
//	var name string
//	var profile *UserProfile
//	var orders []Order
//	var isActive bool
//
//	// Each question returns its specific type
//	count, err := userCount.AnsweredBy(actor)    // int, error
//	name, err := userName.AnsweredBy(actor)      // string, error
//
// Using Questions with Expectations:
//
//	actor.AttemptsTo(
//		ensure.That(userCount, expectations.GreaterThan(0)),
//		ensure.That(userName, expectations.Contains("admin")),
//		ensure.That(userProfile, expectations.HasField("Email", expectations.IsNotEmpty())),
//	)
//

// question implements the Question interface for type-safe system queries.
// This generic implementation allows creating questions that return
// specific types while maintaining the Question interface contract.
//
// Type question is private - use NewQuestion() or Of() factory functions.
type question[T any] struct {
	// description provides a human-readable description of what the question asks
	description string

	// ask is the function that executes when the question is answered
	ask func(actor Actor, ctx context.Context) (T, error)
}

// NewQuestion creates a new question with the given description and ask function.
// This is the primary factory function for creating typed questions.
// Use this when you want explicit control over the question creation.
//
// Type Parameters:
//   - T: The type of answer this question returns
//
// Parameters:
//   - description: Human-readable description of what the question asks
//   - ask: Function that takes an actor and context, returns the typed answer
//
// Returns:
//   - Question[T]: A new question that returns type T when answered
//
// Usage Examples:
//
//	// Simple type question
//	userCount := core.NewQuestion[int]("number of users in system", func(actor core.Actor) (int, error) {
//		db, err := actor.AbilityTo(&database.DatabaseAbility{})
//		if err != nil {
//			return 0, fmt.Errorf("actor needs database ability: %w", err)
//		}
//		return db.QueryRow("SELECT COUNT(*) FROM users").Int()
//	})
//
//	// Complex type question
//	userProfile := core.NewQuestion[*UserProfile]("user profile", func(actor core.Actor) (*UserProfile, error) {
//		db, err := actor.AbilityTo(&database.DatabaseAbility{})
//		if err != nil {
//			return nil, fmt.Errorf("actor needs database ability: %w", err)
//		}
//		return db.GetUserProfile(actor.Name())
//	})
//
//	// Collection question
//	activeOrders := core.NewQuestion[[]Order]("active orders", func(actor core.Actor) ([]Order, error) {
//		api, err := actor.AbilityTo(&api.CallAnAPI{})
//		if err != nil {
//			return nil, fmt.Errorf("actor needs API ability: %w", err)
//		}
//		response, err := api.(api.CallAnAPI).Get("/orders?status=active")
//		if err != nil {
//			return nil, err
//		}
//		return parseOrders(response.Body)
//	})
//
//	// Boolean question
//	isSystemOnline := core.NewQuestion[bool]("system online status", func(actor core.Actor) (bool, error) {
//		health, err := actor.AbilityTo(&monitoring.HealthAbility{})
//		if err != nil {
//			return false, fmt.Errorf("actor needs health check ability: %w", err)
//		}
//		return health.(monitoring.HealthAbility).IsOnline()
//	})
//
// Using Created Questions:
//
//	count, err := userCount.AnsweredBy(actor)
//	if err != nil {
//		return fmt.Errorf("failed to get user count: %w", err)
//	}
//
//	profile, err := userProfile.AnsweredBy(actor)
//	if err != nil {
//		return fmt.Errorf("failed to get user profile: %w", err)
//	}
//
//	// With expectations
//	actor.AttemptsTo(
//		ensure.That(userCount, expectations.GreaterThan(0)),
//		ensure.That(isSystemOnline, expectations.IsTrue()),
//	)
func NewQuestion[T any](description string, ask func(actor Actor, ctx context.Context) (T, error)) Question[T] {
	return &question[T]{
		description: description,
		ask:         ask,
	}
}

// Description returns the question's human-readable description.
// The description is prefixed with "asks " to indicate it's a query.
//
// Returns:
//   - string: Description formatted as "asks [description]"
//
// Example:
//
//	q := core.Of("user count", getUserCount)
//	fmt.Println(q.Description()) // "asks user count"
func (q *question[T]) Description() string {
	return fmt.Sprintf("asks %s", q.description)
}

// AnsweredBy returns the answer when asked by the given actor.
// This method executes the ask function provided to NewQuestion().
//
// Parameters:
//   - actor: The actor asking the question
//   - ctx: Context for cancellation and timeout
//
// Returns:
//   - T: The typed answer to the question
//   - error: Error if the question cannot be answered
//
// Example:
//
//	func (q *userCountQuestion) AnsweredBy(actor core.Actor, ctx context.Context) (int, error) {
//		return q.ask(actor, ctx)
//	}
//
// Usage:
//
//	count, err := question.AnsweredBy(actor, ctx)
//	if err != nil {
//		return fmt.Errorf("failed to answer question '%s': %w", question.Description(), err)
//	}
//	fmt.Printf("Answer: %v\n", count)
func (q *question[T]) AnsweredBy(actor Actor, ctx context.Context) (T, error) {
	return q.ask(actor, ctx)
}

// Of creates a new question with the given description and ask function.
// This is a convenience function that internally calls NewQuestion().
// Use this as a shorter alternative to NewQuestion().
//
// Type Parameters:
//   - T: The type of answer this question returns
//
// Parameters:
//   - description: Human-readable description of what the question asks
//   - ask: Function that takes an actor and returns the typed answer
//
// Returns:
//   - Question[T]: A new question that returns type T when answered
//
// Usage Examples:
//
//	// Simple boolean question
//	isHealthy := core.Of("system health status", func(actor core.Actor) (bool, error) {
//		health := actor.AbilityTo(&monitoring.HealthAbility{})
//		return health.(monitoring.HealthAbility).IsHealthy()
//	})
//
//	// String question
//	currentUser := core.Of("current user name", func(actor core.Actor) (string, error) {
//		session := actor.AbilityTo(&auth.SessionAbility{})
//		return session.(auth.SessionAbility).GetCurrentUser().Name
//	})
//
//	// Integer question with calculation
//	averageResponseTime := core.Of("average response time", func(actor core.Actor) (time.Duration, error) {
//		metrics := actor.AbilityTo(&monitoring.MetricsAbility{})
//		return metrics.(monitoring.MetricsAbility).CalculateAverageResponseTime(time.Hour)
//	})
//
//	// Struct question
//	systemInfo := core.Of("system information", func(actor core.Actor) (*SystemInfo, error) {
//		info := &SystemInfo{}
//		health := actor.AbilityTo(&monitoring.HealthAbility{})
//		metrics := actor.AbilityTo(&monitoring.MetricsAbility{})
//
//		info.IsHealthy = health.(monitoring.HealthAbility).IsHealthy()
//		info.Uptime = metrics.(monitoring.MetricsAbility).GetUptime()
//		info.Version = metrics.(monitoring.MetricsAbility).GetVersion()
//
//		return info, nil
//	})
//
//	// Using in assertions
//	actor.AttemptsTo(
//		ensure.That(isHealthy, expectations.IsTrue()),
//		ensure.That(currentUser, expectations.Not(expectations.IsEmpty())),
//		ensure.That(averageResponseTime, expectations.LessThan(time.Second)),
//		ensure.That(systemInfo, expectations.HasField("Version", expectations.Not(expectations.IsEmpty()))),
//	)
//
// NewQuestion vs Of:
//
//	// Use NewQuestion when you want explicit creation
//	q1 := core.NewQuestion[int]("user count", getUserCount)
//
//	// Use Of for convenience (identical result)
//	q2 := core.Of("user count", getUserCount)
//
//	// Both create the same type of question
//	var q1, q2 core.Question[int]
func Of[T any](description string, ask func(actor Actor, ctx context.Context) (T, error)) Question[T] {
	return NewQuestion(description, ask)
}
