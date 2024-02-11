package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xhyovo.cn/community/pkg/data"
)

func GetOderBy(ctx *gin.Context) (res data.ListSortStrategy) {
	ctx.ShouldBindBodyWith(&res, binding.JSON)
	return
}
