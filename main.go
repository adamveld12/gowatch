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
	shouldTest     = flag.Bool("test", false, "run go test on reload")
	shouldLint     = flag.Bool("lint", false, "run go lint on reload")

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
		buildResult, killProcess := executeBuildSteps(projectPath, *appArgs)
		exit := false
		syncer := make(chan bool)

		go func() {
			for !exit {
				close(syncer)
				select {
				case <-watchNotification:
					select {
					case killProcess <- os.Kill:
					}
					return
				}
			}
		}()
		<-syncer

		if err := <-buildResult; err != nil {
			log.Println("[DEBUG] build result", err)
		}
		exit = true
		time.Sleep(*wait)
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
