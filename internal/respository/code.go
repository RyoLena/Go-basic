package respository

import (
	"Project/internal/respository/cache"
	"context"
)

var ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
var ErrCodeNotFound = cache.ErrUnKnownCode

type CodeCodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
}

type CodeStorage struct {
	cacheCode cache.CodeCache
}

func NewCodeRepo(cacheCode cache.CodeCache) CodeCodeRepository {
	return &CodeStorage{
		cacheCode: cacheCode,
	}
}

func (repo *CodeStorage) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cacheCode.Set(ctx, biz, phone, code)
}

func (repo *CodeStorage) Verify(ctx context.Context, biz, phone, code string) error {
	return repo.cacheCode.Verify(ctx, biz, phone, code)
}
