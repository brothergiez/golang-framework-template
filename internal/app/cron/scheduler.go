package cron

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job interface {
	Name() string
	Run(ctx context.Context) error
}

type Scheduler interface {
	Register(spec string, job Job) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type entry struct {
	interval time.Duration
	job      Job
}

// IntervalScheduler is a lightweight scheduler using fixed intervals.
// Spec format: duration string, e.g. "1m", "30s", "1h".
type IntervalScheduler struct {
	mu      sync.Mutex
	entries []entry
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	running bool
}

func NewIntervalScheduler() *IntervalScheduler {
	return &IntervalScheduler{}
}

func (s *IntervalScheduler) Register(spec string, job Job) error {
	d, err := time.ParseDuration(spec)
	if err != nil {
		return fmt.Errorf("invalid cron spec %q: %w", spec, err)
	}
	if d <= 0 {
		return fmt.Errorf("cron interval must be positive")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry{interval: d, job: job})
	return nil
}

func (s *IntervalScheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("scheduler already running")
	}
	runCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.running = true
	entries := append([]entry(nil), s.entries...)
	s.mu.Unlock()

	for _, e := range entries {
		e := e
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ticker := time.NewTicker(e.interval)
			defer ticker.Stop()
			// run once at start
			_ = e.job.Run(runCtx)
			for {
				select {
				case <-runCtx.Done():
					return
				case <-ticker.C:
					_ = e.job.Run(runCtx)
				}
			}
		}()
	}

	<-runCtx.Done()
	s.wg.Wait()
	return nil
}

func (s *IntervalScheduler) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	cancel := s.cancel
	s.mu.Unlock()
	cancel()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
