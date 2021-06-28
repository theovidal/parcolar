package lib

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

var StandardLogger = Logger{
	Format: "2006-02-01 15:04:05 — ",
	Output: os.Stdout,
	Levels: map[string]Level{
		"debug":   {"💾", color.New(color.FgWhite)},
		"info":    {"▶", color.New(color.FgCyan)},
		"success": {"✅", color.New(color.FgHiGreen)},
		"error":   {"‼", color.New(color.FgRed)},
	},
}

type Level struct {
	Emoji string
	Color *color.Color
}

type Logger struct {
	Format string
	Debug  bool
	Output io.Writer
	Levels map[string]Level
}

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

func LogDebug(format string, a ...interface{}) {
	if StandardLogger.Debug {
		StandardLogger.log("debug", format, a...)
	}
}

func LogInfo(format string, a ...interface{}) {
	StandardLogger.log("info", format, a...)
}

func LogSuccess(format string, a ...interface{}) {
	StandardLogger.log("success", format, a...)
}

func LogError(format string, a ...interface{}) {
	if StandardLogger.Debug {
		StandardLogger.log("error", format, a...)
	}
}

func Fatal(format string, a ...interface{}) {
	StandardLogger.log("error", format, a...)
	os.Exit(1)
}
