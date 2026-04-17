# Curriculum v2 — Go + Distributed Systems (Math-Academy-Way edition)

> **What is this?** A redesign of `implementation-plan.md` applying the Math-Academy-Way principles: knowledge-graph scaffolding, behavior-observable mastery criteria, spaced retrieval, interleaving, deliberate-practice drills, targeted remediation, automaticity targets, and diagnostic placement. The original plan is preserved in `implementation-plan.md`; see `improvements-summary.md` for the diff.
>
> **Who it's for.** A learner who can write JavaScript and/or Python at a working level but has **no prior backend experience** — never built a server, never written raw SQL, never touched a cache or message queue, new to compiled + statically typed languages. The course takes you from there to passing all 6 Fly.io Gossip Glomers challenges. ~26–30 weeks at 2–3 hrs/day.
>
> **How to read it.** Each task carries a machine-readable YAML frontmatter block the tracker parses. The prose underneath is for humans. "Mastery criteria" replaces checkbox-doneness: you progress by passing the test the tracker runs, then re-passing it after a spaced-retrieval interval. TDD is the default workflow throughout — the `verify_test.go` is the source of truth, not your opinion of whether you're done.

---

## How mastery works

Four levels per task, identical across the curriculum:

| Level | Meaning | How you get here |
|---|---|---|
| `unseen` | Not started | default |
| `learning` | Verify test passes at least once | `tracker verify <task-id>` returns green |
| `proficient` | Re-passed after ≥3 days without reopening your old solution | SRS review returns green |
| `automatic` | Associated drills pass within their time target AND one full SRS cycle survived | drill + SRS both green |

SRS schedule defaults to `[3, 7, 21, 60]` days. A missed review bumps mastery down one level — retention is the whole point.

---

## Diagnostic placement

Before each phase, `tracker place phase-N` runs a short quiz/challenge. Passing lets you skip individual tasks within that phase. Don't grind through what you already own.

- `phase-0`: placement = explain HTTP request/response, what a port is, what DNS does, and draw the client/server diagram from memory. Short conceptual check.
- `phase-1`: placement = write a stdin-echo + JSON-roundtrip + table-driven test in under 20 min, stdlib only.
- `phase-2`: placement = write a `net/http` JSON server with context cancellation in under 30 min.
- `phase-3`: placement = containerized Postgres + Go service with one endpoint doing a transactional write in under 45 min.
- `phase-4`: placement = concept-check MCQ (CAP, consistency models, vector clocks, CRDT convergence).
- `phase-5`: no placement — every challenge is run and graded by Maelstrom.

---

## Drill library

Drills are short, timed, reproducible exercises for automaticity — not project work. You'll run them repeatedly. Each has a target time; beating the target counts toward `automatic` mastery on the parent task.

```yaml
drills:
  - id: zero-values-quickfire
    target_seconds: 180
    prompt: "For each Go type, write the zero value from memory: int, string, bool, []int, map[string]int, *T, struct{X int}, chan int. No references."
  - id: err-idiom
    target_seconds: 180
    prompt: "Write a function Parse(s string) (int, error) that returns a wrapped error on bad input. Write the caller with the idiomatic if err != nil block. From scratch."
  - id: http-status-match
    target_seconds: 180
    prompt: "Match each HTTP status code to its meaning: 200, 201, 204, 301, 400, 401, 403, 404, 409, 422, 500, 502, 503. Which ones are client vs server errors?"
  - id: go-fmt-vet-cycle
    target_seconds: 600
    prompt: "New module, failing test, go vet clean, go fmt clean, make it pass. Whole red→green cycle."
  - id: stdin-echo
    target_seconds: 300
    prompt: "Go program that reads stdin line-by-line, prints each to stdout. Handle EOF. Stdlib only."
  - id: struct-json-roundtrip
    target_seconds: 180
    prompt: "Struct with JSON tags, marshal to string, unmarshal back, assert equality."
  - id: table-driven-test
    target_seconds: 300
    prompt: "Add a table-driven test for a function Sum(a, b int) int using t.Run subtests."
  - id: goroutine-channel
    target_seconds: 600
    prompt: "Spawn N goroutines doing work, aggregate results via channel, no leaks."
  - id: mutex-vs-channel
    target_seconds: 900
    prompt: "Given a race-y counter, fix it with sync.Mutex. Then refactor to channels. Explain which you'd ship and why."
  - id: http-handler
    target_seconds: 600
    prompt: "net/http handler returning JSON + proper status codes for GET/POST. No routers."
  - id: context-timeout
    target_seconds: 300
    prompt: "Wrap an outbound HTTP call in context.WithTimeout. Handle cancellation cleanly."
  - id: errgroup-fanout
    target_seconds: 600
    prompt: "errgroup fan-out to 3 URLs. Return first success, cancel the rest."
  - id: tcp-echo
    target_seconds: 600
    prompt: "TCP echo server using net.Listen. Handle concurrent connections."
  - id: ticker-shutdown
    target_seconds: 600
    prompt: "time.Ticker-driven periodic task with graceful shutdown on SIGINT."
  - id: sql-crud-raw
    target_seconds: 600
    prompt: "In psql (no Go): CREATE TABLE, INSERT 3 rows, UPDATE one, DELETE one, SELECT survivors with a WHERE clause."
  - id: sql-crud-txn
    target_seconds: 900
    prompt: "database/sql CRUD with a transaction. Roll back on error. Use prepared statements."
  - id: redis-cli-basics
    target_seconds: 300
    prompt: "In redis-cli (no Go): SET, GET, SET ... EX, INCR, DEL, EXPIRE. Explain what TTL returns for a key that doesn't exist."
  - id: redis-setnx-lock
    target_seconds: 600
    prompt: "Distributed lock with SETNX + TTL. Handle the TOCTOU pitfall on release."
  - id: nats-pubsub
    target_seconds: 300
    prompt: "With the nats CLI: subscribe to a subject in one shell, publish two messages from another, both arrive."
  - id: maelstrom-init
    target_seconds: 600
    prompt: "Handle Maelstrom init message. Respond with init_ok carrying node_id."
  - id: lamport-clock
    target_seconds: 600
    prompt: "Lamport clock: Send/Recv/Tick. Unit-tested."
  - id: g-counter-merge
    target_seconds: 300
    prompt: "G-Counter struct with Inc/Value/Merge. Merge is idempotent and commutative."
```

---

## Phase 0 — Onramp (~2 weeks)

Backend mental model + the Go essentials you can't fake your way past if you're coming from JS/Python. This phase is half concept, half muscle memory. No servers, no databases, no queues yet — just the language and the vocabulary.

```yaml
---
id: 0.1-backend-mental-model
phase: 0
title: What is a backend? (concepts only)
prereqs: []
interleave_with: [0.2-go-toolchain]
mastery_criteria: |
  - Draw the client/server request/response diagram from memory
  - Explain in two sentences each: HTTP, TCP, DNS, port, JSON, stateless, database, cache, queue
  - Given a concrete scenario ("user clicks 'buy'"), trace the request end-to-end through client → HTTP → server → DB, naming each hop
drill_ids: [http-status-match]
review_intervals_days: [3, 7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

No Go code here. Read, watch, and draw. A "backend" in plain English is the machine (or fleet of machines) that answers requests from apps and browsers, stores state that outlives any one user, and coordinates work that shouldn't happen on the user's device. The rest of this curriculum is just "here's how to build one, one primitive at a time." If this feels trivial, `tracker place phase-0` and move on.

```yaml
---
id: 0.2-go-toolchain
phase: 0
title: Go toolchain — the red/green TDD loop
prereqs: []
interleave_with: [0.1-backend-mental-model]
mastery_criteria: |
  - Install Go; `go env` and `go version` work
  - `go mod init` a fresh module, write a failing test, make it pass
  - `go vet` and `go fmt` clean a file you wrote on purpose poorly
  - Articulate why Go has no REPL and why TDD is the natural substitute
drill_ids: [go-fmt-vet-cycle]
review_intervals_days: [3, 7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

Coming from JS/Python you'll miss the REPL. You won't get one. You'll build the habit of writing a test, running it, reading the red, fixing the red, re-running — that's the iteration loop for the next six months.

```yaml
---
id: 0.3-types-zero-values
phase: 0
title: Static types + zero values
prereqs: [0.2-go-toolchain]
interleave_with: [0.4-errors-as-values]
mastery_criteria: |
  - Write the zero value for: int, string, bool, []int, map[string]int, *T, struct, chan — from memory
  - Explain when Go "unitialized" differs from JS undefined / Python None
  - Use var, :=, and const correctly in a short program that passes go vet
drill_ids: [zero-values-quickfire]
review_intervals_days: [3, 7, 21, 60]
remediation: [0.2-go-toolchain]
diagnostic_checkpoint: false
---
```

Every type has a defined zero value. There is no `undefined`. A `map` you haven't made yet is `nil`, and writing to a nil map panics — familiarize yourself with that panic *on purpose*, once.

```yaml
---
id: 0.4-errors-as-values
phase: 0
title: Errors as values, not exceptions
prereqs: [0.3-types-zero-values]
interleave_with: [0.3-types-zero-values]
mastery_criteria: |
  - Write a function that returns (T, error); handle every call with if err != nil
  - Use fmt.Errorf with %w to wrap an underlying error
  - Use errors.Is / errors.As to match a sentinel error in a test
  - Articulate why there's no try/catch in Go and what changes about how you think about failure
drill_ids: [err-idiom]
review_intervals_days: [3, 7, 21, 60]
remediation: [0.3-types-zero-values]
diagnostic_checkpoint: false
---
```

Errors are returned, not thrown. The caller decides. This is the single biggest cognitive shift from JS/Python — it becomes reflex once you've typed `if err != nil { return fmt.Errorf("doing X: %w", err) }` a hundred times.

```yaml
---
id: 0.5-first-http-hit
phase: 0
title: Make your first HTTP request in Go
prereqs: [0.4-errors-as-values, 0.1-backend-mental-model]
interleave_with: []
mastery_criteria: |
  - Use net/http to GET https://httpbin.org/json and print the response body
  - Inspect the status code, response headers, and body length
  - Handle a timeout with context.WithTimeout (5s)
  - Run `curl -v` against the same URL; compare what you see in your Go code to what curl prints
drill_ids: [context-timeout]
review_intervals_days: [3, 7, 21, 60]
remediation: [0.4-errors-as-values, 0.1-backend-mental-model]
diagnostic_checkpoint: true
---
```

First brush with HTTP as a *protocol*, not a library function. You're the client; something out there is the server. By the end of this task you should not be mystified by what curl prints.

**Phase 0 exit:** `tracker place phase-1` must pass before entering Phase 1 proper.

---

## Phase 1 — Go Fundamentals (~3 weeks)

The goal is muscle memory for Go syntax and idioms. Everything is CLI programs, no servers yet. Use [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests) (LGwT) as a TDD companion throughout — chapters are mapped per task.

```yaml
---
id: 1.1-hello-world
phase: 1
title: Setup & stdin-echo
prereqs: [0.2-go-toolchain]
interleave_with: []
mastery_criteria: |
  - go env, go mod init, go run, go test, go vet, go fmt used fluently
  - Can write a stdin→stdout line-echo program from scratch in <5 min without reference
  - Explain when bufio.Scanner is wrong (long lines) and how to fix it
drill_ids: [stdin-echo, go-fmt-vet-cycle]
review_intervals_days: [3, 7, 21, 60]
remediation: [0.2-go-toolchain]
diagnostic_checkpoint: false
lgwt_chapters: [install-go, hello-world]
---
```

Foreshadows Maelstrom: every node you'll write in Phase 5 is a stdin-echo with more logic.

```yaml
---
id: 1.2-types-structs
phase: 1
title: Types, Structs, Interfaces — the Maelstrom node pattern
prereqs: [1.1-hello-world, 0.3-types-zero-values]
interleave_with: [1.3-error-testing]
mastery_criteria: |
  - Build a Message struct with JSON tags, marshal/unmarshal roundtrip clean
  - Define a Handler interface {Handle(Message) (Message, error)}
  - Implement EchoHandler and ReverseHandler
  - Dispatcher routes by "type" field — adding a new handler requires zero dispatcher edits
drill_ids: [struct-json-roundtrip]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.1-hello-world]
diagnostic_checkpoint: false
lgwt_chapters: [integers, iteration, arrays-slices, structs-methods-interfaces]
---
```

This is literally the Maelstrom node pattern — recognize it in Phase 5. Implicit interface satisfaction (no `implements` keyword) is the single Go-vs-JS-class-vs-Java-class concept to internalize here.

```yaml
---
id: 1.3-error-testing
phase: 1
title: Error handling + table-driven tests
prereqs: [1.2-types-structs, 0.4-errors-as-values]
interleave_with: [1.4-pointers-slices-maps]
mastery_criteria: |
  - Dispatcher returns typed errors (malformed JSON, unknown type); errors wrap with %w
  - Table-driven tests cover every handler + every error path
  - Writes raw testing.T tests BEFORE reaching for testify
drill_ids: [table-driven-test]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.2-types-structs, 0.4-errors-as-values]
diagnostic_checkpoint: false
lgwt_chapters: [pointers-errors, error-types, dependency-injection, mocking, why-unit-tests, anti-patterns]
---
```

```yaml
---
id: 1.4-pointers-slices-maps
phase: 1
title: Pointers, slices, maps — KV REPL
prereqs: [1.1-hello-world, 0.3-types-zero-values]
interleave_with: [1.3-error-testing]
mastery_criteria: |
  - In-memory KV store (map) with Get/Set/Delete
  - CLI REPL: SET/GET/DEL commands over stdin
  - Two benchmarks comparing map access patterns
  - Articulate: when you'd pass a *T vs T, and why "pass the pointer if it's big" isn't the only reason
drill_ids: []
review_intervals_days: [3, 7, 21, 60]
remediation: [0.3-types-zero-values]
diagnostic_checkpoint: false
lgwt_chapters: [maps, generics, arrays-slices-with-generics]
---
```

Pointers hide in JS/Python — here they're explicit. Spend the day. This is also where you first feel the *cost* of Go's explicitness pay off.

```yaml
---
id: 1.5-goroutines-channels
phase: 1
title: Goroutines, channels, the intentional race
prereqs: [1.4-pointers-slices-maps]
interleave_with: [1.6-tdd-practice]
mastery_criteria: |
  - Concurrent word counter: chunk → goroutine → channel aggregate
  - Demonstrate a data race with `go test -race`; fix with sync.Mutex
  - Refactor the same fix with channels; articulate which you'd ship and why
  - Contrast: what's different about goroutines vs JavaScript async/await
drill_ids: [goroutine-channel, mutex-vs-channel]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.4-pointers-slices-maps]
diagnostic_checkpoint: false
lgwt_chapters: [concurrency]
---
```

```yaml
---
id: 1.6-tdd-practice
phase: 1
title: TDD muscle — property tests, reflection, SVG clock
prereqs: [1.3-error-testing]
interleave_with: [1.4-pointers-slices-maps, 1.5-goroutines-channels]
mastery_criteria: |
  - Property-based tests with testing/quick (Roman numerals)
  - Acceptance vs unit test separation (SVG clock)
  - reflect.Value walk over a nested struct
drill_ids: [table-driven-test]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.3-error-testing]
diagnostic_checkpoint: false
lgwt_chapters: [property-based-tests, maths-svg-clock, reflection]
---
```

**Phase 1 exit:** `tracker place phase-2` must pass.

---

## Phase 2 — HTTP & the Stdlib (~3 weeks)

Now the protocols underneath the web — seen clearly, built from stdlib, tested before they're trusted.

```yaml
---
id: 2.1-http-protocol
phase: 2
title: HTTP as a protocol (what curl actually sees)
prereqs: [0.5-first-http-hit, 1.2-types-structs]
interleave_with: [2.2-net-http-server]
mastery_criteria: |
  - Read and write a raw HTTP/1.1 request + response by hand (no curl)
  - Use curl -v to see request-line, headers, body, and status clearly
  - Explain: what is a status code, what are headers, what's a content-type, what does Connection: close mean
  - Speak TCP at the right level: HTTP runs over it; ports are TCP's addressing
drill_ids: [http-status-match]
review_intervals_days: [3, 7, 21, 60]
remediation: [0.5-first-http-hit]
diagnostic_checkpoint: false
---
```

Concept-first task. Before you write a server, you need to know what the wire looks like. Open a raw TCP connection to `httpbin.org:80` and type out a real HTTP request. The mystery dies here.

```yaml
---
id: 2.2-net-http-server
phase: 2
title: net/http JSON API — no routers
prereqs: [2.1-http-protocol, 1.3-error-testing, 1.4-pointers-slices-maps]
interleave_with: [2.3-concurrency-patterns]
mastery_criteria: |
  - POST /kv/set, GET /kv/get?key=X, DELETE /kv/del?key=X all working; status codes correct
  - Middleware for request logging + panic recovery composed via http.Handler chain
  - No third-party router — stdlib only. Explain why you'd eventually reach for one and when not to
drill_ids: [http-handler]
review_intervals_days: [3, 7, 21, 60]
remediation: [2.1-http-protocol, 1.4-pointers-slices-maps]
diagnostic_checkpoint: false
lgwt_chapters: [http-server, json-routing-embedding, revisiting-http-handlers]
---
```

```yaml
---
id: 2.3-concurrency-patterns
phase: 2
title: Worker pool + fan-out/fan-in + context
prereqs: [1.5-goroutines-channels, 2.2-net-http-server]
interleave_with: [2.4-networking-serialization]
mastery_criteria: |
  - HTTP-backed worker pool: N goroutines, bounded queue, clean shutdown
  - Fan-out to 3 "upstream" URLs; return the fastest; cancel the stragglers via context
  - sync.WaitGroup vs errgroup chosen deliberately, not reflexively
drill_ids: [context-timeout, errgroup-fanout]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.5-goroutines-channels]
diagnostic_checkpoint: false
lgwt_chapters: [select, sync, context]
---
```

**Why this matters for Phase 5:** Broadcast (Challenge 3) is fan-out to peers.

```yaml
---
id: 2.4-networking-serialization
phase: 2
title: TCP + line-delimited JSON — the Maelstrom twin
prereqs: [1.2-types-structs, 2.2-net-http-server]
interleave_with: [2.5-time-tickers]
mastery_criteria: |
  - TCP echo server; two processes exchanging line-delimited JSON over TCP
  - bufio.Scanner tuned for long lines (understand the default 64KB limit)
  - Explain what Phase 5 changes — same protocol, different transport (stdio)
drill_ids: [tcp-echo]
review_intervals_days: [3, 7, 21, 60]
remediation: [1.2-types-structs]
diagnostic_checkpoint: false
lgwt_chapters: [reading-files, io-sorting, context-aware-reader, os-exec]
---
```

```yaml
---
id: 2.5-time-tickers
phase: 2
title: time.Ticker + graceful shutdown
prereqs: [2.3-concurrency-patterns]
interleave_with: []
mastery_criteria: |
  - Health-check service: ticker-driven, periodic URL pings, per-target status
  - Graceful shutdown with os/signal + ctx cancellation; no goroutine leaks under -race
drill_ids: [ticker-shutdown]
review_intervals_days: [3, 7, 21, 60]
remediation: [2.3-concurrency-patterns]
diagnostic_checkpoint: false
lgwt_chapters: [time, command-line-package-structure]
---
```

```yaml
---
id: 2.6-testing-at-scale
phase: 2
title: Acceptance tests + fakes over mocks
prereqs: [1.3-error-testing, 2.2-net-http-server]
interleave_with: [2.5-time-tickers]
mastery_criteria: |
  - Black-box acceptance tests exercising the HTTP server via net/http.Client
  - Fake-with-contract test-double for one external dependency
  - Articulate when you would choose Testcontainers over a fake
drill_ids: []
review_intervals_days: [3, 7, 21, 60]
remediation: [1.3-error-testing]
diagnostic_checkpoint: true
lgwt_chapters: [acceptance-tests, scaling-acceptance-tests, working-without-mocks]
---
```

**Phase 2 exit:** diagnostic placement for Phase 3.

---

## Phase 3 — Backend Building Blocks (~5 weeks)

Databases, caches, queues. Each tech gets a **concepts-first** task (no Go) before the integration task, so you understand the primitive before wrapping it in a client library.

```yaml
---
id: 3.1-sql-fundamentals
phase: 3
title: SQL fundamentals — no Go
prereqs: [0.1-backend-mental-model]
interleave_with: [3.2-postgres]
mastery_criteria: |
  - In psql (or sqlite3): CREATE TABLE with types + constraints, INSERT, UPDATE, DELETE, SELECT with WHERE
  - INNER JOIN two tables with a COUNT + GROUP BY
  - Explain: what is an index, why adding one helps SELECT and hurts INSERT
  - Explain transactions with ACID in plain language; demonstrate ROLLBACK
drill_ids: [sql-crud-raw]
review_intervals_days: [7, 21, 60]
remediation: [0.1-backend-mental-model]
diagnostic_checkpoint: false
---
```

If you've never written raw SQL, spend a week in psql before touching Go. Every backend engineer's debugging routine eventually reduces to "log in to the DB, run a query, see the truth."

```yaml
---
id: 3.2-postgres
phase: 3
title: Postgres via database/sql + pgx
prereqs: [3.1-sql-fundamentals, 2.2-net-http-server, 2.6-testing-at-scale]
interleave_with: [3.3-caching-fundamentals]
mastery_criteria: |
  - URL shortener: POST /shorten, GET /:slug with raw SQL, no ORM
  - Prepared statements, pooling, BEGIN/COMMIT/ROLLBACK all used intentionally
  - Integration test against a real containerized Postgres passes repeatedly
  - golang-migrate up/down migrations applied in a CI-style script
drill_ids: [sql-crud-txn]
review_intervals_days: [7, 21, 60]
remediation: [3.1-sql-fundamentals, 2.6-testing-at-scale]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 3.3-caching-fundamentals
phase: 3
title: Caching fundamentals — redis-cli, no Go
prereqs: [0.1-backend-mental-model]
interleave_with: [3.4-redis]
mastery_criteria: |
  - In redis-cli: SET/GET, SET with EX, INCR, DEL, EXPIRE, TTL
  - Explain: what problem does a cache solve; when it's the wrong tool
  - Describe cache-aside, TTL-based eviction, and the thundering herd in plain English
drill_ids: [redis-cli-basics]
review_intervals_days: [7, 21, 60]
remediation: [0.1-backend-mental-model]
diagnostic_checkpoint: false
---
```

Before you wrap Redis in a Go client, use it with your fingers. Feel how fast key-value reads are. Watch a key expire.

```yaml
---
id: 3.4-redis
phase: 3
title: Redis from Go — rate limit, lock, cache-aside, pub/sub
prereqs: [3.3-caching-fundamentals, 3.2-postgres]
interleave_with: [3.5-queue-fundamentals]
mastery_criteria: |
  - Rate-limit the shortener via Redis INCR+TTL
  - SETNX-based distributed lock (TOCTOU correctness on release accounted for)
  - Cache-aside for slug lookups with explicit TTL strategy
  - Pub/Sub chat between two processes working end-to-end
drill_ids: [redis-setnx-lock]
review_intervals_days: [7, 21, 60]
remediation: [3.3-caching-fundamentals, 3.2-postgres]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 3.5-queue-fundamentals
phase: 3
title: Message queue fundamentals — nats CLI, no Go
prereqs: [0.1-backend-mental-model]
interleave_with: [3.6-nats]
mastery_criteria: |
  - Install and run a local nats-server; use the nats CLI to subscribe and publish on a subject
  - Explain: pub/sub vs work queue; when you'd want each
  - Define and give one example each: at-most-once, at-least-once, exactly-once (and why exactly-once is a lie)
  - Explain backpressure: what happens when a consumer is slower than producers
drill_ids: [nats-pubsub]
review_intervals_days: [7, 21, 60]
remediation: [0.1-backend-mental-model]
diagnostic_checkpoint: false
---
```

Before you write a Go consumer, you need to know what a queue *is*. Producers push messages onto subjects; consumers pull them off. Everything else is engineering around failure.

```yaml
---
id: 3.6-nats
phase: 3
title: NATS + JetStream from Go
prereqs: [3.5-queue-fundamentals, 2.3-concurrency-patterns]
interleave_with: []
mastery_criteria: |
  - HTTP enqueue → NATS → consumer pipeline; explicit ack/nack handling
  - JetStream stream consumed twice with at-least-once semantics
  - Articulate why NATS subjects map to Kafka topics (prepares Challenge 5)
drill_ids: []
review_intervals_days: [7, 21, 60]
remediation: [3.5-queue-fundamentals, 2.3-concurrency-patterns]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 3.7-mini-system
phase: 3
title: Three-binary link-analytics system
prereqs: [3.2-postgres, 3.4-redis, 3.6-nats]
interleave_with: []
mastery_criteria: |
  - Three Go binaries: shortener, counter, analytics-consumer
  - docker-compose up brings the whole stack up cleanly
  - Graceful shutdown across all three on SIGINT — no half-processed messages
drill_ids: []
review_intervals_days: [21, 60]
remediation: [3.2-postgres, 3.4-redis, 3.6-nats]
diagnostic_checkpoint: true
---
```

---

## Phase 4 — Distributed Systems Theory + Practice (~3 weeks)

```yaml
---
id: 4.1-distsys-theory
phase: 4
title: Read + discuss — DDIA Ch 5/8/9, Young Bloods
prereqs: [3.7-mini-system]
interleave_with: [4.2-gossip-protocol]
mastery_criteria: |
  - Can explain CAP theorem limits (not the myth) in one paragraph
  - Distinguish linearizability, serializability, causal, eventual consistency
  - Name three concrete failure modes distributed systems must handle
drill_ids: []
review_intervals_days: [7, 21, 60]
remediation: []
diagnostic_checkpoint: false
---
```

```yaml
---
id: 4.2-gossip-protocol
phase: 4
title: Implement a gossip protocol standalone
prereqs: [2.3-concurrency-patterns, 4.1-distsys-theory]
interleave_with: [4.3-logical-clocks]
mastery_criteria: |
  - 5-node epidemic protocol inside one process (goroutines + channels)
  - Convergence measurable: prints rounds-to-convergence
  - Simulated 20% message loss still converges; 80% loss explains why it doesn't
drill_ids: []
review_intervals_days: [7, 21, 60]
remediation: [2.3-concurrency-patterns]
diagnostic_checkpoint: false
---
```

Direct prep for Challenge 3.

```yaml
---
id: 4.3-logical-clocks
phase: 4
title: Lamport + vector clocks
prereqs: [4.2-gossip-protocol]
interleave_with: [4.4-crdt]
mastery_criteria: |
  - Lamport clock integrated into the gossip network: events have total order
  - Vector clock correctly identifies concurrent vs causally-ordered events
  - Can hand-trace an example on a whiteboard without looking
drill_ids: [lamport-clock]
review_intervals_days: [7, 21, 60]
remediation: [4.2-gossip-protocol]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 4.4-crdt
phase: 4
title: G-Counter + PN-Counter CRDTs
prereqs: [4.2-gossip-protocol]
interleave_with: [4.3-logical-clocks]
mastery_criteria: |
  - G-Counter merge is idempotent, commutative, associative (property-tested)
  - PN-Counter built on two G-Counters; negative values impossible
  - Integrated into gossip network — final value converges across nodes
drill_ids: [g-counter-merge]
review_intervals_days: [7, 21, 60]
remediation: [4.2-gossip-protocol]
diagnostic_checkpoint: false
---
```

Direct prep for Challenge 4.

```yaml
---
id: 4.5-maelstrom-setup
phase: 4
title: Maelstrom installed + demo echo passes
prereqs: [2.4-networking-serialization, 4.1-distsys-theory]
interleave_with: []
mastery_criteria: |
  - JDK + Graphviz + Gnuplot installed; `maelstrom test -w echo ...` demo passes
  - Read the workloads.md doc end-to-end; can explain init-handshake from memory
drill_ids: [maelstrom-init]
review_intervals_days: [7, 21, 60]
remediation: [2.4-networking-serialization]
diagnostic_checkpoint: true
---
```

---

## Phase 5 — Gossip Glomers (~4–6 weeks)

Maelstrom is the authority; the tracker's `verify` shells out to a Maelstrom run with the correct workload. Claude Code's role here is architecture-only — zero implementation help. If you reach for a solution, take the remediation pointer instead.

```yaml
---
id: 5.1-echo
phase: 5
title: Challenge 1 — Echo
prereqs: [4.5-maelstrom-setup, 1.2-types-structs]
interleave_with: []
mastery_criteria: |
  - Maelstrom echo workload passes at n=1 and n=3
  - init, echo messages handled; msg_id + in_reply_to set correctly
drill_ids: [maelstrom-init]
review_intervals_days: [21, 60]
remediation: [4.5-maelstrom-setup, 1.2-types-structs]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.2-unique-id
phase: 5
title: Challenge 2 — Unique ID generation
prereqs: [5.1-echo]
interleave_with: []
mastery_criteria: |
  - Maelstrom unique-ids workload passes with partitions enabled
  - Can articulate tradeoffs among UUIDv4 / Snowflake / node_id+counter
  - Chosen scheme defended in writing (reflection log entry required)
drill_ids: []
review_intervals_days: [21, 60]
remediation: [5.1-echo]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.3a-broadcast-single
phase: 5
title: Challenge 3a — Single-node broadcast
prereqs: [5.1-echo]
interleave_with: []
mastery_criteria: "broadcast workload passes at n=1"
drill_ids: []
review_intervals_days: [21]
remediation: [5.1-echo]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.3b-broadcast-flooding
phase: 5
title: Challenge 3b — Naive multi-node broadcast
prereqs: [5.3a-broadcast-single, 4.2-gossip-protocol]
interleave_with: []
mastery_criteria: "broadcast workload passes at n=5, no partitions"
drill_ids: []
review_intervals_days: [21]
remediation: [4.2-gossip-protocol]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.3c-broadcast-fault-tolerant
phase: 5
title: Challenge 3c — Fault-tolerant broadcast
prereqs: [5.3b-broadcast-flooding]
interleave_with: []
mastery_criteria: "broadcast workload passes at n=5 with --nemesis partition"
drill_ids: []
review_intervals_days: [21, 60]
remediation: [5.3b-broadcast-flooding, 4.2-gossip-protocol]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.3d-broadcast-efficient
phase: 5
title: Challenge 3d — Efficient broadcast
prereqs: [5.3c-broadcast-fault-tolerant]
interleave_with: []
mastery_criteria: "messages-per-op <= 30, stable latencies under partition"
drill_ids: []
review_intervals_days: [60]
remediation: [5.3c-broadcast-fault-tolerant]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.3e-broadcast-constrained
phase: 5
title: Challenge 3e — Latency + budget constrained broadcast
prereqs: [5.3d-broadcast-efficient]
interleave_with: []
mastery_criteria: "meets the 3e latency SLO and message budget"
drill_ids: []
review_intervals_days: [60]
remediation: [5.3d-broadcast-efficient]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.4-g-counter
phase: 5
title: Challenge 4 — G-Counter with seq-kv
prereqs: [4.4-crdt, 5.3c-broadcast-fault-tolerant]
interleave_with: []
mastery_criteria: "g-counter workload passes at n=3 under partition"
drill_ids: [g-counter-merge]
review_intervals_days: [21, 60]
remediation: [4.4-crdt]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.5a-kafka-single
phase: 5
title: Challenge 5a — Single-node log
prereqs: [5.1-echo]
interleave_with: []
mastery_criteria: "kafka workload passes at n=1"
drill_ids: []
review_intervals_days: [21]
remediation: [5.1-echo]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.5b-kafka-lin-kv
phase: 5
title: Challenge 5b — Multi-node log via lin-kv
prereqs: [5.5a-kafka-single]
interleave_with: []
mastery_criteria: "kafka workload passes at n=3 with lin-kv backing store"
drill_ids: []
review_intervals_days: [21, 60]
remediation: [5.5a-kafka-single]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.5c-kafka-efficient
phase: 5
title: Challenge 5c — Efficient multi-node log
prereqs: [5.5b-kafka-lin-kv]
interleave_with: []
mastery_criteria: "kafka workload passes at n=3 with reduced KV op count"
drill_ids: []
review_intervals_days: [60]
remediation: [5.5b-kafka-lin-kv]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.6a-txn-single
phase: 5
title: Challenge 6a — Single-node transactions
prereqs: [5.1-echo]
interleave_with: []
mastery_criteria: "txn-rw-register passes at n=1"
drill_ids: []
review_intervals_days: [21]
remediation: [5.1-echo]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.6b-txn-read-uncommitted
phase: 5
title: Challenge 6b — Read-uncommitted transactions
prereqs: [5.6a-txn-single]
interleave_with: []
mastery_criteria: "txn-rw-register passes at n=2 with read-uncommitted isolation"
drill_ids: []
review_intervals_days: [21, 60]
remediation: [5.6a-txn-single]
diagnostic_checkpoint: false
---
```

```yaml
---
id: 5.6c-txn-read-committed
phase: 5
title: Challenge 6c — Read-committed transactions
prereqs: [5.6b-txn-read-uncommitted]
interleave_with: []
mastery_criteria: "txn-rw-register passes at n=2 with read-committed isolation"
drill_ids: []
review_intervals_days: [60]
remediation: [5.6b-txn-read-uncommitted]
diagnostic_checkpoint: false
---
```

---

## Spiral review schedule

Spaced retrieval is the engine that moves a task from `learning` → `proficient` → `automatic`. The tracker re-queues review prompts at each interval in `review_intervals_days`. A missed review demotes mastery by one level — ignore the queue and you will feel it.

Every review is a *retrieval* (write from scratch, no peeking), not a re-read.

---

## Appendix: what changed vs. implementation-plan.md

See `improvements-summary.md`.
