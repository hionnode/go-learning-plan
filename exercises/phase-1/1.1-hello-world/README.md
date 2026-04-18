# 1.1 — Setup & stdin-echo

Write a program that reads stdin line-by-line and echoes each line to stdout. Handle EOF gracefully. Stdlib only.

## Files

- `echo.go` — you edit this. Implement `Echo(r io.Reader, w io.Writer) error`.
- `echo_test.go` — the verify harness. **Do not edit** — `go-learn verify 1.1-hello-world` runs these.

You can smoke-test manually: write a tiny `cmd/echo/main.go` under this dir that calls `Echo(os.Stdin, os.Stdout)`, then `echo -e "a\nb" | go run ./cmd/echo`. But the test is the source of truth.

## Goal

This exercise looks trivial. It isn't. The Maelstrom protocol in Phase 5 is literally stdin-echo-with-JSON-message-routing, so build fluent muscle memory here for:

- `bufio.Scanner` on `os.Stdin`
- handling EOF cleanly (no panics, no spurious errors)
- writing to `os.Stdout` with appropriate line terminators

## Mastery criteria

- `go-learn verify 1.1-hello-world` passes
- You can reproduce this from scratch, blank editor, in under 5 min
- You can explain the difference between `bufio.NewReader` and `bufio.NewScanner`

## Don't

- use `fmt.Scanln` — it's the wrong tool for line-delimited stdin
- use panic for EOF
- read lines into a fixed-size byte slice — scanner handles this

## When you're stuck

Run `go doc bufio.Scanner`. Read the example. The API is the hint.
