package main

import (
	"fmt"
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
	// 设置程序使用中国时区
	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error loading China location:", err)
		return
	}
	time.Local = chinaLoc
	log.Init()
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

	err = r.Run(":8080")
	if err != nil {
		log.Errorln(err)
	}

	var users = make([]model.Users, 0)
	model.User().Find(&users)
	var orders = make([]model.Orders, 0)
	for i := range users {
		user := users[i]

		order := model.Orders{
			InviteCode:      user.InviteCode,
			Price:           70,
			Purchaser:       user.ID,
			AcquisitionType: 2,
			Creator:         13,
		}
		orders = append(orders, order)
	}
	model.Order().Save(orders)
}
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}
