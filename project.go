package main

import (
	"log"
	"path/filepath"
)

type project struct {
	Directory    string
	errorLastRun bool
}

func New(workingDirectory string) project {

	workingDirectory, err := filepath.Abs(*dir)
	if err != nil {
		log.Println("-dir not found", err)
	}

	cwd := workingDirectory
	if filepath.Ext(workingDirectory) != "" {
		cwd = filepath.Dir(workingDirectory)
	}

	return project{Directory: cwd}
}
