package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

func getWatch(dir string) (<-chan string, chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	signal, killChan := make(chan string), make(chan bool)

	go func() {
		defer watcher.Close()

		lastEvent := ""
		debouncer := time.AfterFunc(time.Second*2, func() {
			if *debug {
				log.Println("File Updated:", lastEvent)
			}
			lastEvent = ""
			signal <- lastEvent
		})

		debouncer.Stop()

		_, projectName := filepath.Split(dir)

		for {
			select {
			case event := <-watcher.Events:

				_, eventFile := filepath.Split(event.Name)
				if projectName == eventFile {
					continue
				}

				if event.Name == lastEvent {
					debouncer.Reset(time.Second * 2)
				}

				lastEvent = event.Name

			case err := <-watcher.Errors:
				log.Println("\tWatcher error:", err.Error())
			case <-killChan:
				if *debug {
					log.Println("Shit fuck exiting")
				}
				return
			}
		}
	}()

	if *debug {
		log.Println("Starting watcher routine @ ", dir)
		log.Println("\t " + dir + "/.")
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		log.Fatal(err)
	}

	files(dir, func(filePath string) {
		if *debug {
			log.Println("\t " + filePath + "/")
		}

		err := watcher.Add(filePath)
		if err != nil {
			watcher.Close()
			log.Fatal(err)
		}
	})

	return signal, killChan
}

func files(dir string, apply func(string)) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		abs, err := filepath.Abs(filepath.Join(dir, file.Name()))
		shouldContinue := false

		for _, path := range ignorePaths {
			if match, _ := filepath.Match(path, strings.TrimPrefix(abs, *pwd)); match {
				if *debug {
					log.Println("\t ignoring", abs)
				}
				shouldContinue = true
				break
			}
		}

		if shouldContinue {
			continue
		}

		if err != nil {
			log.Fatal(err)
		}

		if file.Name() == ".git" || file.Name() == ".gitignore" {
			continue
		}

		if file.IsDir() {
			apply(abs)
			files(dir+"/"+file.Name(), apply)
		}
	}

}
