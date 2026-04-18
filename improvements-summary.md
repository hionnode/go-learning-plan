# Improvements — Math-Academy-Way audit of `implementation-plan.md`

This is the human-readable changelog between `implementation-plan.md` (v1) and `curriculum-v2.md`. For each gap the audit found: *what was missing, what the Math-Academy-Way principle says, what v2 does instead.*

You don't have to adopt everything. Read it, pick the changes you believe in, reject the rest.

## 0. Learner profile pivot (v2.1)

v1 and the original v2 both assumed an experienced MuleSoft/Java developer transitioning to Go. v2 now assumes a different learner entirely: **someone who can write JavaScript and/or Python at a working level but has never built a backend** — no servers, no SQL, no caches, no queues, new to compiled + statically typed languages. Everything else below still applies; the practical changes from this pivot are:

- **New Phase 0 (~2 weeks).** Five onramp tasks covering the backend mental model (client/server, HTTP, TCP, DNS, JSON, database, cache, queue as *concepts*) and the Go essentials a JS/Python dev stumbles on (compilation vs REPL, zero values, errors-as-values, first HTTP call). No servers yet — just vocabulary and the red/green TDD loop.
- **Phase 3 expanded (~5 weeks, up from ~3–4).** Each tech — SQL, caching, queues — now gets a **concepts-first** task (raw psql, raw redis-cli, raw nats CLI) *before* the Go integration task. You learn what Redis *is* before you wrap it in `go-redis`.
- **Phase 2 gains an HTTP-protocol intro task (2.1).** Before writing a `net/http` server, you write a raw HTTP/1.1 request over a TCP socket by hand. The mystery dies before the library shows up.
- **JS/Python analogies replace Java analogies in `claude.md`.** Async vs goroutines, `undefined`/`None` vs zero values, errors-as-values vs try/catch, pointers vs reference semantics — all framed against what the learner already knows.
- **TDD is reinforced as the default.** Already true in v1, explicit in v2.1: every task ships with a `verify_test.go`; `go-learn`'s `verify` command *is* the grade.

Total timeline expands from ~22 weeks to ~26–30 weeks at 2–3 hrs/day. Task count grows from 22 to 43; drill library grows from 16 to 22 drills.

---

## 1. Done-ness ≠ mastery

**Gap.** v1 tracks progress with Markdown checkboxes. A checkbox answers "did I touch this?" — it does not answer "if I reopened this cold in 6 weeks, could I reproduce it?" Distributed systems problems compound; forgetting Phase 1's channels by Phase 5 is the default outcome, not the surprise.

**Principle.** **Mastery learning** (Bloom): progress is gated by demonstrated proficiency, not exposure. **Automaticity** (Skycak): the sign you've learned something is that you can do it without deliberate thought, quickly, and without reference.

**What v2 does.** Four-level mastery per task — `unseen → learning → proficient → automatic` — advanced only by passing a test `go-learn` runs (`learning`), re-passing the test after a delay (`proficient`), and beating the drill time target (`automatic`). Checkboxes are gone.

---

## 2. No explicit prerequisite graph

**Gap.** v1 orders tasks sequentially but never says *why*. 1.4 (pointers/slices/maps) is listed after 1.3 (errors/testing) even though it only needs 1.1. Worse, nothing tells `go-learn` which earlier task to re-surface when a later one fails.

**Principle.** **Knowledge-graph scaffolding** (Skycak): each skill is a node in a DAG; you earn access to a node when every prerequisite is `proficient` or better. Remediation is the inverse edge — failing node N points you at the prerequisite it exposed.

**What v2 does.** Every task frontmatter carries `prereqs: [...]` and `remediation: [...]`. `go-learn` refuses to let you start a task whose prereqs aren't `learning` yet, and when a verify fails twice it suggests the remediation target. See `internal/graph/dag.go`.

---

## 3. Serial ordering kills retention

**Gap.** v1 treats each phase as a block: finish 1.1–1.6, then start 2.x, then start 3.x. By the time you reach 3.x, 1.5's goroutines are three weeks cold.

**Principle.** **Interleaving + spaced retrieval** (Brown, Roediger): mixing related-but-distinct topics in the same session beats blocking them; spaced re-exposure beats massed practice even when massed practice feels better in the moment.

**What v2 does.**
- Each task has an `interleave_with` hint — tasks you can productively juggle in the same week. E.g. 1.4 ↔ 1.3 ↔ 1.5.
- Each task has `review_intervals_days: [3, 7, 21, 60]`. `go-learn` resurfaces every task at those offsets via the `review` queue. A missed review *demotes* mastery.

---

## 4. No automaticity targets

**Gap.** v1 says "write a worker pool" but never says *how fast* you should be able to do so. You can spend four hours building a worker pool once and call it done. That's exposure, not fluency.

**Principle.** **Deliberate practice** (Ericsson): short, repeated, timed exercises with immediate feedback, targeted at the skill you're weakest at. **Automaticity** comes from doing the same small thing many times until cognitive load drops near zero.

**What v2 does.** A `drills:` library of ~16 short exercises (stdin-echo, struct-json-roundtrip, table-driven-test, http-handler, context-timeout, goroutine-channel, mutex-vs-channel, errgroup-fanout, tcp-echo, ticker-shutdown, sql-crud-txn, redis-setnx-lock, maelstrom-init, lamport-clock, g-counter-merge, go-fmt-vet-cycle). Each has a `target_seconds`. Each parent task lists the drills that gate its `automatic` level. `go-learn`'s `drill` subcommand times you.

---

## 5. No diagnostic placement

**Gap.** v1 assumes everyone starts at 1.1. An experienced Java dev has already internalized some of this; grinding through it wastes weeks and kills motivation.

**Principle.** **Placement / diagnostic assessment** (Math Academy): before committing to N weeks of a phase, show that you *don't* already own it. Skip what you already have; spend the saved time where you're actually weak.

**What v2 does.** `go-learn placement phase-<N>` runs a short timed challenge at phase entry. Passing marks earlier tasks as placement-skipped in `progress.json.placement`. Phase-1–3 placements are concrete coding exercises; Phase 4 is a concept MCQ; Phase 5 has none (Maelstrom is the judge).

---

## 6. No active retrieval structure

**Gap.** v1's only learning artifact is a `notes/log.md` with freeform entries. That's reflection, not retrieval — and reflection alone doesn't build durable memory.

**Principle.** **Active recall + retrieval practice** (Roediger, Karpicke): the act of *producing* information from memory strengthens it more than re-reading or even re-highlighting. Every review should be a closed-book attempt, not a skim.

**What v2 does.**
- Every SRS review is explicitly a *retrieval* (closed book).
- Reflection log is structured with three prompts (**what you built / what clicked / what's still fuzzy**), written after the retrieval attempt, so the fuzzy bits feed back into the next review.
- The task detail page in `go-learn` UI shows the reflection log inline with attempt history.

---

## 7. Phase 5 is too coarse-grained

**Gap.** v1 bundles Broadcast 3a–3e into one bullet list, same for Kafka 5a–5c and Txn 6a–6c. `go-learn` can't tell whether you're at 3b or 3d — the smallest progress unit is "the whole challenge."

**Principle.** **Granular mastery units** (Skycak): break skills down to the smallest piece that can be independently verified. Progress lives in movements of these atoms, not in vague phase-level feelings.

**What v2 does.** Each sub-challenge (`5.3a`, `5.3b`, `5.3c`, `5.3d`, `5.3e`, `5.5a`, `5.5b`, `5.5c`, `5.6a`, `5.6b`, `5.6c`) is its own task with its own Maelstrom invocation as the verify. The prereq graph encodes that 3c blocks 3d, 5b blocks 5c, etc.

---

## 8. No targeted remediation loop

**Gap.** If you fail a Phase 5 challenge, v1 just says "take time, revisit Phase 4.2 notes." That's advice, not a protocol. You have no way to know if re-reading your notes is enough or if you need to re-do the implementation.

**Principle.** **Targeted remediation** (Skycak): a fail points you at a specific prerequisite, and you must re-pass that prerequisite's verify before retrying the failing task. No guessing.

**What v2 does.** Two consecutive `go-learn verify` failures on task N trigger a prompt to run `go-learn verify <remediation-id>` first. The remediation ID comes from the task's frontmatter.

---

## 9. Phase 1.6 is stapled on

**Gap.** v1 collects "TDD extras" into 1.6 as a dedicated section. The chapters are useful, but isolating them from 1.3–1.5 means the LGwT flow gets interrupted.

**Principle.** **Interleaving + integration**: don't quarantine the testing practice from the code it tests.

**What v2 does.** 1.6 remains a task (property tests, reflection, SVG clock) but is explicitly `interleave_with: [1.4, 1.5]` — you do it alongside, not after. Mapped LGwT chapters now live on each task's frontmatter (`lgwt_chapters:`) so `go-learn` can link directly.

---

## 10. No verification harness

**Gap.** v1 expects you to self-assess. Humans are bad at self-assessing fluency, especially on topics they've been reading about — familiarity masquerades as understanding.

**Principle.** **External verification**: the test runs, not your gut. Math Academy has the student solve a graded problem; we have `go test` and Maelstrom.

**What v2 does.** Every task has a concrete verify — unit tests in `exercises/phase-*/<task-id>/verify_test.go` for Phases 1–4, a Maelstrom invocation for Phase 5. `go-learn` CLI runs them and writes the result to `progress.json`. You can't mark yourself complete; the test marks you.

---

## What I explicitly *did not* change

- The five-phase high-level structure. It's sensible and matches the end goal.
- Time budget (2–3 hrs/day, no fixed deadline).
- Resource list in the appendix.
- `CLAUDE.md`'s teaching philosophy — the mentor rules are still the contract.
- The LGwT companion chapters — I kept the mapping; just moved it into frontmatter.

---

## What to skip if this feels like too much

If you only adopt two things: **(1) drills with timed targets** (automaticity) and **(2) spaced-retrieval reviews** (retention). Those two carry ~70% of the benefit. Everything else — placement, DAG, remediation graph, interleaving hints — is polish you can layer on once those are habit.
