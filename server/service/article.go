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
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type ArticleService struct {
}

func (*ArticleService) GetArticleData(id, userId int) (data *model.ArticleData, err error) {
	var a model.Articles
	model.Article().Where("id = ?", id).First(&a)
	var u UserService
	flag, err := u.IsAdmin(userId)
	if err != nil {
		return &model.ArticleData{}, err
	}
	if !flag {
		if (a.ID == 0 && a.UserId != userId && a.State == constant.Draft) || (userId != a.UserId && a.State == constant.PrivateQuestion) {
			return &model.ArticleData{}, errors.New("文章不存在")
		}
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

func (a *ArticleService) PageByClassfily(typeFlag string, tagId []string, article *model.Articles, page data.QueryPage, sort data.ListSortStrategy, currentUserId int) (result []*model.ArticleData, total int64, err error) {
	query := articleDao.GetArticleSql()
	if len(typeFlag) > 0 {
		query.Where("tp.flag_name = ?", typeFlag)
	}
	if article != nil {
		query.Where("articles.state = ?", article.State)
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
	defer rows.Close()

	for rows.Next() {
		item := model.ArticleData{}
		itemType := model.TypeSimple{}
		itemUser := model.UserSimple{}
		tags := ""
		rows.Scan(
			&item.ID, &item.Title, &item.State, &item.Like, &item.CreatedAt, &item.UpdatedAt,
			&itemType.TypeId, &itemType.TypeTitle, &itemType.TypeFlag,
			&itemUser.UName, &itemUser.UId, &itemUser.UAvatar,
			&tags,
		)
		item.UserSimple = itemUser
		item.TypeSimple = itemType
		item.Tags = tags
		if item.State != constant.Published {
			item.StateName = constant.GetArticleName(item.State)
		}
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
	var c1 int64
	var c2 int64
	model.Article().Where("user_id = ? and state = ?", userId, constant.Published).Count(&c1)
	model.Article().Where("user_id = ? and state = ?", userId, constant.Draft).Count(&c2)

	return c1 + c2
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
	types := typeS.GetById(typeO)
	if types.ID == 0 {
		return 0, errors.New("分类不存在")
	}
	if types.ParentId == 0 {
		return 0, errors.New("不能选择一级分类")
	}

	types = typeS.GetById(types.ParentId)
	// 状态是否存在
	state := article.State

	// QA 状态校验
	if (state == constant.Resolved || state == constant.Pending || state == constant.QADraft) && state == constant.Published {
		return 0, errors.New("QA不能选择已发布")
	} else if state == constant.Published || state == constant.Draft {
		// 文章 状态校验
		if state == constant.Pending || state == constant.Resolved || state == constant.PrivateQuestion {
			msg := constant.GetArticleName(state)
			return 0, errors.New("文章不支持该状态:" + msg)
		}
	}
	// 修改
	if id != 0 {
		flag = false
		// 获取老文章
		oldArticle := a.GetById(id)
		oldTypeParentId := typeS.GetById(oldArticle.Type).ParentId
		// 修改 一级分类不能修改,如果parent不同则修改了一级分类
		newTypeParentId := types.ID
		if oldTypeParentId != newTypeParentId {
			return 0, errors.New("修改的分类只能属于同一级分类下")
		}
		// 老文章状态如果为非草稿状态，则新文章不可修改为草稿状态
		if (oldArticle.State != constant.Draft && state == constant.Draft) || (oldArticle.State != constant.QADraft && state == constant.QADraft) {
			return 0, errors.New("旧文章状态不可从非草稿转为草稿")
		}
	}

	articleObject := &model.Articles{
		ID:      article.ID,
		Title:   article.Title,
		Content: article.Content,
		UserId:  article.UserId,
		State:   state,
		Type:    article.Type,
	}
	// 分开写，避免更新 0 值
	if article.ID == 0 {
		mysql.GetInstance().Save(&articleObject)
	} else {
		model.Article().Where("user_id = ? and id = ?", articleObject.UserId, articleObject.ID).Updates(&articleObject)
	}

	id = articleObject.ID
	// 关联关系
	db := model.ArticleTagRelation
	db().Where("article_id = ?", id).Delete(nil)
	var tags []model.ArticleTagRelations
	for i := range article.Tags {
		tags = append(tags, model.ArticleTagRelations{ArticleId: id, TagId: article.Tags[i], UserId: article.UserId})
	}
	db().Create(&tags)
	var subscriptionService SubscriptionService
	var d Draft
	if flag {
		go d.Del(article.UserId, 2)
		var b SubscribeData
		b.UserId = articleObject.UserId
		b.ArticleId = articleObject.ID
		b.CurrentBusinessId = articleObject.ID
		b.SubscribeId = articleObject.UserId
		subscriptionService.Do(event.UserFollowingEvent, b)
		subscriptionService.ConstantAtSend(event.ArticleAt, id, articleObject.Content, b)
	} else {
		go d.Del(article.UserId, 1)
	}

	return id, nil
}

func (a *ArticleService) DeleteByUserId(articleId, userId int) (err error) {

	// 删除文章
	db := mysql.GetInstance()
	tx := db.Where("id = ? and user_id = ?", articleId, userId).Delete(&model.Articles{})
	if tx.RowsAffected == 0 {
		return errors.New("删除文章不存在")
	}
	// 删除文章标签表
	err = db.Where("article_id = ?", articleId).Delete(&model.ArticleTagRelations{}).Error
	return
}

func (a *ArticleService) Delete(articleId int) (err error) {

	// 删除文章
	db := mysql.GetInstance()
	err = db.Where("id = ?", articleId).Delete(&model.Articles{}).Error
	if err != nil {
		return err
	}
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
	model.Article().Limit(limit).Offset((p-1)*limit).Select("id", "created_at", "title", "user_id", "state", "type", "top_number").Order("created_at desc").Find(&articles)
	model.Article().Count(&count)
	if count == 0 {
		return make([]model.ArticleData, 0), 0
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
			TopNumber:  articleO.TopNumber,
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

func (a *ArticleService) Auth(userId, articleId int) bool {
	var count int64
	model.Article().Where("user_id = ? and id = ?", userId, articleId).Count(&count)
	return count == 1
}

func (a *ArticleService) UpdateState(articleId, state int) {
	model.Article().Where("id = ?", articleId).Select("state").Updates(model.Articles{State: state})
}

func (a *ArticleService) QAArticleCount(userId int) (count int64) {
	var c1 int64
	var c2 int64
	var c3 int64
	var c4 int64
	model.Article().Where("user_id = ? and state = ?", userId, constant.Pending).Count(&c1)
	model.Article().Where("user_id = ? and state = ?", userId, constant.Resolved).Count(&c2)
	model.Article().Where("user_id = ? and state = ?", userId, constant.PrivateQuestion).Count(&c3)
	model.Article().Where("user_id = ? and state = ?", userId, constant.QADraft).Count(&c4)
	return c1 + c2 + c3 + c4
}

func (a *ArticleService) PageTopArticle(types string, page, limit int) (articles []model.ArticleData, count int64) {

	query := articleDao.GetArticleSql()
	query.Where("articles.state = ? and tp.flag_name = ?", constant.Top, types)
	query.Group("articles.id").Count(&count)
	rows, err := query.Order("articles.top_number desc").Limit(limit).Offset((page - 1) * limit).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		item := model.ArticleData{}
		itemType := model.TypeSimple{}
		itemUser := model.UserSimple{}
		tags := ""
		rows.Scan(
			&item.ID, &item.Title, &item.State, &item.Like, &item.CreatedAt, &item.UpdatedAt,
			&itemType.TypeId, &itemType.TypeTitle, &itemType.TypeFlag,
			&itemUser.UName, &itemUser.UId, &itemUser.UAvatar,
			&tags,
		)
		item.UserSimple = itemUser
		item.TypeSimple = itemType
		item.Tags = tags
		if item.State != constant.Published {
			item.StateName = constant.GetArticleName(item.State)
		}
		articles = append(articles, item)
	}
	return
}
func (a *ArticleService) UpdateArticleState(article request.TopArticle) error {

	return model.Article().Where("id = ?", article.Id).Updates(&article).Error
}
