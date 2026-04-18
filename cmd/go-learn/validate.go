package main

import (
	"fmt"

	"learning-plan/internal/curriculum"
	"learning-plan/internal/graph"
)

// runValidate parses a curriculum or skill-tree markdown file and reports any
// structural problems: unparseable frontmatter, cyclic prereqs, dangling drill
// refs, dangling remediation targets. No side effects on progress.json — this
// is a pure linter.
func runValidate(ctx *appContext, args []string) error {
	path := ctx.curriculumPath
	if len(args) >= 1 {
		path = args[0]
	}

	c, err := curriculum.Parse(path)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	dag, err := graph.Build(c.Tasks)
	if err != nil {
		return fmt.Errorf("build DAG: %w", err)
	}
	order, err := dag.TopoSort()
	if err != nil {
		return fmt.Errorf("topo sort: %w", err)
	}

	drillIDs := map[string]bool{}
	for _, d := range c.Drills {
		drillIDs[d.ID] = true
	}
	taskIDs := map[string]bool{}
	for _, t := range c.Tasks {
		taskIDs[t.ID] = true
	}

	var problems []string
	for _, t := range c.Tasks {
		for _, did := range t.DrillIDs {
			if !drillIDs[did] {
				problems = append(problems, fmt.Sprintf("task %s references missing drill %s", t.ID, did))
			}
		}
		for _, rid := range t.Remediation {
			if !taskIDs[rid] {
				problems = append(problems, fmt.Sprintf("task %s references missing remediation task %s", t.ID, rid))
			}
		}
		for _, iid := range t.InterleaveWith {
			if !taskIDs[iid] {
				problems = append(problems, fmt.Sprintf("task %s references missing interleave_with task %s", t.ID, iid))
			}
		}
	}

	fmt.Printf("✓ %s\n", path)
	fmt.Printf("  tasks:  %d (topo-sorted: %d)\n", len(c.Tasks), len(order))
	fmt.Printf("  drills: %d\n", len(c.Drills))
	if len(problems) > 0 {
		fmt.Printf("  %d problem(s):\n", len(problems))
		for _, p := range problems {
			fmt.Printf("    - %s\n", p)
		}
		return fmt.Errorf("%d validation problem(s)", len(problems))
	}
	fmt.Println("  no dangling references")
	return nil
}
