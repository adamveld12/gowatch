package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/logutils"
)

var (
	debug          = flag.Bool("debug", false, "")
	wait           = flag.Duration("wait", time.Second*2, "")
	ignore         = flag.String("ignore", ".git/*,node_modules/*", "")
	restartOnExit  = flag.Bool("onexit", true, "")
	restartOnError = flag.Bool("onerror", true, "")
	appArgs        = flag.String("args", "", "")
	shouldTest     = flag.Bool("test", true, "")
	shouldLint     = flag.Bool("lint", true, "")

	ignorePaths = []string{}
)

func init() {
	flag.Usage = func() {
		fmt.Printf("%s options\n%s\n", os.Args[0], strings.Join([]string{
			color.GreenString("  -args") + "=\"\": arguments to pass to the underlying app",
			color.GreenString("  -debug") + "=false: enabled debug print statements",
			color.GreenString("  -ignore") + "=\".git/*,node_modules/*\": comma delimited paths to ignore in the file watcher",
			color.GreenString("  -lint") + "=true: run go lint on reload",
			color.GreenString("  -onerror") + "=true: If the app should restart if a lint/test/build/non-zero exit code occurs",
			color.GreenString("  -onexit") + "=true: If the app sould restart on exit, regardless of exit code",
			color.GreenString("  -test") + "=true: run go test on reload",
			color.GreenString("  -wait") + "=2s: # seconds to wait before restarting",
		}, "\n"))
	}
}

func main() {
	flag.Parse()

	setupLogging()

	projectPath := getAbsPathToProject()
	ignorePaths = setupIgnorePaths(projectPath)

	log.Println("[DEBUG] watching", projectPath)

	watchNotification, _ := startWatch(projectPath)

	for {
		time.Sleep(*wait)
		buildResult, killProcess := executeBuildSteps(projectPath, *appArgs)
		exit, syncer := false, make(chan bool)

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

		if err != nil {
			color.Red(err.Error())
		}

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

func setupIgnorePaths(root string) []string {
	log.Println("[DEBUG] Ignore globs.")
	paths := strings.Split(*ignore, ",")

	expandedPaths := []string{}
	for _, path := range paths {
		abs := filepath.Join(root, path)
		log.Printf("[DEBUG] \t%s\n", abs)
		expandedPaths = append(expandedPaths, abs)
	}

	return expandedPaths
}
