package project

import (
	"sync"

	gwl "github.com/adamveld12/gowatch/log"
	"github.com/fatih/color"
)

func ExecuteBuildSteps(projectDirectory, appArguments string, shouldTest bool, shouldLint bool) *ExecuteHandle {

	handle := &ExecuteHandle{
		sync.Mutex{},
		projectDirectory,
		make(chan StepResult, 1),
		false,
		nil,
		false,
	}

	if !build(projectDirectory) {
		color.Red("Build failed.")
		handle.Kill(ErrorBuildFailed)
	} else if shouldLint && !lint(projectDirectory) {
		color.Red("Linter found errors.")
		handle.Kill(ErrorLintFailed)
	} else if shouldTest && !test(projectDirectory) {
		color.Red("Tests failed.")
		handle.Kill(ErrorTestFailed)
	} else {
		handle.start(run(projectDirectory, appArguments))
	}

	gwl.LogDebug("[DEBUG] build steps completed")

	return handle
}
