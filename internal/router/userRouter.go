package router

import (
	"Project/webBook_git/internal/web"
	"github.com/gin-gonic/gin"
)

// UserGroutine is a function that handles all the user related routes
func UserGroutine(server *gin.Engine) {
	user := web.NewUserHandle()
	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.Login)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
	}
}
