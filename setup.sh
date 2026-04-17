#!/usr/bin/env bash
# setup.sh — verify env, compile tracker, smoke-test. Idempotent.

set -euo pipefail

cd "$(dirname "$0")"

echo "→ checking Go..."
if ! command -v go >/dev/null 2>&1; then
  echo "✗ go not on PATH. install from https://go.dev/dl/ and re-run." >&2
  exit 1
fi

GO_VERSION="$(go env GOVERSION | sed 's/^go//')"      # e.g. 1.23.5
GO_MAJOR="$(echo "$GO_VERSION" | cut -d. -f1)"
GO_MINOR="$(echo "$GO_VERSION" | cut -d. -f2)"
if [ "$GO_MAJOR" -lt 1 ] || { [ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 26 ]; }; then
  echo "✗ need Go ≥ 1.26, found $GO_VERSION. upgrade via https://go.dev/dl/" >&2
  exit 1
fi
echo "  go $GO_VERSION ok"

echo "→ building tracker..."
go build ./cmd/tracker
echo "  ok"

echo "→ smoke-testing internal packages..."
go test ./internal/... >/dev/null
echo "  ok"

cat <<'EOF'

✓ setup complete.

next steps:
  go run ./cmd/tracker serve                      # dashboard on http://localhost:8080
  go run ./cmd/tracker place phase-0              # diagnostic quiz — skip what you know
  go run ./cmd/tracker verify 1.1-hello-world     # run the first exercise

read the learner-facing docs:
  README.md             what this is, quick start
  curriculum-v2.md      the 43-task curriculum
  claude.md             how to pair-mentor with Claude Code
EOF
