package routers

import (
	"github.com/gin-gonic/gin"
	services "xhyovo.cn/community/server/service"
)

var userService services.UserService

// 获取用户信息
func GetInfo(ctx *gin.Context) {
	// todo session get userId

}

func UpdateUser(ctx *gin.Context) {

	// todo session get userId
	var userId uint
	name := ctx.Query("username")
	pswd := ctx.Query("password")
	userService.UpdateUser(userId, name, pswd)
	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}
