package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/fsnotify.v1"
)

func shouldIgnore(file string) bool {
	for _, pattern := range ignorePaths {

		matched, err := filepath.Match(strings.Replace(pattern, "/", "", -1), strings.Replace(file, "/", "", -1))

		if err != nil {
			log.Println("[ERROR]", err)
		}

		if matched && err == nil {
			log.Printf(color.MagentaString("[DEBUG] \tIgnore %s -> %s\n", pattern, file))
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

	_, projectName := filepath.Split(dir)

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if shouldIgnore(filepath.Join(p, info.Name())) {
				return filepath.SkipDir
			} else if err := watcher.Add(p); err != nil {
				log.Println(color.MagentaString("[DEBUG] error adding watched dir", p, info.Name(), err))
				return err
			}
		}
		return nil
	})

	go handleFileEvent(watcher, fileUpdateNotification, projectName)
	go killWatcher(watcher, haltWatch)
	go watcherErrors(watcher.Errors)

	return fileUpdateNotification, haltWatch
}

func killWatcher(watcher *fsnotify.Watcher, haltWatch <-chan bool) {
	<-haltWatch
	log.Println(color.MagentaString("[DEBUG] killing watcher"))
	watcher.Close()
}

func watcherErrors(errorChan <-chan error) {
	for {
		err, ok := <-errorChan
		if !ok {
			log.Println(color.MagentaString("[DEBUG] Closing file watcher error routine"))
			return
		}
		log.Println("[ERROR] watcher error:", err)
	}
}

func handleFileEvent(watcher *fsnotify.Watcher, fileUpdateNotification chan<- string, projectName string) {
	lastEvent, timeSinceLastEvent := "", time.Now().AddDate(-1, 0, 0)

	for {
		event, ok := <-watcher.Events

		if !ok {
			log.Println(color.MagentaString("[DEBUG] closing file event channel"))
			return
		}
		_, filename := filepath.Split(event.Name)

		if projectName == filename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
			log.Println(color.MagentaString("[DEBUG] ignoring go build artifacts"))
			continue
		}

		if event.Name == lastEvent && timeSinceLastEvent.Add(time.Second).After(time.Now()) {
			log.Println(color.MagentaString("[DEBUG] ignoring extra file watch events"))
			timeSinceLastEvent = time.Now()
			continue
		}

		if shouldIgnore(event.Name) {

			log.Println(color.MagentaString("[DEBUG] %s in ignore path ", event.Name))
			continue
		}

		lastEvent = event.Name
		timeSinceLastEvent = time.Now()

		log.Println(color.MagentaString("[DEBUG] updated", event))
		fileUpdateNotification <- event.Name

	}
}
