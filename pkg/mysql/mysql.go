package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"xhyovo.cn/community/pkg/config"
)

var instance *gorm.DB

func GetInstance() *gorm.DB {
	return instance
}

func Init(db *config.DbConfig) {

	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	var err error
	dsn := fmt.Sprintf(d, db.Username, db.Password, db.Address, db.Database)
	instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err.Error())
	}
}
