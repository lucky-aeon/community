package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"log"
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/utils"
)

func main() {

	r := gin.Default()
	r.SetFuncMap(utils.GlobalFunc())
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("./views/**/*")

	store := cookie.NewStore([]byte("123456"))
	// 添加 session
	r.Use(sessions.Sessions("user", store))
	// 添加 auth

	config.Init()
	mysql.Init(&config.GetInstance().DbConfig)
	kodo.Init(&config.GetInstance().KodoConfig)

	routers.InitFrontedRouter(r)
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
