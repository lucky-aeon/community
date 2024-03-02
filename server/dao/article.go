package dao

import (
	"gorm.io/gorm"

	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/server/model"
)

type Article struct {
}

func (a *Article) QuerySingle(article model.Articles) (*model.Articles, error) {
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

func (a *Article) Delete(articleId, userId int) error {
	return model.Article().Model(&model.Articles{}).Delete(&model.Articles{
		ID:     articleId,
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
		ID:     article.ID,
		UserId: article.UserId,
	}).Save(article).Error
}

func (a *Article) CountByTypeId(id int) int64 {
	var count int64
	model.Article().Where("type = ?", id).Count(&count)
	return count
}

func (a *Article) ExistById(id int) bool {
	var count int64
	model.Article().Where("id = ?", id).Count(&count)
	return count == 1
}

func (a *Article) ListByIdsSelectIdTitle(ids []int) []model.Articles {
	var articles []model.Articles
	model.Article().Where("id in ?", ids).Select("id,title").Find(&articles)

	return articles
}

func (a *Article) GetById(id int) model.Articles {
	var article model.Articles
	model.Article().Where("id = ?", id).First(&article)
	return article
}

func (a *Article) CreateLike(articleId, userId int) bool {

	tx := model.ArticleLike().Create(&model.Article_Likes{ArticleId: articleId, UserId: userId})
	affected := tx.RowsAffected
	return affected == 1
}

func (a *Article) DeleteLike(articleId, userId int) {
	model.ArticleLike().Delete("article_id = ? and user_id = ?", articleId, userId)
}

func (a *Article) UpdateCount(articleId, number int) {
	model.Article().Where("id = ?", articleId).Update("like", gorm.Expr("'like' + ?", number))
}
