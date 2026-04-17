package dispatcher

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestMessage_JSONTags(t *testing.T) {
	m := Message{Type: "echo", Body: "hi"}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, `"type":"echo"`) {
		t.Errorf("expected lowercase 'type' tag, got %s", got)
	}
	if !strings.Contains(got, `"body":"hi"`) {
		t.Errorf("expected lowercase 'body' tag, got %s", got)
	}
}

func TestEchoHandler(t *testing.T) {
	got, err := EchoHandler{}.Handle(Message{Type: "echo", Body: "hello"})
	if err != nil {
		t.Fatalf("EchoHandler error: %v", err)
	}
	if got.Body != "hello" {
		t.Errorf("EchoHandler Body = %q, want %q", got.Body, "hello")
	}
}

func TestReverseHandler(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"日本語", "語本日"}, // must reverse by rune, not byte
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ReverseHandler{}.Handle(Message{Body: tc.in})
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if got.Body != tc.want {
				t.Errorf("Reverse(%q) = %q, want %q", tc.in, got.Body, tc.want)
			}
		})
	}
}

func TestDispatcher_Echo(t *testing.T) {
	d := NewDispatcher()
	d.Register("echo", EchoHandler{})
	out, err := d.Dispatch([]byte(`{"type":"echo","body":"hi"}`))
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	var m Message
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if m.Body != "hi" {
		t.Errorf("body = %q", m.Body)
	}
}

func TestDispatcher_Reverse(t *testing.T) {
	d := NewDispatcher()
	d.Register("reverse", ReverseHandler{})
	out, err := d.Dispatch([]byte(`{"type":"reverse","body":"abc"}`))
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	var m Message
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m.Body != "cba" {
		t.Errorf("body = %q", m.Body)
	}
}

func TestDispatcher_UnknownType(t *testing.T) {
	d := NewDispatcher()
	_, err := d.Dispatch([]byte(`{"type":"nope","body":""}`))
	if !errors.Is(err, ErrUnknownType) {
		t.Errorf("expected ErrUnknownType, got %v", err)
	}
}

func TestDispatcher_MalformedJSON(t *testing.T) {
	d := NewDispatcher()
	_, err := d.Dispatch([]byte(`{not json`))
	if !errors.Is(err, ErrMalformedJSON) {
		t.Errorf("expected ErrMalformedJSON, got %v", err)
	}
}

func TestDispatcher_OpenForNewHandlersWithoutCoreEdits(t *testing.T) {
	// Register an ad-hoc handler and check it gets routed — proves the
	// dispatcher is open to new handlers without core changes.
	d := NewDispatcher()
	d.Register("shout", handlerFunc(func(m Message) (Message, error) {
		return Message{Type: m.Type, Body: strings.ToUpper(m.Body)}, nil
	}))
	out, err := d.Dispatch([]byte(`{"type":"shout","body":"hi"}`))
	if err != nil {
		t.Fatalf("dispatch: %v", err)
	}
	var m Message
	_ = json.Unmarshal(out, &m)
	if m.Body != "HI" {
		t.Errorf("body = %q", m.Body)
	}
}

type handlerFunc func(Message) (Message, error)

func (f handlerFunc) Handle(m Message) (Message, error) { return f(m) }
