package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type StepResult error

var (
	BuildFailed = errors.New("Build failed")
	RunFailed   = errors.New("App exited with non-zero exit code")
	TestFailed  = errors.New("Test failed")
	LintFailed  = errors.New("Lint failed")
)

func build(projectDirectory string) bool {
	goPath := os.Getenv("GOPATH")
	cmd := gocmd("build", strings.TrimPrefix(projectDirectory, filepath.Join(goPath, "src/")+"/"), projectDirectory)
	err := cmd.Run()
	return err == nil
}

func test(projectDirectory string) bool {
	return true
}

func lint(projectDirectory string) bool {
	return true
}

func runProject(projectDirectory string, arguments string) (<-chan StepResult, chan<- os.Signal) {
	routineSync, isDone, killApp := make(chan bool), make(chan StepResult), make(chan os.Signal)

	// build
	if !build(projectDirectory) {
		isDone <- BuildFailed
		close(isDone)
		return isDone, killApp
	}

	// lint
	if *shouldLint && !lint(projectDirectory) {
		isDone <- LintFailed
		close(isDone)
		return isDone, killApp
	}

	// test
	if *shouldTest && !test(projectDirectory) {
		isDone <- TestFailed
		close(isDone)
		return isDone, killApp
	}

	cmd := run(projectDirectory, arguments)
	exited := false

	go func() {
		for {
			select {
			case exitSignal := <-killApp:
				if exitSignal != nil {
					cmd.Process.Signal(exitSignal)
				}
				if *debug {
					log.Println("\tSending kill signal...")
				}
				return
			default:
				if exited {
					return
				}
			}
		}

	}()

	go func() {
		close(routineSync)
		err := cmd.Run()

		if *debug {
			log.Println("App has exited", err)
		}

		exited = true
		isDone <- err
		close(isDone)
	}()

	<-routineSync

	return isDone, killApp
}

func run(projectDirectory, arguments string) *exec.Cmd {
	_, command := filepath.Split(projectDirectory)
	cmd := exec.Command("./"+command, arguments)
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func gocmd(command, arg, projectDirectory string) *exec.Cmd {

	cmd := exec.Command("go", command, arg)
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
