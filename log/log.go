package log

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/logutils"
)

// Setup initializes the logging system
func Setup(debug bool) {
	minLevel := logutils.LogLevel("ERROR")

	if debug {
		minLevel = logutils.LogLevel("DEBUG")
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
		MinLevel: minLevel,
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
}

// Errorln prints values with red text terminated with a new line
func Errorln(a ...interface{}) {
	log.Println("[ERROR]", color.RedString(fmt.Sprintln(a...)))
}

// Infoln prints values with cyan text terminated with a new line
func Infoln(a ...interface{}) {
	log.Println("[INFO]", color.CyanString(fmt.Sprintln(a...)))
}

// Debugln prints values with magenta text terminated with a new line
func Debugln(a ...interface{}) {
	log.Println("[DEBUG]", color.MagentaString(fmt.Sprintln(a...)))
}

// LogError prints values with red text with a format string
func LogError(format string, a ...interface{}) {
	log.Println("[ERROR]", color.RedString(format, a...))
}

// LogInfo prints values with Cyan text with a format string
func LogInfo(format string, a ...interface{}) {
	log.Println("[INFO]", color.CyanString(format, a...))
}

// LogDebug prints values with magenta text with a format string
func LogDebug(format string, a ...interface{}) {
	log.Println("[DEBUG]", color.MagentaString(format, a...))
}
