package httplog

import (
	"strings"
	"testing"
	"time"
)

func TestLoggingToBuffer(t *testing.T) {
	l := New("")
	l.Println("hey there")
	if !strings.HasSuffix(l.buf.String(), "hey there\n") {
		t.Fatalf("logging the actual text to the buffer don't work: %s", l.buf.String())
	}
}

func TestEnsureTimeIsLoggedAsRFC3339(t *testing.T) {
	l := New("")
	l.Println("ok")
	x := l.buf.String()
	xs := strings.Split(x, " ")
	_, err := time.Parse(time.RFC3339, strings.TrimSpace(xs[0]))
	if err != nil {
		t.Fatalf("expected first non-space characters to be the RFC3339 timestamp: %s", err)
	}
}

func TestEnsurePrefixCanBeSet(t *testing.T) {
	goodpref := "testing!"
	l := New(goodpref)
	l.Println("ok")
	x := l.buf.String()
	if !strings.HasPrefix(x, goodpref) {
		t.Fatalf("expected prefix %s in %s", goodpref, x)
	}
}

func BenchmarkWriteToBufferWithPrefix(b *testing.B) {
	l := New("testing! ")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Println("this is a test log. log log log log log.")
	}
}

func BenchmarkWriteToBufferWithPrefixPar(b *testing.B) {
	l := New("testing! ")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println("this is a test log. log log log log log.")
		}
	})
}
