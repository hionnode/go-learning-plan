package curriculum

import (
	"strings"
	"testing"
)

func TestParseReader(t *testing.T) {
	src := "```yaml\n" +
		"---\n" +
		"id: 1.1-hello-world\n" +
		"phase: 1\n" +
		"title: Setup\n" +
		"prereqs: []\n" +
		"interleave_with: [1.2-types-structs]\n" +
		"mastery_criteria: |\n" +
		"  - one\n" +
		"  - two\n" +
		"drill_ids: [stdin-echo, go-fmt-vet-cycle]\n" +
		"review_intervals_days: [3, 7, 21]\n" +
		"remediation: []\n" +
		"diagnostic_checkpoint: true\n" +
		"lgwt_chapters: [install-go]\n" +
		"---\n" +
		"```\n\n" +
		"```yaml\n" +
		"drills:\n" +
		"  - id: stdin-echo\n" +
		"    target_seconds: 300\n" +
		"    prompt: \"Echo stdin to stdout\"\n" +
		"  - id: go-fmt-vet-cycle\n" +
		"    target_seconds: 600\n" +
		"    prompt: \"Cycle through fmt/vet\"\n" +
		"```\n"

	c, err := ParseReader(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if len(c.Tasks) != 1 {
		t.Fatalf("want 1 task, got %d", len(c.Tasks))
	}
	task := c.Tasks[0]
	cases := []struct {
		name string
		got  any
		want any
	}{
		{"id", task.ID, "1.1-hello-world"},
		{"phase", task.Phase, 1},
		{"title", task.Title, "Setup"},
		{"prereqs empty", len(task.Prereqs), 0},
		{"interleave_with", strings.Join(task.InterleaveWith, ","), "1.2-types-structs"},
		{"mastery body", task.MasteryCriteria, "- one\n- two"},
		{"drill_ids", strings.Join(task.DrillIDs, ","), "stdin-echo,go-fmt-vet-cycle"},
		{"intervals len", len(task.ReviewIntervalsDays), 3},
		{"intervals[2]", task.ReviewIntervalsDays[2], 21},
		{"remediation empty", len(task.Remediation), 0},
		{"diagnostic", task.DiagnosticCheckpoint, true},
		{"lgwt", strings.Join(task.LgwtChapters, ","), "install-go"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, c.got, c.want)
		}
	}

	if len(c.Drills) != 2 {
		t.Fatalf("want 2 drills, got %d", len(c.Drills))
	}
	if c.Drills[0].ID != "stdin-echo" || c.Drills[0].TargetSeconds != 300 {
		t.Errorf("drill 0 bad: %+v", c.Drills[0])
	}
	if c.Drills[0].Prompt != "Echo stdin to stdout" {
		t.Errorf("drill prompt unquote: %q", c.Drills[0].Prompt)
	}
}

func TestTaskByID(t *testing.T) {
	c := &Curriculum{Tasks: []Task{{ID: "a"}, {ID: "b"}}}
	if c.TaskByID("b").ID != "b" {
		t.Fatal("TaskByID miss")
	}
	if c.TaskByID("zz") != nil {
		t.Fatal("TaskByID spurious hit")
	}
}

func TestParseList(t *testing.T) {
	tests := map[string][]string{
		"[]":              nil,
		"":                nil,
		"[a]":             {"a"},
		"[a, b, c]":       {"a", "b", "c"},
		"[\"a\", \"b\"]":  {"a", "b"},
		"[1, 2, 3]":       {"1", "2", "3"},
		"bare":            {"bare"},
	}
	for in, want := range tests {
		got := parseList(in)
		if len(got) != len(want) {
			t.Errorf("parseList(%q)=%v want %v", in, got, want)
			continue
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("parseList(%q)[%d]=%q want %q", in, i, got[i], want[i])
			}
		}
	}
}
