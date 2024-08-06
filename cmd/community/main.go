package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"time"
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/oss"
	"xhyovo.cn/community/pkg/utils"
)

func main() {

	log.Init()
	// 设置程序使用中国时区
	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	time.Local = chinaLoc
	if err != nil {
		log.Errorf("Error loading China location:", err)
		return
	}
	r := gin.Default()
	r.SetFuncMap(utils.GlobalFunc())
	config.Init()
	appConfig := config.GetInstance()
	db := appConfig.DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)
	ossConfig := appConfig.OssConfig
	oss.Init(ossConfig.Endpoint, ossConfig.AccessKey, ossConfig.SecretKey, ossConfig.Bucket)
	emailConfig := appConfig.EmailConfig
	email.Init(emailConfig.Address, emailConfig.Username, emailConfig.Password, emailConfig.Host, emailConfig.PollCount)
	routers.InitFrontedRouter(r)
	cache.Init()
	log.Info("start web")
	err = r.Run(appConfig.ServerBind)
	if err != nil {
		log.Errorln(err)
	}
}
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}
