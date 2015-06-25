package main

import (
	"flag"
	"log"
	"path/filepath"
	"time"
)

var (
	wait = flag.Duration("wait", time.Second*0, "# seconds to wait before restarting")
	//ignore = flag.Bool("ignore", false, "comma delimited paths to ignore in the file watcher")
	//test           = flag.Bool("test", false, "run go test on reload")
	//lint           = flag.Bool("lint", false, "run go lint on reload")
	debug = flag.Bool("debug", true, "enabled debug print statements")
	dir   = flag.String("dir", ".", "working directory ")
	//restartOnError = flag.Bool("onerror", true, "If a non-zero exit code should restart the lint/build/test/run process")
	//stepUpdates = make(chan bool)
)

func init() {
	flag.Parse()
}

func main() {
	if *debug {
		log.Println("Debug mode enabled.")
	}

	proj := New(*dir)

	cwd := proj.WorkingDirectory()

	if *debug {
		log.Println("Watching", cwd, "for dir changes")
	}

	fileUpdates := getWatch(cwd)
	proj.RunSteps()

	for {
		select {
		case file := <-fileUpdates:
			_, fileName := filepath.Split(file)
			if fileName == proj.Name() {
				continue
			}
			if *debug {
				log.Println("queueing build", file)
			}

			proj.RunSteps()

			time.Sleep(*wait)
		}
	}
}
