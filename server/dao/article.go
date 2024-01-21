package dao

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type Article struct {
}

func (a *Article) QuerySingle(article *model.Articles) (*model.Articles, error) {
	result := &model.Articles{}
	err := model.Article().Model(article).Where(article).First(result).Error
	return result, err
}

func (a *Article) QueryList(article *model.Articles, page, limit int) ([]*model.Articles, error) {
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	userDb := model.Article().Model(article).Where(article)

	articleList := []*model.Articles{}
	userDb.Offset((page - 1) * limit).Limit(limit).Find(&articleList)
	return articleList, userDb.Error
}

func (a *Article) Count() int64 {
	var count int64
	model.Type().Count(&count)
	return count
}

func (a *Article) Delete(articleId, userId uint) error {
	return model.Article().Model(&model.Articles{}).Delete(&model.Articles{
		Model: gorm.Model{
			ID: articleId,
		},
		UserId: userId,
	}).Error
}

func (a *Article) Create(article *model.Articles) error {
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	return model.Article().Model(article).Create(article).Error
}

func (a *Article) Update(article *model.Articles) error {
	article.UpdatedAt = time.Now()
	return model.Article().Model(article).Where(&model.Articles{
		Model: gorm.Model{
			ID: article.ID,
		},
		UserId: article.UserId,
	}).Save(article).Error
}
