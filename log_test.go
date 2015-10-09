package httplog

import (
	"strings"
	"testing"
	"time"
)

func TestLoggingToBuffer(t *testing.T) {
	l := New("", 1000)
	l.Println("hey there")
	if !strings.HasSuffix(l.buf.String(), "hey there\n") {
		t.Fatalf("logging the actual text to the buffer don't work: %s", l.buf.String())
	}
}

func TestEnsureTimeIsLoggedAsRFC3339(t *testing.T) {
	l := New("", 1000)
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
	l := New(goodpref, 1000)
	l.Println("ok")
	x := l.buf.String()
	if !strings.HasPrefix(x, goodpref) {
		t.Fatalf("expected prefix %s in '%s'", goodpref, x)
	}
}

func TestPrint(t *testing.T) {
	l := New("", 1000)
	l.Print("ok")
	x := l.buf.String()
	if !strings.HasSuffix(x, "ok") {
		t.Fatalf("expected suffix %s in '%s'", "ok", x)
	}
}

func TestPrintf(t *testing.T) {
	l := New("", 1000)
	l.Printf("%s %d", "ok", 1)
	x := l.buf.String()
	if !strings.HasSuffix(x, "ok 1") {
		t.Fatalf("expected suffix %s in '%s'", "ok 1", x)
	}
}

func TestPrintln(t *testing.T) {
	l := New("", 1000)
	l.Println("um", "ok")
	x := l.buf.String()
	if !strings.HasSuffix(x, "um ok\n") {
		t.Fatalf("expected suffix %s in '%s'", "um ok", x)
	}
}

func TestTruncateBytes(t *testing.T) {
	l := New("", 30)
	l.Println("AA")
	l.Println("BB")
	x := l.buf.String()
	// we expect the AA line to have been evicted to make room for BB.
	expected := "BB"
	if !strings.HasSuffix(x, expected+"\n") {
		t.Fatalf("expected suffix %s in '%s'", expected, x)
	}
}

func BenchmarkWriteToBufferWithPrefix(b *testing.B) {
	l := New("testing! ", 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Println("this is a test log. log log log log log.")
	}
}

func BenchmarkWriteToBufferWithPrefixPar(b *testing.B) {
	l := New("testing! ", 1000)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println("this is a test log. log log log log log.")
		}
	})
}
