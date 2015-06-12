package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func run(dir string) bool {
	_, file := filepath.Split(dir)

	if *debug {
		log.Println("Running ", file)
	}

	cmd := exec.Command("./"+file, "run")

	cmd.Dir = dir
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	successfulBuild := true
	if err := cmd.Wait(); err != nil {
		successfulBuild = false
	}

	if *debug {
		if successfulBuild {
			log.Println("succeeded.")
		} else {
			log.Println("failed.")
		}
	}

	return successfulBuild

}

// to run, do go build && ./<program>.exe
func build(dir string) bool {
	if *debug {
		log.Println("building ", dir)
	}

	executable, _ := exec.LookPath("go")

	cmd := exec.Command(executable, "build")

	cmd.Dir = dir
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	successfulRun := true
	if err := cmd.Wait(); err != nil {
		successfulRun = false
	}

	if *debug {
		if successfulRun {
			log.Println("Build succeeded.")
		} else {
			log.Println("Build failed.")
		}
	}

	return successfulRun
}
