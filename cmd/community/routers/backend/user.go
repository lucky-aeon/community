package backend

import (
	"github.com/gin-gonic/gin"
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
	result.Ok(page.New(users, count), "").Json(ctx)
}

func updateUser(ctx *gin.Context) {

	user := model.Users{}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	var u services.UserService
	u.UpdateUser(&user)
}
