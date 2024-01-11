package dao

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xhyovo.cn/community/server/config"
)

var DB *sql.DB

func InitDb(config *config.DbConfig) {

	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"

	dsn := fmt.Sprintf(d, config.Username, config.Password, config.Address, config.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	// 配置 sql.DB 连接池参数
	// ref: https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(25)                 // 设置最大的并发连接数（in-use + idle）
	db.SetMaxIdleConns(25)                 // 设置最大的空闲连接数（idle）
	db.SetConnMaxLifetime(5 * time.Minute) // 设置连接的最大生命周期

	// 检查数据库连接
	if err := db.Ping(); err != nil {
		panic(err.Error())
	}
	DB = db
}
