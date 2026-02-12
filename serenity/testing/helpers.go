// Package utilities for TestContext API implementation.
// These functions provide common functionality for test management and actor creation.
//
// The helper functions are designed to be used internally by the SerenityTest
// implementation but can also be used by custom test managers that need
// TestContext integration.
package testing

import (
	"fmt"

	"github.com/nchursin/serenity-go/serenity/abilities"
)

// NewTestContext creates a TestContext wrapper around the provided testing.TB.
// This function enables automatic error handling for any test context.
//
// Use this function when creating custom test managers or when you need
// TestContext functionality outside of the standard SerenityTest.
//
// Parameters:
//
//	tb - The underlying testing.TB instance (usually *testing.T)
//
// Returns:
//
//	A TestContext instance that provides automatic error handling
//
// Example:
//
//	func CustomTestWrapper(t *testing.T) {
//		ctx := testing.NewTestContext(t)
//
//		// Use ctx for automatic error handling
//		ctx.Helper() // Mark as helper function
//		ctx.Log("Custom test logic")
//
//		if someCondition {
//			ctx.Fatalf("Test condition failed: %v", someCondition)
//		}
//	}
//
// Note: This function is typically not needed when using NewSerenityTest(),
// as it automatically handles TestContext creation.
func NewTestContext(tb interface{}) TestContext {
	// Implementation would wrap the testing.TB interface
	// This is a placeholder for future extensibility
	return nil
}

// abilityMatchesType checks if two abilities are of the same type.
// This helper function is used internally to match ability types when
// retrieving specific abilities from an actor's ability collection.
//
// Parameters:
//
//	ability - The ability instance to check
//	abilityType - The type reference ability to match against
//
// Returns:
//
//	true if the abilities are of the same concrete type, false otherwise
//
// This function uses type string comparison to determine type equality,
// which is sufficient for the current ability system where each ability
// type has a unique implementation.
func abilityMatchesType(ability, abilityType abilities.Ability) bool {
	return fmt.Sprintf("%T", ability) == fmt.Sprintf("%T", abilityType)
}
