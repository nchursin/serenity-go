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
