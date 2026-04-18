# explorations/

Pedagogical **skill trees** for third-party codebases — beginner-oriented prerequisite DAGs for understanding how a repo actually works.

Each file layers two frameworks:

- **[Math-Academy-Way](../curriculum-v2.md):** 4-level mastery, explicit prereqs, spaced retrieval intervals, timed drills for automaticity.
- **[`codebase-study-guide` skill](https://skills.sh/petekp/agent-skills/codebase-study-guide):** purpose-before-structure, 3 threshold concepts front-loaded, Mermaid system map, PRIMM (Predict-Run-Investigate-Modify) active-learning prompts.

Unlike `curriculum-v2.md`, these are **not** tracker-driven — there's no `verify_test.go` we control in a third-party repo. Treat them as reading orders you self-assess against.

## Existing trees

| File | Repo | Nodes | Drills |
|---|---|---|---|
| [netbird-skill-tree.md](netbird-skill-tree.md) | [netbirdio/netbird](https://github.com/netbirdio/netbird) — WireGuard mesh VPN | 15 | 14 |

## Add one for a new repo

### Prereqs (once per machine)

```sh
npx skills add petekp/agent-skills@codebase-study-guide -g -y
```

The skill's pedagogy lives at `~/.agents/skills/codebase-study-guide/SKILL.md`. Re-read it when designing a new tree.

### The five-step process

**1. Shallow-clone the target**
```sh
git clone --depth=1 <URL> /tmp/<repo>
```

**2. Explore via a subagent.** Launch one `Explore` agent (thoroughness: "very thorough") with a prompt that asks for:

- one-paragraph product summary (what it is, what problem it solves, who uses it)
- top-level directory map with one-line purposes
- entry points (every `main.go`, what binary each builds)
- 2–3 primary end-to-end flows traced file-by-file (e.g. "user logs in", "two peers connect")
- **3–5 threshold concepts** — non-obvious architectural insights that make the rest click
- external protocols/deps the reader needs to grok (gRPC, WireGuard, ICE, etc.) with pointers to where they're wired in
- testing strategy (what the tests reveal about intended behavior)
- suggested learning order — ~10–15 atomic concepts/subsystems with rough prereq hints

**3. Draft `explorations/<repo>-skill-tree.md`** using the structure below. Use the Explore agent's learning-order list as the node backbone; enrich each node with mastery criteria, prereqs, a drill ID, and a PRIMM task.

**4. Validate.**
```sh
go run ./cmd/tracker validate explorations/<repo>-skill-tree.md
```
This parses the frontmatter, builds the DAG, topo-sorts it (rejects cycles), and reports any dangling `drill_ids` / `remediation` / `interleave_with` refs. Fix problems until it prints `no dangling references`.

**5. Clean up and commit.**
```sh
rm -rf /tmp/<repo>
git add explorations/<repo>-skill-tree.md
git commit -m "docs: skill tree for <repo>"
```

## Required structure

Sections in this exact order:

1. **Header** — one-line description + who it's for + "fun-experiment status" note.
2. **Product in one paragraph** — what it is, the user's problem, the shape of the system.
3. **Threshold concepts** — exactly 3 paragraphs. Each is a transformative, integrative, non-obvious idea. These go *before* the diagram, not after.
4. **System map** — one Mermaid `graph TD` with the 4–8 primary components and their communication arrows.
5. **How mastery works** — terse table: what each of `unseen`/`learning`/`proficient`/`automatic` means.
6. **The tree** — 10–20 nodes across 4–6 stages. Each node has:
   - YAML frontmatter block (schema below)
   - Prose **reading list** with specific file paths in the target repo
   - **Active-learning task** using PRIMM phrasing ("Predict: … then read the code and check yourself.")
7. **Drill library** — a single `drills:` YAML block with every `drill_id` referenced above. Each drill has `target_seconds` and a `prompt` that's reproducible without the codebase open.
8. **Where to read first** — a small role-based table (architect / client-dev / server-dev / security / devops).
9. **Exit criterion** — 3–4 observable behaviors that mean "you've got it."

### Node frontmatter schema

```yaml
---
id: <repo>.<stage>.<n>-<slug>          # e.g. nb.3.2-sync-stream
phase: <int>                            # 1..5, used as DAG column
title: <short title>
prereqs: [<task-id>, ...]               # strict DAG — no cycles
interleave_with: [<task-id>, ...]       # siblings you can study alongside
mastery_criteria: |
  - <behavior-observable bullet 1>
  - <behavior-observable bullet 2>
  - <behavior-observable bullet 3>
drill_ids: [<drill-id>, ...]
review_intervals_days: [3, 7, 21, 60]   # keep the default unless you have a reason
remediation: [<task-id>, ...]           # where to go back to if this task fails
diagnostic_checkpoint: false            # true on the last node of each stage
---
```

### Drill format

```yaml
drills:
  - id: <drill-id>
    target_seconds: <int>                # 180 for recall, 300-600 for code tracing
    prompt: "One sentence. Timed. Reproducible without opening the codebase."
```

## Quality checklist

Before committing, run through these:

**Structure (mechanical):**
- [ ] `tracker validate` prints `no dangling references`
- [ ] Every node has YAML frontmatter; every frontmatter has all 10 fields
- [ ] Every `drill_id` on a task resolves in the drill library
- [ ] DAG has no cycles (validator enforces, but sanity-check the prereqs read right)

**Pedagogy (judgment):**
- [ ] **Purpose before structure** — can a reader skim sections 1–3 and understand *why this codebase exists*?
- [ ] **Threshold concepts front-loaded** — exactly 3, before any diagrams or file paths
- [ ] **Dual coding** — at least one Mermaid diagram, plus one per complex flow if warranted
- [ ] **Active, not passive** — every node ends with a PRIMM-style task, not just a description
- [ ] **Progressive disclosure** — a reader who stops after stage 2 still has useful understanding
- [ ] **Patterns named + linked** — ICE, WireGuard, gRPC streaming, etc. link to a canonical resource

**Scoping:**
- [ ] 10–20 nodes total. Fewer feels shallow; more is overwhelming.
- [ ] 4–6 stages. Each stage is 2–5 nodes.
- [ ] Each drill is timeboxed and achievable without the codebase open.

## Tips from the first tree (netbird)

- **Start prerequisites outside the repo.** Stage 1 should be the network/crypto/protocol concepts the learner needs *before* they can read the Go. Don't embed "learn WireGuard" inside a NetBird-specific node — make it its own prereq.
- **One threshold concept per diagram.** Don't try to put "management is source of truth," "control plane ≠ data plane," and "WireGuard keys = identity" into one diagram. Use three.
- **PRIMM prompts should fail safely.** "Predict what happens if the same setup key is used from two machines" is a good prompt because the learner can be wrong and learn something. "Summarize what this code does" is a bad prompt because there's no prediction.
- **Drills at 180–600 seconds.** Recall drills sit at 180s; file-tracing drills at 300–600s. Anything past 10 minutes isn't a drill, it's a subtask.
- **Name one exit criterion they'll feel.** "You can redraw the system map from memory" is testable. "You understand the architecture" isn't.
