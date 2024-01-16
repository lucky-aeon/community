package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	services "xhyovo.cn/community/server/service"
)

func Login(c *gin.Context) {

	user, err := services.Login(c.Query("account"), c.Query("password"))
	if err != nil {
		R.Error().setMsg(err.Error()).Res(c)
		return
	}
	user.Password = ""
	R.Ok().setData(user).Res(c)

}

func Register(c *gin.Context) {
	code, err := strconv.Atoi(c.Query("inviteCode"))
	if err != nil {
		R.Error().setMsg("序列化邀请码失败,请检查邀请码是否为数字").Res(c)
		return
	}

	if err := services.Register(c.Query("account"),
		c.Query("password"), c.Query("name"), uint16(code)); err != nil {
		fmt.Printf(err.Error())
		R.Error().setMsg(err.Error()).Res(c)
		return
	}
	R.Ok().setMsg("注册成功,你可以成长咯").Res(c)

}
