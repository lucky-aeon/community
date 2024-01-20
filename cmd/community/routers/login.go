package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	services "xhyovo.cn/community/server/service"
)

func InitLoginRegisterRouter(ctx *gin.Engine) {
	group := ctx.Group("/community")
	group.POST("/login", Login)
	group.POST("/register", Register)
}

func Login(c *gin.Context) {

	user, err := services.Login(c.Query("account"), c.Query("password"))
	if err != nil {
		c.JSON(500, &R{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	user.Password = ""
	c.JSON(200, &R{
		Code: 200,
		Msg:  "登录成功",
		Data: user,
	})

}

func Register(c *gin.Context) {
	code, err := strconv.Atoi(c.Query("inviteCode"))
	if err != nil {
		c.JSON(500, &R{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}

	if err := services.Register(c.Query("account"),
		c.Query("password"), c.Query("name"), uint16(code)); err != nil {
		c.JSON(500, &R{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	c.JSONP(200, &R{
		Code: 200,
		Msg:  "注册成功,你可以成长了",
	})

}
