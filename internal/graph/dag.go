// Package graph builds and queries the curriculum prerequisite DAG.
//
// Nodes are task IDs. Edges go prereq → dependent. Two operations matter:
//   (1) TopoSort for a phase-by-phase, prereq-respecting display order;
//   (2) Ready(id, state) to tell whether a learner has earned access to a task.
package graph

import (
	"fmt"

	"learning-plan/internal/curriculum"
	"learning-plan/internal/progress"
)

type DAG struct {
	Nodes    []string            // task IDs in declared order
	Phase    map[string]int      // id -> phase
	Prereqs  map[string][]string // id -> prereqs
	Children map[string][]string // id -> dependents
}

func Build(tasks []curriculum.Task) (*DAG, error) {
	d := &DAG{
		Phase:    map[string]int{},
		Prereqs:  map[string][]string{},
		Children: map[string][]string{},
	}
	for _, t := range tasks {
		d.Nodes = append(d.Nodes, t.ID)
		d.Phase[t.ID] = t.Phase
		d.Prereqs[t.ID] = append([]string(nil), t.Prereqs...)
	}
	for _, t := range tasks {
		for _, p := range t.Prereqs {
			if _, ok := d.Phase[p]; !ok {
				return nil, fmt.Errorf("task %q prereq %q is unknown", t.ID, p)
			}
			d.Children[p] = append(d.Children[p], t.ID)
		}
	}
	if _, err := d.TopoSort(); err != nil {
		return nil, err
	}
	return d, nil
}

// TopoSort returns node IDs in an order where every prereq precedes its dependents.
// Returns an error if there is a cycle.
func (d *DAG) TopoSort() ([]string, error) {
	indeg := map[string]int{}
	for _, id := range d.Nodes {
		indeg[id] = len(d.Prereqs[id])
	}
	var ready []string
	for _, id := range d.Nodes {
		if indeg[id] == 0 {
			ready = append(ready, id)
		}
	}
	var out []string
	for len(ready) > 0 {
		id := ready[0]
		ready = ready[1:]
		out = append(out, id)
		for _, child := range d.Children[id] {
			indeg[child]--
			if indeg[child] == 0 {
				ready = append(ready, child)
			}
		}
	}
	if len(out) != len(d.Nodes) {
		return nil, fmt.Errorf("prereq cycle detected")
	}
	return out, nil
}

// Ready reports whether all prereqs of id are at Learning mastery or better.
// A task is "ready" to be worked on, which is weaker than its prereqs being Proficient —
// you can start 1.5 once 1.4 is Learning, but you shouldn't expect to make it to
// Automatic without 1.4 also reaching Proficient.
func (d *DAG) Ready(id string, state *progress.State) bool {
	for _, p := range d.Prereqs[id] {
		tp, ok := state.Tasks[p]
		if !ok || tp.Mastery == "" || tp.Mastery == progress.Unseen {
			return false
		}
	}
	return true
}

// NextFocus returns the first ready, non-Automatic task (prereq order). If
// none are ready, returns "".
func (d *DAG) NextFocus(state *progress.State) string {
	order, _ := d.TopoSort()
	for _, id := range order {
		tp := state.Tasks[id]
		if tp != nil && tp.Mastery == progress.Automatic {
			continue
		}
		if d.Ready(id, state) {
			return id
		}
	}
	return ""
}
