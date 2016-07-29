package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamveld12/gowatch"
	"github.com/fatih/color"
)

var (
	debug          = flag.Bool("debug", false, "")
	wait           = flag.Duration("wait", time.Second*2, "")
	ignore         = flag.String("ignore", ".git/*,node_modules/*", "")
	restartOnExit  = flag.Bool("onexit", true, "")
	restartOnError = flag.Bool("onerror", true, "")
	outputName     = flag.String("output", "", "")
	appArgs        = flag.String("args", "", "")
	shouldTest     = flag.Bool("test", true, "")
	shouldLint     = flag.Bool("lint", true, "")
)

func init() {
	flag.Usage = func() {
		fmt.Println("gowatch [options] <path to main package>")
		fmt.Printf("options\n%s\n", strings.Join([]string{
			color.GreenString(" -output") + "=\"\": the name of the program to output",
			color.GreenString(" -args") + "=\"\": arguments to pass to the underlying app",
			color.GreenString(" -debug") + "=false: enabled debug print statements",
			color.GreenString(" -ignore") + "=\".git/*,node_modules/*\": comma delimited paths to ignore in the file watcher",
			color.GreenString(" -lint") + "=true: run go lint on reload",
			color.GreenString(" -onerror") + "=true: If the app should restart if a lint/test/build/non-zero exit code occurs",
			color.GreenString(" -onexit") + "=true: If the app sould restart on exit, regardless of exit code",
			color.GreenString(" -test") + "=true: run go test on reload",
			color.GreenString(" -wait") + "=2s: # seconds to wait before restarting",
		}, "\n"))
	}
}

func main() {
	flag.Parse()

	packagePath := flag.Arg(0)

	watchConfig := gowatch.Config{
		packagePath,
		*outputName,
		*appArgs,
		*shouldLint,
		*shouldTest,
		*restartOnError,
		*restartOnExit,
		*wait,
		buildIgnorePaths(packagePath, *ignore),
		*debug,
	}

	gowatch.Start(watchConfig)
}

func buildIgnorePaths(root string, rawIgnorePaths string) []string {
	paths := strings.Split(rawIgnorePaths, ",")

	expandedPaths := []string{}
	for _, path := range paths {
		abs := filepath.Join(root, path)
		expandedPaths = append(expandedPaths, abs)
	}

	return expandedPaths
}
