// Package srs implements a Leitner-style spaced retrieval scheduler.
//
// The scheduler is deliberately simple: a "box" is an index into the task's
// review_intervals_days slice. Passing a review advances to the next box and
// schedules the next review that many days out. Missing a review demotes
// mastery one level and resets the box. No exponential magic, no personalized
// difficulty — the curriculum author picks intervals per task.
package srs

import (
	"time"

	"learning-plan/internal/progress"
)

// OnVerifyResult updates the TaskProgress after a verify attempt.
// Passing: bumps Learning→Proficient if eligible (≥ minProficientDays since first pass,
// otherwise stays Learning). Automatic is gated on drills elsewhere.
// Failing: demotes one level, but never below Learning if the task has ever passed.
func OnVerifyResult(tp *progress.TaskProgress, intervals []int, passed bool, now time.Time) {
	if tp.Mastery == "" {
		tp.Mastery = progress.Unseen
	}
	if passed {
		if tp.Mastery == progress.Unseen {
			tp.Mastery = progress.Learning
			tp.ReviewBox = 0
		} else if tp.Mastery == progress.Learning && eligibleForProficient(tp, now) {
			tp.Mastery = progress.Proficient
			tp.ReviewBox = advance(tp.ReviewBox, intervals)
		} else if tp.Mastery == progress.Proficient {
			tp.ReviewBox = advance(tp.ReviewBox, intervals)
		}
		// Automatic level is promoted elsewhere (drills), not here.
		tp.LastVerifiedAt = timePtr(now)
		tp.NextReviewAt = nextReview(tp.ReviewBox, intervals, now)
		return
	}
	// failed
	tp.Mastery = demote(tp.Mastery)
	tp.ReviewBox = 0
	tp.LastVerifiedAt = timePtr(now)
	tp.NextReviewAt = nextReview(0, intervals, now)
}

// OnReviewMiss is called when a scheduled review lapses — the learner didn't
// retrieve the task by nextReviewAt. Demote mastery one level and reset box.
func OnReviewMiss(tp *progress.TaskProgress, intervals []int, now time.Time) {
	tp.Mastery = demote(tp.Mastery)
	tp.ReviewBox = 0
	tp.NextReviewAt = nextReview(0, intervals, now)
}

// DueTasks returns task IDs whose NextReviewAt is at or before now, sorted by
// how overdue they are (most overdue first).
func DueTasks(state *progress.State, now time.Time) []string {
	type due struct {
		id   string
		late time.Duration
	}
	var out []due
	for id, tp := range state.Tasks {
		if tp.NextReviewAt == nil {
			continue
		}
		if tp.Mastery == progress.Unseen {
			continue
		}
		if !tp.NextReviewAt.After(now) {
			out = append(out, due{id, now.Sub(*tp.NextReviewAt)})
		}
	}
	// simple insertion sort by most overdue first
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].late > out[j-1].late; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	ids := make([]string, len(out))
	for i := range out {
		ids[i] = out[i].id
	}
	return ids
}

// PromoteToAutomatic upgrades mastery if a drill was completed within target.
// Call this after a successful drill tied to the task.
func PromoteToAutomatic(tp *progress.TaskProgress) {
	if tp.Mastery == progress.Proficient {
		tp.Mastery = progress.Automatic
	}
}

func advance(box int, intervals []int) int {
	if len(intervals) == 0 {
		return 0
	}
	next := box + 1
	if next >= len(intervals) {
		next = len(intervals) - 1
	}
	return next
}

func nextReview(box int, intervals []int, now time.Time) *time.Time {
	if len(intervals) == 0 {
		return nil
	}
	if box >= len(intervals) {
		box = len(intervals) - 1
	}
	days := intervals[box]
	t := now.Add(time.Duration(days) * 24 * time.Hour)
	return &t
}

func demote(m progress.Mastery) progress.Mastery {
	switch m {
	case progress.Automatic:
		return progress.Proficient
	case progress.Proficient:
		return progress.Learning
	case progress.Learning:
		return progress.Learning // don't drop a once-learned task back to Unseen
	default:
		return progress.Unseen
	}
}

// A task can move Learning→Proficient once enough time has elapsed since the
// first verify — we want to see it survive a gap, not just repeat back-to-back.
const minProficientDays = 3

func eligibleForProficient(tp *progress.TaskProgress, now time.Time) bool {
	// Find earliest passed attempt.
	var first *time.Time
	for _, a := range tp.Attempts {
		if a.Passed {
			a := a
			first = &a.At
			break
		}
	}
	if first == nil {
		return false
	}
	return now.Sub(*first) >= minProficientDays*24*time.Hour
}

func timePtr(t time.Time) *time.Time { return &t }
