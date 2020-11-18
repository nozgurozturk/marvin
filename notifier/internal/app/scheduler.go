package app

import (
	"context"
	"sync"
	"time"
)

type Job func(ctx context.Context)

type Scheduler struct {
	wg      *sync.WaitGroup
	cancels []context.CancelFunc
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:      new(sync.WaitGroup),
		cancels: make([]context.CancelFunc, 0),
	}
}

func (s *Scheduler) Add(ctx context.Context, j Job, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancels = append(s.cancels, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, interval)
}

func (s *Scheduler) process(ctx context.Context, j Job, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			j(ctx)
		case <-ctx.Done():
			s.wg.Done()
			return
		}
	}
}
