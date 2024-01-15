package dao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(username, password, address, database string) *gorm.DB {

	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"

	dsn := fmt.Sprintf(d, username, password, address, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	return db
}
