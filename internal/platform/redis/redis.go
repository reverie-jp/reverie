package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"reverie.jp/reverie/internal/config"
)

type Client = redis.Client

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return c, nil
}
