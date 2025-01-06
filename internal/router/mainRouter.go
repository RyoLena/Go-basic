package router

import (
	"Project/webBook_git/internal/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// MainGroute MainRouter 管理着所有的router
func MainGroute() *gin.Engine {
	server := gin.Default()
	//跨域问题
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		ExposeHeaders:    []string{"Content-Type", "Authorization", "x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "localhost") {
				//这里添加开发环境

				return true
			}
			return strings.Contains(origin, "your-company.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	//限流
	redisClien := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "Ryo19120705",
	})
	server.Use(ratelimit.NewBuilder(redisClien, time.Minute, 20).Build())

	UserGroutine(server)

	return server
}
