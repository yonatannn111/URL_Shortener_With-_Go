package worker

import (
	"fmt"
	"sync"
	"time"

	"github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
)

// Worker manages analytics for shortened URLs
type Worker struct {
	store  *storage.Store
	stop   chan bool
	mu     sync.Mutex
	clicks map[string]int // stores click counts for each short code
}

// NewWorker initializes a new Worker
func NewWorker(store *storage.Store) *Worker {
	return &Worker{
		store:  store,
		stop:   make(chan bool),
		clicks: make(map[string]int),
	}
}

// Start runs a background analytics job (optional, can be extended)
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.stop:
				return
			default:
				// Example analytics job; can extend to persist analytics to DB
				fmt.Println("worker running background analytics...")
				time.Sleep(10 * time.Second)
			}
		}
	}()
}

// Stop stops the background worker
func (w *Worker) Stop() {
	w.stop <- true
}

// AddAnalytics records a click/analytics event for a given code
func (w *Worker) AddAnalytics(code, ip, country, city string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Increment click count safely
	w.clicks[code]++

	// Log analytics; safe to show empty values for localhost
	fmt.Printf("Analytics: code=%s, ip=%s, country=%s, city=%s\n", code, ip, country, city)
}

// GetClicks returns the total clicks for a given code
func (w *Worker) GetClicks(code string) int {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Return 0 if code does not exist
	if count, ok := w.clicks[code]; ok {
		return count
	}
	return 0
}
