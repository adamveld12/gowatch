package project

import (
	"github.com/golang/lint"
)

lintEnabled := false


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


