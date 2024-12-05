package io

import (
	"fmt"
	"os"
	"strings"

	opts "blockbuffer/internal/settings"
)

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
	Fatal LogLevel = "FATAL"
	Panic LogLevel = "PANIC"
	Nuxt  LogLevel = "NUXT"
)

const (
	ResetColor  = "\033[0m"
	DebugPrefix = "\033[0;32m[DEBUG] "
	InfoPrefix  = "\033[0;36m[INFO] "
	WarnPrefix  = "\033[0;33m[WARN] "
	ErrorPrefix = "\033[0;31m[ERROR] "
	FatalPrefix = "\033[0;31m[FATAL] "
	PanicPrefix = "\033[0;31m[PANIC] "
	NuxtPrefix  = "\033[0;35m[NUXT]\033[0m "
)

func CheckError(err error) {
	if err != nil {
		Log(err.Error(), Error)
	}
}

func severityPrefix(severity LogLevel) string {
	// Return the prefix for the log message with color
	switch severity {
	case Debug:
		return fmt.Sprintf("%s%s", DebugPrefix, ResetColor)
	case Info:
		return fmt.Sprintf("%s%s", InfoPrefix, ResetColor)
	case Warn:
		return fmt.Sprintf("%s%s", WarnPrefix, ResetColor)
	case Error:
		return fmt.Sprintf("%s%s", ErrorPrefix, ResetColor)
	case Fatal:
		return fmt.Sprintf("%s%s", FatalPrefix, ResetColor)
	case Panic:
		return fmt.Sprintf("%s%s", PanicPrefix, ResetColor)
	}
	return ""
}

func filterSeverity(severity LogLevel) bool {
	// need to ensure opts.LogLevel is not nil and is all-caps
	target := LogLevel(strings.ToUpper(*opts.LogLevel))
	if target == "" {
		target = Info
	}

	// rank the severity and filter based on the log level
	sRank := map[LogLevel]int{
		Debug: 0,
		Info:  1,
		Warn:  2,
		Error: 3,
	}

	if sRank[severity] < sRank[target] {
		return true
	}
	return false
}

func Log(msg string, severity ...LogLevel) {
	// Default to info level
	var sv LogLevel = Info
	if len(severity) > 0 {
		sv = LogLevel(severity[0])
	}

	if filterSeverity(sv) {
		return
	}

	s := fmt.Sprintf("%s%s\n", severityPrefix(sv), msg)
	fmt.Printf(s)
	if sv == Fatal {
		os.Exit(1)
	}
	if sv == Panic {
		panic(s)
	}
}

func Logf(msg string, severity LogLevel, args ...interface{}) {
	s := fmt.Sprintf("%s%s\n", severityPrefix(severity), msg)
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
