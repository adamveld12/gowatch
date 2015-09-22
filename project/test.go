package project

import (
	"os"
	"os/exec"

	gwl "github.com/adamveld12/gowatch/log"
)

func test(projectDirectory string) bool {
	gwl.LogDebug("testing code...")

	cmd := exec.Command("go", "test")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		gwl.LogDebug("test failures: ", err)
		return false
	}

	return true
}
