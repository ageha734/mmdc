package logger

import (
	"fmt"
	"os"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

type Interface interface {
	Info(msg string)
	Success(msg string)
	Warn(msg string)
	ErrorExit(msg string)
}

type Logger struct {
	quiet    bool
	exitFunc func(int)
}

func New(quiet bool) *Logger {
	return &Logger{quiet: quiet, exitFunc: os.Exit}
}

func NewWithExitFunc(quiet bool, exitFunc func(int)) *Logger {
	return &Logger{quiet: quiet, exitFunc: exitFunc}
}

func (l *Logger) Info(msg string) {
	if !l.quiet {
		fmt.Println(msg)
	}
}

func (l *Logger) Success(msg string) {
	if !l.quiet {
		fmt.Printf("%s%s%s\n", colorGreen, msg, colorReset)
	}
}

func (l *Logger) Warn(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colorYellow, msg, colorReset)
}

func (l *Logger) ErrorExit(msg string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colorRed, msg, colorReset)
	l.exitFunc(1)
}
