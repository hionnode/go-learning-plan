package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"learning-plan/internal/progress"
)

// Placement quizzes are deliberately short (3–5 items per phase). The point
// is to signal "you already own this" — a full pretest would defeat the
// purpose. Phase 5 has no placement; Maelstrom is the judge.
var placementQuizzes = map[string]placementQuiz{
	"phase-0": {
		title:   "Phase 0 placement — backend mental model + Go toolchain",
		passPct: 0.8,
		items: []placementItem{
			{q: "When your browser loads example.com, DNS does what? (a) decrypts TLS (b) translates the hostname into an IP address (c) caches the HTML (d) compresses the response", answer: "b"},
			{q: "HTTP runs on top of which transport protocol? (a) UDP (b) TCP (c) ICMP (d) HTTP is its own transport", answer: "b"},
			{q: "A Go function that can fail typically returns: (a) (T, error) and never panics (b) T, and throws on error (c) Promise<T> (d) nothing — Go has exceptions", answer: "a"},
			{q: "The zero value of a Go int is: (a) nil (b) 0 (c) undefined (d) panics on read", answer: "b"},
			{q: "Why is there no Go REPL in daily workflow? (a) Go compiles ahead of time; tests are the scratchpad (b) Go does have one, called goplay (c) security policy (d) Go hates iteration", answer: "a"},
		},
		coverage: []string{"0.1-backend-mental-model", "0.2-go-toolchain", "0.3-types-zero-values", "0.4-errors-as-values", "0.5-first-http-hit"},
	},
	"phase-1": {
		title:   "Phase 1 placement — Go fundamentals",
		passPct: 0.8,
		items: []placementItem{
			{q: "Which command initializes a new Go module? (a) go init (b) go mod init (c) go new mod (d) go create", answer: "b"},
			{q: "The zero value of a map is: (a) an empty map (b) nil — reads ok, writes panic (c) undefined (d) a runtime error at declaration", answer: "b"},
			{q: "To check for a data race at test time you pass: (a) -race to go test (b) -verbose (c) -detect (d) nothing, Go handles it", answer: "a"},
			{q: "A Go interface is satisfied: (a) with implements keyword (b) by matching method set, implicitly (c) only by pointer receivers (d) never across packages", answer: "b"},
		},
		coverage: []string{"1.1-hello-world", "1.2-types-structs", "1.4-pointers-slices-maps", "1.5-goroutines-channels"},
	},
	"phase-2": {
		title:   "Phase 2 placement — HTTP & stdlib",
		passPct: 0.75,
		items: []placementItem{
			{q: "An HTTP request's status code is on the: (a) request line you send (b) response status line the server sends back (c) body (d) there are no status codes in HTTP/1.1", answer: "b"},
			{q: "To cancel an outbound HTTP call after 500ms you use: (a) time.AfterFunc (b) context.WithTimeout (c) http.Client.Timeout only (d) a goroutine + select", answer: "b"},
			{q: "Fan-out to N urls and cancel stragglers on first success is idiomatic with: (a) sync.WaitGroup (b) errgroup + context cancel (c) mutex + channel (d) time.Ticker", answer: "b"},
			{q: "Acceptance tests are distinguished from unit tests by: (a) using testify (b) black-box vs white-box scope (c) being slower (d) running in Docker", answer: "b"},
		},
		coverage: []string{"2.1-http-protocol", "2.2-net-http-server", "2.3-concurrency-patterns", "2.5-time-tickers", "2.6-testing-at-scale"},
	},
	"phase-3": {
		title:   "Phase 3 placement — backend building blocks",
		passPct: 0.75,
		items: []placementItem{
			{q: "In SQL, an index on a column primarily speeds up: (a) SELECT with WHERE on that column (b) INSERT (c) DELETE (d) every operation equally", answer: "a"},
			{q: "In database/sql, prepared statements help with: (a) SQL injection AND plan caching (b) only syntax highlighting (c) nothing (d) connection pooling only", answer: "a"},
			{q: "A Redis distributed lock with SETNX must also set: (a) a TTL (b) a password (c) a DB index (d) nothing else", answer: "a"},
			{q: "A message queue provides, at minimum: (a) exactly-once delivery for free (b) decoupling between producers and consumers (c) persistent storage (d) SQL querying", answer: "b"},
		},
		coverage: []string{"3.1-sql-fundamentals", "3.2-postgres", "3.3-caching-fundamentals", "3.4-redis", "3.5-queue-fundamentals", "3.6-nats"},
	},
	"phase-4": {
		title:   "Phase 4 placement — distsys concepts",
		passPct: 0.75,
		items: []placementItem{
			{q: "CAP theorem: under a network partition a system must sacrifice: (a) consistency or availability (b) latency or throughput (c) correctness or uptime (d) durability or security", answer: "a"},
			{q: "Vector clocks detect: (a) wall-clock drift (b) causal vs concurrent events (c) disk failures (d) process restarts", answer: "b"},
			{q: "A G-Counter's merge is correct because it is: (a) idempotent, commutative, associative (b) transactional (c) linearizable (d) lock-free", answer: "a"},
		},
		coverage: []string{"4.1-distsys-theory", "4.3-logical-clocks", "4.4-crdt"},
	},
}

type placementItem struct {
	q      string
	answer string
}

type placementQuiz struct {
	title    string
	passPct  float64
	items    []placementItem
	coverage []string // tasks to mark skipped if score >= passPct
}

func runPlace(ctx *appContext, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: tracker place <phase-id>  (e.g. phase-1)")
	}
	phaseID := args[0]
	quiz, ok := placementQuizzes[phaseID]
	if !ok {
		return fmt.Errorf("no placement quiz for %q", phaseID)
	}

	store, state, err := ctx.loadState()
	if err != nil {
		return err
	}

	fmt.Println(quiz.title)
	fmt.Println("type a, b, c, or d for each. type 'skip' to abandon.")
	fmt.Println()

	r := bufio.NewReader(os.Stdin)
	correct := 0
	for i, item := range quiz.items {
		fmt.Printf("Q%d. %s\n> ", i+1, item.q)
		ans, err := r.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading answer: %w", err)
		}
		ans = strings.ToLower(strings.TrimSpace(ans))
		if ans == "skip" {
			fmt.Println("abandoned.")
			return nil
		}
		if ans == item.answer {
			correct++
			fmt.Println("  ✓")
		} else {
			fmt.Printf("  ✗ (answer: %s)\n", item.answer)
		}
		fmt.Println()
	}

	score := float64(correct) / float64(len(quiz.items))
	result := &progress.PlacementResult{At: time.Now().UTC(), Score: score}
	fmt.Printf("score: %d/%d (%.0f%%)\n", correct, len(quiz.items), score*100)
	if score >= quiz.passPct {
		result.SkippedTasks = append(result.SkippedTasks, quiz.coverage...)
		// mark each covered task as Learning (placement-skipped) so the DAG doesn't block progress
		for _, id := range quiz.coverage {
			tp := state.TaskOrInit(id)
			if tp.Mastery == progress.Unseen {
				tp.Mastery = progress.Learning
			}
		}
		fmt.Printf("placement passed — marked %d tasks as placement-skipped.\n", len(quiz.coverage))
	} else {
		fmt.Println("below threshold — no tasks skipped. work through the phase normally.")
	}
	state.Placement[phaseID] = result
	return store.Save(state)
}
