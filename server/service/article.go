package services

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/service_context"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
)

type ArticleService struct {
}

// 获取文章，分页，类型
func (a *ArticleService) Page(ctx *service_context.BaseContext) gin.H {

	c := ctx.Ctx
	cur, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))
	typeId, _ := strconv.Atoi(c.DefaultQuery("typeId", ""))

	articles, err := articleDao.QueryList(&model.Articles{Type: uint(typeId)}, cur, limit)
	if err != nil {
		return nil
	}
	// 设置用户信息
	count := a.Count()
	userDao.SetUserInfo(articles)
	// 设置 label 信息 todo

	baseUrl := c.Request.RequestURI

	pageObj := page.New(int(count), limit, cur, baseUrl)

	return gin.H{"articles": articles, "page": pageObj, "type": typeId}

}

func (a *ArticleService) Count() int64 {
	return articleDao.Count()
}
