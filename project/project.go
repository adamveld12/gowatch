package project

import (
	"sync"

	gwl "github.com/adamveld12/gowatch/log"
)

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

func Execute(projectDirectory, outputName, appArguments string, shouldTest, shouldLint bool) *ExecuteHandle {
	handle := &ExecuteHandle{
		sync.Mutex{},
		projectDirectory,
		make(chan StepResult, 1),
		false,
		nil,
		false,
		nil,
	}

	if !build(projectDirectory, outputName) {
		handle.Kill(ErrorBuildFailed)
	} else if shouldLint && !lint(projectDirectory) {
		handle.Kill(ErrorLintFailed)
	} else if shouldTest && !test(projectDirectory) {
		handle.Kill(ErrorTestFailed)
	} else {
		handle.start(run(projectDirectory, outputName, appArguments))
	}

	gwl.Debug("build steps completed")

	return handle
}
