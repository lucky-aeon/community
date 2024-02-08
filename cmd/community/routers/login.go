package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	services "xhyovo.cn/community/server/service"
)

type registerForm struct {
	Code     int    `binding:"required" form:"code" msg:"code不能为空" `
	Account  string `binding:"required,email" form:"account" msg:"邮箱格式不正确"`
	Name     string `binding:"required" form:"name" msg:"用户名不能为空"`
	Password string `binding:"required" form:"password" msg:"密码不能为空"`
}

type loginForm struct {
	Account  string `binding:"required" json:"account" msg:"账号不能为空"`
	Password string `binding:"required" json:"password" msg:"密码不能为空"`
}

func InitLoginRegisterRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	group.POST("/login", Login)
	group.POST("/register", Register)
}

func Login(c *gin.Context) {
	var form loginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		result.Err(utils.GetValidateErr(form, err)).Json(c)
		return
	}
	user, err := services.Login(form.Account, form.Password)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	user.Password = ""

	token, _ := middleware.GenerateToken(user.ID, user.Name)

	result.Ok(map[string]string{"token": token}, "登录成功").Json(c)
}

func Register(c *gin.Context) {
	var form registerForm

	err := c.ShouldBindJSON(&form)

	if err != nil {
		result.Err(utils.GetValidateErr(form, err)).Json(c)
		return
	}

	err = services.Register(form.Account, form.Password, form.Name, form.Code)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(form, "注册成功").Json(c)
}
