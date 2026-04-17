package drills

import (
	"bytes"
	"strings"
	"testing"

	"learning-plan/internal/curriculum"
)

func TestRun_HappyPath(t *testing.T) {
	// Two newlines: one to start, one to finish.
	in := strings.NewReader("\n\n")
	var out bytes.Buffer
	drill := curriculum.Drill{ID: "test-drill", TargetSeconds: 300, Prompt: "do the thing"}
	res, err := Run(drill, in, &out)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if res.DurationMs < 0 {
		t.Errorf("negative duration: %d", res.DurationMs)
	}
	if !strings.Contains(out.String(), "drill: test-drill") {
		t.Errorf("missing id in output:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "do the thing") {
		t.Errorf("missing prompt in output")
	}
}

func TestRun_RejectsEmptyID(t *testing.T) {
	_, err := Run(curriculum.Drill{}, strings.NewReader("\n\n"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("want error for empty id")
	}
}
