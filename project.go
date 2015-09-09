package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	linter "github.com/golang/lint"
)

// StepResult the result of a build step
type StepResult error

var (
	// ErrorBuildFailed indicates build failure
	ErrorBuildFailed = errors.New("Build failed")

	// ErrorRunFailed indicates app exited with error
	ErrorRunFailed = errors.New("App exited with non-zero exit code")

	// ErrorTestFailed indicates one or more tests failed
	ErrorTestFailed = errors.New("Test failed")

	// ErrorLintFailed indicates Linter errors
	ErrorLintFailed = errors.New("Lint failed")

	errorProcessAlreadyFinished = errors.New("process already finished")
)

func build(projectDirectory string) bool {
	cmd := exec.Command("go", "build")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("[DEBUG] building code...")
	if err := cmd.Run(); err != nil {
		log.Println("[ERROR]", err)
		return false
	}

	return true
}

func test(projectDirectory string) bool {
	log.Println("[DEBUG] testing code...")

	cmd := exec.Command("go", "test")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("[DEBUG] test failures:", err)
		return false
	}

	return true
}

func lint(projectDirectory string) bool {
	log.Println("[DEBUG] linting code...")
	lint := &linter.Linter{}

	files := make(map[string]map[string][]byte)
	filepath.Walk(projectDirectory, func(p string, info os.FileInfo, err error) error {
		if filepath.Ext(p) == ".go" {
			fileWithPackage := strings.TrimPrefix(p, projectDirectory)
			packageName := strings.Trim(strings.TrimSuffix(fileWithPackage, info.Name()), "/")

			if packageName == "" {
				packageName = "main"
			}

			files[packageName] = make(map[string][]byte)

			f, err := os.Open(p)
			if err != nil {
				return err
			}

			if files[packageName][p], err = ioutil.ReadAll(f); err != nil {
				return err
			}
		}
		return nil
	})

	lintErrors := false
	for k, v := range files {
		log.Println("[DEBUG] linting package", k)

		problems, err := lint.LintFiles(v)

		if err != nil {
			color.Red("[ERROR]", err)
			lintErrors = true
		} else if len(problems) > 0 {

			log.Println("[DEBUG] lint issues found")
			color.Yellow("%d lint issue(s) found in %s\n\n", len(problems), k)
			linterConfidenceThresholdReached := false
			for i, p := range problems {
				position := p.Position
				fileWithPackage := strings.Trim(strings.TrimPrefix(position.Filename, projectDirectory), "/")
				lintInfo := strings.Split(p.String(), "\n")

				lintLineOutput := fmt.Sprintf("\t%d. %s line %d - %s\n\t%s\n\n",
					i+1, fileWithPackage, position.Line, lintInfo[2], lintInfo[0])

				if p.Confidence > 0.5 {
					color.Red(lintLineOutput)
					linterConfidenceThresholdReached = true
				} else {
					color.Yellow(lintLineOutput)
				}
			}
			lintErrors = linterConfidenceThresholdReached
		}
	}

	return !lintErrors
}

func executeBuildSteps(projectDirectory, appArguments string) (<-chan StepResult, chan<- os.Signal) {

	isDone, killApp := make(chan StepResult, 1), make(chan os.Signal)

	if !build(projectDirectory) {
		isDone <- ErrorBuildFailed
	} else if *shouldLint && !lint(projectDirectory) {
		color.Red("Linter found errors.")
		isDone <- ErrorLintFailed
	} else if *shouldTest && !test(projectDirectory) {
		color.Red("Tests failed.")
		isDone <- ErrorTestFailed
	} else {
		color.Green("Starting...")
		return runProject(projectDirectory, appArguments)
	}
	log.Println("[DEBUG] build steps completed")

	close(isDone)
	return isDone, killApp
}

func runProject(projectDirectory string, arguments string) (<-chan StepResult, chan<- os.Signal) {
	routineSync, isDone, killApp := make(chan bool), make(chan StepResult), make(chan os.Signal)
	cmd := run(projectDirectory, arguments)
	exited := false

	go func() {
		for {
			select {
			case <-killApp:
				log.Println("[DEBUG] killing app")
				if err := cmd.Process.Kill(); err != nil && err.Error() != errorProcessAlreadyFinished.Error() {
					log.Fatal("[DEBUG] wow this sucks", err)
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
		log.Println("[DEBUG] app has exited", err)

		if err != nil {
			color.Red("exited with", err)
		} else {
			color.Green("exited successfully (0)")
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
