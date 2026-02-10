package testing

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/nchursin/serenity-go/serenity/reporting"
	"github.com/nchursin/serenity-go/serenity/reporting/console_reporter"
	reportingMocks "github.com/nchursin/serenity-go/serenity/reporting/mocks"
	"github.com/nchursin/serenity-go/serenity/testing/mocks"
)

func TestSerenityTestWithConsoleReporter(t *testing.T) {
	// Create a SerenityTest with console reporter
	test := NewSerenityTestWithReporter(t, console_reporter.NewConsoleReporter())
	defer test.Shutdown()

	actor := test.ActorCalled("TestActor")
	require.NotNil(t, actor)

	// Verify that reporter is configured
	adapter := test.GetReporterAdapter()
	require.NotNil(t, adapter)
	require.IsType(t, &console_reporter.ConsoleReporter{}, adapter.GetReporter())
}

func TestNewSerenityTestUsesConsoleReporter(t *testing.T) {
	test := NewSerenityTest(t)
	defer test.Shutdown()

	adapter := test.GetReporterAdapter()
	require.NotNil(t, adapter)

	// Verify it's a ConsoleReporter
	reporter := adapter.GetReporter()
	_, isConsole := reporter.(*console_reporter.ConsoleReporter)
	require.True(t, isConsole, "Expected ConsoleReporter")
}

func TestSerenityTestLifecycleReporting(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReporter := reportingMocks.NewMockReporter(ctrl)
	mockTestContext := mocks.NewMockTestContext(ctrl)

	// Expect test lifecycle events
	mockReporter.EXPECT().OnTestStart("TestExample")
	mockReporter.EXPECT().OnTestFinish(gomock.Any()).Do(func(result reporting.TestResult) {
		require.Equal(t, "TestExample", result.Name())
		require.Equal(t, reporting.StatusPassed, result.Status())
		require.True(t, result.Duration() >= 0)
		require.NoError(t, result.Error())
	})

	mockTestContext.EXPECT().Name().Return("TestExample")
	mockTestContext.EXPECT().Failed().Return(false)

	test := NewSerenityTestWithReporter(mockTestContext, mockReporter)

	// Simulate test end
	test.Shutdown()
}

func TestSerenityTestLifecycleReportingFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReporter := reportingMocks.NewMockReporter(ctrl)
	mockTestContext := mocks.NewMockTestContext(ctrl)

	// Expect test lifecycle events for failed test
	mockReporter.EXPECT().OnTestStart("FailedTest")
	mockReporter.EXPECT().OnTestFinish(gomock.Any()).Do(func(result reporting.TestResult) {
		require.Equal(t, "FailedTest", result.Name())
		require.Equal(t, reporting.StatusFailed, result.Status())
		require.True(t, result.Duration() >= 0)
		require.Error(t, result.Error())
		require.Equal(t, "test failed", result.Error().Error())
	})

	mockTestContext.EXPECT().Name().Return("FailedTest")
	mockTestContext.EXPECT().Failed().Return(true)

	test := NewSerenityTestWithReporter(mockTestContext, mockReporter)

	// Simulate test end
	test.Shutdown()
}
