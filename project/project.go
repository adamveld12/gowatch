package project

import (
	"sync"

	gwl "github.com/adamveld12/gowatch/log"
)

func ExecuteBuildSteps(projectDirectory, appArguments string, shouldTest bool, shouldLint bool) *ExecuteHandle {

	handle := &ExecuteHandle{
		sync.Mutex{},
		projectDirectory,
		make(chan StepResult, 1),
		false,
		nil,
		false,
		nil,
	}

	if !build(projectDirectory) {
		handle.Kill(ErrorBuildFailed)
	} else if shouldLint && !lint(projectDirectory) {
		handle.Kill(ErrorLintFailed)
	} else if shouldTest && !test(projectDirectory) {
		handle.Kill(ErrorTestFailed)
	} else {
		handle.start(run(projectDirectory, appArguments))
	}

	gwl.LogDebug("[DEBUG] build steps completed")

	return handle
}
