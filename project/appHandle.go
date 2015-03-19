package project

func CreateAppHandle(mainFile, cwd string) *AppHandle {
	handle := &AppHandle{mainFile, cwd, false}
	return handle
}

type AppHandle struct {
	main    string
	cwd     string
	running bool
}

func (a *AppHandle) Running() bool {
	return a.running
}

func (a *AppHandle) Halt() {
	a.running = false
}

func (a *AppHandle) Start() bool {
	a.running = true
	if err := command(a.cwd, "go", "run", a.main); err != nil {
		return false
	}
	a.running = false
	return true
}
