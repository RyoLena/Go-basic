//go:build wireinject

package usersWire

import (
	"Project/webBook_git/internal/DB"
	"Project/webBook_git/internal/respository"
	"Project/webBook_git/internal/respository/cache"
	"Project/webBook_git/internal/respository/dao"
	"Project/webBook_git/internal/service"
	"Project/webBook_git/internal/web"
	"github.com/google/wire"
)

func InitWebService() *web.UserHandle {
	wire.Build(DB.InitDB, DB.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache,
		respository.NewUserRepo,
		service.NewUserService,
		//InitFakerSMS,
		web.NewUserHandle,
	)
	return new(web.UserHandle)
}
