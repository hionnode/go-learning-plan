package main

import (
	"fmt"
	"os"
)

// subcommandHelp holds the full help text for each subcommand. Shown via
//   go-dojo <cmd> --help
//   go-dojo <cmd> -h
//   go-dojo help <cmd>
var subcommandHelp = map[string]string{
	"serve": `go-dojo serve [addr]

Start the web dashboard + task browser.

Arguments:
  [addr]   Listen address. Default: :8080

Examples:
  go-dojo serve
  go-dojo serve :9090
  go-dojo serve 127.0.0.1:8080

What you get:
  /           mastery heatmap + current focus + due count
  /graph      prereq DAG rendered as server-side SVG
  /review     SRS queue for today
  /task/<id>  per-task detail + reflection form
  /drill      deliberate-practice drill list + runner

Shutdown: SIGINT (Ctrl-C). Graceful, 5s timeout.
`,

	"verify": `go-dojo verify <task-id>

Run the task's go test suite and update mastery state.

Arguments:
  <task-id>   ID from curriculum-v2.md (e.g. 1.1-hello-world)

Examples:
  go-dojo verify 1.1-hello-world
  go-dojo verify 2.2-net-http-server

What it does:
  Shells out to 'go test ./exercises/phase-N/<task-id>/...', captures
  output, writes the attempt to progress.json. On pass, advances mastery
  and schedules the next SRS review. On fail, suggests the remediation
  target from the task's frontmatter.

Mastery transitions on pass:
  unseen     → learning    (first time)
  learning   → proficient  (after ≥3 days since first pass)
  proficient → (stays)     (re-passing re-extends the review interval)

Tip: to list task IDs, 'grep ^id: curriculum-v2.md'.
`,

	"drill": `go-dojo drill <drill-id>

Run a timed deliberate-practice drill.

Arguments:
  <drill-id>   ID from the drills: library in curriculum-v2.md
               (e.g. stdin-echo, http-handler, goroutine-channel)

Examples:
  go-dojo drill stdin-echo
  go-dojo drill http-handler

How it works:
  Prints the prompt, waits for 'enter' to start, times your attempt,
  waits for 'enter' to stop. Records best time. Beating the target time
  promotes any parent task's mastery from 'proficient' to 'automatic'.

Tip: to list drill IDs, 'go-dojo drill' with no arg prints the usage.
     Or: 'grep "- id:" curriculum-v2.md'.
`,

	"review": `go-dojo review

Print the list of tasks whose SRS review is due today.

What it does:
  Reads progress.json, finds every task with nextReviewAt ≤ now, sorts
  by how overdue (most overdue first), prints them to stdout.

Every review is a closed-book retrieval — re-solve from memory, then
re-verify with 'go-dojo verify'. Missed reviews demote mastery one level.

Examples:
  go-dojo review
`,

	"placement": `go-dojo placement <phase-id>

Run a diagnostic placement quiz at phase entry. Passing marks earlier
tasks in that phase as placement-skipped so you don't grind through
what you already own.

Arguments:
  <phase-id>   phase-0 | phase-1 | phase-2 | phase-3 | phase-4

Examples:
  go-dojo placement phase-0
  go-dojo placement phase-2

Notes:
  - phase-5 (Gossip Glomers) has no placement — Maelstrom is the judge.
  - Placement sets covered tasks to 'learning' mastery if they were
    'unseen', so the prereq DAG unblocks downstream work.
`,

	"validate": `go-dojo validate [path]

Parse a curriculum or skill-tree markdown file and check its integrity.

Arguments:
  [path]   Path to .md file. Default: curriculum-v2.md at repo root.

Examples:
  go-dojo validate
  go-dojo validate explorations/netbird-skill-tree.md

Checks:
  - YAML frontmatter parses
  - Prereq DAG is acyclic (topological sort succeeds)
  - Every drill_id on a task resolves to a drill in the library
  - Every remediation / interleave_with task-id resolves

Exit code 0 if clean; non-zero and a list of problems otherwise.
`,

	"version": `go-dojo version

Print version, commit, build date, and runtime info.

Also accepts: go-dojo --version
`,

	"help": `go-dojo help [subcommand]

Print help. With no argument, print the top-level help. With a subcommand
name, print detailed help for that subcommand.

Examples:
  go-dojo help
  go-dojo help verify
  go-dojo verify --help     # equivalent
`,
}

// printSubcommandHelp writes the help for name (or the top-level usage if
// name isn't recognized) to stderr.
func printSubcommandHelp(name string) {
	if text, ok := subcommandHelp[name]; ok {
		fmt.Fprint(os.Stderr, text)
		return
	}
	fmt.Fprintf(os.Stderr, "no help for %q\n\n", name)
	usage()
}

// hasHelpFlag returns true if --help or -h appears anywhere in args.
func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}
