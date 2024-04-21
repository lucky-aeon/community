package main

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
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

	// 获取所有评论
	var comments []model.Comments
	model.Comment().Find(&comments)
	set := mapset.NewSet[int]()
	for i := range comments {
		set.Add(comments[i].BusinessId)
	}
	// 获取所有文章
	var articles []model.Articles
	model.Article().Where("id in ?", set.ToSlice()).Find(&articles)
	var m = make(map[int]int)
	for i := range articles {
		m[articles[i].ID] = articles[i].UserId
	}
	// 建立关系
	for i := range comments {
		comments[i].BusinessUserId = m[comments[i].BusinessId]
	}

	model.Comment().Where("business_id in ?", set.ToSlice()).Save(&comments)
	err := r.Run(":8080")
	if err != nil {
		log.Errorln(err)
	}
}
