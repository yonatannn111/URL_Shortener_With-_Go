package worker

import (
	"context"
	"log"
)

type ClickEvent struct {
	URLID  int
	Code   string
	IP     string
	Country string
	City   string
	UserAgent string
	Referer   string
}

type Worker struct {
	Store    *storage.Store
	ClickCh  chan ClickEvent
	Ctx      context.Context
	Cancel   context.CancelFunc
}

func NewWorker(store *storage.Store) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		Store: store,
		ClickCh: make(chan ClickEvent, 1000),
		Ctx: ctx, Cancel: cancel,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case ev := <-w.ClickCh:
				err := w.Store.InsertClick(w.Ctx, ev.URLID, ev.Code, ev.IP, ev.Country, ev.City, ev.UserAgent, ev.Referer)
				if err != nil {
					log.Printf("worker: failed to insert click: %v", err)
				}
			case <-w.Ctx.Done():
				return
			}
		}
	}()
}
