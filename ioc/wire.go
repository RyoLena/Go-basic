//go:build wireinject

package ioc

import (
	"Project/internal/respository"
	"Project/internal/respository/cache"
	"Project/internal/respository/dao"
	"Project/internal/service"
	"Project/internal/web"
	myjwt "Project/internal/web/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebService() *gin.Engine {
	wire.Build(
		InitDBS,
		InitRedisDBS,
		dao.NewUserDao,
		cache.NewUserCache, cache.NewCodeCache, // 提供 *cache.CodeRedisCache
		respository.NewUserRepo, respository.NewCodeRepo, // 提供 *respository.CodeRepository
		service.NewUserService, service.NewCodeService, // 提供 *service.CodeServiceImpl
		InitFakeSMS,
		web.NewUserHandle, // 提供 *web.UserHandle
		myjwt.NewRedisJWTHandle,
		InitGin,
		InitMiddlewares,
	)
	return new(gin.Engine)
}
