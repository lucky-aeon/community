package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitIndexRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	// 社区首页
	group.GET("/user/count", getUserCount)
	group.GET("/rate/page", pageRate)
	group.GET("/labels/", getKnowedgeLabels)
}

func pageRate(ctx *gin.Context) {

	p, limit := page.GetPage(ctx)
	var noteService services.RateService
	state, notes := noteService.Page(p, limit)

	result.Ok(map[string]interface{}{
		"state": state,
		"data":  notes,
	}, "").Json(ctx)
	return
}

func getUserCount(ctx *gin.Context) {
	var count int64
	model.User().Count(&count)
	result.Ok(count, "").Json(ctx)
}

func getKnowedgeLabels(ctx *gin.Context) {

	labels := []string{"javase", "juc", "jvm", "mysql", "redis", "mq", "多线程", "反射", "字节码", "设计模式", "spring", "springmvc", "mybatis", "springboot", "dubbo", "分布式", "微服务", "zookeeper", "计算机网络", "操作系统"}
	result.Ok(labels, "").Json(ctx)
}
