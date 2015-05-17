package project

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

var testingEnabled bool = false

func EnableTesting() {
	testingEnabled = true
}

func Test(mainFile, cwd string) bool {
	if !testingEnabled {
		return true
	}

	testsPassed := false

	cmd := exec.Command("go", "test", mainFile)
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

	return testsPassed
}
