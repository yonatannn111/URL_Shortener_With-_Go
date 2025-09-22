package worker

import (
	"fmt"
	"time"

	"github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
)

type Worker struct {
	store *storage.Store
	stop  chan bool
}

func NewWorker(store *storage.Store) *Worker {
	return &Worker{store: store, stop: make(chan bool)}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.stop:
				return
			default:
				// Example analytics job
				fmt.Println("worker running background analytics...")
				time.Sleep(10 * time.Second)
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.stop <- true
}
