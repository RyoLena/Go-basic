//go:build wireinject

package usersWire

import (
	"Project/internal/DB"
	"Project/internal/respository"
	"Project/internal/respository/cache"
	"Project/internal/respository/dao"
	"Project/internal/service"
	"Project/internal/web"
	"github.com/google/wire"
)

func InitWebService() *web.UserHandle {
	wire.Build(
		DB.InitDB, DB.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache, cache.NewCodeCache, // 提供 *cache.CodeRedisCache
		respository.NewUserRepo, respository.NewCodeRepo, // 提供 *respository.CodeRepository
		service.NewUserService, service.NewCodeService, // 提供 *service.CodeServiceImpl
		InitFakerSMS,
		web.NewUserHandle, // 提供 *web.UserHandle
	)
	return new(web.UserHandle)
}
