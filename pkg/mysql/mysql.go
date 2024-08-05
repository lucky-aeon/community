package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
	"xhyovo.cn/community/pkg/log"
)

var instance *gorm.DB

func GetInstance() *gorm.DB {
	return instance
}

func Init(username, password, address, database string) {
	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=10s"
	var err error
	dsn := fmt.Sprintf(d, username, password, address, database)
	instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Info),
	})
	db, err := instance.DB()
	db.SetMaxIdleConns(10)  //空闲连接数
	db.SetMaxOpenConns(100) //最大连接数
	db.SetConnMaxLifetime(time.Minute)
	if err != nil {
		panic("连接数据库失败！")
	}
	if err != nil {
		log.Errorf("初始化 db 失败,err: %s", err.Error())
		panic(err.Error())
	}
}
