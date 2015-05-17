package main

import (
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"path/filepath"
)

func getWatch(dir string) <-chan string {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	signal := make(chan string)
	go func() {
		if *debug {
			log.Println("Starting watcher routine")
		}
		defer watcher.Close()
		for {
			select {
			// need to debounce this because it could potentially fire a ton of events at once.
			// maybe wait 1000ms after last channel event before passing an event.Name to the signal channel
			case event := <-watcher.Events:
				if *debug {
					log.Println("Event:", event.Op)
				}
				signal <- event.Name
				// if event.Op&fsnotify.Write == fsnotify.Write /*&& strings.HasSuffix(event.Name, ".go")*/ {
				// 	signal <- event.Name
				// }
			case err := <-watcher.Errors:
				log.Println("Error:", err.Error())
			}
		}
	}()

	watcher.Add(dir)
	files(dir, func(fileName string) {
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

		if file.Name() == ".git" {
			continue
		}

		if file.IsDir() {
			if *debug {
				log.Println("Watching", file.Name()+"/")
			}
			apply(abs)
			files(dir+"/"+file.Name(), apply)
		}
	}

}
