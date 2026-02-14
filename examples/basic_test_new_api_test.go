package examples

import (
	"context"
	"testing"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestJSONPlaceholderBasicsNewAPI demonstrates basic API testing with JSONPlaceholder using new TestContext API
func TestJSONPlaceholderBasicsNewAPI(t *testing.T) {
	ctx := context.Background()
	test := serenity.NewSerenityTestWithContext(ctx, t)

	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// Test GET posts - should return existing posts
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("title")),
	)

	// Test GET users - should return existing users
	apiTester.AttemptsTo(
		api.SendGetRequest("/users"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
		ensure.That(api.LastResponseBody{}, expectations.Contains("email")),
	)

	// The console output will now show detailed step-by-step execution
	// with emojis, timing, and activity tracking thanks to ConsoleReporter
}
