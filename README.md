# learning-plan

A ~26-week self-paced curriculum to take you from **JS/Python working knowledge + zero backend experience** to independently solving all six [Fly.io Gossip Glomers](https://fly.io/dist-sys/) distributed-systems challenges in Go.

It ships with a local tracker (Go, stdlib-only) that enforces **Math-Academy-Way** learning: mastery levels instead of checkboxes, spaced retrieval, prerequisite DAG, timed drills for automaticity, phase-entry placement quizzes.

## Who this is for

You can write JavaScript and/or Python at a working level. You have **never**:

- built an HTTP server from scratch
- written raw SQL
- used Redis or a message queue
- worked in a compiled, statically-typed language

You want to build real backends and understand distributed systems. You're willing to commit 2–3 hrs/day for ~6 months.

## Quick start

```sh
./setup.sh                                      # check go, compile tracker, smoke-test
go run ./cmd/tracker serve                      # dashboard on http://localhost:8080
go run ./cmd/tracker verify 1.1-hello-world     # run the first exercise's tests
```

Open `exercises/phase-1/1.1-hello-world/` — fill in `echo.go`, re-run the verify. When it goes green, your mastery on `1.1-hello-world` becomes `learning` and a review is scheduled 3 days out.

## How the tracker works

- **Four mastery levels** per task: `unseen → learning → proficient → automatic`. You advance only when the task's `verify_test.go` passes — no self-report.
- **Spaced retrieval** at `[3, 7, 21, 60]` days. Missing a review demotes mastery — retention is the whole point.
- **Prereq DAG.** Every task declares its prerequisites; the tracker refuses to surface a task whose prereqs aren't at least `learning`. Failed verifies route you to the remediation target.
- **Drills** are short, timed deliberate-practice exercises (e.g. "write a net/http handler in <10 min"). Beating target time promotes the parent task to `automatic`.
- **Diagnostic placement** at phase entry lets you skip what you already own: `tracker place phase-0`.

## Curriculum at a glance

| Phase | Focus | Weeks | Tasks |
|---|---|---|---|
| 0 | Onramp: backend mental model + Go toolchain | ~2 | 5 |
| 1 | Go fundamentals (stdin, structs, errors, goroutines) | ~3 | 6 |
| 2 | HTTP + stdlib (protocol, `net/http`, concurrency, TCP) | ~3 | 6 |
| 3 | Backend building blocks (SQL, Redis, NATS) | ~5 | 7 |
| 4 | Distributed-systems theory + practice (gossip, CRDTs, Maelstrom) | ~3 | 5 |
| 5 | Gossip Glomers (6 challenges / 14 sub-tasks) | ~4–6 | 14 |
| | **Total** | **~26–30** | **43** |

Each technology in Phase 3 gets a **concepts-first** task (raw psql, raw redis-cli, raw nats CLI) before the Go integration. You learn what a cache *is* before you wrap Redis in a client.

## Repo layout

```
.
├── README.md                 # you are here
├── setup.sh                  # one-shot env check + first build
├── claude.md                 # teaching rules for Claude Code (mentor mode)
├── curriculum-v2.md          # the 43-task curriculum with YAML frontmatter
├── implementation-plan.md    # v1 (frozen, historical)
├── improvements-summary.md   # why v2 changed vs v1
├── go.mod
├── cmd/tracker/              # the tracker binary (serve/verify/drill/review/place)
├── internal/                 # parser, progress store, SRS, DAG, drill runner
└── exercises/
    └── phase-1/
        ├── 1.1-hello-world/        # seeded scaffold + verify_test.go
        └── 1.2-types-structs/      # seeded scaffold + verify_test.go
```

More exercise scaffolds are added lazily as you reach them.

## Tracker CLI

```sh
go run ./cmd/tracker serve                # dashboard, graph, review queue, drills
go run ./cmd/tracker verify <task-id>     # run that task's tests, update mastery
go run ./cmd/tracker drill  <drill-id>    # timed deliberate-practice exercise
go run ./cmd/tracker review               # what's due for retrieval today
go run ./cmd/tracker place  <phase-id>    # diagnostic quiz (phase-0 … phase-4)
```

Progress lives in `progress.json` (gitignored — your state, not repo state).

## Where to read next

- **`claude.md`** — how Claude Code should pair-mentor you (scaffolds over solutions, JS/Python → Go analogies, concept-before-code).
- **`curriculum-v2.md`** — the full plan, every task with mastery criteria and prereqs.
- **`improvements-summary.md`** — the ten gaps in the original v1 curriculum and what v2 does about them.

## Requirements

- Go **1.26** or newer — `setup.sh` checks this for you.
- Later phases require: Docker (Phase 3), JDK + Graphviz + Gnuplot + Maelstrom (Phases 4–5). Installed when you get there, not now.

## License

Personal learning project. Do whatever you want with it.
