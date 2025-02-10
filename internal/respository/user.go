package respository

import (
	"Project/internal/domain"
	"Project/internal/respository/cache"
	"Project/internal/respository/dao"
	"context"
	"database/sql"
	"errors"
)

var (
	RepErrUserDuplicated = dao.ErrUserDuplicateEmail
	RepoErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByID(ctx context.Context, id int64) (domain.User, error)
}

type UserStorage struct {
	dao   *dao.UserDataAccess
	cache *cache.UserRedisCache
}

func NewUserRepo(db *dao.UserDataAccess, c *cache.UserRedisCache) *UserStorage {
	return &UserStorage{
		dao:   db,
		cache: c,
	}
}

func (repo *UserStorage) Create(ctx context.Context, user domain.User) error {
	//存数据
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}

func (repo *UserStorage) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	return repo.entityToDomain(u), err
}

func (repo *UserStorage) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	return repo.entityToDomain(u), err
}

func (repo *UserStorage) FindByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := repo.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	if errors.Is(err, cache.ErrUserNotFound) {
		//去数据库中查找
	}
	dbUser, err := repo.dao.FindByID(ctx, id)
	u = repo.entityToDomain(dbUser)
	go func() {
		err := repo.cache.Set(ctx, u)
		if err != nil {
			//打个日志做监控就行
			//缓存失败不一定是redis崩溃 也有timeout
		}
	}()
	return u, err
}

func (repo *UserStorage) entityToDomain(ud dao.User) domain.User {
	return domain.User{
		ID:       ud.ID,
		Email:    ud.Email.String,
		Phone:    ud.Phone.String,
		Password: ud.Password,
		Ctime:    ud.Ctime,
	}
}

func (repo *UserStorage) domainToEntity(ud domain.User) dao.User {
	return dao.User{
		ID: ud.ID,
		Email: sql.NullString{
			String: ud.Email,
			Valid:  ud.Email != "",
		},
		Phone: sql.NullString{
			String: ud.Phone,
			Valid:  ud.Phone != "",
		},
		Password: ud.Password,
		Ctime:    ud.Ctime,
	}
}
