package router

import (
	"Project/webBook_git/internal/respository"
	"Project/webBook_git/internal/respository/dao"
	"Project/webBook_git/internal/service"
	"Project/webBook_git/internal/web"
	"Project/webBook_git/internal/web/Middleware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	err = dao.InitTable(db)
	if err != nil {
		fmt.Println("服务器创建失败")
		return
	}
	daoDb := dao.NewUserDao(db)
	repo := respository.NewUserRepo(daoDb)
	svc := service.NewUserService(repo)
	user := web.NewUserHandle(svc)

	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("session_id", store))
	server.Use(Middleware.Build().CheckLogin())

	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.Login)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
	}

}
