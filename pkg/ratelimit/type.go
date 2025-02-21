package ratelimit

import "context"

type Limiter interface {
	// Limit 通过key来限流，key--->ip地址
	//bool 判断是否开始限流
	Limit(ctx context.Context, key string) (bool, error)
}
