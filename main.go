package main

import (
	"flag"
	"log"
	"path/filepath"
)

var (
	wait   = flag.Int("wait", 0, "# seconds to wait before restarting")
	ignore = flag.Bool("ignore", false, "comma delimited paths to ignore in the file watcher")
	test   = flag.Bool("test", false, "run go test on reload")
	lint   = flag.Bool("lint", false, "run go lint on reload")
	debug  = flag.Bool("debug", true, "enabled debug print statements")
	file   = flag.String("file", "main.go", "main file")
)

func init() {
	flag.Parse()
}

func main() {
	mainFile, err := filepath.Abs(*file)
	if err != nil {
		log.Fatal(err)
	}

	cwd := filepath.Dir(mainFile)

	if *debug {
		log.Println("running", mainFile)
		log.Println("Watching", cwd, "for file changes")
	}

	fileUpdates := getWatch(cwd)

	for {
		select {
		case updateType := <-fileUpdates:
			if *debug {
				log.Println(updateType)
			}
		}
	}
}
