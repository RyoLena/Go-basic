package Middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

func Build() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要校验的路由
		if ctx.Request.URL.Path == "/user/signup" ||
			ctx.Request.URL.Path == "/user/login" {
			ctx.Next()
		}
		//校验
		sess := sessions.Default(ctx)
		if sess.Get("userID") == nil {
			ctx.JSON(401, gin.H{
				"msg": "未登录",
			})
			ctx.Abort()
		}
	}
}
