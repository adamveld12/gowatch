package main

import (
	"flag"
	"log"
	"strings"
	"time"
)

var (
	wait   = flag.Duration("wait", time.Second*0, "# seconds to wait before restarting")
	ignore = flag.String("ignore", "", "comma delimited paths to ignore in the file watcher")
	debug  = flag.Bool("debug", true, "enabled debug print statements")
	dir    = flag.String("dir", ".", "working directory ")
	//test           = flag.Bool("test", false, "run go test on reload")
	//lint           = flag.Bool("lint", false, "run go lint on reload")
	//restartOnError = flag.Bool("onerror", true, "If a non-zero exit code should restart the lint/build/test/run process")
	//stepUpdates = make(chan bool)
	ignorePaths = []string{}
)

func init() {
	flag.Parse()
}

func main() {
	ignorePaths = strings.Split(*ignore, ",")

	if *debug {
		log.Println("Debug mode enabled.")
		for _, files := range ignorePaths {
			log.Println("\tignoring", files)
		}
	}

	proj := createProject(*dir)
	cwd := proj.WorkingDirectory()
	fileUpdates := getWatch(cwd)

	proj.RunSteps()

	for {
		select {
		case file := <-fileUpdates:
			if *debug {
				log.Println("\tbuilding", file)
			}

			proj.RunSteps()

			time.Sleep(*wait)
		}
	}
}
