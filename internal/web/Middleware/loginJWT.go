package Middleware

import (
	myjwt "Project/internal/web/jwt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
)

type LoginMiddlewareBuilderJWT struct {
	ignorePath []string
	jwtHdl     myjwt.Handle
}

func (l *LoginMiddlewareBuilderJWT) IgnorePathJWT(path string) *LoginMiddlewareBuilderJWT {
	l.ignorePath = append(l.ignorePath, path)
	return l
}

func NewLoginJWTMiddleware(jwtHdl myjwt.Handle) *LoginMiddlewareBuilderJWT {
	return &LoginMiddlewareBuilderJWT{
		jwtHdl: jwtHdl,
	}
}

func (l *LoginMiddlewareBuilderJWT) BuildJWT() gin.HandlerFunc {
	fmt.Printf("BuildJWT 中 jwtHdl 地址: %p\n", l.jwtHdl)
	return func(ctx *gin.Context) {
		//不需要校验的路由
		for _, path := range l.ignorePath {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		if l.jwtHdl == nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Println("jwtHdl 未初始化")
			return
		}
		tokenHead := l.jwtHdl.ExtractToken(ctx)
		if tokenHead == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Head_未登录"))
			return
		}
		//在这里拿到token
		claims := &myjwt.UserClaims{}

		parse, err := jwt.ParseWithClaims(tokenHead, claims, func(token *jwt.Token) (interface{}, error) {
			return myjwt.AtKey, nil
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
		err = l.jwtHdl.CheckSession(ctx, claims.SsID)
		if err != nil {
			_ = ctx.AbortWithError(401, fmt.Errorf("redis出错或者退出登录"))
			return
		}
		ctx.Set("claims", claims)
	}
}
