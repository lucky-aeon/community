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
	"xhyovo.cn/community/server/model"
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
	log.Info("start web :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Errorln(err)
	}
}

// 给所有用户的 expire_time 当前时间 + 1年
func initUserExpireTime() {

	// 计算一年后的时间
	oneYearLater := time.Now().AddDate(1, 0, 0)

	// 使用 Gorm 的 Update 方法更新所有记录
	result := model.User().Where("id > ?", 0).Update("expire_time", oneYearLater)
	if result.Error != nil {
		log.Infof("Failed to update expire_time: %v", result.Error)
	} else {
		log.Infof("Successfully updated %v user(s)", result.RowsAffected)
	}

}
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}
