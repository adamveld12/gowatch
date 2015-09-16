package project

import (
	"log"
	"os"
	"os/exec"
)

func test(projectDirectory string) bool {
	log.Println("[DEBUG] testing code...")

	cmd := exec.Command("go", "test")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("[DEBUG] test failures:", err)
		return false
	}

	return true
}
