package gowatch

import (
	"os"
	"path/filepath"
	"time"

	gwl "github.com/adamveld12/gowatch/log"
	"github.com/adamveld12/gowatch/watch"

	"github.com/adamveld12/gowatch/project"

	"github.com/fatih/color"
)

// Config has several configuration options used to setup a new Watch
type Config struct {
	PackagePath    string
	OutputName     string
	Arguments      string
	Lint           bool
	Test           bool
	RestartOnError bool
	RestartOnExit  bool
	Wait           time.Duration
	IgnorePaths    []string
	Debug          bool
}

// Start starts running GoWatch on a specified package directory
func Start(config Config) {
	config.PackagePath = defaultPackagePath(config.PackagePath)
	config.OutputName = defaultArtifactOutputName(config.PackagePath, config.OutputName)
	config.Wait = defaultWait(config.Wait)

	startWatch(config.PackagePath,
		config.OutputName,
		config.Arguments,
		config.Lint,
		config.Test,
		config.RestartOnError,
		config.RestartOnExit,
		config.Wait,
		config.IgnorePaths,
		config.Debug)
}

func startWatch(projectPath, outputName, appArgs string,
	shouldLint, shouldTest, restartOnError, restartOnExit bool,
	wait time.Duration,
	ignorePaths []string,
	debug bool) {

	gwl.Setup(debug)

	watchHandle := watch.StartWatch(projectPath, outputName, ignorePaths)

	if outputName == "" {
		outputName = filepath.Base(projectPath)
		gwl.Debug("using", outputName, "as the output name")
	}

	for {
		gwl.Debug("---Starting app monitor---")
		time.Sleep(wait)
		execHandle := project.Execute(projectPath,
			outputName,
			appArgs,
			shouldTest,
			shouldLint)

		gwl.Debug("---Setting up watch cb---")
		watchHandle.Subscribe(func(fileName string) {
			if !execHandle.Halted() {
				gwl.Error("attempting to kill process")
				execHandle.Kill(nil)
				gwl.Debug("exiting file watch routine in main")
			}
		})

		gwl.Debug("waiting on app to exit")
		err := execHandle.Error()
		gwl.Debug("---App exited---")
		watchHandle.Subscribe(nil)

		exitedSuccessfully := err == nil || err == project.ErrorAppKilled

		if exitedSuccessfully {
			color.Green("exited successfully\n")
		} else {
			color.Red("%s\n", err.Error())
		}

		sync := make(chan bool)
		if (!exitedSuccessfully && !restartOnError) || (!restartOnExit && err != project.ErrorAppKilled) {
			watchHandle.Subscribe(func(fileName string) {
				close(sync)
				watchHandle.Subscribe(nil)
			})
			gwl.Debug("waiting on file notification")
			<-sync
		}
	}
}

func defaultArtifactOutputName(packagePath, outputName string) string {
	if outputName == "" {
		return filepath.Base(packagePath)
	}

	return outputName
}

func defaultWait(waitTime time.Duration) time.Duration {

	if waitTime/time.Millisecond < 500 {
		return time.Millisecond * 500
	} else if waitTime.Seconds() > 10 {
		return time.Second * 10
	}

	return waitTime
}

func defaultPackagePath(packagePath string) string {
	finalDir := packagePath
	if packagePath == "" {
		if dir, err := os.Getwd(); err == nil {
			finalDir = dir
		}
	} else if dir, err := filepath.Abs(packagePath); err == nil {
		finalDir = dir
	}

	return finalDir
}
