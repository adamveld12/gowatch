package main

import (
	"errors"
	"log"
	"path/filepath"
)

type project struct {
	directory string
	name      string
	kill      chan<- bool
	isDone    <-chan bool
}

type StepResult error

var (
	BuildFailed = errors.New("Build failed")
	RunFailed   = errors.New("App exited with non-zero exit code")
	TestFailed  = errors.New("Test failed")
	LintFailed  = errors.New("Lint failed")
)

func (p *project) RunSteps() StepResult {
	if p.kill != nil {
		if *debug {
			log.Println("\tkilling process and restarting")
		}
		p.kill <- true
	}

	if built := build(p.directory); built {

		finish, kill := run(p.directory, p.name)
		p.kill = kill
		p.isDone = finish

		go func() {
			if <-finish {
				p.kill = nil
				p.isDone = nil
			}
		}()

	} else {
		return BuildFailed
	}

	return nil
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
