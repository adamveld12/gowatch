package main

import (
	"flag"
	"gopkg.in/fsnotify.v1"
	"log"
	"os/exec"
	"strings"
)

var ignore, restartMode, mainFile string
var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.StringVar(&mainFile, "main", "main.go", "The go file to run")
	flag.StringVar(&ignore, "ignore", "", "Ignores a file path based on the glob specified")
	flag.StringVar(&restartMode, "restart-mode", "", "how to handle application restarting. error = restart only on error, build = restart only on successful build, save = restart only on save")
}

func main() {
	flag.Parse()

	globs := strings.Split(ignore, ",")
	files := expandGlobs(globs)

	watch(files)
}

func build() error {
	cmd := exec.Command("go", "run", mainFile)
	cmd.Dir = ""

	err := cmd.Start()
	if err != nil {
		// restart build
		return err
	}
	return nil
}

func watch(paths []string) chan string {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()
	output := make(chan string)
	go func() {
		var err error = nil
		for {
			select {
			case event := <-watcher.Events:

				if debug {
					log.Println("event:", event)
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if debug {
						log.Println("modified file:", event.Name)
					}
				}

				err = build()
				for err != nil && (restartMode == "error" || restartMode == "") {
					err = build()
				}

			case err := <-watcher.Errors:
				if debug {
					log.Println("error:", err)
				}
			}
		}
	}()

	err = watcher.Add("./")

	if err != nil {
		log.Fatal(err)
	}

	return output
}

func expandGlobs(globs []string) []string {
	return []string{}
}
