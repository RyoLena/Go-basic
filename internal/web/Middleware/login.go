package Middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	ignorePath []string
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{
		ignorePath: append(l.ignorePath, path),
	}
}

func Build() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		//不需要校验的路由
		for _, path := range l.ignorePath {
			if ctx.Request.URL.Path == path {
				ctx.Next()
			}
		}
		//校验
		sess := sessions.Default(ctx)
		id := sess.Get("userID")
		if id == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "未登录",
			})
			ctx.Abort()
		}

		//延长时间
		update := sess.Get("update_time")
		sess.Set("userID,", id)
		sess.Options(sessions.Options{
			MaxAge: 30,
		})
		err := sess.Save()
		if err != nil {
			fmt.Println("会进这里么")
			return
		}
		now := time.Now()
		if update == nil {
			sess.Set("update_time", now)
			sess.Options(sessions.Options{
				MaxAge: 30,
			})
			err = sess.Save()
			if err != nil {
				fmt.Println("Session保存失败")
				return
			}
		}
		updatetimeVal := update.(time.Time)
		if now.Sub(updatetimeVal) > time.Second*20 {
			sess.Set("update_time", now)
			err = sess.Save()
			if err != nil {
				fmt.Println("这里的session保存错误")
				return
			}
		}
	}
}
