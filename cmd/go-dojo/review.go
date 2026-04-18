package main

import (
	"fmt"
	"time"

	"learning-plan/internal/srs"
)

func runReview(ctx *appContext, args []string) error {
	_ = args
	curr, err := ctx.loadCurriculum()
	if err != nil {
		return err
	}
	_, state, err := ctx.loadState()
	if err != nil {
		return err
	}
	due := srs.DueTasks(state, time.Now().UTC())
	if len(due) == 0 {
		fmt.Println("nothing due for review today.")
		return nil
	}
	fmt.Printf("due for review (%d):\n", len(due))
	for _, id := range due {
		tp := state.Tasks[id]
		task := curr.TaskByID(id)
		title := id
		if task != nil {
			title = task.Title
		}
		late := time.Since(*tp.NextReviewAt)
		fmt.Printf("  %-30s  %-12s  %s overdue  — %s\n",
			id, tp.Mastery, roundDuration(late), title)
	}
	fmt.Println("\nrun: go-dojo verify <task-id>   for each — retrieval, not re-reading.")
	return nil
}

func roundDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}
