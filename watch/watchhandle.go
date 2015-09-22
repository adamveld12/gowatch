package watch

import (
	"path/filepath"
	"sync"
	"time"

	gwl "github.com/adamveld12/gowatch/log"
	"gopkg.in/fsnotify.v1"
)

type WatchHandle struct {
	sync.Mutex
	halted       bool
	watcher      *fsnotify.Watcher
	ignorePaths  []string
	fileUpdateCb func(string)
}

func (w *WatchHandle) Halt() {
	w.Lock()
	defer w.Unlock()
	if !w.halted {
		gwl.LogDebug("Halting watcher")
		w.halted = true
		w.fileUpdateCb = nil
		go func() {
			w.watcher.Close()
		}()
	}
}

func (handle *WatchHandle) Subscribe(cb func(string)) {
	handle.Lock()
	handle.fileUpdateCb = cb
	handle.Unlock()
}

func (handle *WatchHandle) handleFileEvent(projectName string) {
	lastEvent, timeSinceLastEvent := "", time.Now().AddDate(-1, 0, 0)

	watcher := handle.watcher
	ignorePaths := handle.ignorePaths
	errorChan := watcher.Errors

	for {
		select {
		case <-time.After(time.Second):
			continue

		case err := <-errorChan:
			if err != nil {
				gwl.LogError("Closing file watcher error routine %s", err.Error())
			}

		case event, ok := <-watcher.Events:

			if !ok {
				gwl.LogDebug("closing file event channel")
				return
			}
			_, filename := filepath.Split(event.Name)

			// ignore any files that have the same name as the package
			if projectName == filename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
				gwl.LogDebug("ignoring go build artifacts")
				continue
			}

			//ignores individual files that may match any ignore paths
			if shouldIgnore(event.Name, ignorePaths) {
				gwl.LogDebug("%s in ignore path ", event.Name)
				continue
			}

			// debounces file events
			if event.Name == lastEvent && timeSinceLastEvent.Add(time.Second).After(time.Now()) {
				gwl.LogDebug("ignoring extra file watch events")
				timeSinceLastEvent = time.Now()
				continue
			}

			lastEvent = event.Name
			timeSinceLastEvent = time.Now()

			if handle.halted {
				gwl.LogDebug("\tfilewatcher halteds")
				return
			} else if handle.fileUpdateCb != nil {
				gwl.LogDebug("\tinvoking file callback with %s", event.Name)
				handle.Lock()
				handle.fileUpdateCb(event.Name)
				handle.Unlock()
				gwl.LogDebug("\tfile callback complete")
			}
		}
	}
}
