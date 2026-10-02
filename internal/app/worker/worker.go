package worker

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID      string
	Type    string
	Payload []byte
}

type Consumer interface {
	Handle(ctx context.Context, message Message) error
}

type Worker interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type InMemoryQueue struct {
	ch chan Message
}

func NewInMemoryQueue(size int) *InMemoryQueue {
	if size <= 0 {
		size = 100
	}
	return &InMemoryQueue{ch: make(chan Message, size)}
}

func (q *InMemoryQueue) Publish(msg Message) {
	q.ch <- msg
}

func (q *InMemoryQueue) Chan() <-chan Message {
	return q.ch
}

type Runner struct {
	consumer    Consumer
	queue       *InMemoryQueue
	concurrency int
	poll        time.Duration
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex
	running     bool
}

func NewRunner(consumer Consumer, queue *InMemoryQueue, concurrency int, poll time.Duration) *Runner {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &Runner{
		consumer:    consumer,
		queue:       queue,
		concurrency: concurrency,
		poll:        poll,
	}
}

func (r *Runner) Start(ctx context.Context) error {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return fmt.Errorf("worker already running")
	}
	runCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	r.running = true
	r.mu.Unlock()

	for i := 0; i < r.concurrency; i++ {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case <-runCtx.Done():
					return
				case msg, ok := <-r.queue.Chan():
					if !ok {
						return
					}
					_ = r.consumer.Handle(runCtx, msg)
				}
			}
		}()
	}

	<-runCtx.Done()
	r.wg.Wait()
	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return nil
	}
	cancel := r.cancel
	r.mu.Unlock()
	cancel()

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
