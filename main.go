package main

import (
	"flag"
	"log"
	"time"
)

var (
	wait           = flag.Duration("wait", time.Second*0, "# seconds to wait before restarting")
	ignore         = flag.Bool("ignore", false, "comma delimited paths to ignore in the file watcher")
	test           = flag.Bool("test", false, "run go test on reload")
	lint           = flag.Bool("lint", false, "run go lint on reload")
	debug          = flag.Bool("debug", true, "enabled debug print statements")
	dir            = flag.String("dir", ".", "working directory ")
	restartOnError = flag.Bool("onerror", true, "If a non-zero exit code should restart the lint/build/test/run process")
	stepUpdates    = make(chan bool)
)

func init() {
	flag.Parse()
}

func main() {
	if *debug {
		log.Println("Debug mode enabled.")
	}

	proj := New(*dir)
	if *debug {
		log.Println("CWD: ", proj.Directory)
	}

	if *debug {
		log.Println("Watching", proj.Directory, "for dir changes")
	}

	fileUpdates := getWatch(proj.Directory)
	buildTestRun(proj.Directory)

	for {
		select {
		case file := <-fileUpdates:
			if *debug {
				log.Println("queueing build", file)
			}

			buildTestRun(proj.Directory)
		}
	}
}

func buildTestRun(cwd string) {
	if buildSucceeded := build(cwd); buildSucceeded {
		if exitedSuccessfully := run(cwd); exitedSuccessfully {
			log.Println("Exited successfully")
		} else if *restartOnError {
			log.Println("Exit fail.")
		}
	}

	time.Sleep(*wait)
}
