package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func command(cwd, command, flags, args string) error {
	cmd := exec.Command(command, flags, args)
	cmd.Dir = cwd

	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	stdoutOutput, err := ioutil.ReadAll(stdout)
	stderrOutput, err := ioutil.ReadAll(stderr)

	if processErr := cmd.Wait(); processErr != nil {
		fmt.Println(string(stderrOutput))
		return processErr
	}

	fmt.Println(string(stdoutOutput))
	return nil
}
