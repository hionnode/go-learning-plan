// Package drills runs a timed deliberate-practice exercise.
//
// Interface:
//   Run displays the prompt, waits for a "start" signal, timers the learner's
//   work, and records the elapsed time. No correctness checking — the drill
//   is an honor-system commitment. This matches how Math Academy grades drills:
//   speed is the target, the learner verifies their own output.
package drills

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"learning-plan/internal/curriculum"
)

type Result struct {
	DurationMs int64
	MetTarget  bool
}

// Run drives an interactive drill session against stdin/stdout.
//   in:  reader for the learner's keystrokes
//   out: writer for prompts
func Run(d curriculum.Drill, in io.Reader, out io.Writer) (*Result, error) {
	if d.ID == "" {
		return nil, fmt.Errorf("empty drill id")
	}
	r := bufio.NewReader(in)

	fmt.Fprintf(out, "drill: %s\n", d.ID)
	fmt.Fprintf(out, "target: %ds\n\n", d.TargetSeconds)
	fmt.Fprintf(out, "prompt:\n  %s\n\n", d.Prompt)
	fmt.Fprint(out, "press enter when ready to start > ")
	if _, err := r.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("reading start: %w", err)
	}
	start := time.Now()

	fmt.Fprint(out, "\ntimer running. press enter when done > ")
	if _, err := r.ReadString('\n'); err != nil {
		return nil, fmt.Errorf("reading stop: %w", err)
	}
	elapsed := time.Since(start)
	ms := elapsed.Milliseconds()
	target := time.Duration(d.TargetSeconds) * time.Second
	met := elapsed <= target

	fmt.Fprintf(out, "\nelapsed: %.1fs  (target %ds) — %s\n",
		elapsed.Seconds(), d.TargetSeconds, hitMiss(met))

	return &Result{DurationMs: ms, MetTarget: met}, nil
}

func hitMiss(met bool) string {
	if met {
		return "met target"
	}
	return "over target"
}
