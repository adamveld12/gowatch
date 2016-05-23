package log

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/logutils"
)

var (
	errC   = color.New(color.FgRed)
	debugC = color.New(color.FgMagenta)
	infoC  = color.New(color.FgCyan)

	errSprintFunc   = errC.SprintFunc()
	infoSprintFunc  = infoC.SprintFunc()
	debugSprintFunc = debugC.SprintFunc()

	errSprintFmt   = errC.SprintfFunc()
	infoSprintFmt  = infoC.SprintfFunc()
	debugSprintFmt = debugC.SprintfFunc()
	logger         *log.Logger
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

	logger = log.New(filter, "", log.Ltime)
}

// Error prints values with red text with a format string
func Error(a ...interface{}) {
	write("[ERROR]", errSprintFunc(a...))
}

// Info prints values with Cyan text with a format string
func Info(a ...interface{}) {
	write("[INFO]", infoSprintFunc(a...))
}

// Debug prints values with magenta text with a format string
func Debug(a ...interface{}) {
	write("[DEBUG]", debugSprintFunc(a...))
}

// Errorf prints values with red text with a format string
func Errorf(format string, a ...interface{}) {
	write("[ERROR]", errSprintFmt(format, a...))
}

// Infof prints values with Cyan text with a format string
func Infof(format string, a ...interface{}) {
	write("[INFO]", infoSprintFmt(format, a...))
}

// Debugf prints values with magenta text with a format string
func Debugf(format string, a ...interface{}) {
	write("[DEBUG]", debugSprintFmt(format, a...))
}

func write(level, msg string) {
	logger.Println(level, msg)
}
