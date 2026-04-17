package graph

import (
	"fmt"
	"html"
	"io"
	"strings"

	"learning-plan/internal/progress"
)

// Layout constants — tuned to fit the full 34-node curriculum on a single page
// without overlap. Phases run left-to-right; within a phase, tasks stack
// vertically in declaration order.
const (
	nodeW    = 200
	nodeH    = 56
	colW     = 260
	rowH     = 80
	padX     = 40
	padY     = 40
	arrowLen = 8
)

// Color per mastery level — readable against dark-mode cards, also ok on light.
var masteryColor = map[progress.Mastery]string{
	progress.Unseen:     "#3a3f4b",
	progress.Learning:   "#d4a017",
	progress.Proficient: "#4a9a6e",
	progress.Automatic:  "#2d6cdf",
}

// RenderSVG writes the full DAG as one SVG. Caller supplies the progress state
// so nodes can be colored by mastery.
func (d *DAG) RenderSVG(w io.Writer, state *progress.State) error {
	type pos struct{ x, y int }
	positions := map[string]pos{}

	// Column per phase (phase 0 → column 0); row by declaration order within phase.
	minPhase, maxPhase := 0, 0
	for _, id := range d.Nodes {
		p := d.Phase[id]
		if p < minPhase {
			minPhase = p
		}
		if p > maxPhase {
			maxPhase = p
		}
	}
	rowByPhase := map[int]int{}
	for _, id := range d.Nodes {
		phase := d.Phase[id]
		col := phase - minPhase
		row := rowByPhase[phase]
		rowByPhase[phase] = row + 1
		positions[id] = pos{
			x: padX + col*colW,
			y: padY + row*rowH,
		}
	}

	totalCols := maxPhase - minPhase + 1
	maxRows := 0
	for _, r := range rowByPhase {
		if r > maxRows {
			maxRows = r
		}
	}
	width := padX*2 + totalCols*colW
	height := padY*2 + maxRows*rowH

	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="100%%" style="font-family: ui-monospace, Menlo, monospace; font-size: 12px;">`, width, height)
	fmt.Fprintln(w, `<defs>
  <marker id="arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
    <path d="M 0 0 L 10 5 L 0 10 z" fill="#6b7380"/>
  </marker>
</defs>`)

	// Edges first so nodes sit on top.
	for _, id := range d.Nodes {
		p := positions[id]
		for _, pr := range d.Prereqs[id] {
			q, ok := positions[pr]
			if !ok {
				continue
			}
			x1 := q.x + nodeW
			y1 := q.y + nodeH/2
			x2 := p.x
			y2 := p.y + nodeH/2
			fmt.Fprintf(w, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#6b7380" stroke-width="1.2" marker-end="url(#arrow)" opacity="0.55"/>`, x1, y1, x2, y2)
		}
	}

	for _, id := range d.Nodes {
		p := positions[id]
		tp := state.Tasks[id]
		mast := progress.Unseen
		if tp != nil && tp.Mastery != "" {
			mast = tp.Mastery
		}
		fill := masteryColor[mast]
		if fill == "" {
			fill = "#3a3f4b"
		}
		fmt.Fprintf(w, `<a href="/task/%s"><g>`, html.EscapeString(id))
		fmt.Fprintf(w, `<rect x="%d" y="%d" width="%d" height="%d" rx="8" fill="%s" stroke="#1a1d23" stroke-width="1"/>`, p.x, p.y, nodeW, nodeH, fill)
		fmt.Fprintf(w, `<text x="%d" y="%d" fill="#f5f5f7" font-weight="600">%s</text>`, p.x+12, p.y+22, html.EscapeString(id))
		fmt.Fprintf(w, `<text x="%d" y="%d" fill="#d7d9de" opacity="0.85">%s</text>`, p.x+12, p.y+40, html.EscapeString(truncate(titleOrID(id, d), 26)))
		fmt.Fprintf(w, `</g></a>`)
	}

	// Column labels — one per phase actually present.
	for p := minPhase; p <= maxPhase; p++ {
		col := p - minPhase
		x := padX + col*colW
		fmt.Fprintf(w, `<text x="%d" y="%d" fill="#8c909a" font-size="11" letter-spacing="1.5">PHASE %d</text>`, x, padY-14, p)
	}

	_, err := fmt.Fprintln(w, `</svg>`)
	return err
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return strings.TrimSpace(s[:n-1]) + "…"
}

func titleOrID(id string, d *DAG) string {
	// Fallback: strip phase prefix and dashes.
	// The DAG doesn't store titles directly; keep it simple.
	_ = d
	if i := strings.Index(id, "-"); i > 0 && i < len(id)-1 {
		return strings.ReplaceAll(id[i+1:], "-", " ")
	}
	return id
}
