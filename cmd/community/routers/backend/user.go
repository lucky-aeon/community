package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitUserRouters(r *gin.Engine) {

	group := r.Group("/community/admin/user")
	group.GET("", listUser)
	group.POST("", updateUser)
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
