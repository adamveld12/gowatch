package log

import (
	"log"

	"github.com/fatih/color"
)

func LogError(format string, a ...interface{}) {
	log.Println("[ERROR]", color.RedString(format, a...))
}

func LogInfo(format string, a ...interface{}) {
	log.Println("[INFO]", color.CyanString(format, a...))
}

func LogDebug(format string, a ...interface{}) {
	log.Println("[DEBUG]", color.MagentaString(format, a...))
}
