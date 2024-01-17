package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/dao"
)

var (
	commentDao = &dao.Comment{}
)

func InitCommentRouter(r *gin.Engine) {
	r.GET("/comments", commentList)
	r.POST("/comments", commentAdd)
	r.DELETE("/comments", commentDeleted)
	r.PUT("/comments", commentUpdate)
}

func commentList(c *gin.Context) {
	// commentDao.GetCommentsByArticleID()
}

func commentAdd(c *gin.Context) {

}

func commentDeleted(c *gin.Context) {

}

func commentUpdate(c *gin.Context) {

}
