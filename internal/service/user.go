package service

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/respository"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repo     *respository.UserRepo
	Email    string
	Password string
}

func NewUserService(repo *respository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx *gin.Context, u domain.User) error {
	//1.密码加密
	//2.存储
	return svc.repo.Create(ctx, u)
}
