package watch

import (
	"log"
	"path/filepath"
	"sync"

	"gopkg.in/fsnotify.v1"
)

func StartWatch(dir string, ignorePaths []string) *WatchHandle {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal("Could not create file watcher", err)
	}

	addFilesToWatch(dir, ignorePaths, watcher)

	handle := &WatchHandle{
		sync.Mutex{},
		false,
		watcher,
		ignorePaths,
		nil,
	}

	_, projectName := filepath.Split(dir)
	go handle.handleFileEvent(projectName)

	return handle
}
