package project

import (
//"github.com/golang/lint"
)

var lintEnabled bool = false

func EnableLinting() {
	lintEnabled = true
}

func RunLint(mainFile, cwd string) bool {
	if !lintEnabled {
		return true
	}

	lintSuccessful := false

	return lintSuccessful
}
