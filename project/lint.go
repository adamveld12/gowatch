package project

import (
//"github.com/golang/lint"
)

var lintEnabled bool = false

func EnableLinting() {
	lintEnabled = true
}

func RunLint() bool {
	if !lintEnabled {
		return true
	}

	lintSuccessful := true

	return lintSuccessful
}
