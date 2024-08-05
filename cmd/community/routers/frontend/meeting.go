package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/request"
	services "xhyovo.cn/community/server/service"
)

func InitMeetingRouters(ctx *gin.Engine) {

	group := ctx.Group("/community/meeting")
	group.GET("", pageMeeting)
	group.GET("/:id", getMeeting)
	group.GET("/manager", managerMeeting)
	group.GET("/inMeetingState/:id", inMeetingState)
	group.GET("/:id/inMeetingUserAvatar", getJoinMeetingUsers)

	group.Use(middleware.OperLogger())
	group.DELETE(":id", deleteMeeting)
	group.POST("/:id/join", joinMeeting)
	group.POST("/:id/quit", quitMeeting)
	group.POST("/apply", applyMeeting)
}

func managerMeeting(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	userId := middleware.GetUserId(ctx)
	var meetingService services.MeetingService
	meetings, count := meetingService.Page(p, limit, userId)
	result.Page(meetings, count, nil).Json(ctx)
}

// 申请会议
func applyMeeting(ctx *gin.Context) {
	var reqMeeting request.ReqMeeting
	if err := ctx.ShouldBindJSON(&reqMeeting); err != nil {
		msg := utils.GetValidateErr(&reqMeeting, err)
		result.Err(msg).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	var meeting model.Meetings
	meeting.Id = reqMeeting.Id
	meeting.Title = reqMeeting.Title
	meeting.Description = reqMeeting.Description
	meeting.InitiatorId = middleware.GetUserId(ctx)
	meeting.InitiatorTime = reqMeeting.InitiatorTime
	if err := meetingService.Save(meeting); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "申请成功,等待管理员审核").Json(ctx)
}

func pageMeeting(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var meetingService services.MeetingService
	meetings, count := meetingService.Page(p, limit, 0)
	result.Page(meetings, count, nil).Json(ctx)
}
func getMeeting(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	meeting := meetingService.GetById(idInt)
	result.Ok(meeting, "").Json(ctx)

}

func deleteMeeting(ctx *gin.Context) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	var meetingService services.MeetingService
	if err = meetingService.DeleteById(idInt, userId); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 加入会议
func joinMeeting(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	userId := middleware.GetUserId(ctx)
	if err = meetingService.JoinMeeting(id, userId); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "加入成功").Json(ctx)
}

func quitMeeting(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	userId := middleware.GetUserId(ctx)
	if err = meetingService.QuitJoinMeeting(id, userId); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "退出成功").Json(ctx)
}

func inMeetingState(ctx *gin.Context) {

	idInt, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var meetingService services.MeetingService
	userId := middleware.GetUserId(ctx)
	state, err := meetingService.InMeetingState(idInt, userId)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Ok(state, "").Json(ctx)
}

func getJoinMeetingUsers(ctx *gin.Context) {
	idInt, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error())
		return
	}
	var meetingService services.MeetingService
	avatars := meetingService.GetJoinMeetingUserSelectAvatar(idInt)

	result.Ok(avatars, "").Json(ctx)
	return
}
