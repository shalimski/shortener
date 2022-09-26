package cache

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v9"
	"github.com/shalimski/shortener/config"
	"github.com/shalimski/shortener/internal/domain"
	"github.com/shalimski/shortener/internal/ports"
)

type cache struct {
	rdb *redis.Client
}

// NewCache create instance of redis server
func NewCache(cfg *config.Config) ports.Cacher {
	return &cache{
		rdb: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.DSN,
			Password: cfg.Redis.Password,
			DB:       0,
		}),
	}
}

// Set value by key
func (c *cache) Set(ctx context.Context, key, value string) error {
	return c.rdb.Set(ctx, key, value, 0).Err()
}

// Get value by key
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

// Delete value by key
func (c *cache) Del(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
