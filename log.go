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
	// rw-lock for buf
	m *sync.RWMutex
	// internal buffer that we periodically flush to the HTTP endpoint
	buf *bytes.Buffer

	// prefix for every line
	prefix string

	// contains temporary buffers for writing log lines to, before
	// flushing to buf.
	bufPool *sync.Pool
}

// New creates a new instance of the httplog.Logger. All Loggers will
// log to the same http endpoint.
func New(prefix string) *Logger {
	lenPrefix := len(prefix)
	buf := &bytes.Buffer{}
	return &Logger{
		prefix: prefix,
		buf:    buf,
		m:      new(sync.RWMutex),
		bufPool: &sync.Pool{
			New: func() interface{} {
				b := &bytes.Buffer{}
				b.Grow(25 + 1 + 1 + lenPrefix + 1024)
				return b
			},
		},
	}
}

func (l *Logger) setupBuffer() *bytes.Buffer {
	rfc3339now := time.Now().Format(time.RFC3339)
	buf := l.bufPool.Get().(*bytes.Buffer)
	if l.prefix != "" {
		buf.WriteString(l.prefix)
	}
	buf.WriteString(rfc3339now)
	buf.WriteString(" ")
	return buf
}

func (l *Logger) printBuffer(buf *bytes.Buffer) {
	l.m.Lock()
	io.Copy(l.buf, buf)
	l.m.Unlock()

	buf.Reset()
	l.bufPool.Put(buf)
}

// Println prints a line to the logger.
func (l *Logger) Println(v ...interface{}) {
	b := l.setupBuffer()
	fmt.Fprintln(b, v...)
	l.printBuffer(b)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	b := l.setupBuffer()
	fmt.Fprintf(b, format, v...)
	l.printBuffer(b)
}

func (l *Logger) Print(v ...interface{}) {
	b := l.setupBuffer()
	fmt.Fprint(b, v...)
	l.printBuffer(b)
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
