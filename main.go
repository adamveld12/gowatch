package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	wait           = flag.Duration("wait", time.Second*2, "# seconds to wait before restarting")
	ignore         = flag.String("ignore", "", "comma delimited paths to ignore in the file watcher")
	debug          = flag.Bool("debug", true, "enabled debug print statements")
	pwd            = flag.String("dir", ".", "working directory ")
	restartOnExit  = flag.Bool("onexit", true, "If the app sould restart on exit, regardless of exit code")
	restartOnError = flag.Bool("onerror", true, "If the app should restart if a lint/test/build/non-zero exit code occurs")
	appArgs        = flag.String("args", "", "arguments to pass to the underlying app")
	shouldTest     = flag.Bool("test", false, "run go test on reload")
	shouldLint     = flag.Bool("lint", false, "run go lint on reload")

	ignorePaths = []string{}
)

func init() {
	flag.Parse()
}

func main() {
	ignorePaths = strings.Split(*ignore, ",")

	if *debug {
		log.Println("Debug mode enabled.")
		if !*restartOnError {
			log.Println("\tRestart on error disabled")
		}
		for _, files := range ignorePaths {
			log.Println("\tignoring", files)
		}
	}

	*pwd, _ = filepath.Abs(*pwd)
	fileUpdates, killWatcher := getWatch(*pwd)

	done := false
	appStopped, killApp := make(<-chan StepResult), make(chan<- os.Signal)

	defer func() {
		killWatcher <- true
		if killApp != nil {
			killApp <- os.Kill
		}
		done = true
	}()

	go func() {
		for !done {
			select {
			case <-fileUpdates:
				if killApp != nil {
					killApp <- os.Kill
				}
			}
		}

	}()

	for {
		appStopped, killApp = runProject(*pwd, *appArgs)
		exitError := <-appStopped
		log.Println(exitError)
		time.Sleep(*wait)
	}
}
