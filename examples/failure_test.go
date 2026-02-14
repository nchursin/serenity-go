package examples

import (
	"testing"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
)

// TestFailureHandling demonstrates how different failure modes work
func TestFailureHandling(t *testing.T) {
	test := serenity.NewSerenityTest(t)

	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// This should pass - all assertions are correct
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(200)),
	)

	// This should fail and stop the test - wrong status code
	// apiTester.AttemptsTo(
	// 	api.SendGetRequest("/posts"),
	// 	ensure.That(api.LastResponseStatus{}, expectations.Equals(202)),
	// 	ensure.That(api.LastResponseStatus{}, expectations.Equals(404)), // This will fail
	// )
}
