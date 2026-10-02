package redis

import (
	"context"
	"time"

	"github.com/lion/go-framework/internal/config"
	"github.com/lion/go-framework/internal/infrastructure/cache"
	goredis "github.com/redis/go-redis/v9"
)

type Cache struct {
	client *goredis.Client
}

func New(cfg config.RedisConfig) (*Cache, error) {
	opts, err := goredis.ParseURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	opts.DB = cfg.Database
	client := goredis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return &Cache{client: client}, nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == goredis.Nil {
		return nil, cache.ErrMiss
	}
	return val, err
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
