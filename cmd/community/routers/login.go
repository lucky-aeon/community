package routers

import (
	"github.com/gin-gonic/gin"
	context2 "xhyovo.cn/community/pkg/service_context"
	"xhyovo.cn/community/pkg/utils"
	services "xhyovo.cn/community/server/service"
)

type registerForm struct {
	Code     int    `binding:"required" form:"code" msg:"code不能为空" `
	Account  string `binding:"required" form:"account" msg:"账号不能为空"`
	Name     string `binding:"required" form:"name" msg:"用户名不能为空"`
	Password string `binding:"required" form:"password" msg:"密码不能为空"`
}

type loginForm struct {
	Account  string `binding:"required" form:"account" msg:"账号不能为空"`
	Password string `binding:"required" form:"password" msg:"密码不能为空"`
}

func InitLoginRegisterRouter(ctx *gin.Engine) {
	group := ctx.Group("/community")
	group.POST("/login", Login)
	group.POST("/register", Register)
}

func Login(c *gin.Context) {
	context := context2.DataContext(c)
	var form loginForm
	if err := c.ShouldBind(&form); err != nil {
		context.To("/login").WithError(utils.GetValidateErr(form, err)).Redirect()
		return
	}
	user, err := services.Login(form.Account, form.Password)
	if err != nil {
		context.To("/login").WithError(err).Redirect()
		return
	}
	user.Password = ""
	context.SetAuth(user)
	context.To("/").Redirect()
}

func Register(c *gin.Context) {
	var form registerForm

	context := context2.DataContext(c)
	err := c.ShouldBind(&form)

	if err != nil {
		context.To("/register").WithError(utils.GetValidateErr(form, err)).Redirect()
		return
	}

	err = services.Register(form.Account, form.Password, form.Name, uint16(form.Code))
	if err != nil {
		context.To("/register").WithError(err).Redirect()
		return
	}

	context.WithMsg("注册成功").To("/login").Redirect()
}
