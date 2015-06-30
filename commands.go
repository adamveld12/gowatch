package main

import (
	"log"
	"os"
	"os/exec"
)

func run(dir, file string) (<-chan bool, chan<- bool) {
	if *debug {
		log.Println("running", file)
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
	return gocmd(dir, "build", "")
}

func test(dir string) bool {
	return gocmd(dir, "test", "")
}

func gocmd(cwd, command, target string) bool {
	if *debug {
		log.Println("running", command)
	}

	executable, _ := exec.LookPath("go")

	if target != "" {
		target = " " + target
	}
	cmd := exec.Command(executable, command+target)

	cmd.Dir = cwd
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	success := true
	if err := cmd.Wait(); err != nil {
		success = false
		log.Println("\t", command, "failed.")
		log.Println("\t", err.Error())
	}

	if *debug {
		if success {
			log.Println("\t", command, "succeeded.")
		}
	}

	return success
}
