// go-learn is the learning-plan CLI + local web dashboard.
//
// Subcommands:
//   serve      — run HTTP server on :8080
//   verify     — run a task's go-test suite, update mastery
//   drill      — run a timed deliberate-practice drill
//   review     — print today's spaced-retrieval queue
//   placement  — run a phase diagnostic placement quiz
//   validate   — lint a curriculum or skill-tree file
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

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

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
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `go-learn — learning-plan Math-Academy-Way tracker

Usage:
  go-learn serve                   start dashboard on http://localhost:8080
  go-learn verify <task-id>        run the task's verify tests, update mastery
  go-learn drill  <drill-id>       run a timed drill
  go-learn review                  list tasks due for review today
  go-learn placement <phase-id>    run a phase diagnostic (phase-0 … phase-4)
  go-learn validate [path]         parse + DAG-check a curriculum / skill-tree
                                   file (defaults to curriculum-v2.md)

Examples:
  go run ./cmd/go-learn serve
  go run ./cmd/go-learn verify 1.1-hello-world
  go run ./cmd/go-learn drill stdin-echo
  go run ./cmd/go-learn review
  go run ./cmd/go-learn placement phase-0
  go run ./cmd/go-learn validate explorations/netbird-skill-tree.md
`)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "go-learn: "+format+"\n", args...)
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
