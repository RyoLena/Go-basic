package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handle interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractToken(ctx *gin.Context) string
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	SsID string
	Uid  int64
}

type UserClaims struct {
	jwt.RegisteredClaims
	//想要获取什么从这里添加
	Uid  int64
	SsID string
	//Email     string
	UserAgent string
}
