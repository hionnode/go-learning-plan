package progress

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMissingReturnsFresh(t *testing.T) {
	dir := t.TempDir()
	st, err := Open(filepath.Join(dir, "nope.json")).Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if st.SchemaVersion != SchemaVersion {
		t.Errorf("schema: %d", st.SchemaVersion)
	}
	if st.Tasks == nil || st.Drills == nil || st.Placement == nil {
		t.Errorf("nil maps in fresh state")
	}
}

func TestSaveLoadRoundtrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "progress.json")
	s := Open(p)

	now := time.Now().UTC().Truncate(time.Second)
	st := NewState()
	tp := st.TaskOrInit("1.1-hello-world")
	tp.Mastery = Learning
	tp.LastVerifiedAt = &now
	tp.Attempts = []Attempt{{At: now, Passed: true, DurationMs: 420}}
	tp.Reflections = []Reflection{{At: now, Built: "echo", Clicked: "modules", Fuzzy: "imports"}}

	dp := st.DrillOrInit("stdin-echo")
	dp.BestMs = 250_000
	dp.History = append(dp.History, DrillAttempt{At: now, DurationMs: 250_000, MetTarget: true})

	st.Placement["phase-1"] = &PlacementResult{At: now, Score: 0.8, SkippedTasks: []string{"1.1-hello-world"}}

	if err := s.Save(st); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("stat: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Tasks["1.1-hello-world"].Mastery != Learning {
		t.Errorf("mastery: %s", got.Tasks["1.1-hello-world"].Mastery)
	}
	if got.Drills["stdin-echo"].BestMs != 250_000 {
		t.Errorf("drill bestMs: %d", got.Drills["stdin-echo"].BestMs)
	}
	if got.Placement["phase-1"].Score != 0.8 {
		t.Errorf("placement score: %v", got.Placement["phase-1"].Score)
	}
	if len(got.Tasks["1.1-hello-world"].Attempts) != 1 {
		t.Errorf("attempts: %d", len(got.Tasks["1.1-hello-world"].Attempts))
	}
}

func TestTaskOrInitUnseen(t *testing.T) {
	st := NewState()
	tp := st.TaskOrInit("new-task")
	if tp.Mastery != Unseen {
		t.Errorf("fresh task mastery: %s", tp.Mastery)
	}
}
