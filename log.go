// Package httplog contains a logger compatible with the stdlib log
// package. This logger will send output to an HTTP endpoint which can
// then be queries to obtain information.
package httplog

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Logger struct {
	// internal buffer that we periodically flush to the HTTP endpoint
	buf *bytes.Buffer
	m   *sync.RWMutex

	// TODO: prefixing support
}

// New creates a new instance of the httplog.Logger. All Loggers will
// log to the same http endpoint.
func New() *Logger {
	buf := &bytes.Buffer{}
	return &Logger{buf: buf, m: new(sync.RWMutex)}
}

// Println prints a line to the logger.
func (l *Logger) Println(text string) {
	l.m.Lock()
	l.buf.WriteString(time.Now().Format(time.RFC3339))
	l.buf.WriteString(" ")
	l.buf.WriteString(text)
	l.buf.WriteString("\n")
	l.m.Unlock()
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.m.RLock()
	n, err := io.Copy(w, l.buf)
	l.m.RUnlock()

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error: only written %d bytes due to: %s", n, err)))
		return
	}
}
