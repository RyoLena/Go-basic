package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var slideWindowScript string

type RedisSlideWindowLimiter struct {
	// TODO
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlideWindowLimiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	vals, err := r.cmd.Eval(ctx, slideWindowScript, []string{key},
		r.interval, r.rate).Bool()
	if err != nil {
		return false, err
	}
	return vals, nil
}
