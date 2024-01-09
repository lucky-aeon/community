package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"xhyovo.cn/community/config"
)

var db gorm.DB

func InitDb(dbconfig config.Db) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)%scharset=utf8mb4&parseTime=True&loc=Local",
		dbconfig.Username, dbconfig.Password, dbconfig.Address, dbconfig.Database)
	dbObject, err := gorm.Open(mysql.Open(dsn), &gorm.Config{

		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		panic(err.Error())
	}
	db = *dbObject

}
