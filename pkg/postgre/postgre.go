package postgre

import (
	"fmt"
	"gorm.io/driver/postgres"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		address, username, password, database)
	var err error
	instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Errorf("连接 postgres 数据库失败！,err: %s", err.Error())
		//panic("连接 postgres 数据库失败！")
	}

	db, err := instance.DB()
	db.SetMaxIdleConns(10)  //空闲连接数
	db.SetMaxOpenConns(100) //最大连接数
	db.SetConnMaxLifetime(time.Minute)

}
