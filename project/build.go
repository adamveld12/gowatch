package project

import (
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func build(projectDirectory string) bool {
	cmd := exec.Command("go", "build")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		errors, _ := cmd.CombinedOutput()
		color.Red(string(errors))
		return false
	}

	return true
}
