package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamveld12/gowatch/project"
	"github.com/adamveld12/gowatch/watch"
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

	log.Println("[DEBUG] watching", projectPath)

	handleWatch(projectPath, setupIgnorePaths(projectPath))
}

func handleWatch(projectPath string, ignorePaths []string) {
	watchHandle := watch.StartWatch(projectPath, ignorePaths)

	for {
		time.Sleep(*wait)
		execHandle := project.ExecuteBuildSteps(projectPath, *appArgs, *shouldTest, *shouldLint)

		go func() {
			for execHandle.Running() {
				select {
				case <-watchHandle.FileNotifier():
					log.Println(color.MagentaString("[ERROR] attempting to kill process"))
					execHandle.Kill(nil)
					log.Println("[DEBUG] exiting file watch routine in main")
					return
				default:
					if execHandle.Halted() {
						return
					}
				}
			}
		}()

		log.Println(color.MagentaString("[DEBUG] Starting forealsies"))
		err := execHandle.Error()

		exitedSuccessfully := err == nil || err == project.ErrorAppKilled

		if exitedSuccessfully {
			color.Green("exited successfully\n")
		} else {
			color.Red("%s\n", err.Error())
		}

		if (!exitedSuccessfully && !*restartOnError) || (!*restartOnExit && err != project.ErrorAppKilled) {
			log.Println("[DEBUG] waiting on file notification ", err.Error())
			<-watchHandle.FileNotifier()
		}

		log.Println("[DEBUG] exiting routine")
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
