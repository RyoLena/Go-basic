package router

import "github.com/gin-gonic/gin"

// MainRouter 管理着所有的router
func MainGroute() *gin.Engine {
	server := gin.Default()
	UserGroutine(server)
	return server
}
