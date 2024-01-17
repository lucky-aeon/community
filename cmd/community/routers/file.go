package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	services "xhyovo.cn/community/server/service"
)

var file services.File

var kodo services.Kodo

func InitFileRouter(ctx *gin.Engine) {
	group := ctx.Group("/community/file")
	group.GET("/getUploadToken", GetUploadToken)
	group.POST("/saveFile", SaveFile)
}

// 获取上传token
func GetUploadToken(ctx *gin.Context) {
	ctx.JSON(200, &R{
		Code: 200,
		Data: kodo.GetToken(),
	})
}

func SaveFile(ctx *gin.Context) {
	articleId := ctx.Query("articleId")
	fileKey := ctx.Query("fileKey")
	atoi, err := strconv.Atoi(articleId)
	if err != nil {
		ctx.JSON(500, &R{
			Code: 500,
			Msg:  "序列化文章id失败,请检查是否为数字",
		})
		return
	}
	// todo get userId

	err = file.Save(0, uint(atoi), fileKey)
	if err != nil {
		ctx.JSON(500, &R{
			Code: 500,
			Msg:  "在我们空间中没有该文件",
		})
		return
	}

}
