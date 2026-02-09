package core

// FailureMode defines how activities should handle failures
type FailureMode int

const (
	// FailFast stops execution immediately when an activity fails
	FailFast FailureMode = iota

	// ErrorButContinue logs the error but continues with remaining activities
	ErrorButContinue

	// Ignore completely ignores the failure and continues
	Ignore
)

// Critical returns a failure mode that stops execution on failure
func Critical() FailureMode { return FailFast }

// NonCritical returns a failure mode that logs errors but continues
func NonCritical() FailureMode { return ErrorButContinue }

// Optional returns a failure mode that ignores errors completely
func Optional() FailureMode { return Ignore }
