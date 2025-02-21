package ioc

import (
	"Project/config"
	"github.com/redis/go-redis/v9"
)

func InitRedisDBS() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
	})
	return rdb
}
