package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitLogRouters(r *gin.Engine) {
	group := r.Group("/community/admin")
	group.GET("/oper/log", listLogs)
	group.DELETE("/oper/log", deleteLogs)
}

func listLogs(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	logSearch := model.LogSearch{}
	if err := ctx.ShouldBindQuery(&logSearch); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	if (logSearch.StartTime != "" && logSearch.EndTime == "") || (logSearch.StartTime == "" && logSearch.EndTime != "") {
		result.Err("选择范围时间，开始时间和结束时间必须同时有值").Json(ctx)
		return
	}
	var logsS services.LogServices
	logs, count := logsS.GetPageOperLog(p, limit, logSearch)
	result.Ok(page.New(logs, count), "").Json(ctx)
}

func deleteLogs(ctx *gin.Context) {
	var ids []int
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var logsS services.LogServices
	logsS.DeletesOperLogs(ids)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
