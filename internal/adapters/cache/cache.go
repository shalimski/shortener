package cache

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
	"github.com/shalimski/shortener/config"
	"github.com/shalimski/shortener/internal/domain"
)

type cache struct {
	rdb *redis.Client
}

func NewCache(cfg *config.Config) *cache {
	return &cache{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.DSN,
			Password: cfg.Redis.Password,
			DB:       0,
		}),
	}
}

func (c *cache) Set(ctx context.Context, key, value string) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	url, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrNotFound
	}

	if err != nil {
		return "", err
	}

	return url, nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

func (c *cache) Shutdown(ctx context.Context) {
	c.rdb.Shutdown(ctx)
}
