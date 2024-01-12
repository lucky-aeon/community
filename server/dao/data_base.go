package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"xhyovo.cn/community/server/config"
)

// dao层使用
var db *gorm.DB

func InitDb(config *config.DbConfig) {

	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"

	dsn := fmt.Sprintf(d, config.Username, config.Password, config.Address, config.Database)
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db = db1
}
