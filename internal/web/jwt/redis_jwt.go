package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	AtKey = []byte("sUvca2dpn7veAV4odb4xQNwYFV0EescZ")
	RtKey = []byte("sUvca2dpn7veAV4odb4xQNwYFV0EescT")
)

type RedisJWTHandle struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandle(cmd redis.Cmdable) Handle {
	return RedisJWTHandle{
		cmd: cmd,
	}
}

func (r RedisJWTHandle) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, uid, ssid)
	return err
}

func (r RedisJWTHandle) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		SsID:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		ctx.String(http.StatusOK, "jwt加密系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r RedisJWTHandle) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	c, ok := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return errors.New("系统错误")
	}
	err := r.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.SsID), "", time.Hour*24*7).Err()

	return err
}

func (r RedisJWTHandle) CheckSession(ctx *gin.Context, ssid string) error {
	cnt, err := r.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	if err != nil || cnt > 0 {
		_ = ctx.AbortWithError(401, fmt.Errorf("redis出错或者退出登录"))
		return err
	}
	return nil
}

func (r RedisJWTHandle) ExtractToken(ctx *gin.Context) string {
	tokenHead := ctx.GetHeader("Authorization")
	segs := strings.SplitN(tokenHead, " ", 2)
	if len(segs) != 2 || segs[0] != "Bearer" {
		return ""
	}
	return segs[1]
}

func (r RedisJWTHandle) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		SsID: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		ctx.String(http.StatusOK, "jwt加密系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}
