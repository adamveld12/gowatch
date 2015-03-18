package project

import ()

var testingEnabled bool = false

func EnableTesting() {
	testingEnabled = true
}

func Test(mainFile, cwd string) bool {
	if !testingEnabled {
		return true
	}

	testsPassed := false

	return testsPassed
}
