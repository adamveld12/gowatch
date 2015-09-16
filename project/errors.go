package project

import "errors"

// StepResult the result of a build step
type StepResult error

var (
	// ErrorBuildFailed indicates build failure
	ErrorBuildFailed = errors.New("Build failed")

	// ErrorRunFailed indicates app exited with error
	ErrorRunFailed = errors.New("App exited with non-zero exit code")

	// ErrorTestFailed indicates one or more tests failed
	ErrorTestFailed = errors.New("Test failed")

	// ErrorLintFailed indicates Linter errors
	ErrorLintFailed = errors.New("Lint failed")

	// ErrorAppKilled indicates the file watcher killed the process
	ErrorAppKilled = errors.New("File watcher killed")

	errorProcessAlreadyFinished = errors.New("os: process already finished")
)
