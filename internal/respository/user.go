package respository

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/respository/dao"
	"context"
)

type UserRepo struct {
	dao *dao.UserDao
}

func NewUserRepo(db *dao.UserDao) *UserRepo {
	return &UserRepo{dao: db}
}

func (repo *UserRepo) Create(ctx context.Context, user domain.User) error {
	//存数据
	return repo.dao.Insert(ctx, dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
