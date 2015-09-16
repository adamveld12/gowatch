package watch

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/fsnotify.v1"
)

func addFilesToWatch(dir string, ignorePaths []string, watcher *fsnotify.Watcher) {
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if shouldIgnore(filepath.Join(p, info.Name()), ignorePaths) {
				return filepath.SkipDir
			} else if err := watcher.Add(p); err != nil {
				log.Println(color.MagentaString("[DEBUG] error adding watched dir", p, info.Name(), err.Error()))
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
			log.Println("[ERROR]", err.Error())
		}

		if matched && err == nil {
			log.Printf(color.MagentaString("[DEBUG] \tIgnore %s -> %s\n", pattern, file))
			return true
		}
	}

	return false
}
