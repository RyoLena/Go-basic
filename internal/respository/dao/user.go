package dao

import (
	"context"
	"database/sql"
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

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type UserDataAccess struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDataAccess {
	return &UserDataAccess{db: db}
}

func (ud *UserDataAccess) Insert(ctx context.Context, user User) error {
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

func (ud *UserDataAccess) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	if errors.Is(err, ErrUserNotFound) {
		return u, ErrUserNotFound
	}
	return u, err
}

func (ud *UserDataAccess) FindByID(ctx context.Context, id int64) (User, error) {
	user := User{
		ID: id,
	}
	return user, ud.db.WithContext(ctx).Where("email=?", user.ID).First(&user).Error
}

func (ud *UserDataAccess) FindByPhone(ctx context.Context, phone string) (User, error) {
	user := User{
		Phone: sql.NullString{
			String: phone,
			Valid:  phone != "",
		},
	}
	return user, ud.db.WithContext(ctx).Where("phone=?", user.Phone).First(&user).Error
}

type User struct {
	ID       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string

	Phone sql.NullString `gorm:"unique"`

	Ctime int64
	Utime int64
}
