// Package progress owns the tracker's persistent state in progress.json.
//
// Writes are atomic: we write to a temp file in the same directory, fsync, and
// rename over the target. That matches how battle-tested CLIs (git, vim) avoid
// corrupting state on crash — the rename is the commit point.
package progress

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const SchemaVersion = 1

type Mastery string

const (
	Unseen     Mastery = "unseen"
	Learning   Mastery = "learning"
	Proficient Mastery = "proficient"
	Automatic  Mastery = "automatic"
)

type Attempt struct {
	At         time.Time `json:"at"`
	Passed     bool      `json:"passed"`
	DurationMs int64     `json:"durationMs"`
	Output     string    `json:"output,omitempty"`
}

type Reflection struct {
	At     time.Time `json:"at"`
	Built  string    `json:"built"`
	Clicked string   `json:"clicked"`
	Fuzzy  string    `json:"fuzzy"`
}

type TaskProgress struct {
	Mastery        Mastery      `json:"mastery"`
	LastVerifiedAt *time.Time   `json:"lastVerifiedAt,omitempty"`
	NextReviewAt   *time.Time   `json:"nextReviewAt,omitempty"`
	ReviewBox      int          `json:"reviewBox"` // 0-based index into ReviewIntervalsDays
	Attempts       []Attempt    `json:"attempts,omitempty"`
	Reflections    []Reflection `json:"reflections,omitempty"`
}

type DrillAttempt struct {
	At         time.Time `json:"at"`
	DurationMs int64     `json:"durationMs"`
	MetTarget  bool      `json:"metTarget"`
}

type DrillProgress struct {
	BestMs  int64          `json:"bestMs,omitempty"`
	History []DrillAttempt `json:"history,omitempty"`
}

type PlacementResult struct {
	At           time.Time `json:"at"`
	Score        float64   `json:"score"`
	SkippedTasks []string  `json:"skippedTasks,omitempty"`
}

type State struct {
	SchemaVersion int                         `json:"schemaVersion"`
	Tasks         map[string]*TaskProgress    `json:"tasks"`
	Drills        map[string]*DrillProgress   `json:"drills"`
	Placement     map[string]*PlacementResult `json:"placement"`
}

func NewState() *State {
	return &State{
		SchemaVersion: SchemaVersion,
		Tasks:         map[string]*TaskProgress{},
		Drills:        map[string]*DrillProgress{},
		Placement:     map[string]*PlacementResult{},
	}
}

// TaskOrInit returns the TaskProgress entry, creating an Unseen one if absent.
func (s *State) TaskOrInit(id string) *TaskProgress {
	if s.Tasks == nil {
		s.Tasks = map[string]*TaskProgress{}
	}
	tp, ok := s.Tasks[id]
	if !ok {
		tp = &TaskProgress{Mastery: Unseen}
		s.Tasks[id] = tp
	}
	if tp.Mastery == "" {
		tp.Mastery = Unseen
	}
	return tp
}

// DrillOrInit returns the DrillProgress entry, creating a blank one if absent.
func (s *State) DrillOrInit(id string) *DrillProgress {
	if s.Drills == nil {
		s.Drills = map[string]*DrillProgress{}
	}
	dp, ok := s.Drills[id]
	if !ok {
		dp = &DrillProgress{}
		s.Drills[id] = dp
	}
	return dp
}

// Store is a single-file, process-local progress store with a coarse mutex.
type Store struct {
	path string
	mu   sync.Mutex
}

// Open returns a Store for the given path. The file does not need to exist yet.
func Open(path string) *Store {
	return &Store{path: path}
}

// Load reads and decodes the state. If the file is missing, a fresh state is returned.
func (s *Store) Load() (*State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return NewState(), nil
		}
		return nil, fmt.Errorf("reading %s: %w", s.path, err)
	}
	if len(data) == 0 {
		return NewState(), nil
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", s.path, err)
	}
	if st.SchemaVersion == 0 {
		st.SchemaVersion = SchemaVersion
	}
	if st.SchemaVersion != SchemaVersion {
		return nil, fmt.Errorf("schemaVersion %d not supported (want %d)", st.SchemaVersion, SchemaVersion)
	}
	if st.Tasks == nil {
		st.Tasks = map[string]*TaskProgress{}
	}
	if st.Drills == nil {
		st.Drills = map[string]*DrillProgress{}
	}
	if st.Placement == nil {
		st.Placement = map[string]*PlacementResult{}
	}
	return &st, nil
}

// Save atomically writes the state to disk.
func (s *Store) Save(state *State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state.SchemaVersion = SchemaVersion
	buf, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding state: %w", err)
	}
	dir := filepath.Dir(s.path)
	tmp, err := os.CreateTemp(dir, ".progress-*.json.tmp")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(buf); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("writing temp: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("syncing temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("closing temp: %w", err)
	}
	if err := os.Rename(tmpPath, s.path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming temp to target: %w", err)
	}
	return nil
}
