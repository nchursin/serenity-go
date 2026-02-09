package examples

import (
	"testing"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestIntentionalFailure demonstrates error handling with wrong assertion
func TestIntentionalFailure(t *testing.T) {
	test := serenity.NewSerenityTest(t)
	defer test.Shutdown()

	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// This should fail - wrong status code (expecting 404 but getting 200)
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(404)), // This will fail
	)
}
