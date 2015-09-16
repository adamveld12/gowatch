package watch

import (
	"os"
	"path/filepath"
	"strings"

	gwl "github.com/adamveld12/gowatch/log"
	"gopkg.in/fsnotify.v1"
)

func addFilesToWatch(dir string, ignorePaths []string, watcher *fsnotify.Watcher) {
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if shouldIgnore(filepath.Join(p, info.Name()), ignorePaths) {
				return filepath.SkipDir
			} else if err := watcher.Add(p); err != nil {
				gwl.LogDebug("error adding watched dir", p, info.Name(), err.Error())
				return err
			}
		}
		return nil
	})
}

func shouldIgnore(file string, ignorePaths []string) bool {
	for _, pattern := range ignorePaths {

		matched, err := filepath.Match(strings.Replace(pattern, "/", "", -1), strings.Replace(file, "/", "", -1))

		if err != nil {
			gwl.LogError(err.Error())
		}

		if matched && err == nil {
			gwl.LogDebug("\tIgnore %s -> %s\n", pattern, file)
			return true
		}
	}

	return false
}
