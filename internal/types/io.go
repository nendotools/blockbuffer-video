package types

import (
	"io"
	"log"
	"strings"
)

type FilterWriter struct {
	Writer io.Writer
}

func (fw *FilterWriter) Write(p []byte) (n int, err error) {
	// get string from byte slice
	// filter out nuxt.js logs
	// write to writer
	s := string(p)
	if !strings.Contains(s, "WARN  Deprecation") {
		var prefix = "[NUXT] "
		return fw.Writer.Write([]byte(prefix + s))
	}
	return len(p), nil
}

func Writer() io.Writer {
	return &FilterWriter{Writer: log.Writer()}
}
