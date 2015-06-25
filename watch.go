package main

import (
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func getWatch(dir string) <-chan string {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	signal := make(chan string)
	go func() {
		defer watcher.Close()

		lastEvent := ""
		debouncer := time.AfterFunc(time.Second*2, func() {
			if *debug {
				log.Println("File Updated:", lastEvent)
			}

			signal <- lastEvent
		})
		debouncer.Stop()

		// so we ignore the go build artifact
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
			}
		}
	}()

	if *debug {
		log.Println("Starting watcher routine @ ", dir)
		log.Println("\t\t" + dir + "/.")
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		log.Fatal(err)
	}

	files(dir, func(fileName string) {
		if *debug {
			log.Println("\t\t" + fileName + "/")
		}
		err := watcher.Add(fileName)
		if err != nil {
			watcher.Close()
			log.Fatal(err)
		}
	})

	return signal
}

func files(dir string, apply func(string)) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		abs, err := filepath.Abs(dir + "/" + file.Name())

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
