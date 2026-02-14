package examples

import (
	"context"
	"testing"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestNewAPIDemonstration demonstrates the new TestContext API without require.NoError
func TestNewAPIDemonstration(t *testing.T) {
	// Create SerenityTest context - no more manual error handling!
	ctx := context.Background()
	test := serenity.NewSerenityTestWithContext(ctx, t)

	// Create actor through test context
	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Chain activities without require.NoError - errors are handled automatically!
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
	)

	// Multiple actors are supported
	user := test.ActorCalled("RegularUser").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	user.AttemptsTo(
		api.SendGetRequest("/users"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("email")),
	)
}
