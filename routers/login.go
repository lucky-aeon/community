package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/services"
	"xhyovo.cn/community/utils"
)

func Login(c *gin.Context) {

	user := services.Login(c.Query("account"), c.Query("password"))
	user.Password = ""
	utils.Ok().Data(user).Res(c)
}

func Register(c *gin.Context) {
	services.
}
