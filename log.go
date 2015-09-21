// Package httplog contains a logger compatible with the stdlib log
// package. This logger will send output to an HTTP endpoint which can
// then be queries to obtain information.
package httplog

type Logger struct {
	// internal buffer that we periodically flush to the HTTP endpoint
	buf []byte
	// TODO: prefixing support
}

// New creates a new instance of the httplog.Logger. All Loggers will
// log to the same http endpoint.
func New() *Logger {
	return &Logger{
		buf: make([]byte, 0),
	}
}

func (l *Logger) Println(text string) {
	l.buf = append(l.buf, []byte(text)...)
	l.buf = append(l.buf, []byte("\n")...)
}
