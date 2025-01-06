package router

import (
	"Project/webBook_git/internal/respository"
	"Project/webBook_git/internal/respository/dao"
	"Project/webBook_git/internal/service"
	"Project/webBook_git/internal/web"
	"Project/webBook_git/internal/web/Middleware"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// UserGroutine is a function that handles all the user related routes
func UserGroutine(server *gin.Engine) {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:13306)/webook"))
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

	//使用session的方式
	//store := cookie.NewStore([]byte("secret"))
	//server.Use(sessions.Sessions("session_id", store))
	//server.Use(Middleware.Build().CheckLogin())

	//使用内存/redis保存sess_id的方式
	store, err := redis.NewStore(16, "tcp", "localhost:6379",
		"Ryo19120705",
		[]byte("sUvca2dpn7veAV4odb4xQNwYFV0EescZ"),
		[]byte("zYkJFgYaKEEgDTgQLLpomR028ZuQc6BE"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("session_id", store))
	//server.Use(Middleware.Build().
	//	IgnorePath("/user/signup").
	//	IgnorePath("/user/login").CheckLogin())

	//换成JWT实现的检测登录的中间件
	server.Use(Middleware.Build().
		IgnorePathJWT("/user/signup").
		IgnorePath("/user/login").CheckLoginJWT())

	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.LoginJWT)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
	}

}
