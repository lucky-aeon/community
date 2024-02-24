package frontend

import (
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"

	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

var (
	userService services.UserService
)

type editUserForm struct {
	Name string `binding:"required" form:"name" msg:"用户名不可为空"`
	Desc string `form:"desc"`
}

type editPasswordForm struct {
	OldPassword     string `form:"oldPassword" binding:"required" msg:"旧密码不能为空"`
	NewPassword     string `form:"newPassword" binding:"required" msg:"新密码不能为空"`
	ConfirmPassword string `form:"confirmPassword" binding:"required" msg:"确认密码不能为空"`
}

func InitUserRouters(r *gin.Engine) {
	group := r.Group("/community/user")
	group.GET("/info", getUserInfo)
	group.POST("/edit/:tab", updateUser)
	group.GET("/menu", getUserMenu)
	group.GET("/statistics", statistics)
}

func getUserMenu(ctx *gin.Context) {
	result.Ok(userService.GetUserMenu(), "ok").Json(ctx)
}

// 获取用户信息
func getUserInfo(ctx *gin.Context) {

	var userService services.UserService
	userId := ctx.Query("userId")
	uId, err := strconv.Atoi(userId)
	if err != nil {
		uId = middleware.GetUserId(ctx)
	}
	user := userService.GetUserById(uId)

	result.Ok(user, "").Json(ctx)
}

func updateUser(ctx *gin.Context) {

	userId := middleware.GetUserId(ctx)
	t := ctx.Param("tab")
	switch t {
	case "info":
		form := editUserForm{}
		err := ctx.ShouldBind(&form)

		if err != nil {
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		if len(form.Desc) > 200 {
			result.Err("描述长度不可超过200字").Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Name: form.Name, Desc: form.Desc, ID: userId})
	case "pass":
		form := editPasswordForm{}
		err := ctx.ShouldBind(&form)
		if err != nil {
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		// check 旧密码
		if form.OldPassword != userService.GetUserById(userId).Password {
			result.Err("旧密码不一致").Json(ctx)
			return
		}
		// check 新密码
		if form.NewPassword != form.ConfirmPassword {
			result.Err("两次新密码不一致").Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Password: form.ConfirmPassword, ID: userId})
	case "avatar":
		type avatar struct {
			Avatar string `json:"avatar" binding:"required" msg:"头像不能为空"`
		}
		object := &avatar{}
		if err := ctx.ShouldBindJSON(&object); err != nil {
			result.Err(utils.GetValidateErr(object, err)).Json(ctx)
			return
		}
		// 更改用户信息
		userService.UpdateUser(&model.Users{ID: userId, Avatar: object.Avatar})
	}
	result.OkWithMsg(nil, "修改成功").Json(ctx)
}

// 数据统计
func statistics(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)

	m := userService.Statistics(userId)
	result.Ok(m, "").Json(ctx)
}
