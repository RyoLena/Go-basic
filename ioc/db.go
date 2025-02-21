package ioc

import (
	"Project/config"
	"Project/internal/respository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDBS() *gorm.DB {
	userDB, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		return nil
	}
	err = dao.InitTable(userDB)
	if err != nil {
		return nil
	}
	return userDB
}
