package testing

import (
	"fmt"

	"github.com/nchursin/serenity-go/serenity/abilities"
)

// Helper function to check if ability matches type
func abilityMatchesType(ability, abilityType abilities.Ability) bool {
	return fmt.Sprintf("%T", ability) == fmt.Sprintf("%T", abilityType)
}
