package main

import (
	"flag"
	"github.com/adamveld12/gowatch/project"
	"log"
	"path/filepath"
)

/*
 * options to run:
 *  goimport
 *  golint
 *  go test
 */

var ignore, mainFile, cwd string
var shouldTest, shouldLint, noRestartOnErr, noRestartOnExit bool

func init() {
	flag.StringVar(&mainFile, "main", "main.go", "A go file with func main().")

	flag.BoolVar(&quiet, "quiet", true, "Controls debug printing.")

	flag.BoolVar(&shouldTest, "test", false, "Run tests.")
	flag.BoolVar(&shouldLint, "lint", false, "Run lint.")

	flag.BoolVar(&noRestartOnErr, "error", true, "No restart on error.")
	flag.BoolVar(&noRestartOnExit, "exit", true, "No restart on exit, regardless of exit code.")

	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified.")
}

var appHandle project.AppHandle = nil

func main() {
	flag.Parse()

	mainFile, err := filepath.Abs(mainFile)

	if err != nil {
		log.Fatal(err)
	}

	cwd = filepath.Dir(mainFile)
	fileUpdates := project.Watch(cwd)

	if shouldLint {
		project.EnableLinting()
	}

	if shouldTest {
		project.EnableTesting()
	}

	appHandle = project.CreateAppHandle(mainFile, cwd)

	for {
		buildSuccess := false

		select {
		case updateType := <-fileUpdates:
			log.Println("Update:", updateType)
			buildSuccess = runSteps()

		default:
			if !appHandle.Running() {
				buildSuccess = runSteps()
			}
		}

		if buildSuccess {
			runApp()
		}

	}
}

func runSteps() bool {

	if !quiet {
		log.Println("Linting", mainFile)
	}
	lintSuccess := project.Lint(mainFile, cwd)

	testsPassed := false
	if lintSuccess {

		if !quiet {
			log.Println("Testing", mainFile)
		}
		testsPassed = project.Test(mainFile, cwd)
	}

	buildSuccessful := false
	if lintSuccess && testsPassed {

		if !quiet {
			log.Println("Building", mainFile)
		}

		buildSuccessful = project.Build(mainFile, cwd)
	}

	return lintSuccess && testsPassed && buildSuccessful
}

func runApp() {
	if !quiet {
		log.Println("Restarting", mainFile)
	}

	// halt if running
	if appHandle.Running() {
		appHandle.Halt()
	}

	// run and don't block
	appHandle.Run()
}
