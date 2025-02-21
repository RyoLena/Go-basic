package auth

import (
	"Project/internal/service/ShortMessage"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	svc ShortMessage.Service
	key string
}

func (a AuthService) Sends(ctx context.Context, tpl string, args []string, number ...string) error {
	var authClaims AuthClaims
	token, err := jwt.ParseWithClaims(tpl, &authClaims, func(token *jwt.Token) (interface{}, error) {
		return a.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("token无效")
	}
	return a.svc.Sends(ctx, authClaims.TplID, args, number...)
}

type AuthClaims struct {
	jwt.RegisteredClaims
	TplID string
}
