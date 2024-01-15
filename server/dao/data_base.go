package dao

import (
	"fmt"
	"gorm.io/gorm/schema"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init(username, password, address, database string) {

	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	var err error
	dsn := fmt.Sprintf(d, username, password, address, database)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	if err != nil {
		panic(err.Error())
	}

}
