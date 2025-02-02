package respository

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/respository/cache"
	"Project/webBook_git/internal/respository/dao"
	"context"
	"errors"
)

var (
	RepErrUserDuplicated = dao.ErrUserDuplicateEmail
	RepoErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepo struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepo(db *dao.UserDao, c *cache.UserCache) *UserRepo {
	return &UserRepo{
		dao:   db,
		cache: c,
	}
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

func (repo *UserRepo) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	if errors.Is(err, cache.ErrUserNotFound) {
		//去数据库中查找
	}
	dbUser, err := repo.dao.FindByID(ctx, id)
	u = domain.User{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}
	go func() {
		err := repo.cache.Set(ctx, u)
		if err != nil {
			//打个日志做监控就行
			//缓存失败不一定是redis崩溃 也有timeout
		}
	}()
	return u, err
}
