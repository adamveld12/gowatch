package project

import (
	"os"
	"os/exec"
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
		handle.start(buildCmd(projectDirectory, "./"+outputName, appArguments))
	}

	gwl.Debug("build steps completed")

	return handle
}

func buildCmd(projectDirectory, command string, arguments ...string) *exec.Cmd {
	cmd := exec.Command(command, arguments...)

	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func runCmd(pwd, command string, args ...string) bool {
	cmd := buildCmd(pwd, command, args...)

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
