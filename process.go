package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func command(cwd, command, args string) (success bool) {
	cmd := exec.Command(command, args)
	cmd.Dir = cwd

	stdin, err := cmd.StdinPipe()
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	return success
}
