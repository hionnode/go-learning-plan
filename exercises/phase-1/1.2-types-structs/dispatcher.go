package dispatcher

import "errors"

// Message is the wire format. TODO: add JSON tags so this round-trips
// to/from JSON like {"type": "...", "body": "..."}.
type Message struct {
	Type string
	Body string
}

// Handler consumes one Message and returns a response or error.
// Any struct that has a Handle method with this signature satisfies
// Handler — implicitly, no "implements" keyword.
type Handler interface {
	Handle(msg Message) (Message, error)
}

// EchoHandler returns the input unchanged.
type EchoHandler struct{}

// TODO: implement Handle on EchoHandler.
func (EchoHandler) Handle(msg Message) (Message, error) {
	return Message{}, errors.New("EchoHandler: not implemented")
}

// ReverseHandler reverses the Body field.
type ReverseHandler struct{}

// TODO: implement Handle on ReverseHandler.
func (ReverseHandler) Handle(msg Message) (Message, error) {
	return Message{}, errors.New("ReverseHandler: not implemented")
}

// Dispatcher routes messages to handlers by Type.
type Dispatcher struct {
	handlers map[string]Handler
}

// NewDispatcher returns an empty Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: map[string]Handler{}}
}

// Register wires a type string to the handler that should receive it.
func (d *Dispatcher) Register(typ string, h Handler) {
	d.handlers[typ] = h
}

// Sentinel errors — use errors.Is / errors.As against these in tests.
var (
	ErrMalformedJSON = errors.New("malformed json")
	ErrUnknownType   = errors.New("unknown message type")
)

// Dispatch takes raw JSON, routes it, and returns the marshaled response.
//
// TODO: implement.
//   1. json.Unmarshal raw into a Message
//   2. look up the handler for msg.Type
//   3. call handler.Handle(msg)
//   4. json.Marshal the response
// Error handling:
//   - invalid JSON: return wrapped ErrMalformedJSON
//   - missing handler: return wrapped ErrUnknownType
//   - handler error: wrap and return
func (d *Dispatcher) Dispatch(raw []byte) ([]byte, error) {
	return nil, errors.New("Dispatch: not implemented")
}
