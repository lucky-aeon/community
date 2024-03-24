package services

import (
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/server/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/data"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type ArticleService struct {
}

func (*ArticleService) GetArticleData(id, userId int) (data *model.ArticleData, err error) {
	var a model.Articles
	model.Article().Where("id = ?", id).First(&a)
	if a.ID == 0 && a.UserId != userId && a.State == constant.Draft {
		return nil, errors.New("文章不存在")
	}
	var tags []*model.ArticleTagSimple
	model.ArticleTag().Joins("LEFT JOIN article_tag_relations as atr ON atr.tag_id = article_tags.id").
		Where("atr.article_id = ?", a.ID).Find(&tags)
	us, err := userDao.QueryUserSimple(&model.Users{ID: a.UserId})
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
		Desc:       a.Content,
		UserSimple: us,
		TypeSimple: typeData,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}, err
}

func (a *ArticleService) PageByClassfily(tagId []int, article *model.Articles, page data.QueryPage, sort data.ListSortStrategy) (result []*model.ArticleData, total int64, err error) {
	query := mysql.GetInstance().Table("articles").
		Select("articles.id, articles.title, articles.state, articles.`like`, "+
			"tp.id as type_id, tp.title as type_title, tp.flag_name as type_flag, "+
			"u.name as u_name, u.id as u_id, u.avatar as u_avatar, "+
			"articles.created_at, articles.updated_at, GROUP_CONCAT(DISTINCT atg.tag_name) as tags").
		Joins("LEFT JOIN article_tag_relations as atr on atr.article_id = articles.id").
		Joins("LEFT JOIN article_tags as atg on atg.id = atr.tag_id").
		Joins("LEFT JOIN types as tp on tp.id = articles.type").
		Joins("LEFT JOIN users as u on u.id = articles.user_id").
		Where("articles.state != ?", constant.Draft)
	if article != nil {
		if article.Type > 0 {
			query.Where("articles.type = ?", article.Type)
		}
		if len(article.Title) > 0 {
			query.Where("articles.title like ?", "%"+article.Title+"%")
		}
		if len(article.Content) > 0 {
			query.Where("articles.content like ?", "%"+article.Content+"%")
		}
		if article.UserId > 0 {
			query.Where("articles.user_id = ?", article.UserId)
		}
	}
	if len(tagId) > 0 {
		query.Where("atg.tag_name in ?", tagId)
	}
	query.Group("articles.id").Count(&total)
	if len(sort.OrderBy) > 0 {
		query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: "articles." + sort.OrderBy},
			Desc:   sort.DescOrder,
		})
	}
	rows, err := query.Offset((page.Page - 1) * page.Limit).
		Limit(page.Limit).
		Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		item := model.ArticleData{}
		itemType := model.TypeSimple{}
		itemUser := model.UserSimple{}
		tags := ""
		rows.Scan(
			&item.ID, &item.Title, &item.State, &item.Like,
			&itemType.TypeId, &itemType.TypeTitle, &itemType.TypeFlag,
			&itemUser.UName, &itemUser.UId, &itemUser.UAvatar,
			&item.CreatedAt, &item.UpdatedAt, &tags,
		)
		item.UserSimple = itemUser
		item.TypeSimple = itemType
		item.Tags = tags
		result = append(result, &item)
	}
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

func (a *ArticleService) GetById(id int) model.Articles {
	article := articleDao.GetById(id)
	user := userDao.GetById(article.UserId)
	article.Users = user
	return article
}

// 点赞/取消点赞文章
func (a *ArticleService) Like(articleId, userId int) bool {

	// 点赞
	err := mysql.GetInstance().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.Article_Likes{ArticleId: articleId, UserId: userId}).Error; err != nil {
			return err
		}

		return tx.Model(&model.Articles{}).Where("id = ?", articleId).Update("like", gorm.Expr("`like` + ?", 1)).Error
	})

	if err != nil {
		mysql.GetInstance().Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("article_id = ? and user_id = ?", articleId, userId).Delete(&model.Article_Likes{}).Error; err != nil {
				return err
			}
			return tx.Model(&model.Articles{}).Where("id = ?", articleId).Update("like", gorm.Expr("`like` + ?", -1)).Error
		})
	}
	return err == nil
}

func (a *ArticleService) PublishArticleCount(userId int) (count int64) {
	model.Article().Where("user_id = ?", userId).Count(&count)
	return
}

func (a *ArticleService) PublishArticlesSelectId(userId int) (id []int) {
	model.Article().Where("user_id = ?", userId).Select("id").Find(&id)
	return
}

// 获取文章的点赞次数
func (a *ArticleService) ArticlesLikeCount(ids []int) (count int64) {
	model.ArticleLike().Where("article_id  in ? ", ids).Count(&count)
	return
}

func (a *ArticleService) SaveArticle(article request.ReqArticle) (int, error) {

	id := article.ID
	typeO := article.Type
	flag := true
	var typeS TypeService
	// 分类是否存在
	if !typeS.Exist(typeO) {
		return 0, errors.New("分类不存在")
	}
	// 状态是否存在
	oldState := article.State
	if oldState < 0 || oldState > 5 {
		return 0, errors.New("状态不存在")
	}

	// 根据分类选择状态：QA分类没有发布,普通分类只有草稿和发布
	if typeO == 1 && oldState == constant.Published {
		return 0, errors.New("QA分类状态不能选择已发布")
	} else {
		// 普通分类校验状态
		if oldState == constant.Pending || oldState == constant.Resolved || oldState == constant.PrivateQuestion {
			msg := constant.GetArticleName(oldState)
			return 0, errors.New("普通分类不支持该状态:" + msg)
		}
	}
	// 修改
	if id != 0 {
		flag = false
		// 获取老文章
		oldArticle := a.GetById(id)
		oldTypeParentId := typeS.GetById(oldArticle.Type).ParentId
		// 修改 一级分类不能修改,如果parent不同则修改了一级分类
		newTypeParentId := typeS.GetById(typeO).ParentId
		if oldTypeParentId != newTypeParentId {
			return 0, errors.New("修改的分类只能属于同一级分类下")
		}
		// 老文章状态如果为非草稿状态，则新文章不可修改为草稿状态
		if oldState != constant.Draft && article.State == constant.Draft {
			return 0, errors.New("旧文章状态不可从非草稿转为草稿")
		}
	}

	articleObject := &model.Articles{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		UserId:    article.UserId,
		State:     article.State,
		Type:      article.Type,
		UpdatedAt: time.Now(),
	}
	mysql.GetInstance().Save(&articleObject)
	id = article.ID
	// 关联关系
	db := model.ArticleTagRelation
	db().Where("article_id = ?", id).Delete(nil)
	var tags []model.ArticleTagRelations
	for i := range article.Tags {
		tags = append(tags, model.ArticleTagRelations{ArticleId: id, TagId: article.Tags[i], UserId: article.UserId})
	}
	db().Create(&tags)
	var subscriptionService SubscriptionService
	if flag {
		var b SubscribeData
		b.UserId = article.UserId
		b.ArticleId = article.ID
		b.CurrentBusinessId = article.ID
		b.SubscribeId = article.UserId
		subscriptionService.Do(event.UserFollowingEvent, b)
		subscriptionService.ConstantAtSend(event.ArticleAt, id, article.Content, b)
	}
	return id, nil
}

func (a *ArticleService) DeleteByUserId(articleId, userId int) (err error) {

	// 删除文章
	db := mysql.GetInstance()
	err = db.Where("id = ? and user_id = ?", articleId, userId).Delete(&model.Articles{}).Error
	// 删除文章标签表
	err = db.Where("article_id = ?", articleId).Delete(&model.ArticleTagRelations{}).Error
	return
}

func (a *ArticleService) Delete(articleId int) (err error) {

	// 删除文章
	db := mysql.GetInstance()
	err = db.Where("id = ?", articleId).Delete(&model.Articles{}).Error
	// 删除文章标签表
	err = db.Where("article_id = ?", articleId).Delete(&model.ArticleTagRelations{}).Error
	return
}

func (a *ArticleService) GetLikeState(articleId, userId int) bool {
	var count int64
	model.ArticleLike().Where("article_id  = ? and user_id = ?", articleId, userId).Count(&count)
	return count == 1
}

func (a *ArticleService) PageArticles(p, limit int) (articleList []model.ArticleData, count int64) {
	var articles []model.Articles
	model.Article().Limit(limit).Offset((p-1)*limit).Select("id", "created_at", "title", "user_id", "state", "type").Order("created_at desc").Find(&articles)
	model.Article().Count(&count)
	if count == 0 {
		return
	}
	// 找到文章的userId
	userIds := mapset.NewSet[int]()
	typeIds := mapset.NewSet[int]()

	for i := range articles {
		articleO := articles[i]
		aritcle := model.ArticleData{
			ID:         articleO.ID,
			Title:      articleO.Title,
			State:      articleO.State,
			CreatedAt:  articleO.CreatedAt,
			UserSimple: model.UserSimple{UId: articleO.UserId},
			TypeSimple: model.TypeSimple{TypeId: articleO.Type},
		}
		typeIds.Add(articleO.Type)
		userIds.Add(articleO.UserId)
		articleList = append(articleList, aritcle)
	}

	// 填充分类,状态名称

	var u UserService
	var t TypeService
	userMap := u.ListByIdsToMap(userIds.ToSlice())
	typeMap := t.ListByIdToMap(typeIds.ToSlice())
	for i := range articles {
		articleList[i].StateName = constant.GetArticleName(articleList[i].State)
		articleList[i].TypeSimple.TypeTitle = typeMap[articleList[i].TypeSimple.TypeId]
		articleList[i].UserSimple.UName = userMap[articleList[i].UserSimple.UId].Name
	}
	return
}
