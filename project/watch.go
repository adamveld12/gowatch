package project

import (
	"gopkg.in/fsnotify.v1"
	"log"
	"strings"
)

func Watch(dir string) <-chan string {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	signal := make(chan string)
	go func() {
		for {
			select {
			// need to debounce this because it could potentially fire a ton of events at once.
			// maybe wait 1000ms after last channel event before passing an event.Name to the signal channel
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".go") {
					signal <- event.Name
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err.Error())
			}
		}
	}()

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	return signal
}
