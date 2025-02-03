package service

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/respository"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var SVCErrUserDuplicated = respository.RepErrUserDuplicated
var ErrInvalidUserOrPassword = errors.New("账号或者密码不对")

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
	//用bcrypt加密
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPwd)
	//2.存储
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	//1.查询
	u, err := svc.repo.FindByEmail(ctx, user.Email)
	//用户不存在
	if errors.Is(err, respository.RepoErrUserNotFound) {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	//查询出错
	if err != nil {
		return domain.User{}, err
	}

	//2.密码比对
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	//密码不对
	if err != nil {
		fmt.Println(u.Password, user.Password)
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if !errors.Is(err, respository.RepoErrUserNotFound) {
		return domain.User{}, err
	}

	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})

	if err != nil {
		return u, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}
