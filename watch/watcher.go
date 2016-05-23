package watch

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"

	gwl "github.com/adamveld12/gowatch/log"
)

func StartWatch(dir, outputName string, ignorePaths []string) *WatchHandle {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Fatal("Could not create file watcher", err)
	}

	for _, filteredDir := range filterDirectories(dir, ignorePaths) {
		if err := watcher.Add(filteredDir); err != nil {
			gwl.Debug("error adding watched dir", filteredDir, err.Error())
		}
	}

	handle := &WatchHandle{
		sync.Mutex{},
		false,
		watcher,
		ignorePaths,
		nil,
	}

	go handle.handleFileEvent(outputName)

	return handle
}

// TODO needs tests
func filterDirectories(dir string, ignorePaths []string) []string {
	dirList := []string{}

	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if shouldIgnore(filepath.Join(p, info.Name()), ignorePaths) {
				return filepath.SkipDir
			} else {
				dirList = append(dirList, p)
			}
		}
		return nil
	})

	return dirList
}

// TODO needs tests
func shouldIgnore(file string, ignorePaths []string) bool {
	for _, pattern := range ignorePaths {

		matched, err := filepath.Match(strings.Replace(pattern, "/", "", -1), strings.Replace(file, "/", "", -1))

		if err != nil {
			gwl.Error(err.Error())
		}

		if matched && err == nil {
			gwl.Debug("\tIgnore %s -> %s\n", pattern, file)
			return true
		}
	}

	return false
}
