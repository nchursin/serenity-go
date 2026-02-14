package examples

import (
	"context"
	"fmt"
	"testing"

	"github.com/nchursin/serenity-go/serenity/core"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestCoreDoFunction demonstrates the new core.Do function for quick activity creation
func TestCoreDoFunction(t *testing.T) {
	ctx := context.Background()
	test := serenity.NewSerenityTestWithContext(ctx, t)
	defer test.Shutdown()

	actor := test.ActorCalled("TestActor")

	// Test the new core.Do function with FailFast mode
	actor.AttemptsTo(
		core.Do("#actor performs a simple action", func(actor core.Actor, ctx context.Context) error {
			// Simple test action
			t.Logf("Actor %s is performing a custom action", actor.Name())
			return nil
		}),
	)

	// Test core.Do with access to actor abilities
	actor.AttemptsTo(
		core.Do("#actor accesses actor information", func(actor core.Actor, ctx context.Context) error {
			// Verify we can access actor properties
			if actor.Name() != "TestActor" {
				return fmt.Errorf("expected actor name 'TestActor', got '%s'", actor.Name())
			}
			t.Logf("Successfully accessed actor: %s", actor.Name())
			return nil
		}),
	)
}
