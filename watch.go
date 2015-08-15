package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/fsnotify.v1"
)

func shouldIgnore(file string) bool {
	for _, pattern := range ignorePaths {

		matched, err := filepath.Match(strings.Replace(pattern, "/", "", -1), strings.Replace(file, "/", "", -1))

		if err != nil {
			log.Println("[ERROR]", err)
		}

		if matched && err == nil {
			log.Printf("[DEBUG] \tIgnore %s -> %s\n", pattern, file)
			return true
		}
	}

	return false
}

func startWatch(dir string) (<-chan string, chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Could not create file watcher", err)
	}

	fileUpdateNotification, haltWatch := make(chan string), make(chan bool)

	lastEvent, timeSinceLastEvent := "", time.Now().AddDate(-1, 0, 0)
	_, projectName := filepath.Split(dir)

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if shouldIgnore(filepath.Join(p, info.Name())) {
				return filepath.SkipDir
			} else if err := watcher.Add(p); err != nil {
				log.Println("[DEBUG] error adding watched dir", p, info.Name(), err)
				return err
			}
		}
		return nil
	})

	go func() {
		defer watcher.Close()

		for {
			select {
			case event := <-watcher.Events:
				_, filename := filepath.Split(event.Name)

				if projectName == filename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("[DEBUG] skipping restart")
					continue
				}

				if event.Name == lastEvent && timeSinceLastEvent.Add(time.Second).After(time.Now()) {
					log.Println("[DEBUG] skipping restart")
					timeSinceLastEvent = time.Now()
					continue
				}

				if shouldIgnore(event.Name) {
					log.Println("[DEBUG] ignoring update to ", event.Name)
					continue
				}

				lastEvent = event.Name
				timeSinceLastEvent = time.Now()

				log.Println("[DEBUG] updated", event)
				fileUpdateNotification <- event.Name

			case err := <-watcher.Errors:
				log.Println("[ERROR] watcher error:", err)

			case <-haltWatch:
				log.Println("[DEBUG] killing watcher")
				return
			default:
			}
		}
	}()

	return fileUpdateNotification, haltWatch
}
