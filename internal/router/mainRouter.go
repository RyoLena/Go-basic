package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		ExposeHeaders:    []string{"Content-Type", "Authorization"},
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
	//session

	UserGroutine(server)

	return server
}
