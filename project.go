package main

import (
	"log"
	"path/filepath"
)

type project struct {
	directory    string
	name         string
	errorLastRun bool
}

func (p *project) RunSteps() bool {
	if buildSucceeded := build(p.directory); buildSucceeded {
		if runFailed := run(p.directory, p.name); runFailed {
			return true
		}
		return false
	}

	return false
}

func (p *project) WorkingDirectory() string {
	return p.directory
}

func (p *project) Name() string {
	return p.name
}

func New(workingDirectory string) *project {

	workingDirectory, err := filepath.Abs(*dir)
	if err != nil {
		log.Println("-dir not found", err)
	}

	cwd := workingDirectory
	if filepath.Ext(workingDirectory) != "" {
		cwd = filepath.Dir(workingDirectory)
	}

	_, projectName := filepath.Split(cwd)

	return &project{directory: cwd, name: projectName}
}
