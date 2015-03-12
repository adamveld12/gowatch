package main

import (
	"flag"
	//"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	//"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
 * options to run:
 *  goimport
 *  golint
 *  go test
 */

var ignore, restartMode, mainFile string
var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.StringVar(&mainFile, "main", "main.go", "A go file with func main()")
	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified")
	flag.StringVar(&restartMode, "restart-mode", "", "How to handle application restarting. error = restart only on error, build = restart only on successful build, save = restart only on save")
}

func main() {
	flag.Parse()

	mainFile, err := filepath.Abs(mainFile)

	if err != nil {
		log.Fatal(err)
	}

	debugLog(mainFile)

	cwd := filepath.Dir(mainFile)
	execCommand("go", "run", mainFile, cwd)

	//watch(files)
}

func execCommand(command, flags, args, cwd string) {
	cmd := exec.Command(command, flags, args)
	cmd.Dir = cwd

	stdout, err := cmd.StdoutPipe()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	output, err := ioutil.ReadAll(stdout)

	cmd.Wait()
	log.Println(string(output))
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
