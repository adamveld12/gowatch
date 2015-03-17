package project

import (
	"gopkg.in/fsnotify.v1"
	"log"
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
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
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
