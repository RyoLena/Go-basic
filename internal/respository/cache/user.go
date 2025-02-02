package cache

import (
	"Project/webBook_git/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrUserNotFound = errors.New("key 不存在")

type UserCache struct {
	//将接口定义成一个字段 ---- 面向接口编程
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client: client,
		//这一段可以通 client一样外部谁调用谁传进来
		expiration: time.Minute * 10,
	}
}

func (cache *UserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		fmt.Println("fail to marshal")
		return err
	}
	key := cache.key(user.ID)
	cache.client.Set(ctx, key, val, cache.expiration)
	return nil
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		fmt.Println("cache Get fail")
		return domain.User{}, ErrUserNotFound
	}
	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
