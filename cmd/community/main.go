package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/utils"
)

func main() {

	r := gin.Default()
	r.SetFuncMap(utils.GlobalFunc())
	config.Init()
	appConfig := config.GetInstance()
	db := appConfig.DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)
	kodoConfig := appConfig.KodoConfig
	kodo.Init(kodoConfig.AccessKey, kodoConfig.SecretKey, kodoConfig.Bucket, kodoConfig.Domain)
	emailConfig := appConfig.EmailConfig
	email.Init(emailConfig.Address, emailConfig.Username, emailConfig.Password, emailConfig.Host, emailConfig.PollCount)
	routers.InitFrontedRouter(r)
	err := r.Run("127.0.0.1:8080")
	if err != nil {
		log.Fatalln(err)
	}
}
