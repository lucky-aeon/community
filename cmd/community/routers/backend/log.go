package backend

import (
	"github.com/gin-gonic/gin"
	"strings"
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

	logs, count := logsS.GetPageOperLog(p, limit, logSearch)
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

type DeviceInfo struct {
	Device  string
	Browser string
}

func parseUserAgent(userAgent string) *DeviceInfo {
	deviceInfo := &DeviceInfo{}

	// 解析设备类型
	if strings.Contains(userAgent, "Mobile") {
		deviceInfo.Device = "Mobile"
	} else if strings.Contains(userAgent, "Tablet") {
		deviceInfo.Device = "Tablet"
	} else {
		deviceInfo.Device = "Desktop"
	}

	// 解析浏览器信息
	if strings.Contains(userAgent, "Chrome") {
		deviceInfo.Browser = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		deviceInfo.Browser = "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		deviceInfo.Browser = "Safari"
	} else if strings.Contains(userAgent, "Edge") {
		deviceInfo.Browser = "Edge"
	} else {
		deviceInfo.Browser = "Unknown"
	}

	return deviceInfo
}
