package respository

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/respository/dao"
	"context"
)

var (
	RepErrUserDuplicated = dao.ErrUserDuplicateEmail
	RepoErrUserNotFound  = dao.ErrUserNotFound
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

func (repo *UserRepo) FindByEmail(ctx context.Context, user domain.User) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, dao.User{
		Email: user.Email,
	})
	return domain.User{
		Email:    u.Email,
		Password: u.Password,
	}, err
}
