package io

import (
	"fmt"
	"os"
)

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
	Fatal LogLevel = "FATAL"
	Panic LogLevel = "PANIC"
)

func CheckError(err error) {
	if err != nil {
		Log(err.Error(), Error)
	}
}

func Log(msg string, severity ...LogLevel) {
	// Default to info level
	var sv LogLevel = Info
	if len(severity) > 0 {
		sv = LogLevel(severity[0])
	}

	s := fmt.Sprintf("[%s] %s\n", sv, msg)
	fmt.Println(s)
	if sv == Fatal {
		os.Exit(1)
	}
	if sv == Panic {
		panic(s)
	}
}

func Logf(msg string, severity LogLevel, args ...interface{}) {
	s := fmt.Sprintf("[%s] %s\n", severity, msg)
	fmt.Printf(s, args...)
	if severity == Fatal {
		os.Exit(1)
	}
	if severity == Panic {
		panic(s)
	}
}

func Panicf(v ...any) {
	s := fmt.Sprint(v...)
	Log(s, Panic)
}
