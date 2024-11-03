package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var logsS services.LogServices

func InitLogRouters(r *gin.Engine) {
	group := r.Group("/community/admin")
	group.GET("/oper/log", listOperLogs)
	group.GET("/login/log", listLoginLogs)
	group.GET("/file/log", listFileLogs)
	group.GET("/question/log", listQuestionLOgs)
}

func listOperLogs(ctx *gin.Context) {
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

	logs, count := logsS.GetPageOperLog(p, limit, logSearch, true)
	result.Page(logs, count, nil).Json(ctx)
}

func listLoginLogs(ctx *gin.Context) {
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
	logs, count := logsS.GetPageLoginPage(p, limit, logSearch)
	result.Page(logs, count, nil).Json(ctx)
	return
}

func listFileLogs(ctx *gin.Context) {
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
	logs, count := logsS.GetPageOperLog(p, limit, logSearch, false)
	result.Page(logs, count, nil).Json(ctx)
	return
}

func listQuestionLOgs(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	logs, count := logsS.GetPageQuestionLogs(p, limit)
	result.Page(logs, count, nil).Json(ctx)
	return
}

type DeviceInfo struct {
	Device  string
	Browser string
}
