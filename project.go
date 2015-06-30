package main

import (
	"log"
	"path/filepath"
)

type project struct {
	directory string
	name      string
	kill      chan<- bool
	isDone    <-chan bool
}

func (p *project) RunSteps() {
	if p.kill != nil {
		if *debug {
			log.Println("\tkilling process and restarting")
		}
		p.kill <- true
	}

	if built := build(p.directory); built {
		finish, kill := run(p.directory, p.name)
		p.kill = kill
		go func() {
			if <-finish {
				kill = nil
			}
		}()
	}
}

func (p *project) WorkingDirectory() string {
	return p.directory
}

func (p *project) Name() string {
	return p.name
}

func createProject(workingDirectory string) *project {

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
