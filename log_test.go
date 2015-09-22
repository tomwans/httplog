package httplog

import (
	"strings"
	"testing"
	"time"
)

func TestLoggingToBuffer(t *testing.T) {
	l := New()
	l.Println("hey there")
	if !strings.HasSuffix(l.buf.String(), "hey there\n") {
		t.Fatalf("logging to buffer don't work")
	}
}

func TestEnsureTimeIsLoggedAsRFC3339(t *testing.T) {
	l := New()
	l.Println("ok")
	x := l.buf.String()
	xs := strings.Split(x, " ")
	_, err := time.Parse(time.RFC3339, xs[0])
	if err != nil {
		t.Fatalf("expected first non-space characters to be the RFC3339 timestamp: %s", err)
	}
}
