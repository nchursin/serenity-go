package testing

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/nchursin/serenity-go/serenity/core"
	coreMocks "github.com/nchursin/serenity-go/serenity/core/testing/mocks"
	"github.com/nchursin/serenity-go/serenity/reporting"
	reportingMocks "github.com/nchursin/serenity-go/serenity/reporting/mocks"
	testingMocks "github.com/nchursin/serenity-go/serenity/testing/mocks"
)

func TestTestActorAttemptsToWithReporting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockReporter := reportingMocks.NewMockReporter(ctrl)
	mockTestContext := testingMocks.NewMockTestContext(ctrl)

	// Expect OnStepStart and OnStepFinish for activity
	mockReporter.EXPECT().OnStepStart("Send GET request to /posts").Times(1)
	mockReporter.EXPECT().OnStepFinish(gomock.Any()).Times(1)

	// Expect no error from test context
	mockTestContext.EXPECT().Failed().Return(false).AnyTimes()

	// Create test actor with mock reporter
	adapter := reporting.NewTestRunnerAdapter(mockReporter)
	testCtx := context.Background()
	test := &serenityTest{
		testCtx: mockTestContext,
		ctx:     testCtx,
		actors:  make(map[string]core.Actor),
		adapter: adapter,
	}

	actor := test.ActorCalled("TestActor")

	// Create mock activity
	mockActivity := coreMocks.NewMockActivity(ctrl)
	mockActivity.EXPECT().PerformAs(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockActivity.EXPECT().Description().Return("Send GET request to /posts").Times(1)
	mockActivity.EXPECT().FailureMode().Return(core.FailFast).AnyTimes()

	// Execute activity
	actor.AttemptsTo(mockActivity)
}
