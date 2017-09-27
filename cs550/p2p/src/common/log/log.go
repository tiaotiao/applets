package log

import (
	"fmt"
	"time"
)

var LevelDebug = true

var ModuleName = ""

func Debug(format string, args ...interface{}) {
	if LevelDebug {
		output("Debug", format, args...)
	}
}

func Info(format string, args ...interface{}) {
	output("Info", format, args...)
}

func Warning(format string, args ...interface{}) {
	output("Warning", format, args...)
}

func Error(format string, args ...interface{}) {
	output("ERROR", format, args...)
}

func output(tag string, format string, args ...interface{}) {
	now := time.Now()
	t := now.Format("15:04:05")
	s := fmt.Sprintf(format, args...)
	if ModuleName != "" {
		s = fmt.Sprintf("[%s] %s", ModuleName, s)
	}
	fmt.Printf("[%s] %-6s - %s\n", t, tag, s)
}
