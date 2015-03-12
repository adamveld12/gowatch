package main

import (
	"./exec"
	"flag"
	//"gopkg.in/fsnotify.v1"
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

var ignore, restartMode, mainFile, cwd string
var debug, shouldTest, shouldLint bool

func init() {
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.BoolVar(&shouldTest, "test", false, "run tests")
	flag.BoolVar(&shouldLint, "lint", false, "run lint")
	flag.StringVar(&mainFile, "main", "main.go", "A go file with func main()")
	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified")
	flag.StringVar(&restartMode, "restart-mode", "", "How to handle application restarting. error = restart only on error, build = restart only on successful build, save = restart only on save")
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

	runSteps()

	//watch(files)
}

func runSteps() {

	build()
}

func build() {
	exec.Command("go", "run", mainFile, cwd)
}

// func watch(paths []string) chan string {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	defer watcher.Close()
// 	output := make(chan string)
// 	go func() {
// 		var err error = nil
// 		for {
// 			select {
// 			case event := <-watcher.Events:
//
// 				if debug {
// 					log.Println("event:", event)
// 				}
//
// 				if event.Op&fsnotify.Write == fsnotify.Write {
// 					if debug {
// 						log.Println("modified file:", event.Name)
// 					}
// 				}
//
// 				err = build()
// 				for err != nil && (restartMode == "error" || restartMode == "") {
// 					err = build()
// 				}
//
// 			case err := <-watcher.Errors:
// 				if debug {
// 					log.Println("error:", err)
// 				}
// 			}
// 		}
// 	}()
//
// 	err = watcher.Add("./")
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	return output
// }

func debugLog(args ...string) {
	if debug {
		log.Println("DEBUG:", strings.Join(args, " "))
	}
}

func expandGlobs(globs []string) []string {
	return []string{}
}
