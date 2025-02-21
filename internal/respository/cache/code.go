package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var (
	ErrSetCodeFreq            = errors.New("发送验证码太平凡")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnKnownCode            = errors.New("我也不知道是什么错误，在Code这里")
)

// 通过go的语法嵌入lua脚本到luaSetCode中
//
//go:embed Lua/set_code.lua
var luaSetCode string

//go:embed Lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, expeditedCode string) error
}

type CodeRedisCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	if client == nil {
		log.Fatal("Redis client is nil")
	}
	return &CodeRedisCache{
		client: client,
	}
}

func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int64()

	if err != nil {
		fmt.Println("出错了嘛")
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrSetCodeFreq
	default:
		return errors.New("系统出错")
	}
}

func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, expeditedCode string) error {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, expeditedCode).Int64()

	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeVerifyTooManyTimes
	case -2:
		return nil
	default:
		return ErrUnKnownCode
	}

}

func (c *CodeRedisCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s,%s", biz, phone)
}
