package DB

import (
	"Project/config"
	"Project/internal/respository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	userDB, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(userDB)
	if err != nil {
		panic(err)
	}
	return userDB
}

func InitRedis() redis.Cmdable {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
	})
	return rdb
}
