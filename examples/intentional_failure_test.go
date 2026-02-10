package examples

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/nchursin/serenity-go/serenity/abilities/api"
	"github.com/nchursin/serenity-go/serenity/expectations"
	"github.com/nchursin/serenity-go/serenity/expectations/ensure"
	serenity "github.com/nchursin/serenity-go/serenity/testing"
	"github.com/nchursin/serenity-go/serenity/testing/mocks"
)

// TestIntentionalFailure demonstrates error handling with wrong assertion using mock TestContext
func TestIntentionalFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock TestContext
	mockCtx := mocks.NewMockTestContext(ctrl)

	// Expect Name() to be called during test initialization
	mockCtx.EXPECT().Name().Return("TestIntentionalFailure")

	// Expect Failed() to be called during shutdown
	mockCtx.EXPECT().Failed().Return(true)

	// Expect Errorf to be called when the assertion fails
	// We expect it to be called exactly once with any format string and arguments
	mockCtx.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

	// Create SerenityTest with mock context
	test := serenity.NewSerenityTest(mockCtx)
	defer test.Shutdown()

	apiTester := test.ActorCalled("APITester").WhoCan(api.CallAnApiAt("https://jsonplaceholder.typicode.com"))

	// This should trigger Errorf call - wrong status code (expecting 404 but getting 200)
	apiTester.AttemptsTo(
		api.SendGetRequest("/posts"),
		ensure.That(api.LastResponseStatus{}, expectations.Equals(404)), // This will fail and call Errorf
	)

	// The test will pass because we're using a mock, but the expectation will be verified
	// automatically when ctrl.Finish() is called, ensuring Errorf was called exactly once.
}
