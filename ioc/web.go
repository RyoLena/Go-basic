package ioc

import (
	"Project/internal/web"
	"Project/internal/web/Middleware"
	myjwt "Project/internal/web/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandle) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.Register(server)
	return server
}

func InitMiddlewares(jwtHdl myjwt.Handle) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Recovery(),
		func(ctx *gin.Context) {
			log.Printf("请求进入中间件链: %s", ctx.FullPath())
			ctx.Next()
		},
		corsHdl(),
		Middleware.NewLoginJWTMiddleware(jwtHdl).
			IgnorePathJWT("/user/signup").
			IgnorePathJWT("/user/login").
			IgnorePathJWT("/user/login_sms/code/send").
			IgnorePathJWT("/user/login_sms").
			IgnorePathJWT("/user/refresh_token").BuildJWT(),

		//全局限流在这里

	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		ExposeHeaders: []string{"Content-Type", "Authorization",
			"x-jwt-token", "x-refresh-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "localhost") {
				//这里添加开发环境

				return true
			}
			return strings.Contains(origin, "your-company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
