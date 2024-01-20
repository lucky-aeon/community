package dao

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type Article struct {
}

func (a *Article) QuerySingle(article *model.Article) (*model.Article, error) {
	result := &model.Article{}
	err := db.Model(article).Where(article).First(result).Error
	return result, err
}

func (a *Article) QueryList(article *model.Article, page, limit int) (*[]model.Article, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	userDb := db.Model(article).Where(article)
	pageSize := int64(0)
	articleList := &[]model.Article{}
	userDb.Count(&pageSize)

	if pageSize == 0 {
		return articleList, userDb.Error
	}

	userDb.Offset((page - 1) * limit).Limit(limit).Find(articleList)
	return articleList, userDb.Error
}

func (a *Article) Delete(articleId, userId uint) error {
	return db.Model(&model.Article{}).Delete(&model.Article{
		Model: gorm.Model{
			ID: articleId,
		},
		UserId: userId,
	}).Error
}

func (a *Article) Create(article *model.Article) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return db.Model(article).Create(article).Error
}

func (a *Article) Update(article *model.Article) error {
	article.UpdatedAt = time.Now()
	return db.Model(article).Where(&model.Article{
		Model: gorm.Model{
			ID: article.ID,
		},
		UserId: article.UserId,
	}).Save(article).Error
}
