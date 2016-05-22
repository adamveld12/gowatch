package gowatch

import (
	"log"
	"path/filepath"
)

func AbsPathToProject(workingDirectory string) string {
	if workingDirectory == "" {
		workingDirectory = "."
	}

	directoryPath, err := filepath.Abs(workingDirectory)

	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	return directoryPath
}
