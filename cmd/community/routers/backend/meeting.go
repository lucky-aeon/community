package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/request"
	services "xhyovo.cn/community/server/service"
)

func InitMeetingRouters(r *gin.Engine) {
	group := r.Group("/community/admin/meeting")
	group.Use(middleware.OperLogger())
	group.POST("/approve", approve)
	group.POST("/pass", pass)
	group.POST("/record", record)
	group.DELETE("/:id", deleteMeeting)
}

func approve(ctx *gin.Context) {
	var reqProveMeeting request.ReqApproveMeeting
	if err := ctx.ShouldBindJSON(&reqProveMeeting); err != nil {
		msg := utils.GetValidateErr(reqProveMeeting, err)
		result.Err(msg).Json(ctx)
		return
	}

	var meetingService services.MeetingService
	if err := meetingService.Approve(reqProveMeeting); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "审核通过").Json(ctx)
}

func pass(ctx *gin.Context) {
	var reqPassMeeting request.ReqPassMeeting
	if err := ctx.ShouldBindJSON(&reqPassMeeting); err != nil {
		msg := utils.GetValidateErr(reqPassMeeting, err)
		result.Err(msg).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	if err := meetingService.Pass(reqPassMeeting); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "已PASS").Json(ctx)
}

func record(ctx *gin.Context) {
	var reqRecordMeeting request.ReqRecordMeeting
	if err := ctx.ShouldBindJSON(&reqRecordMeeting); err != nil {
		msg := utils.GetValidateErr(&reqRecordMeeting, err)
		result.Err(msg).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	if err := meetingService.Record(reqRecordMeeting); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "已记录").Json(ctx)
}

func deleteMeeting(ctx *gin.Context) {
	idInt, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	if err = meetingService.DeleteById(idInt, 0); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
