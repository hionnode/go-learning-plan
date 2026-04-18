// go-dojo is the learning-plan CLI + local web dashboard.
//
// Subcommands:
//   serve      — run HTTP server on :8080
//   verify     — run a task's go-test suite, update mastery
//   drill      — run a timed deliberate-practice drill
//   review     — print today's spaced-retrieval queue
//   placement  — run a phase diagnostic placement quiz
//   validate   — lint a curriculum or skill-tree file
//   version    — print version info
//   help       — show help (top-level or per-subcommand)
//
// All state lives in progress.json at the repo root; the curriculum is loaded
// from curriculum-v2.md. No external dependencies — stdlib only.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"learning-plan/internal/curriculum"
	"learning-plan/internal/progress"
)

// banner is the ASCII gate rendered on serve boot, no-args, and `version`.
// The crossbeam strokes (────┤ ... ├────) evoke a torii — the entrance to a dojo.
const banner = `
          ┌─────────────────────────────────────┐
       ───┤              go-dojo                ├───
          │      math-academy way tracker       │
          └─────────────────────────────────────┘
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, banner)
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	// Top-level --version / -v flag → version subcommand.
	if cmd == "--version" || cmd == "-v" {
		cmd = "version"
	}

	// Per-subcommand help: `go-dojo <cmd> --help` or `go-dojo help <cmd>`.
	if hasHelpFlag(args) {
		printSubcommandHelp(cmd)
		return
	}
	if cmd == "help" {
		if len(args) >= 1 {
			printSubcommandHelp(args[0])
		} else {
			fmt.Fprint(os.Stderr, banner)
			usage()
		}
		return
	}
	if cmd == "-h" || cmd == "--help" {
		fmt.Fprint(os.Stderr, banner)
		usage()
		return
	}

	root, err := findRepoRoot()
	if err != nil {
		die("finding repo root: %v", err)
	}
	ctx := newContext(root)

	switch cmd {
	case "serve":
		if err := runServe(ctx, args); err != nil {
			die("serve: %v", err)
		}
	case "verify":
		if err := runVerify(ctx, args); err != nil {
			die("verify: %v", err)
		}
	case "drill":
		if err := runDrill(ctx, args); err != nil {
			die("drill: %v", err)
		}
	case "review":
		if err := runReview(ctx, args); err != nil {
			die("review: %v", err)
		}
	case "placement":
		if err := runPlacement(ctx, args); err != nil {
			die("placement: %v", err)
		}
	case "validate":
		if err := runValidate(ctx, args); err != nil {
			die("validate: %v", err)
		}
	case "version":
		if err := runVersion(ctx, args); err != nil {
			die("version: %v", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `go-dojo — learning-plan Math-Academy-Way tracker

Usage:
  go-dojo serve                    start dashboard on http://localhost:8080
  go-dojo verify <task-id>         run the task's verify tests, update mastery
  go-dojo drill  <drill-id>        run a timed drill
  go-dojo review                   list tasks due for review today
  go-dojo placement <phase-id>     run a phase diagnostic (phase-0 … phase-4)
  go-dojo validate [path]          parse + DAG-check a curriculum / skill-tree
                                   file (defaults to curriculum-v2.md)
  go-dojo version                  print version, commit, build info
  go-dojo help [subcommand]        per-subcommand help

Shortcuts:
  go-dojo <subcommand> --help      detailed help for one subcommand
  go-dojo --version                same as 'go-dojo version'

Examples:
  go-dojo serve
  go-dojo verify 1.1-hello-world
  go-dojo drill stdin-echo
  go-dojo placement phase-0
  go-dojo validate explorations/netbird-skill-tree.md
`)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "go-dojo: "+format+"\n", args...)
	os.Exit(1)
}

type appContext struct {
	root           string
	curriculumPath string
	progressPath   string
}

func newContext(root string) *appContext {
	return &appContext{
		root:           root,
		curriculumPath: filepath.Join(root, "curriculum-v2.md"),
		progressPath:   filepath.Join(root, "progress.json"),
	}
}

func (a *appContext) loadCurriculum() (*curriculum.Curriculum, error) {
	return curriculum.Parse(a.curriculumPath)
}

func (a *appContext) loadState() (*progress.Store, *progress.State, error) {
	s := progress.Open(a.progressPath)
	st, err := s.Load()
	return s, st, err
}

// findRepoRoot walks up from cwd looking for curriculum-v2.md.
func findRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for dir := wd; ; {
		if _, err := os.Stat(filepath.Join(dir, "curriculum-v2.md")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("curriculum-v2.md not found in %s or any parent", wd)
		}
		dir = parent
	}
}
