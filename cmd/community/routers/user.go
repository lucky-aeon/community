package routers

import (
	"github.com/gin-gonic/gin"
	"io"
	"xhyovo.cn/community/pkg/kodo"
	context2 "xhyovo.cn/community/pkg/service_context"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var userService services.UserService

var fileService services.FileService

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
	group.GET("/:id", getUserInfo)
	group.POST("/edit", updateUser)
}

// 获取用户信息
func getUserInfo(ctx *gin.Context) {
	context := context2.DataContext(ctx)
	id := context.Auth().ID
	user := userService.GetUserById(id)
	context.View("user.edit", user)

}

func updateUser(ctx *gin.Context) {

	context := context2.DataContext(ctx)
	id := context.Auth().ID
	t := ctx.DefaultQuery("tab", "info")
	switch t {
	case "info":
		form := editUserForm{}
		err := ctx.ShouldBind(&form)
		if err != nil {
			context.To("/community/user/edit?tab=info").WithError(utils.GetValidateErr(form, err)).Redirect()
			return
		}
		userService.UpdateUser(&model.Users{Name: form.Name, Desc: form.Desc, ID: id})
	case "pass":
		form := editPasswordForm{}
		err := ctx.ShouldBind(&form)
		if err != nil {
			context.To("/community/user/edit?tab=pass").WithError(utils.GetValidateErr(form, err)).Redirect()
			return
		}
		// check 旧密码
		if form.OldPassword != userService.GetUserById(id).Password {
			context.To("/community/user/edit?tab=pass").WithError("旧密码不一致").Redirect()
			return
		}
		// check 新密码
		if form.NewPassword != form.ConfirmPassword {
			context.To("/community/user/edit?tab=pass").WithError("两次新密码不一致").Redirect()
			return
		}
		userService.UpdateUser(&model.Users{Password: form.ConfirmPassword, ID: id})
	case "avatar":

		multipart, err := ctx.FormFile("avatar")
		if err != nil {
			context.To("/community/user/edit?tab=avatar").WithError("请选择上传的图片").Redirect()
			return
		}
		filename := multipart.Filename
		if err != nil {
			context.To("/community/user/edit?tab=avatar").WithError("上传头像失败").Redirect()
			return
		}
		if !utils.IsImage(filename) {
			context.To("/community/user/edit?tab=avatar").WithError("请上传图片类型").Redirect()
			return
		}
		// 打开上传的文件
		uploadedFile, err := multipart.Open()
		if err != nil {
			context.To("/community/user/edit?tab=avatar").WithError(err).Redirect()
			return
		}
		defer uploadedFile.Close()

		// 读取文件的二进制数据
		fileBytes, err := io.ReadAll(uploadedFile)
		if err != nil {
			context.To("/community/user/edit?tab=avatar").WithError(err).Redirect()
			return
		}

		// 上传oss
		fileKey := utils.BuildFileKey(id)
		if _, err = kodo.Upload(fileBytes, fileKey); err != nil {
			context.To("/community/user/edit?tab=avatar").WithError(err).Redirect()
			return
		}
		// 保存文件表
		fileService.Save(id, 1, fileKey)

		// 更改用户信息
		userService.UpdateUser(&model.Users{ID: id, Avatar: fileKey})
	}

	context.Refresh(userService.GetUserById(id))
	context.Back().WithMsg("修改成功").Redirect()
}
