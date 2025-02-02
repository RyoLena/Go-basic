package respository

import (
	"Project/webBook_git/internal/respository/cache"
	"context"
)

var ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
var ErrCodeNotFound = cache.ErrUnKnownCode

type CodeRepo struct {
	cacheCode cache.CodeCache
}

func NewCodeRepo(cacheCode cache.CodeCache) *CodeRepo {
	return &CodeRepo{
		cacheCode: cacheCode,
	}
}

func (repo *CodeRepo) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cacheCode.Set(ctx, biz, phone, code)
}

func (repo *CodeRepo) Verify(ctx context.Context, biz, phone, code string) error {
	return repo.cacheCode.Verify(ctx, biz, phone, code)
}
