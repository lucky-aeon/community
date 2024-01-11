package routers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	services "xhyovo.cn/community/server/service"
)

// todo: feat login
func Login(c *gin.Context) {

	user, _ := services.Login(c.Query("account"), c.Query("password"))
	user.Password = ""
	R.Ok().Data(user).Res(c)
}

func Register(c *gin.Context) {
	code, err := strconv.Atoi(c.Query("inviteCode"))
	if err != nil {
		R.Error().Msg("序列化邀请码失败,请检查邀请码是否为数字")
	}
	if err := services.Register(c.Query("account"),
		c.Query("password"), c.Query("name"), code); err != nil {
		R.Error().Msg(err.Error()).Res(c)
	}
	R.Ok().Msg("注册成功,你可以成长咯").Res(c)
}
