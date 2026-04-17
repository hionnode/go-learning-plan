package hello

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"single line", "hello\n", "hello\n"},
		{"two lines", "hello\nworld\n", "hello\nworld\n"},
		{"no trailing newline", "tail", "tail\n"},
		{"empty", "", ""},
		{"blank line preserved", "a\n\nb\n", "a\n\nb\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			err := Echo(strings.NewReader(tc.in), &out)
			if err != nil {
				t.Fatalf("Echo error: %v", err)
			}
			if got := out.String(); got != tc.want {
				t.Errorf("Echo(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// errReader returns an error after the first read to make sure Echo propagates
// I/O errors instead of swallowing them.
type errReader struct{ data string; read bool }

func (e *errReader) Read(p []byte) (int, error) {
	if e.read {
		return 0, io.ErrUnexpectedEOF
	}
	n := copy(p, e.data)
	e.read = true
	return n, nil
}

func TestEcho_PropagatesReadError(t *testing.T) {
	var out bytes.Buffer
	err := Echo(&errReader{data: "partial"}, &out)
	if err == nil {
		t.Errorf("Echo should propagate read error, got nil")
	}
}
