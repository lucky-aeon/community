package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitMemberRouters(r *gin.Engine) {
	group := r.Group("/community/admin/member")
	group.GET("", listMembers)
	group.POST("", saveMember)
	group.DELETE("/:id", deleteMember)
}

func listMembers(ctx *gin.Context) {
	var m services.MemberInfoService
	members := m.ListMember()
	result.Ok(page.New(members, int64(len(members))), "").Json(ctx)
}

func saveMember(ctx *gin.Context) {
	var m services.MemberInfoService
	var member model.MemberInfos
	if err := ctx.ShouldBindJSON(&member); err != nil {
		log.Warnf("用户id: %d 添加等级参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	m.SaveMember(&member)
	result.OkWithMsg(nil, "操作成功").Json(ctx)
}

func deleteMember(ctx *gin.Context) {
	var m services.MemberInfoService
	id := ctx.Param("id")
	atoi, _ := strconv.Atoi(id)
	if err := m.DeleteMember(atoi); err != nil {
		log.Warnf("用户id: %d 删除等级参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "操作成功").Json(ctx)
}
