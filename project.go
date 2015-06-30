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
	isDone    <-chan error
}

type StepResult error

var (
	BuildFailed = errors.New("Build failed")
	RunFailed   = errors.New("App exited with non-zero exit code")
	TestFailed  = errors.New("Test failed")
	LintFailed  = errors.New("Lint failed")
)

func (p *project) RunSteps() <-chan error {
	if p.kill != nil {
		if *debug {
			log.Println("\tkilling process and restarting")
		}
		p.kill <- true
	}
	isDoneSender := make(chan error, 1)

	if built := build(p.directory); built {

		finish, kill := run(p.directory, p.name)
		p.kill = kill
		p.isDone = isDoneSender

		go func() {
			if exitError := <-finish; exitError != nil {
				p.kill = nil
				isDoneSender <- exitError
				p.isDone = nil
			}
		}()

	} else {
		isDoneSender <- BuildFailed
	}

	return isDoneSender
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
