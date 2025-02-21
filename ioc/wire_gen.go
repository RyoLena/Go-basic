// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ioc

import (
	"Project/internal/respository"
	"Project/internal/respository/cache"
	"Project/internal/respository/dao"
	"Project/internal/service"
	"Project/internal/web"
	"Project/internal/web/jwt"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func InitWebService() *gin.Engine {
	cmdable := InitRedisDBS()
	handle := jwt.NewRedisJWTHandle(cmdable)
	v := InitMiddlewares(handle)
	db := InitDBS()
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(cmdable)
	userRepository := respository.NewUserRepo(userDao, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewCodeCache(cmdable)
	codeCodeRepository := respository.NewCodeRepo(codeCache)
	shortMessageService := InitFakeSMS()
	codeService := service.NewCodeService(codeCodeRepository, shortMessageService)
	userHandle := web.NewUserHandle(userService, codeService, handle)
	engine := InitGin(v, userHandle)
	return engine
}
