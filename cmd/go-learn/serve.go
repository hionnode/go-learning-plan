package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"learning-plan/internal/curriculum"
	"learning-plan/internal/graph"
	"learning-plan/internal/progress"
	"learning-plan/internal/srs"
)

//go:embed templates/*.html
var templateFS embed.FS

func runServe(ctx *appContext, args []string) error {
	addr := ":8080"
	if len(args) > 0 {
		addr = args[0]
	}

	srv := &server{ctx: ctx}
	if err := srv.loadTemplates(); err != nil {
		return fmt.Errorf("loading templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", srv.handleDashboard)
	mux.HandleFunc("/graph", srv.handleGraph)
	mux.HandleFunc("/review", srv.handleReview)
	mux.HandleFunc("/task/", srv.handleTask)
	mux.HandleFunc("/drill", srv.handleDrillIndex)
	mux.HandleFunc("/drill/", srv.handleDrillDetail)

	httpSrv := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("go-learn listening on http://localhost%s\n", addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-stop:
		fmt.Println("\nshutting down…")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

type server struct {
	ctx   *appContext
	pages map[string]*template.Template
}

func (s *server) loadTemplates() error {
	funcs := template.FuncMap{
		"masteryColor": func(m progress.Mastery) string {
			switch m {
			case progress.Learning:
				return "#d4a017"
			case progress.Proficient:
				return "#4a9a6e"
			case progress.Automatic:
				return "#2d6cdf"
			default:
				return "#3a3f4b"
			}
		},
		"masteryLabel": func(m progress.Mastery) string {
			if m == "" {
				return "unseen"
			}
			return string(m)
		},
		"joinCSV": func(ss []string) string { return strings.Join(ss, ", ") },
		"humanDuration": func(d time.Duration) string {
			if d < 0 {
				d = -d
			}
			if d < time.Minute {
				return fmt.Sprintf("%ds", int(d.Seconds()))
			}
			if d < time.Hour {
				return fmt.Sprintf("%dm", int(d.Minutes()))
			}
			if d < 24*time.Hour {
				return fmt.Sprintf("%.1fh", d.Hours())
			}
			return fmt.Sprintf("%.1fd", d.Hours()/24)
		},
		"fmtTime": func(t *time.Time) string {
			if t == nil {
				return "—"
			}
			return t.Local().Format("2006-01-02 15:04")
		},
		"relTime": func(t *time.Time) string {
			if t == nil {
				return "—"
			}
			d := time.Until(*t)
			if d < 0 {
				return fmt.Sprintf("%s overdue", roundDuration(-d))
			}
			return "in " + roundDuration(d)
		},
		"phasePill": func(p int) string { return fmt.Sprintf("P%d", p) },
		"percent":   func(n, d int) int { if d == 0 { return 0 }; return 100 * n / d },
		"timePtr":   func(t time.Time) *time.Time { return &t },
		"msToSec": func(ms int64) string {
			if ms <= 0 {
				return "—"
			}
			return fmt.Sprintf("%.1fs", float64(ms)/1000)
		},
	}
	// Each page is built as its own template set: layout.html + one page.
	// Sharing a single set would collide on the "content" define.
	pages := []string{
		"dashboard.html",
		"graph.html",
		"review.html",
		"task.html",
		"drill.html",
		"drill_detail.html",
	}
	s.pages = map[string]*template.Template{}
	for _, p := range pages {
		t, err := template.New("layout").Funcs(funcs).ParseFS(templateFS, "templates/layout.html", "templates/"+p)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", p, err)
		}
		s.pages[p] = t
	}
	return nil
}

type pageData struct {
	Title      string
	ActiveNav  string
	Curriculum *curriculum.Curriculum
	State      *progress.State
	Body       any
}

func (s *server) render(w http.ResponseWriter, name string, data pageData) {
	tmpl, ok := s.pages[name]
	if !ok {
		http.Error(w, "unknown template: "+name, http.StatusInternalServerError)
		return
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}

func (s *server) load() (*curriculum.Curriculum, *progress.Store, *progress.State, error) {
	c, err := s.ctx.loadCurriculum()
	if err != nil {
		return nil, nil, nil, err
	}
	store, state, err := s.ctx.loadState()
	return c, store, state, err
}

type dashboardView struct {
	Phases       []phaseSummary
	Focus        *curriculum.Task
	FocusState   *progress.TaskProgress
	DueCount     int
	TotalTasks   int
	MasteredPct  int
	LearningPct  int
	AutomaticPct int
}

type phaseSummary struct {
	Number      int
	Tasks       []taskRow
	Counts      map[progress.Mastery]int
	PercentDone int
}

type taskRow struct {
	ID        string
	Title     string
	Mastery   progress.Mastery
	NextDue   *time.Time
	Prereqs   []string
	Ready     bool
	Phase     int
}

func (s *server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	c, _, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	dag, err := graph.Build(c.Tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	counts := map[progress.Mastery]int{}
	phases := map[int][]taskRow{}
	for _, t := range c.Tasks {
		tp := state.Tasks[t.ID]
		m := progress.Unseen
		var next *time.Time
		if tp != nil {
			if tp.Mastery != "" {
				m = tp.Mastery
			}
			next = tp.NextReviewAt
		}
		counts[m]++
		phases[t.Phase] = append(phases[t.Phase], taskRow{
			ID:      t.ID,
			Title:   t.Title,
			Mastery: m,
			NextDue: next,
			Prereqs: t.Prereqs,
			Ready:   dag.Ready(t.ID, state),
			Phase:   t.Phase,
		})
	}
	var phaseNums []int
	for n := range phases {
		phaseNums = append(phaseNums, n)
	}
	sort.Ints(phaseNums)
	var phaseSummaries []phaseSummary
	for _, n := range phaseNums {
		rows := phases[n]
		c := map[progress.Mastery]int{}
		for _, row := range rows {
			c[row.Mastery]++
		}
		done := c[progress.Proficient] + c[progress.Automatic]
		phaseSummaries = append(phaseSummaries, phaseSummary{
			Number:      n,
			Tasks:       rows,
			Counts:      c,
			PercentDone: 100 * done / len(rows),
		})
	}

	focusID := dag.NextFocus(state)
	var focus *curriculum.Task
	var focusState *progress.TaskProgress
	if focusID != "" {
		focus = c.TaskByID(focusID)
		focusState = state.Tasks[focusID]
	}

	due := srs.DueTasks(state, time.Now().UTC())
	total := len(c.Tasks)

	view := dashboardView{
		Phases:       phaseSummaries,
		Focus:        focus,
		FocusState:   focusState,
		DueCount:     len(due),
		TotalTasks:   total,
		MasteredPct:  100 * (counts[progress.Proficient] + counts[progress.Automatic]) / total,
		LearningPct:  100 * counts[progress.Learning] / total,
		AutomaticPct: 100 * counts[progress.Automatic] / total,
	}

	s.render(w, "dashboard.html", pageData{
		Title:      "learning-plan",
		ActiveNav:  "dashboard",
		Curriculum: c,
		State:      state,
		Body:       view,
	})
}

type graphView struct {
	SVG template.HTML
}

func (s *server) handleGraph(w http.ResponseWriter, r *http.Request) {
	c, _, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	dag, err := graph.Build(c.Tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var buf bytes.Buffer
	if err := dag.RenderSVG(&buf, state); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.render(w, "graph.html", pageData{
		Title:      "prereq graph",
		ActiveNav:  "graph",
		Curriculum: c,
		State:      state,
		Body:       graphView{SVG: template.HTML(buf.String())},
	})
}

type reviewView struct {
	Items []reviewItem
}

type reviewItem struct {
	ID      string
	Title   string
	Mastery progress.Mastery
	Due     *time.Time
	Late    time.Duration
}

func (s *server) handleReview(w http.ResponseWriter, r *http.Request) {
	c, _, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	now := time.Now().UTC()
	ids := srs.DueTasks(state, now)
	var items []reviewItem
	for _, id := range ids {
		task := c.TaskByID(id)
		title := id
		if task != nil {
			title = task.Title
		}
		tp := state.Tasks[id]
		items = append(items, reviewItem{
			ID:      id,
			Title:   title,
			Mastery: tp.Mastery,
			Due:     tp.NextReviewAt,
			Late:    now.Sub(*tp.NextReviewAt),
		})
	}
	s.render(w, "review.html", pageData{
		Title:      "review queue",
		ActiveNav:  "review",
		Curriculum: c,
		State:      state,
		Body:       reviewView{Items: items},
	})
}

type taskView struct {
	Task       curriculum.Task
	State      *progress.TaskProgress
	DAG        *graph.DAG
	Ready      bool
	AllDrills  []curriculum.Drill
	DrillState map[string]*progress.DrillProgress
}

func (s *server) handleTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/task/")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	c, store, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	task := c.TaskByID(id)
	if task == nil {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		tp := state.TaskOrInit(id)
		tp.Reflections = append(tp.Reflections, progress.Reflection{
			At:      time.Now().UTC(),
			Built:   strings.TrimSpace(r.FormValue("built")),
			Clicked: strings.TrimSpace(r.FormValue("clicked")),
			Fuzzy:   strings.TrimSpace(r.FormValue("fuzzy")),
		})
		if err := store.Save(state); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/task/"+id, http.StatusSeeOther)
		return
	}

	dag, _ := graph.Build(c.Tasks)
	var drills []curriculum.Drill
	for _, did := range task.DrillIDs {
		if d := c.DrillByID(did); d != nil {
			drills = append(drills, *d)
		}
	}

	view := taskView{
		Task:       *task,
		State:      state.TaskOrInit(id),
		DAG:        dag,
		Ready:      dag.Ready(id, state),
		AllDrills:  drills,
		DrillState: state.Drills,
	}
	s.render(w, "task.html", pageData{
		Title:      task.Title,
		ActiveNav:  "dashboard",
		Curriculum: c,
		State:      state,
		Body:       view,
	})
}

type drillIndexView struct {
	Drills  []curriculum.Drill
	Records map[string]*progress.DrillProgress
}

func (s *server) handleDrillIndex(w http.ResponseWriter, r *http.Request) {
	c, _, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.render(w, "drill.html", pageData{
		Title:      "drills",
		ActiveNav:  "drill",
		Curriculum: c,
		State:      state,
		Body: drillIndexView{
			Drills:  c.Drills,
			Records: state.Drills,
		},
	})
}

type drillDetailView struct {
	Drill   curriculum.Drill
	Record  *progress.DrillProgress
	Message string
}

func (s *server) handleDrillDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/drill/")
	if id == "" {
		http.Redirect(w, r, "/drill", http.StatusFound)
		return
	}
	c, store, state, err := s.load()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	d := c.DrillByID(id)
	if d == nil {
		http.NotFound(w, r)
		return
	}

	message := ""
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		elapsedMs, err := parseMs(r.FormValue("elapsedMs"))
		if err != nil {
			http.Error(w, "bad elapsedMs: "+err.Error(), 400)
			return
		}
		target := time.Duration(d.TargetSeconds) * time.Second
		met := time.Duration(elapsedMs)*time.Millisecond <= target
		dp := state.DrillOrInit(id)
		dp.History = append(dp.History, progress.DrillAttempt{
			At:         time.Now().UTC(),
			DurationMs: elapsedMs,
			MetTarget:  met,
		})
		if dp.BestMs == 0 || elapsedMs < dp.BestMs {
			dp.BestMs = elapsedMs
		}
		if met {
			for _, t := range c.Tasks {
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
			http.Error(w, err.Error(), 500)
			return
		}
		if met {
			message = fmt.Sprintf("logged %.1fs — met target.", float64(elapsedMs)/1000)
		} else {
			message = fmt.Sprintf("logged %.1fs — over target (%ds).", float64(elapsedMs)/1000, d.TargetSeconds)
		}
	}

	s.render(w, "drill_detail.html", pageData{
		Title:      d.ID,
		ActiveNav:  "drill",
		Curriculum: c,
		State:      state,
		Body: drillDetailView{
			Drill:   *d,
			Record:  state.Drills[id],
			Message: message,
		},
	})
}

func parseMs(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, time.Since(start))
	})
}
