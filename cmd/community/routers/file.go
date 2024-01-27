package routers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/pkg/result"
)

func InitFileRouter(ctx *gin.Engine) {
	group := ctx.Group("/community/file")
	group.GET("/getUploadToken", GetUploadToken)
	group.POST("/saveFile", SaveFile)
}

// 获取上传token
func GetUploadToken(ctx *gin.Context) {
	result.Ok(kodo.GetToken(), "").Json(ctx)
}

func SaveFile(ctx *gin.Context) {

	articleId := ctx.Query("articleId")
	fileKey := ctx.Query("fileKey")
	atoi, err := strconv.Atoi(articleId)
	if err != nil {
		result.Err("序列化文章id失败,请检查是否为数字").Json(ctx)
		return
	}
	// todo get userId

	err = file.Save(0, uint(atoi), fileKey)
	if err != nil {
		result.Err("在我们空间中没有该文件").Json(ctx)
		return
	}

}
