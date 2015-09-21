package httplog

import (
	"reflect"
	"testing"
)

func TestLoggingToBuffer(t *testing.T) {
	l := New()
	l.Println("hey there")
	if !reflect.DeepEqual(l.buf, []byte("hey there\n")) {
		t.Fatalf("logging to buffer don't work")
	}
}
