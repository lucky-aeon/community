package services

import (
	"gorm.io/gorm/clause"
	"xhyovo.cn/community/pkg/data"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type ArticleService struct {
}

func (*ArticleService) GetArticleData(id int) (data *model.ArticleData, err error) {
	a, err := articleDao.QuerySingle(model.Articles{ID: id})
	if err != nil {
		return nil, err
	}
	var tags []*model.ArticleTagSimple
	model.ArticleTag().Joins("LEFT JOIN article_tag_relations as atr ON atr.tag_id = article_tags.id").
		Where("atr.article_id = ?", a.ID).Find(&tags)
	us, err := userDao.QueryUserSimple(&model.Users{ID: a.ID})
	if err != nil {
		us = model.UserSimple{
			UId:     0,
			UName:   "未知用户",
			UDesc:   "未知啦",
			UAvatar: "",
		}
	}
	var typeData model.TypeSimple
	model.Type().Where("id = ?", a.Type).First(&typeData)
	return &model.ArticleData{
		ID:         a.ID,
		Title:      a.Title,
		State:      a.State,
		Like:       a.Like,
		Tags:       tags,
		UserSimple: us,
		TypeSimple: typeData,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}, err
}

func (a *ArticleService) PageByClassfily(tagId []int, article *model.Articles, page data.QueryPage, sort data.ListSortStrategy) (result []*model.ArticleData, total int64, err error) {
	query := mysql.GetInstance().Table("articles").
		Select("articles.id, articles.title, articles.state, articles.`like`, tp.title as type_title, tp.flag_name as type_flag, u.name as u_name, u.id as u_id, articles.created_at, articles.updated_at, GROUP_CONCAT(DISTINCT atg.tag_name) as tags").
		Joins("LEFT JOIN article_tag_relations as atr on atr.article_id = articles.id").
		Joins("LEFT JOIN article_tags as atg on atg.id = atr.tag_id").
		Joins("LEFT JOIN types as tp on tp.id = articles.type").
		Joins("LEFT JOIN users as u on u.id = articles.user_id")
	if article != nil {
		if article.Type > 0 {
			query.Where("type = ?", article.Type)
		}
		if len(article.Title) > 0 {
			query.Where("title like ?", article.Title)
		}
		if len(article.Desc) > 0 {
			query.Where("desc like ?", article.Desc)
		}
	}
	if len(tagId) > 0 {
		query.Where("atr.tag_id in ?", tagId)
	}
	query.Group("articles.id").Count(&total)
	if len(sort.OrderBy) > 0 {
		query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: "articles." + sort.OrderBy},
			Desc:   sort.DescOrder,
		})
	}
	query.Offset((page.Page - 1) * page.Limit).
		Limit(page.Limit).
		Find(&result)
	return
}

func (a *ArticleService) Count() int64 {
	return articleDao.Count()
}

func (a *ArticleService) CountByTypeId(typeId int) int64 {
	return articleDao.CountByTypeId(typeId)
}

func (a *ArticleService) ListByIdsSelectIdTitleMap(id []int) map[int]string {

	m := make(map[int]string)
	articles := articleDao.ListByIdsSelectIdTitle(id)
	for i := range articles {
		v := articles[i]
		m[v.ID] = v.Title
	}
	return m
}
