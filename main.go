package main

import (
	"flag"
	"github.com/adamveld12/goad/exec"
	"gopkg.in/fsnotify.v1"
	"log"
	"path/filepath"
	"strings"
)

/*
 * options to run:
 *  goimport
 *  golint
 *  go test
 */

var ignore, mainFile, cwd string
var debug, shouldTest, shouldLint, restartOnErr, restartOnSave, restartOnBuildErr bool

func init() {
	flag.BoolVar(&debug, "debug", false, "run in debug mode.")

	flag.StringVar(&mainFile, "main", "main.go", "A go file with func main().")

	flag.BoolVar(&shouldTest, "test", false, "Run tests.")
	flag.BoolVar(&shouldLint, "lint", false, "Run lint.")
	flag.BoolVar(&restartOnSave, "save", true, "Restart app on save.")
	flag.BoolVar(&restartOnErr, "error", true, "Restart app if an error occurs during executions.")
	flag.BoolVar(&restartOnBuildErr, "build", true, "Restart app if an error occurs during compilation.")

	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified.")
}

func main() {
	flag.Parse()

	var err error = nil
	if mainFile, err = filepath.Abs(mainFile); err != nil {
		log.Fatal(err)
	}
	cwd = filepath.Dir(mainFile)

	debugLog("main file:", mainFile)
	debugLog("current working directory:", cwd)

	if restartOnSave {
		//watch(files)
	}
}

func runSteps() bool {

	lintedSuccessfully := !shouldLint
	if shouldLint {
		lintedSuccessfully = lint()
	}

	passedTests := !shouldTest
	if shouldTest && lintedSuccessfully {
		passedTests = test()
	}

	var builtSuccessfully bool = false
	if lintedSuccessfully && passedTests {
		builtSuccessfully = build()
	}

	buildFailed := restartOnBuildErr && (!lintedSuccessfully || !passedTests || !builtSuccessfully)

	restart := restartOnBuildErr && (!lintedSuccessfully || !passedTests || runState.BuildFailed())
	restart = restart || (restartOnErr && !runState.ExitedSuccessfully())

	return restart
}

func lint() bool {
	debugLog("Linting code")
	return true
}

func test() bool {
	debugLog("Testing code")
	return true
}

func build() bool {
	return true
}

func run() bool {
	debugLog("Building and execing code")
	exec.Command("go", "run", mainFile, cwd)
	return true
}

func watch(paths []string) chan string {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	signal := make(chan string)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					debugLog("modified file:", event.Name)
					signal <- event.Name
				}
			case err := <-watcher.Errors:
				debugLog("Error:", err.Error())
			}
		}
	}()

	err = watcher.Add("./")

	if err != nil {
		log.Fatal(err)
	}

	return signal
}

func debugLog(args ...string) {
	if debug {
		log.Println("DEBUG:", strings.Join(args, " "))
	}
}
