package project

import (
	"log"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func build(projectDirectory string) bool {
	cmd := exec.Command("go", "build")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("[DEBUG] building code...")
	if err := cmd.Run(); err != nil {
		log.Println("[DEBUG]", color.RedString(err.Error()))
		return false
	}

	return true
}
