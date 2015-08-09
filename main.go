package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/logutils"
)

var (
	debug          = flag.Bool("debug", false, "enabled debug print statements")
	wait           = flag.Duration("wait", time.Second*2, "# seconds to wait before restarting")
	ignore         = flag.String("ignore", "", "comma delimited paths to ignore in the file watcher")
	restartOnExit  = flag.Bool("onexit", true, "If the app sould restart on exit, regardless of exit code")
	restartOnError = flag.Bool("onerror", true, "If the app should restart if a lint/test/build/non-zero exit code occurs")
	appArgs        = flag.String("args", "", "arguments to pass to the underlying app")
	shouldTest     = flag.Bool("test", true, "run go test on reload")
	shouldLint     = flag.Bool("lint", true, "run go lint on reload")

	ignorePaths = []string{}
)

func setupLogging() {

	minLevel := logutils.LogLevel("ERROR")

	if *debug {
		minLevel = logutils.LogLevel("DEBUG")
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "ERROR"},
		MinLevel: minLevel,
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
}

func main() {
	flag.Parse()

	setupLogging()

	projectPath := getAbsPathToProject()

	log.Println("[DEBUG] watching", projectPath)

	watchNotification, _ := startWatch(projectPath)

	for {
		time.Sleep(*wait)
		buildResult, killProcess := executeBuildSteps(projectPath, *appArgs)
		exit := false
		syncer := make(chan bool)

		go func() {
			close(syncer)
			for !exit {
				select {
				default:
				case <-watchNotification:
					select {
					case killProcess <- os.Kill:
					}
					return
				}
			}
			log.Println("[DEBUG] exiting routine")
		}()

		<-syncer
		err := <-buildResult
		exit = true
		if err != nil && *restartOnError && *restartOnExit {
			log.Println("[DEBUG] build result", err)
		} else if !*restartOnError || !*restartOnExit {
			log.Println("[DEBUG] waiting on file notification")
			<-watchNotification
		}
	}

}

func getAbsPathToProject() string {
	pwd := "."

	if flag.Arg(0) != "" {
		pwd = flag.Arg(0)
	}

	directoryPath, err := filepath.Abs(pwd)

	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	return directoryPath
}
