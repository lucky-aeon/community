package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/utils"
)

func main() {

	r := gin.Default()
	r.SetFuncMap(utils.GlobalFunc())
	// r.Static("/assets", "assets")
	// r.LoadHTMLGlob("views/**/*")

	// 添加 auth

	config.Init()
	db := &config.GetInstance().DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)
	kodo.Init(&config.GetInstance().KodoConfig)

	routers.InitFrontedRouter(r)
	err := r.Run("127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
}
