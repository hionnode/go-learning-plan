package srs

import (
	"testing"
	"time"

	"learning-plan/internal/progress"
)

func TestOnVerifyResult_FirstPass(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	tp := &progress.TaskProgress{Mastery: progress.Unseen}
	OnVerifyResult(tp, []int{3, 7, 21}, true, now)
	if tp.Mastery != progress.Learning {
		t.Errorf("want Learning, got %s", tp.Mastery)
	}
	if tp.NextReviewAt == nil {
		t.Fatal("NextReviewAt nil")
	}
	want := now.Add(72 * time.Hour)
	if !tp.NextReviewAt.Equal(want) {
		t.Errorf("NextReviewAt=%v want %v", tp.NextReviewAt, want)
	}
}

func TestOnVerifyResult_PromoteAfterGap(t *testing.T) {
	first := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)
	now := first.Add(5 * 24 * time.Hour)
	tp := &progress.TaskProgress{
		Mastery: progress.Learning,
		Attempts: []progress.Attempt{
			{At: first, Passed: true},
		},
	}
	OnVerifyResult(tp, []int{3, 7, 21}, true, now)
	if tp.Mastery != progress.Proficient {
		t.Errorf("want Proficient, got %s", tp.Mastery)
	}
}

func TestOnVerifyResult_NoPromoteSameDay(t *testing.T) {
	first := time.Date(2026, 4, 17, 0, 0, 0, 0, time.UTC)
	now := first.Add(2 * time.Hour)
	tp := &progress.TaskProgress{
		Mastery: progress.Learning,
		Attempts: []progress.Attempt{
			{At: first, Passed: true},
		},
	}
	OnVerifyResult(tp, []int{3, 7}, true, now)
	if tp.Mastery != progress.Learning {
		t.Errorf("stay Learning until ≥3 days, got %s", tp.Mastery)
	}
}

func TestOnVerifyResult_Fail_DemoteButFloorAtLearning(t *testing.T) {
	now := time.Now().UTC()
	tp := &progress.TaskProgress{Mastery: progress.Proficient, ReviewBox: 2}
	OnVerifyResult(tp, []int{3, 7, 21}, false, now)
	if tp.Mastery != progress.Learning {
		t.Errorf("want Learning, got %s", tp.Mastery)
	}
	if tp.ReviewBox != 0 {
		t.Errorf("box not reset: %d", tp.ReviewBox)
	}
}

func TestDueTasks_SortedByOverdue(t *testing.T) {
	now := time.Date(2026, 4, 17, 12, 0, 0, 0, time.UTC)
	old := now.Add(-48 * time.Hour)
	mid := now.Add(-6 * time.Hour)
	future := now.Add(24 * time.Hour)
	st := &progress.State{Tasks: map[string]*progress.TaskProgress{
		"late-most":     {Mastery: progress.Proficient, NextReviewAt: &old},
		"late-a-little": {Mastery: progress.Proficient, NextReviewAt: &mid},
		"not-due":       {Mastery: progress.Proficient, NextReviewAt: &future},
		"unseen":        {Mastery: progress.Unseen, NextReviewAt: &old},
	}}
	due := DueTasks(st, now)
	if len(due) != 2 {
		t.Fatalf("due=%v", due)
	}
	if due[0] != "late-most" {
		t.Errorf("first should be late-most, got %q", due[0])
	}
}

func TestOnReviewMiss_Demotes(t *testing.T) {
	now := time.Now().UTC()
	tp := &progress.TaskProgress{Mastery: progress.Automatic, ReviewBox: 3}
	OnReviewMiss(tp, []int{3, 7, 21, 60}, now)
	if tp.Mastery != progress.Proficient {
		t.Errorf("demote: %s", tp.Mastery)
	}
	if tp.ReviewBox != 0 {
		t.Errorf("box reset: %d", tp.ReviewBox)
	}
}

func TestPromoteToAutomatic_GatedOnProficient(t *testing.T) {
	tp := &progress.TaskProgress{Mastery: progress.Learning}
	PromoteToAutomatic(tp)
	if tp.Mastery != progress.Learning {
		t.Errorf("Learning must not jump to Automatic: %s", tp.Mastery)
	}
	tp.Mastery = progress.Proficient
	PromoteToAutomatic(tp)
	if tp.Mastery != progress.Automatic {
		t.Errorf("Proficient should become Automatic: %s", tp.Mastery)
	}
}
