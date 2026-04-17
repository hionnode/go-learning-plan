package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"learning-plan/internal/progress"
	"learning-plan/internal/srs"
)

func runVerify(ctx *appContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: tracker verify <task-id>")
	}
	taskID := args[0]

	curr, err := ctx.loadCurriculum()
	if err != nil {
		return fmt.Errorf("loading curriculum: %w", err)
	}
	task := curr.TaskByID(taskID)
	if task == nil {
		return fmt.Errorf("unknown task %q", taskID)
	}

	dir := taskDir(ctx.root, taskID)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("no exercise scaffold at %s — seed it first", dir)
	}

	store, state, err := ctx.loadState()
	if err != nil {
		return fmt.Errorf("loading state: %w", err)
	}

	start := time.Now()
	passed, output, _ := runGoTest(context.Background(), dir)
	duration := time.Since(start)

	tp := state.TaskOrInit(taskID)
	tp.Attempts = append(tp.Attempts, progress.Attempt{
		At:         start.UTC(),
		Passed:     passed,
		DurationMs: duration.Milliseconds(),
		Output:     output,
	})
	srs.OnVerifyResult(tp, task.ReviewIntervalsDays, passed, start.UTC())

	if err := store.Save(state); err != nil {
		return fmt.Errorf("saving state: %w", err)
	}

	fmt.Println(output)
	fmt.Println(strings.Repeat("─", 60))
	if passed {
		fmt.Printf("✓ %s passed — mastery now %s\n", taskID, tp.Mastery)
		if tp.NextReviewAt != nil {
			fmt.Printf("  next review: %s\n", tp.NextReviewAt.Format(time.RFC3339))
		}
	} else {
		fmt.Printf("✗ %s failed — mastery %s\n", taskID, tp.Mastery)
		if len(task.Remediation) > 0 {
			fmt.Printf("  remediation: %s\n", strings.Join(task.Remediation, ", "))
			fmt.Printf("  try: tracker verify %s\n", task.Remediation[0])
		}
	}
	return nil
}

// taskDir maps a task ID like "1.1-hello-world" to
// "<root>/exercises/phase-1/1.1-hello-world".
func taskDir(root, id string) string {
	if id == "" {
		return ""
	}
	phase := string(id[0])
	return filepath.Join(root, "exercises", "phase-"+phase, id)
}

func runGoTest(ctx context.Context, dir string) (bool, string, error) {
	cmd := exec.CommandContext(ctx, "go", "test", "-v", "-count=1", "./...")
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	output := buf.String()
	if err != nil {
		return false, output, err
	}
	return true, output, nil
}
