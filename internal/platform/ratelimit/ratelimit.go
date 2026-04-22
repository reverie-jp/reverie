// Package ratelimit wraps Redis SETNX with TTL into a tiny "claim a key for
// cooldown" primitive. Used to prevent notification spam from rapid follow /
// unfollow toggles without depending on DB state that callers may delete.
package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter interface {
	// TryAcquire returns true if the key was previously free and is now
	// reserved for ttl. Returns false when the key is still held (caller
	// should treat the action as rate-limited). Errors are transport-level
	// only — callers typically fail open (proceed as if acquired) to avoid
	// hard-coupling to Redis uptime.
	TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

type redisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) Limiter {
	return &redisLimiter{client: client}
}

func (l *redisLimiter) TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return l.client.SetNX(ctx, key, 1, ttl).Result()
}
