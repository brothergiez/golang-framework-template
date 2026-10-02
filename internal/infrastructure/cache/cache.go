package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Ping(ctx context.Context) error
	Close() error
}

type NoopCache struct{}

func NewNoop() *NoopCache {
	return &NoopCache{}
}

func (n *NoopCache) Get(context.Context, string) ([]byte, error) { return nil, ErrMiss }
func (n *NoopCache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}
func (n *NoopCache) Delete(context.Context, string) error { return nil }
func (n *NoopCache) Ping(context.Context) error           { return nil }
func (n *NoopCache) Close() error                         { return nil }

var ErrMiss = errMiss("cache miss")

type errMiss string

func (e errMiss) Error() string { return string(e) }
