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
var shouldTest, shouldLint, noRestartOnErr, noRestartOnExit, quiet bool

func init() {
	flag.StringVar(&mainFile, "main", "main.go", "A go file with func main().")

	flag.BoolVar(&quiet, "quiet", true, "Controls debug printing.")

	flag.BoolVar(&shouldTest, "test", false, "Run tests.")
	flag.BoolVar(&shouldLint, "lint", false, "Run lint.")

	flag.BoolVar(&noRestartOnErr, "error", false, "No restart on error.")
	flag.BoolVar(&noRestartOnExit, "exit", false, "No restart on exit, regardless of exit code.")

	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified.")
}

var appHandle *project.AppHandle

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

	buildSuccess := runSteps()
	exitedSuccessfully := runApp()

	for {
		select {
		case updateType := <-fileUpdates:
			if !quiet {
				log.Println("File Update:", updateType)
			}
			buildSuccess = runSteps()
		default:
			if !quiet {
				log.Println("App running?", appHandle.Running())
				log.Println("Should restart on err?", (!exitedSuccessfully && noRestartOnErr))
				log.Println("Should restart on exit?", (!appHandle.Running() && noRestartOnExit))
			}
			if (!exitedSuccessfully && noRestartOnErr) || (!appHandle.Running() && noRestartOnExit) {
				continue
			} else if !appHandle.Running() {
				buildSuccess = runSteps()
			}
		}
		if buildSuccess {
			exitedSuccessfully = runApp()
		}
	}
}

func runSteps() bool {
	if !quiet && shouldLint {
		log.Println("Linting", mainFile)
	}

	lintSuccess := project.Lint(mainFile, cwd)

	testsPassed := false
	if lintSuccess {

		if !quiet && shouldTest {
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

func runApp() bool {
	if !quiet {
		log.Println("Restarting", mainFile)
	}

	// halt if running
	if appHandle.Running() {
		appHandle.Halt()
	}

	// run and don't block
	return appHandle.Start()
}
