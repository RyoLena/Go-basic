package router

import (
	"Project/webBook_git/internal/respository"
	"Project/webBook_git/internal/respository/dao"
	"Project/webBook_git/internal/service"
	"Project/webBook_git/internal/web"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// UserGroutine is a function that handles all the user related routes
func UserGroutine(server *gin.Engine) {

	db, err := gorm.Open(mysql.Open("root:210912@tcp(127.0.0.1:13306)/webook"))
	if err != nil {
		panic(err)
	}
	daoDb := dao.NewUserDao(db)
	repo := respository.NewUserRepo(daoDb)
	svc := service.NewUserService(repo)
	user := web.NewUserHandle(svc)

	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.Login)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
	}
}
