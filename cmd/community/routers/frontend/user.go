package frontend

import (
	"io"
	"xhyovo.cn/community/cmd/community/middleware"

	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/kodo"
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
}

func getUserMenu(ctx *gin.Context) {
	result.Ok(userService.GetUserMenu(), "ok").Json(ctx)
}

// 获取用户信息
func getUserInfo(ctx *gin.Context) {

	var userService services.UserService
	user := userService.GetUserById(3)

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

		multipart, err := ctx.FormFile("avatar")
		if err != nil {
			result.Err("请选择上传的图片").Json(ctx)
			return
		}
		filename := multipart.Filename
		if err != nil {
			result.Err("上传头像失败").Json(ctx)
			return
		}
		if !utils.IsImage(filename) {
			result.Err("请上传图片类型").Json(ctx)
			return
		}
		// 打开上传的文件
		uploadedFile, err := multipart.Open()
		if err != nil {
			result.Err(err.Error()).Json(ctx)
			return
		}
		defer uploadedFile.Close()

		// 读取文件的二进制数据
		fileBytes, err := io.ReadAll(uploadedFile)
		if err != nil {
			result.Err(err.Error()).Json(ctx)
			return
		}

		// 上传oss
		fileKey := utils.BuildFileKey(userId)
		if _, err = kodo.Upload(fileBytes, fileKey); err != nil {
			result.Err(err.Error()).Json(ctx)
			return
		}
		// 保存文件表
		var fileService services.FileService
		fileService.Save(userId, 1, fileKey)

		// 更改用户信息
		userService.UpdateUser(&model.Users{ID: userId, Avatar: fileKey})
	}

	result.Ok(nil, "修改成功").Json(ctx)
}
