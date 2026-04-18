package main

import (
	"fmt"
	"os"
	"time"

	"learning-plan/internal/drills"
	"learning-plan/internal/progress"
	"learning-plan/internal/srs"
)

func runDrill(ctx *appContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: go-learn drill <drill-id>")
	}
	id := args[0]

	curr, err := ctx.loadCurriculum()
	if err != nil {
		return fmt.Errorf("loading curriculum: %w", err)
	}
	d := curr.DrillByID(id)
	if d == nil {
		return fmt.Errorf("unknown drill %q", id)
	}

	store, state, err := ctx.loadState()
	if err != nil {
		return err
	}

	result, err := drills.Run(*d, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	dp := state.DrillOrInit(id)
	dp.History = append(dp.History, progress.DrillAttempt{
		At:         time.Now().UTC(),
		DurationMs: result.DurationMs,
		MetTarget:  result.MetTarget,
	})
	if dp.BestMs == 0 || result.DurationMs < dp.BestMs {
		dp.BestMs = result.DurationMs
	}

	// If the drill met target, promote any parent task from Proficient → Automatic.
	if result.MetTarget {
		for _, t := range curr.Tasks {
			for _, did := range t.DrillIDs {
				if did == id {
					if tp, ok := state.Tasks[t.ID]; ok {
						srs.PromoteToAutomatic(tp)
					}
				}
			}
		}
	}

	if err := store.Save(state); err != nil {
		return err
	}
	if dp.BestMs == result.DurationMs {
		fmt.Println("new personal best.")
	} else {
		fmt.Printf("best so far: %.1fs\n", float64(dp.BestMs)/1000)
	}
	return nil
}
