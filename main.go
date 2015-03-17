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

	flag.BoolVar(&shouldTest, "test", false, "Run tests.")
	flag.BoolVar(&shouldLint, "lint", false, "Run lint.")

	flag.BoolVar(&noRestartOnErr, "error", true, "No restart on error.")
	flag.BoolVar(&noRestartOnExit, "exit", true, "No restart on exit, regardless of exit code.")

	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified.")
}

func main() {
	flag.Parse()

	mainFile, err := filepath.Abs(mainFile)

	if err != nil {
		log.Fatal(err)
	}

	cwd = filepath.Dir(mainFile)

	fileUpdates := project.Watch(cwd)

	for {
		select {
		case updateType := <-fileUpdates:
			log.Println("Update:", updateType)
		default:
		}
	}
}
