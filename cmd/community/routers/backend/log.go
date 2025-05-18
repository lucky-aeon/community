package backend

import (
	"strconv"

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
	group.GET("/question/log", listQuestionLogs)
	group.GET("/courses/log", coursesLogs)
	group.GET("/courses/log/trend", coursesLogsTrend)
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

func listQuestionLogs(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	logs, count := logsS.GetPageQuestionLogs(p, limit)
	result.Page(logs, count, nil).Json(ctx)
	return
}

func coursesLogs(ctx *gin.Context) {
	coursesStats := logsS.GetCoursesStatistics()
	result.Ok(coursesStats, "").Json(ctx)
}

// coursesLogsTrend 获取课程访问时间序列数据
func coursesLogsTrend(ctx *gin.Context) {
	// 获取请求参数
	courseId, _ := strconv.Atoi(ctx.Query("courseId"))
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	granularity := ctx.Query("granularity")

	// 获取时间序列数据
	timeSeriesData := logsS.GetCoursesTimeSeries(courseId, startDate, endDate, granularity)

	// 返回数据
	result.Ok(timeSeriesData, "").Json(ctx)
}

type DeviceInfo struct {
	Device  string
	Browser string
}
