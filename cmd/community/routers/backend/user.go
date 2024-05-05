package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitUserRouters(r *gin.Engine) {

	group := r.Group("/community/admin/user")
	group.GET("", listUser)
	group.Use(middleware.OperLogger())
	group.POST("", updateUser)
	group.DELETE("/:id", deleteUser)
}

func listUser(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)

	var u services.UserService

	users, count := u.PageUsers(p, limit)
	result.Page(users, count, nil).Json(ctx)
}

type updateUserInfo struct {
	model.Users
	Tags []int `json:"tags"`
}

func updateUser(ctx *gin.Context) {

	user := updateUserInfo{}

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	var u services.UserService
	u.UpdateUser(&user.Users)
	user.Tags = userTagS.AssignUserLabel(user.ID, user.Tags)

	result.OkWithMsg(user, "修改成功").Json(ctx)
}

func deleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	if userId == id {
		result.Err("不能删除自己").Json(ctx)
		return
	}
	var u services.UserService
	u.DeleteUser(id)
	log.Infof("用户id: %d,删除用户: %d", userId, id)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
