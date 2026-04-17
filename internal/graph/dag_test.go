package graph

import (
	"strings"
	"testing"

	"learning-plan/internal/curriculum"
	"learning-plan/internal/progress"
)

func TestBuildAndTopoSort(t *testing.T) {
	tasks := []curriculum.Task{
		{ID: "a", Phase: 1},
		{ID: "b", Phase: 1, Prereqs: []string{"a"}},
		{ID: "c", Phase: 2, Prereqs: []string{"b"}},
		{ID: "d", Phase: 2, Prereqs: []string{"a"}},
	}
	d, err := Build(tasks)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	order, err := d.TopoSort()
	if err != nil {
		t.Fatalf("topo: %v", err)
	}
	pos := map[string]int{}
	for i, id := range order {
		pos[id] = i
	}
	if !(pos["a"] < pos["b"] && pos["b"] < pos["c"] && pos["a"] < pos["d"]) {
		t.Errorf("bad order: %v", order)
	}
}

func TestBuildCycleRejected(t *testing.T) {
	tasks := []curriculum.Task{
		{ID: "a", Prereqs: []string{"b"}},
		{ID: "b", Prereqs: []string{"a"}},
	}
	_, err := Build(tasks)
	if err == nil {
		t.Fatal("cycle accepted")
	}
}

func TestBuildUnknownPrereq(t *testing.T) {
	tasks := []curriculum.Task{
		{ID: "a", Prereqs: []string{"missing"}},
	}
	_, err := Build(tasks)
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("want unknown prereq error, got %v", err)
	}
}

func TestReadyAndNextFocus(t *testing.T) {
	tasks := []curriculum.Task{
		{ID: "a", Phase: 1},
		{ID: "b", Phase: 1, Prereqs: []string{"a"}},
		{ID: "c", Phase: 2, Prereqs: []string{"b"}},
	}
	d, err := Build(tasks)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	st := progress.NewState()
	if d.Ready("b", st) {
		t.Error("b shouldn't be ready yet")
	}
	if d.NextFocus(st) != "a" {
		t.Errorf("focus: %s", d.NextFocus(st))
	}
	st.TaskOrInit("a").Mastery = progress.Learning
	if !d.Ready("b", st) {
		t.Error("b should be ready now")
	}
	if d.NextFocus(st) != "a" {
		// a is still learning, not Automatic, so it remains the focus
		t.Errorf("focus: %s", d.NextFocus(st))
	}
	st.TaskOrInit("a").Mastery = progress.Automatic
	if d.NextFocus(st) != "b" {
		t.Errorf("after a automatic, focus=%s", d.NextFocus(st))
	}
}
