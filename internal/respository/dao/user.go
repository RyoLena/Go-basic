package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("email has been registered")
	// ErrUserNotFound 邮箱不在数据库中
	ErrUserNotFound = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud *UserDao) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now

	//获取邮件冲突的信息
	err := ud.db.WithContext(ctx).Create(&user).Error
	var emilErr *mysql.MySQLError
	if errors.As(err, &emilErr) {
		var uniqueEmailErr = 1062
		if emilErr.Number == uint16(uniqueEmailErr) {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (ud *UserDao) FindByEmail(ctx context.Context, user User) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email=?", user.Email).First(&u).Error
	if errors.Is(err, ErrUserNotFound) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (ud *UserDao) FindByID(ctx context.Context, id int64) (User, error) {
	user := User{
		ID: id,
	}
	return user, ud.db.WithContext(ctx).Where("email=?", user.ID).First(&user).Error
}

type User struct {
	ID       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Ctime int64
	Utime int64
}
