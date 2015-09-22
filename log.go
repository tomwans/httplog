// Package httplog contains a logger compatible with the stdlib log
// package. This logger will send output to an HTTP endpoint which can
// then be queries to obtain information.
package httplog

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Logger struct {
	// internal buffer that we periodically flush to the HTTP endpoint
	buf *bytes.Buffer
	// TODO: prefixing support
}

// New creates a new instance of the httplog.Logger. All Loggers will
// log to the same http endpoint.
func New() *Logger {
	buf := &bytes.Buffer{}
	return &Logger{buf: buf}
}

// Println prints a line to the logger.
func (l *Logger) Println(text string) {
	// l.buf = append(l.buf, []byte(text)...)
	// l.buf = append(l.buf, []byte("\n")...)
	l.buf.WriteString(time.Now().Format(time.RFC3339))
	l.buf.WriteString(" ")
	l.buf.WriteString(text)
	l.buf.WriteString("\n")
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, l.buf)
}
