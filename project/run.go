package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	gwl "github.com/adamveld12/gowatch/log"
)

type ExecuteHandle struct {
	sync.Mutex
	projectDirectory string
	result           chan StepResult
	halted           bool
	cmd              *exec.Cmd
	running          bool
}

func (h *ExecuteHandle) Running() bool {
	h.Lock()
	defer h.Unlock()
	return h.running
}

func (h *ExecuteHandle) Error() StepResult {
	return <-h.result
}

// Kill kills the underlying application if its started
func (h *ExecuteHandle) Kill(reason StepResult) {

	if reason == nil {
		reason = ErrorAppKilled
	}

	h.Lock()
	if h.running {
		cmd := h.cmd
		proc := cmd.Process

		gwl.LogDebug("Killing")

		if proc != nil {
			if err := proc.Kill(); err != nil && err.Error() != errorProcessAlreadyFinished.Error() {
				gwl.LogDebug("process didn't seem to exit gracefully", err)
				h.writeError(err)
			} else {
				h.writeError(reason)
			}
		} else {
			h.writeError(reason)
		}

		h.running = false
		h.halted = true
		close(h.result)
	} else {
		gwl.LogDebug("process never started", reason.Error())
		h.writeError(reason)
	}

	h.Unlock()
}

func (h *ExecuteHandle) writeError(reason StepResult) {
	if h.running {
		gwl.LogDebug("sending error")
		h.result <- reason
	}
}

func (h *ExecuteHandle) Halted() bool {
	h.Lock()
	defer h.Unlock()
	return h.halted
}

func (h *ExecuteHandle) start(cmd *exec.Cmd) {
	h.Lock()
	h.cmd = cmd
	err := cmd.Start()
	h.running = true
	h.Unlock()

	if err != nil {
		h.Kill(err)
	}

	waiter := make(chan bool)
	go func() {
		close(waiter)
		if err := cmd.Wait(); err != nil {
			gwl.LogDebug("app exited prematurely")
			h.Kill(err)
		}
	}()
	<-waiter
}

func run(projectDirectory, arguments string) *exec.Cmd {
	_, command := filepath.Split(projectDirectory)
	cmd := exec.Command("./"+command, arguments)
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}
