package project

import (
	"os"
	"os/exec"
)

func build(projectDirectory string, outputName string) bool {
	cmd := exec.Command("go", "build", "-o", outputName)
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
