package exec

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

type processHandler struct {
	output    <-chan string
	input     chan<- string
	exitError error
	done      bool
	process   *os.Process
}

func (p *processHandler) exited(err error) {
	p.exitError = err
	p.done = true
}

func (p *processHandler) Output() <-chan string {
	return p.output
}

func (p *processHandler) Input() chan<- string {
	return p.input
}

func (p *processHandler) Stop() {
	p.process.Signal(os.Kill)
	p.Wait()
}

func (p *processHandler) Exited() bool {
	return p.done
}

func (p *processHandler) Wait() error {
	status, err := p.process.Wait()
	p.done = status.Exited()
	return err
}

// ProcessHandle is a facade for interacting with a running process
type ProcessHandle interface {
	// Stops the underlying process with SIGKILL
	Stop()
	// A channel to the underlying process's output stream
	// An aggregate of the stdin and stderr
	Output() <-chan string
	// A channel to the underlying process's input stream
	Input() chan<- string
	// True if the process exited successfully
	Exited() bool
	// Waits until the underlying process exits
	Wait() error
}

func New(cwd, command, flags, args string) ProcessHandle {
	cmd := exec.Command(command, flags, args)
	cmd.Dir = cwd

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	// stdin, _ := cmd.StdinPipe()

	outReader := bufio.NewReader(stdout)
	errReader := bufio.NewReader(stderr)
	// stdinWriter := bufio.NewWriter(stdin)

	processOutput := make(chan string)
	processInput := make(chan string)

	procHandle := &processHandler{processOutput, processInput, nil, false, cmd.Process}

	// pump the output to std out
	go func() {
		for {
			output, err := outReader.ReadString(byte('\n'))
			fmt.Println(output)
			if err != nil {
				break
			}
			errOutput, err := errReader.ReadString(byte('\n'))
			fmt.Println(errOutput)
			if err != nil {
				break
			}
		}
	}()

	// so that we don't block
	go func() {
		if processErr := cmd.Wait(); processErr != nil {
			procHandle.exited(processErr)
			return
		}

		procHandle.exited(nil)
	}()

	return procHandle
}
