// Package curriculum loads curriculum-v2.md and parses its YAML frontmatter
// blocks into Task and Drill values.
//
// The parser intentionally supports only the subset of YAML the curriculum
// file actually uses: scalars, inline lists, block scalars (|), bools, and
// a single list-of-maps shape for drills. Pulling in a full YAML library
// would violate the stdlib-only constraint and isn't needed.
package curriculum

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Task is one curriculum unit parsed from a frontmatter block.
type Task struct {
	ID                   string
	Phase                int
	Title                string
	Prereqs              []string
	InterleaveWith       []string
	MasteryCriteria      string
	DrillIDs             []string
	ReviewIntervalsDays  []int
	Remediation          []string
	DiagnosticCheckpoint bool
	LgwtChapters         []string
}

// Drill is one entry from the drills: list.
type Drill struct {
	ID            string
	TargetSeconds int
	Prompt        string
}

// Curriculum is the fully parsed result.
type Curriculum struct {
	Tasks  []Task
	Drills []Drill
}

// TaskByID returns the task with the given ID, or nil if absent.
func (c *Curriculum) TaskByID(id string) *Task {
	for i := range c.Tasks {
		if c.Tasks[i].ID == id {
			return &c.Tasks[i]
		}
	}
	return nil
}

// DrillByID returns the drill with the given ID, or nil if absent.
func (c *Curriculum) DrillByID(id string) *Drill {
	for i := range c.Drills {
		if c.Drills[i].ID == id {
			return &c.Drills[i]
		}
	}
	return nil
}

// Parse reads a curriculum markdown file.
func Parse(path string) (*Curriculum, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening curriculum %q: %w", path, err)
	}
	defer f.Close()
	return ParseReader(f)
}

// ParseReader is the core parser.
func ParseReader(r io.Reader) (*Curriculum, error) {
	blocks, err := extractYAMLBlocks(r)
	if err != nil {
		return nil, fmt.Errorf("extracting yaml blocks: %w", err)
	}
	c := &Curriculum{}
	for idx, b := range blocks {
		switch {
		case isDrillsBlock(b):
			drills, err := parseDrills(b)
			if err != nil {
				return nil, fmt.Errorf("block %d drills: %w", idx, err)
			}
			c.Drills = append(c.Drills, drills...)
		case isTaskBlock(b):
			t, err := parseTask(b)
			if err != nil {
				return nil, fmt.Errorf("block %d task: %w", idx, err)
			}
			c.Tasks = append(c.Tasks, *t)
		}
	}
	return c, nil
}

func extractYAMLBlocks(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var blocks []string
	var cur strings.Builder
	inBlock := false
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		if !inBlock && trimmed == "```yaml" {
			inBlock = true
			cur.Reset()
			continue
		}
		if inBlock && trimmed == "```" {
			inBlock = false
			blocks = append(blocks, cur.String())
			continue
		}
		if inBlock {
			cur.WriteString(line)
			cur.WriteByte('\n')
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return blocks, nil
}

func isDrillsBlock(block string) bool {
	for _, ln := range strings.Split(block, "\n") {
		t := strings.TrimSpace(ln)
		if t == "" || t == "---" {
			continue
		}
		return strings.HasPrefix(t, "drills:")
	}
	return false
}

func isTaskBlock(block string) bool {
	for _, ln := range strings.Split(block, "\n") {
		t := strings.TrimSpace(ln)
		if t == "" || t == "---" {
			continue
		}
		return strings.HasPrefix(t, "id:")
	}
	return false
}

func parseTask(block string) (*Task, error) {
	lines := strings.Split(block, "\n")
	t := &Task{}
	for i := 0; i < len(lines); {
		ln := lines[i]
		tr := strings.TrimSpace(ln)
		if tr == "" || tr == "---" {
			i++
			continue
		}
		if leadingSpaces(ln) > 0 {
			// orphan continuation line, skip (block scalar body is consumed inline below)
			i++
			continue
		}
		key, val, ok := splitKeyValue(ln)
		if !ok {
			i++
			continue
		}
		if val == "|" {
			i++
			body, consumed := readBlockScalar(lines[i:])
			if err := assignTaskField(t, key, body); err != nil {
				return nil, err
			}
			i += consumed
			continue
		}
		if err := assignTaskField(t, key, val); err != nil {
			return nil, err
		}
		i++
	}
	if t.ID == "" {
		return nil, fmt.Errorf("task block missing id")
	}
	return t, nil
}

func readBlockScalar(rest []string) (string, int) {
	baseIndent := -1
	var sb strings.Builder
	consumed := 0
	for _, ln := range rest {
		if strings.TrimSpace(ln) == "" {
			sb.WriteByte('\n')
			consumed++
			continue
		}
		ind := leadingSpaces(ln)
		if ind == 0 {
			break
		}
		if baseIndent < 0 {
			baseIndent = ind
		}
		if ind < baseIndent {
			break
		}
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		if ind >= baseIndent {
			sb.WriteString(ln[baseIndent:])
		} else {
			sb.WriteString(strings.TrimLeft(ln, " "))
		}
		consumed++
	}
	return strings.TrimRight(sb.String(), "\n"), consumed
}

func assignTaskField(t *Task, key, val string) error {
	switch key {
	case "id":
		t.ID = unquote(val)
	case "phase":
		n, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			return fmt.Errorf("phase: %w", err)
		}
		t.Phase = n
	case "title":
		t.Title = unquote(val)
	case "prereqs":
		t.Prereqs = parseList(val)
	case "interleave_with":
		t.InterleaveWith = parseList(val)
	case "mastery_criteria":
		t.MasteryCriteria = val
	case "drill_ids":
		t.DrillIDs = parseList(val)
	case "review_intervals_days":
		for _, s := range parseList(val) {
			n, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("review_intervals_days %q: %w", s, err)
			}
			t.ReviewIntervalsDays = append(t.ReviewIntervalsDays, n)
		}
	case "remediation":
		t.Remediation = parseList(val)
	case "diagnostic_checkpoint":
		t.DiagnosticCheckpoint = strings.TrimSpace(val) == "true"
	case "lgwt_chapters":
		t.LgwtChapters = parseList(val)
	}
	return nil
}

func parseDrills(block string) ([]Drill, error) {
	lines := strings.Split(block, "\n")
	var drills []Drill
	var cur *Drill
	for _, ln := range lines {
		tr := strings.TrimSpace(ln)
		if tr == "" || tr == "drills:" {
			continue
		}
		if strings.HasPrefix(tr, "- ") {
			if cur != nil {
				drills = append(drills, *cur)
			}
			cur = &Drill{}
			tr = strings.TrimPrefix(tr, "- ")
		}
		if cur == nil {
			continue
		}
		key, val, ok := splitKeyValue(tr)
		if !ok {
			continue
		}
		switch key {
		case "id":
			cur.ID = unquote(val)
		case "target_seconds":
			n, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return nil, fmt.Errorf("drill target_seconds: %w", err)
			}
			cur.TargetSeconds = n
		case "prompt":
			cur.Prompt = unquote(val)
		}
	}
	if cur != nil {
		drills = append(drills, *cur)
	}
	return drills, nil
}

func splitKeyValue(raw string) (string, string, bool) {
	ln := strings.TrimLeft(raw, " ")
	colon := strings.Index(ln, ":")
	if colon < 0 {
		return "", "", false
	}
	key := strings.TrimSpace(ln[:colon])
	val := strings.TrimSpace(ln[colon+1:])
	return key, val, true
}

func parseList(val string) []string {
	val = strings.TrimSpace(val)
	if val == "" || val == "[]" {
		return nil
	}
	if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
		inner := strings.TrimSpace(val[1 : len(val)-1])
		if inner == "" {
			return nil
		}
		parts := strings.Split(inner, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			out = append(out, unquote(strings.TrimSpace(p)))
		}
		return out
	}
	return []string{unquote(val)}
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func leadingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r == ' ' {
			n++
			continue
		}
		break
	}
	return n
}
