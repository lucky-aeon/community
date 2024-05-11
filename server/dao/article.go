package dao

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"

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

func (a *Article) GetArticleSql() *gorm.DB {
	query := mysql.GetInstance().Table("articles").
		Select("articles.id, articles.title, LEFT(articles.content, 100) as `desc`, articles.cover,articles.state, articles.`like`, articles.created_at,articles.updated_at," +
			"tp.id as type_id, tp.title as type_title, tp.flag_name as type_flag, " +
			"u.name as u_name, u.id as u_id, u.avatar as u_avatar, " +
			"  GROUP_CONCAT(DISTINCT atg.tag_name) as tags").
		Joins("LEFT JOIN article_tag_relations as atr on atr.article_id = articles.id").
		Joins("LEFT JOIN article_tags as atg on atg.id = atr.tag_id").
		Joins("LEFT JOIN types as tp on tp.id = articles.type").
		Joins("LEFT JOIN users as u on u.id = articles.user_id").
		Where("articles.deleted_at is null")
	return query
}

// 用这个
func (a *Article) GetQueryArticleSql() *gorm.DB {
	query := mysql.GetInstance().Table("articles").
		Select("articles.id, articles.title, articles.abstract,articles.cover," +
			"articles.state,articles.`like`,articles.created_at,articles.updated_at,types.id as type_id,types.title as type_title," +
			"types.flag_name as type_flag,users.`name` as u_name,users.id as u_id,users.avatar as u_avatar, ( SELECT COUNT(*) FROM comments WHERE comments.business_id = articles.id and tenant_id = 0) AS comments, GROUP_CONCAT(at.tag_name SEPARATOR ', ') AS tags").
		Joins("LEFT JOIN article_tag_relations atr ON articles.id = atr.article_id").
		Joins("LEFT JOIN article_tags at ON atr.tag_id = at.id").
		Joins("join users on users.id = articles.user_id").
		Joins("JOIN types on types.id = articles.type").
		Group("articles.id, articles.title").
		Where("articles.deleted_at is null")
	return query
}
