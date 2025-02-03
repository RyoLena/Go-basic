package router

import (
	"Project/webBook_git/config"
	"Project/webBook_git/internal/web/Middleware"
	"Project/webBook_git/ioc/usersWire"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// UserGroutine is a function that handles all the user related routes
func UserGroutine(server *gin.Engine) {
	user := usersWire.InitWebService()

	//使用session的方式
	//store := cookie.NewStore([]byte("secret"))
	//server.Use(sessions.Sessions("session_id", store))
	//server.Use(Middleware.Build().CheckLogin())

	//使用内存/redis保存sess_id的方式
	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr,
		config.Config.Redis.Password,
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
	server.Use(Middleware.BuildJWT().
		IgnorePathJWT("/user/signup").
		IgnorePathJWT("/user/login").
		IgnorePathJWT("/user/login_sms/code/send").
		IgnorePathJWT("/login_sms").CheckLoginJWT())

	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.LoginJWT)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
		userHandle.POST("/login_sms/code/send", user.SendSMSCode)
		userHandle.POST("/login_sms", user.LoginBySMS)
	}

}
