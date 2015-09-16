package watch

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"gopkg.in/fsnotify.v1"
)

type WatchHandle struct {
	sync.Mutex
	halted       bool
	fileNotifier chan string
	watcher      *fsnotify.Watcher
	ignorePaths  []string
}

func (w *WatchHandle) Halt() {
	log.Println("[DEBUG] Hitting watcher halt lock")
	if !w.halted {
		w.halted = true
		close(w.fileNotifier)
		go func() {
			w.watcher.Close()
		}()
	}
	log.Println("[DEBUG] exited watcher halt lock")
}

func (w *WatchHandle) FileNotifier() <-chan string {
	return w.fileNotifier
}

func StartWatch(dir string, ignorePaths []string) *WatchHandle {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal("Could not create file watcher", err)
	}

	addFilesToWatch(dir, ignorePaths, watcher)

	handle := &WatchHandle{
		sync.Mutex{},
		false,
		make(chan string),
		watcher,
		ignorePaths,
	}

	_, projectName := filepath.Split(dir)
	go handle.handleFileEvent(projectName)

	return handle
}

func (handle *WatchHandle) handleFileEvent(projectName string) {
	lastEvent, timeSinceLastEvent := "", time.Now().AddDate(-1, 0, 0)

	watcher := handle.watcher
	fileUpdateNotification := handle.fileNotifier
	ignorePaths := handle.ignorePaths
	errorChan := watcher.Errors

	for {

		select {
		case err := <-errorChan:
			if err != nil {
				log.Fatal(color.RedString("[DEBUG] Closing file watcher error routine %s", err.Error()))
			}

		case event, ok := <-watcher.Events:

			if !ok {
				log.Println(color.MagentaString("[DEBUG] closing file event channel"))
				return
			}
			_, filename := filepath.Split(event.Name)

			// ignore any files that have the same name as the package
			if projectName == filename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Println(color.MagentaString("[DEBUG] ignoring go build artifacts"))
				continue
			}

			// debounces file events
			if event.Name == lastEvent && timeSinceLastEvent.Add(time.Second).After(time.Now()) {
				log.Println(color.MagentaString("[DEBUG] ignoring extra file watch events"))
				timeSinceLastEvent = time.Now()
				continue
			}

			//ignores individual files that may match any ignore paths
			if shouldIgnore(event.Name, ignorePaths) {
				log.Println(color.MagentaString("[DEBUG] %s in ignore path ", event.Name))
				continue
			}

			lastEvent = event.Name
			timeSinceLastEvent = time.Now()

			if !handle.halted {
				fileUpdateNotification <- event.Name
				log.Println(color.MagentaString("[DEBUG] updated %s", event.Name))
			} else {
				return
			}
		}
	}
}
