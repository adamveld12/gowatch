package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func Build(mainFile, cwd string) bool {
	cmd := exec.Command("go", "build", mainFile)
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
		return false
	}

	fmt.Println(string(stdoutOutput))
	return true
}
