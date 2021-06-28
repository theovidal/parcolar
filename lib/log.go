package lib

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

// StandardLogger is the default logger for the app, with predetermined parameters and levels
var StandardLogger = Logger{
	Format: "2006-01-02 15:04:05 — ",
	Output: os.Stdout,
	Levels: map[string]Level{
		"debug":   {"💾", color.New(color.FgWhite)},
		"info":    {"▶", color.New(color.FgCyan)},
		"success": {"✅", color.New(color.FgHiGreen)},
		"error":   {"‼", color.New(color.FgRed)},
	},
}

// Level describes a log level with an emoji and color for console display
type Level struct {
	Emoji string
	Color *color.Color
}

// Logger is an instance that displays logs with a defined style in the console
type Logger struct {
	// Format used for the date, uses `time` package
	Format string
	// Print all debug and error messages or not
	Debug bool
	// Where to print the logs
	Output io.Writer
	// List of the log levels (debug, info, success and error)
	Levels map[string]Level
}

// log prints a log message in the console using a specific level
func (l *Logger) log(levelName, format string, a ...interface{}) {
	level := l.Levels[levelName]
	_, err := fmt.Fprintln(
		l.Output,
		level.Color.Sprintf(
			time.Now().Format(l.Format)+level.Emoji+" "+format,
			a...,
		),
	)
	if err != nil {
		panic(err)
	}
}

// LogDebug prints a debug log using the StandardLogger (if debug mode is enabled)
func LogDebug(format string, a ...interface{}) {
	if StandardLogger.Debug {
		StandardLogger.log("debug", format, a...)
	}
}

// LogInfo prints an information log using the StandardLogger
func LogInfo(format string, a ...interface{}) {
	StandardLogger.log("info", format, a...)
}

// LogSuccess prints a success log using the StandardLogger
func LogSuccess(format string, a ...interface{}) {
	StandardLogger.log("success", format, a...)
}

// LogError prints an error log using the StandardLogger (if debug mode is enabled)
func LogError(format string, a ...interface{}) {
	if StandardLogger.Debug {
		StandardLogger.log("error", format, a...)
	}
}

// Fatal prints an error log using the StandardLogger and exists the program
func Fatal(format string, a ...interface{}) {
	StandardLogger.log("error", format, a...)
	os.Exit(1)
}
