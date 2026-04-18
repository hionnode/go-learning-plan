# 1.2 — Types, Structs, Interfaces (the Maelstrom node pattern)

You are building the same pattern every Maelstrom node in Phase 5 uses: parse an incoming JSON message, route by `type`, run the handler, return a response.

## Files

- `dispatcher.go` — you edit this. Fill in JSON tags, handler bodies, and `Dispatch`.
- `dispatcher_test.go` — the verify harness. **Do not edit** — `go-learn verify 1.2-types-structs` runs these.

## What to build

1. Add JSON tags to `Message` so `type` and `body` fields round-trip correctly.
2. `EchoHandler.Handle` — return the input unchanged.
3. `ReverseHandler.Handle` — return a new message with `Body` reversed.
4. `Dispatcher.Dispatch([]byte) ([]byte, error)`:
   - Unmarshal input
   - Look up handler by `Type`
   - Call the handler
   - Marshal the response
   - Return typed errors for: malformed JSON, unknown type

## Why it matters

Every Maelstrom workload dispatches messages by `type` — `init`, `echo`, `topology`, `broadcast`, `read`. Get this interface pattern fluent now and Phase 5 Challenge 1 will be an hour of work instead of an afternoon.

## Go-vs-Java note

Go interfaces are satisfied *implicitly*. You do not write `EchoHandler implements Handler`. The compiler checks method sets. If `EchoHandler` has a `Handle(Message) (Message, error)` method, it satisfies the `Handler` interface. If it doesn't, you find out at the assignment site, not the declaration.

## Don't

- reach for generics yet
- use reflection — you don't need it
- over-abstract with middleware / handler chains; that's Phase 2

## Mastery criteria

- `go-learn verify 1.2-types-structs` passes
- You can explain why the JSON tags matter (hint: what happens if you forget them)
- You can add a new handler with zero changes to `Dispatcher`
