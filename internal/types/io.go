package types

import (
	appIO "blockbuffer/internal/io"
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
		return fw.Writer.Write([]byte(appIO.NuxtPrefix + s + "\n"))
	}
	return len(p), nil
}

func Writer() io.Writer {
	return &FilterWriter{Writer: log.Writer()}
}
