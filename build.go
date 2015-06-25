package main

import (
	"log"
	"os"
	"os/exec"
)

func run(dir, file string) (<-chan bool, chan<- bool) {

	if *debug {
		log.Println("Running", file)
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

	procSignal := make(chan bool)
	killSignal := make(chan bool)
	killed := false
	go func() {
		if err := cmd.Wait(); err != nil {
			procSignal <- false
		} else {
			procSignal <- true
		}
		close(procSignal)
		killed = true

	}()

	go func() {
		for !killed {
			select {
			case <-killSignal:
				cmd.Process.Kill()
				procSignal <- true
				killed = true
			}
		}
	}()

	return procSignal, killSignal
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

	successfulBuild := true
	if err := cmd.Wait(); err != nil {
		successfulBuild = false
		log.Println("\tbuild error.")
		log.Println("\t", err)
	}

	if *debug {
		if successfulBuild {
			log.Println("\tsucceeded.")
		} else {
			log.Println("\tbuild failed.")
		}
	}

	return successfulBuild

}
