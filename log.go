package main

// Logger holds information to log information to files and the console.
type Logger struct {
}

const (

	// LogError level is used for critical errors that could lead to data loss
	// or panic that would not be returned to a calling function.
	LogError int = iota

	// LogWarning level is used for very abnormal events and errors that are
	// also returned to a calling function.
	LogWarning

	// LogInformational level is used for normal non-error activity
	LogInformational

	// LogDebug level is for very detailed non-error activity.  This is
	// very spammy and will impact performance.
	LogDebug
)
