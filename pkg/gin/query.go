package ginutils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/data"
)

func GetPage(ctx *gin.Context) data.QueryPage {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 15
	}
	return data.QueryPage{Page: page, Limit: limit}
}
