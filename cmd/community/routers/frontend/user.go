package frontend

import (
	"strconv"

	"xhyovo.cn/community/pkg/log"

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
	Name      string `binding:"required" form:"name" msg:"用户名不可为空"`
	Desc      string `form:"desc"`
	Subscribe int    `form:"subscribe"`
}

type editPasswordForm struct {
	OldPassword     string `form:"oldPassword" binding:"required" msg:"旧密码不能为空"`
	NewPassword     string `form:"newPassword" binding:"required" msg:"新密码不能为空"`
	ConfirmPassword string `form:"confirmPassword" binding:"required" msg:"确认密码不能为空"`
}

func InitUserRouters(r *gin.Engine) {
	group := r.Group("/community/user")
	group.GET("/info", getUserInfo)
	group.GET("/menu", getUserMenu)
	group.GET("/statistics", statistics)
	group.GET("", listUsers)
	group.GET("/tags/:userId", getTagsByUserId)
	group.Use(middleware.OperLogger())
	group.POST("/edit/:tab", updateUser)
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
	user := userService.GetUserSimpleById(uId)

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
			log.Warnf("用户id: %d 修改信息参数解析失败,err: %s", userId, err.Error())
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		if len(form.Desc) > 200 {
			var msg = "描述长度不可超过200字"
			log.Warnf("用户id: %d 修改信息失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Name: form.Name, Desc: form.Desc, ID: userId, Subscribe: form.Subscribe})
	case "pass":
		form := editPasswordForm{}
		err := ctx.ShouldBind(&form)
		if err != nil {
			log.Warnf("用户id: %d 修改密码参数解析失败,err: %s", userId, err.Error())
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		// check 旧密码
		if form.OldPassword != userService.GetUserById(userId).Password {
			var msg = "旧密码不一致"
			log.Warnf("用户id: %d 修改密码失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		// check 新密码
		if form.NewPassword != form.ConfirmPassword {
			var msg = "两次新密码不一致"
			log.Warnf("用户id: %d 修改密码失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Password: form.ConfirmPassword, ID: userId})
	case "avatar":
		type avatar struct {
			Avatar string `json:"avatar" binding:"required" msg:"头像不能为空"`
		}
		object := &avatar{}
		if err := ctx.ShouldBindJSON(&object); err != nil {
			log.Warnf("用户id: %d 修改头像参数解析失败,err: %s", userId, err.Error())
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

func listUsers(ctx *gin.Context) {
	name := ctx.Query("name")
	users := userService.ListUsers(name)
	result.Ok(users, "").Json(ctx)
}

var userTagS services.UserTag

func getTagsByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		log.Warnf("获取用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	tagNames := userTagS.GetTagsByUserId(userId)
	result.Ok(tagNames, "").Json(ctx)
}
