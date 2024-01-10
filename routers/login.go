package routers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/services"
	"xhyovo.cn/community/utils"
)

// todo: feat login
func Login(c *gin.Context) {

	user := services.Login(c.Query("account"), c.Query("password"))
	user.Password = ""
	utils.Ok().Data(user).Res(c)
}

// todo: feat register
func Register(c *gin.Context) {
	code, err := strconv.Atoi(c.Query("inviteCode"))
	if err != nil {
		utils.Error().Msg("序列化邀请码失败,请检查邀请码是否为数字")
	}
	services.Register(c.Query("account"),
		c.Query("password"), c.Query("name"), code)
}
