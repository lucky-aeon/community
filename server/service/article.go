package services

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/model"
)

type ArticleService struct {
}

func (*ArticleService) Get(id uint) (*model.Articles, error) {

	return articleDao.QuerySingle(model.Articles{ID: id})
}

// 获取文章，分页，类型
func (a *ArticleService) Page() gin.H {
	/*
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

		return gin.H{"articles": articles, "page": pageObj, "type": typeId}*/
	return nil

}

func (a *ArticleService) Count() int64 {
	return articleDao.Count()
}

func (a *ArticleService) CountByTypeId(typeId int) int64 {
	return articleDao.CountByTypeId(typeId)
}
