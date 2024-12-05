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
	s := strings.TrimSpace(string(p))
	if !strings.Contains(s, "WARN  Deprecation") && s != "" {
		var prefix = "[NUXT] "
		return fw.Writer.Write([]byte(prefix + s + "\n"))
	}
	return len(p), nil
}

func Writer() io.Writer {
	return &FilterWriter{Writer: log.Writer()}
}
