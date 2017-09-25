package log

import (
	"fmt"
	"time"
)

func Debug(format string, args ...interface{}) {
	output("Debug\t", format, args)
}

func Info(format string, args ...interface{}) {
	output("Info\t", format, args)
}

func Warning(format string, args ...interface{}) {
	output("Warning\t", format, args)
}

func Error(format string, args ...interface{}) {
	output("ERROR\t", format, args)
}

func output(tag string, format string, args ...interface{}) {
	now := time.Now()
	t := now.Format("15:04:05")
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s - %s\n", t, tag, s)
}
