package Middleware

import (
	"Project/webBook_git/internal/web"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilderJWT struct {
	ignorePath []string
}

func (l *LoginMiddlewareBuilder) IgnorePathJWT(path string) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		ignorePath: append(l.ignorePath, path),
	}
}

func BuildJWT() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) CheckLoginJWT() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要校验的路由
		for _, path := range l.ignorePath {
			if ctx.Request.URL.Path == path {
				ctx.Next()
				return
			}
		}
		tokenHead := ctx.GetHeader("Authorization")
		if tokenHead == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Head_未登录"))
			return
		}
		//在这里拿到token
		segs := strings.Split(tokenHead, " ")
		claims := &web.UserClaims{}
		if len(segs) != 2 || segs[0] != "Bearer" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Len_未登录"))
			return
		}
		token := segs[1]
		parse, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("sUvca2dpn7veAV4odb4xQNwYFV0EescZ"), nil
		})
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("parse_未登录"))
			return
		}
		if parse == nil || !parse.Valid || claims.Uid == -1 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("parse.Valid_未登录"))
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			_ = ctx.AbortWithError(401, fmt.Errorf("UserAgent-未登录"))
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*10 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * 30))
			tokenStr, err := parse.SignedString([]byte("sUvca2dpn7veAV4odb4xQNwYFV0EescZ"))
			if err != nil {
				log.Println("刷新token失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
	}
}
