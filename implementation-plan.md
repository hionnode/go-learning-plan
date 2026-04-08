# Implementation Plan — Go + Distributed Systems Curriculum

## Overview

**Learner profile:** Experienced dev (MuleSoft/Java), new to Go, limited DB/infra depth.
**Time commitment:** 2-3 hrs/day, no fixed deadline.
**End goal:** Independently complete all 6 Fly.io Gossip Glomers challenges.
**Method:** Learn by doing with Claude Code as a pairing mentor, not a solution generator.

---

## Phase 1 — Go Fundamentals (~2 weeks)

The goal is muscle memory for Go syntax and idioms. Everything is CLI programs, no servers yet. Use [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests) (LGwT) as a TDD companion throughout — chapters are mapped below.

### 1.1 Setup & Hello World
- [ ] Install Go, configure GOPATH/GOROOT
- [ ] `go mod init`, understand modules
- [ ] Write a program that reads stdin line-by-line and echoes it back (foreshadows Maelstrom's stdin/stdout protocol)
- 📖 **LGwT:** Install Go, Hello World — learn the red/green/refactor TDD cycle from the start

### 1.2 Types, Structs, Interfaces
- [ ] Build a `Message` struct with JSON tags — marshal/unmarshal from a JSON string
- [ ] Define a `Handler` interface with a `Handle(msg Message) (Message, error)` method
- [ ] Implement 2 handlers: `EchoHandler` and `ReverseHandler`
- [ ] Write a dispatcher that routes messages to handlers based on a `type` field
- **Why this early:** This is literally the Maelstrom node pattern. They'll recognize it in Phase 5.
- 📖 **LGwT:** Integers, Iteration, Arrays and slices, Structs/methods/interfaces — build fluency with Go's type system via TDD

### 1.3 Error Handling & Testing
- [ ] Refactor the dispatcher to return proper errors (malformed JSON, unknown type)
- [ ] Write table-driven tests for all handlers
- [ ] Use `testify` only after writing raw `testing.T` tests first
- 📖 **LGwT:** Pointers & errors, Error types, Dependency Injection, Mocking — learn idiomatic error handling and test doubles
- 📖 **LGwT:** Why unit tests and how to make them work for you, Anti-patterns, Refactoring Checklist — read these for testing philosophy

### 1.4 Pointers, Slices, Maps
- [ ] Build an in-memory key-value store (map-based) with Get/Set/Delete
- [ ] Expose it via a CLI REPL — `SET foo bar`, `GET foo`, `DEL foo`
- [ ] Write benchmarks (`testing.B`) comparing map access patterns
- 📖 **LGwT:** Maps, Generics, Revisiting arrays/slices with generics — practice with collections and type parameters

### 1.5 Goroutines & Channels — Intro
- [ ] Build a concurrent word counter: split a file into chunks, count words per chunk in goroutines, combine via channels
- [ ] Demonstrate a race condition intentionally, then fix it with `sync.Mutex`
- [ ] Redo the fix using channels instead of mutex — compare both approaches
- 📖 **LGwT:** Concurrency — goroutines, channels, and anonymous functions via TDD

### 1.6 TDD Practice (LGwT extras)
- [ ] Work through Intro to property based tests (Roman numerals) — `testing/quick`
- [ ] Work through Maths (SVG clock face) — acceptance vs unit test separation
- [ ] Work through Reflection — understand `reflect` package, recursive struct traversal
- **Why a dedicated section:** These chapters don't map to a specific project task above but build important Go + TDD muscle. Do them as side exercises between tasks or as warm-ups.

**Phase 1 checkpoint:** Can write, test, and run a Go program from scratch. Understands interfaces, error handling, basic concurrency. Comfortable with red/green/refactor TDD workflow.

---

## Phase 2 — Stdlib Deep Dive (~2-3 weeks)

Servers, serialization, and concurrency patterns — still no external deps. LGwT's "Build an Application" section is a parallel project you can work through alongside or interleaved with the tasks below.

### 2.1 net/http Server
- [ ] Build a JSON API server with `net/http` only (no router libraries)
- [ ] Endpoints: `POST /kv/set`, `GET /kv/get?key=X`, `DELETE /kv/del?key=X`
- [ ] Proper status codes, JSON error responses, request validation
- [ ] Middleware: request logging, panic recovery
- 📖 **LGwT:** HTTP server, JSON/routing/embedding, Revisiting HTTP Handlers — TDD-driven approach to building the same kind of server

### 2.2 Concurrency Patterns
- [ ] Worker pool: accept jobs via HTTP, process in N goroutines, return results
- [ ] Fan-out/fan-in: query 3 "upstream" URLs concurrently, return fastest response
- [ ] Use `context.Context` for timeouts and cancellation
- [ ] `sync.WaitGroup`, `errgroup` usage
- **Why this matters:** Gossip Glomers broadcast challenge requires fan-out to peer nodes.
- 📖 **LGwT:** Select (race concurrent requests), Sync (`sync.Mutex`, `sync.WaitGroup`), Context (cancellation with `context.Done()` + `select`)

### 2.3 Networking & Serialization
- [ ] Build a TCP echo server using `net` package
- [ ] Implement a simple line-delimited JSON protocol over TCP (stdin/stdout simulation)
- [ ] Two instances talking to each other over TCP — send messages back and forth
- **Why this matters:** This is Maelstrom's communication model, just over a different transport.
- 📖 **LGwT:** Reading files (io/fs, testing/fstest), IO and sorting (io.Reader/Writer, file persistence), Context-aware Reader, OS Exec

### 2.4 Time, Tickers, and Periodic Tasks
- [ ] Build a health-check service: periodically ping a list of URLs, report status
- [ ] Use `time.Ticker` for periodic gossip (will reuse in broadcast challenge)
- [ ] Graceful shutdown with `os/signal` and context cancellation
- 📖 **LGwT:** Time (scheduling with `time.AfterFunc`), Command line & package structure (multi-binary projects, `cmd/` layout)

### 2.5 Testing at Scale (LGwT)
- [ ] Work through Introduction to acceptance tests — black-box testing of HTTP servers, graceful shutdown
- [ ] Work through Scaling acceptance tests — Docker/Testcontainers, specification-driver-system pattern
- [ ] Work through Working without mocks — fakes + contracts over mocks/stubs
- [ ] Optionally: Templating, WebSockets chapters if time allows
- **Why a dedicated section:** These testing chapters teach patterns you'll rely on heavily in Phase 3 (integration tests against real Postgres/Redis) and Phase 5 (Maelstrom verification).

**Phase 2 checkpoint:** Can build a concurrent HTTP service in pure stdlib. Understands context, cancellation, worker pools, fan-out. Has a mature TDD practice including acceptance tests and test architecture.

---

## Phase 3 — Backend Building Blocks (~3-4 weeks)

Databases, caches, queues. External deps allowed now. Each topic is a mini-project.

### 3.1 PostgreSQL with database/sql
- [ ] Set up Postgres locally (Docker)
- [ ] Build a URL shortener: `POST /shorten` stores URL+slug, `GET /:slug` redirects
- [ ] Raw SQL with `database/sql` + `pgx` driver — no ORM
- [ ] Connection pooling, prepared statements, transactions
- [ ] Write integration tests using a test database
- [ ] Add migrations with `golang-migrate`

### 3.2 Redis
- [ ] Add rate limiting to the URL shortener using Redis (`go-redis`)
- [ ] Implement a distributed lock using Redis SETNX
- [ ] Cache frequently accessed URLs — cache-aside pattern with TTL
- [ ] Pub/Sub: build a simple chat where messages go through Redis channels

### 3.3 Message Queues — NATS
- [ ] Set up NATS locally (Docker)
- [ ] Build an async task processor: HTTP endpoint enqueues work, NATS consumer processes it
- [ ] Implement at-least-once delivery with acknowledgments
- [ ] JetStream for persistence and replay
- **Why NATS over Kafka:** Lightweight, Go-native, and the concepts transfer directly. Gossip Glomers Challenge 5 is Kafka-style logs — understanding the concepts matters more than the specific tool.

### 3.4 Putting It Together — Mini System
- [ ] Build a link analytics service: shorten URLs (Postgres), count clicks (Redis), emit click events (NATS), consume events to build analytics (separate service)
- [ ] 3 separate Go binaries talking to each other
- [ ] Graceful shutdown for all services
- [ ] Docker Compose for the full stack

**Phase 3 checkpoint:** Can build multi-service systems with persistent storage, caching, and async messaging. Comfortable with Docker, connection management, and integration testing.

---

## Phase 4 — Distributed Systems Theory + Practice (~2-3 weeks)

Read, understand, then implement. This phase bridges the gap to Gossip Glomers.

### 4.1 Core Concepts (Read + Discuss)
- [ ] Read: [Designing Data-Intensive Applications](https://dataintensive.net/) Chapters 5, 8, 9 (Replication, Faults, Consistency)
- [ ] Read: [Notes on Distributed Systems for Young Bloods](https://www.somethingsimilar.com/2013/01/14/what-is-distributed-systems/)
- [ ] Discuss with Claude Code: CAP theorem, consistency models, failure modes
- [ ] No code yet — just build mental models

### 4.2 Implement a Gossip Protocol (Standalone)
- [ ] Build 5 nodes as goroutines within one process communicating via channels
- [ ] Implement epidemic/gossip protocol: each node periodically picks random peers and shares state
- [ ] Measure convergence time — how many rounds until all nodes agree?
- [ ] Introduce simulated failures (drop messages randomly) and observe behavior
- **Direct prep for:** Gossip Glomers Challenge 3 (Broadcast)

### 4.3 Implement a Logical Clock
- [ ] Lamport timestamps: add to the gossip protocol
- [ ] Vector clocks: extend the implementation
- [ ] Use clocks to order events across nodes
- **Direct prep for:** Understanding ordering in Challenges 4, 5, 6

### 4.4 Implement a CRDT
- [ ] G-Counter (grow-only counter) as a CRDT — each node has its own counter, merge by taking max per node
- [ ] PN-Counter (increment + decrement)
- [ ] Test with the simulated gossip network from 4.2
- **Direct prep for:** Gossip Glomers Challenge 4 (G-Counter)

### 4.5 Maelstrom Setup
- [ ] Install Maelstrom (Java/Clojure tooling)
- [ ] Read the Maelstrom docs and protocol specification
- [ ] Run the demo echo node from the Maelstrom repo
- [ ] Understand the stdin/stdout JSON protocol (they've already built this in Phase 2!)

**Phase 4 checkpoint:** Understands gossip protocols, consistency models, CRDTs, logical clocks. Maelstrom is installed and working. Ready for the real challenges.

---

## Phase 5 — Gossip Glomers (~4-6 weeks)

The main event. Each challenge builds on the previous. Claude Code's role is minimal here — architecture guidance only, no implementation help.

### Challenge 1: Echo
- [ ] Build a Maelstrom node in Go that responds to echo messages
- [ ] Understand the init handshake, node IDs, message IDs
- [ ] **Estimated time:** 1-2 sessions

### Challenge 2: Unique ID Generation
- [ ] Generate globally unique IDs without coordination
- [ ] Consider: UUIDs? Snowflake IDs? Node-ID + counter?
- [ ] Trade-offs: simplicity vs. sortability vs. size
- [ ] **Estimated time:** 1-2 sessions

### Challenge 3: Broadcast (3a → 3b → 3c → 3d → 3e)
- [ ] 3a: Single-node broadcast — just store and return
- [ ] 3b: Multi-node with naive flooding
- [ ] 3c: Fault-tolerant broadcast (handle network partitions)
- [ ] 3d: Efficient broadcast (reduce message count)
- [ ] 3e: Efficient broadcast under constraints (latency + message budget)
- [ ] **This is the hardest challenge.** Take time. Revisit Phase 4.2 notes.
- [ ] **Estimated time:** 1-2 weeks

### Challenge 4: Grow-Only Counter
- [ ] Implement a G-Counter using Maelstrom's seq-kv service
- [ ] Handle eventual consistency — reads may be stale
- [ ] Revisit CRDT implementation from Phase 4.4
- [ ] **Estimated time:** 2-3 sessions

### Challenge 5: Kafka-Style Log (5a → 5b → 5c)
- [ ] 5a: Single-node log
- [ ] 5b: Multi-node with Maelstrom's lin-kv service
- [ ] 5c: Efficient multi-node (reduce KV operations)
- [ ] **Estimated time:** 1-2 weeks

### Challenge 6: Totally-Available Transactions (6a → 6b → 6c)
- [ ] 6a: Single-node transactions
- [ ] 6b: Read-uncommitted isolation
- [ ] 6c: Read-committed isolation
- [ ] This requires understanding isolation levels, write conflicts, and anti-entropy
- [ ] **Estimated time:** 1-2 weeks

**Phase 5 checkpoint:** All 6 challenges passing Maelstrom verification. The learner can explain the trade-offs in their solutions and why they made specific design choices.

---

## Appendix: Resources

### Books
- *The Go Programming Language* — Donovan & Kernighan (Phase 1-2 reference)
- *[Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)* — Chris James (Phase 1-2 TDD companion, worked through alongside each task)
- *Designing Data-Intensive Applications* — Martin Kleppmann (Phase 4 required reading, Ch 5/8/9)
- *Understanding Distributed Systems* — Roberto Vitillo (lighter alternative to DDIA)

### Online
- [Go by Example](https://gobyexample.com/) — quick syntax reference
- [Effective Go](https://go.dev/doc/effective_go) — idiomatic patterns
- [Maelstrom docs](https://github.com/jepsen-io/maelstrom/blob/main/doc/01-getting-ready/index.md) — protocol reference
- [Fly.io dist-sys challenges](https://fly.io/dist-sys/) — the challenges themselves
- [Jon Gjengset's Gossip Glomers stream](https://www.youtube.com/watch?v=gboGyccRVXI) — see someone else work through them (watch AFTER attempting)

### Tools
- `go`, `go test`, `go vet`, `go fmt` — daily drivers
- `dlv` — debugger for concurrency issues
- Docker + Docker Compose — for Postgres, Redis, NATS
- Maelstrom — challenge verification platform
