package services

import (
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"strings"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type ArticleTagService struct {
}

// get hot top 10 item
func (*ArticleTagService) QueryHotTags(limit int) (result []*model.ArticleTags, err error) {
	result = make([]*model.ArticleTags, 0)
	d := model.ArticleTag().Order("updated_at").Limit(limit).Find(&result)
	return result, d.Error
}

func (*ArticleTagService) QueryList(page, limit int, title string) (result map[string]interface{}, err error) {
	var count int64
	list := make([]*model.ArticleTags, 0)
	err = model.ArticleTag().Where("tag_name like ?", "%"+title+"%").Count(&count).Offset((page - 1) * limit).Limit(limit).Find(&list).Error
	result = map[string]interface{}{
		"list":  list,
		"total": count,
	}
	return
}

func (*ArticleTagService) CreateTag(tag model.ArticleTags) (result *model.ArticleTags, err error) {
	db := model.ArticleTag()
	// 暂时先count代替查重,后面用索引去重
	tagName := tag.TagName
	tagName = strings.ToLower(tagName)
	db.Where("tag_name = ?", tagName).First(&tag)
	if tag.Id != 0 {
		return &tag, nil
	}
	model.ArticleTag().Save(&tag)
	model.ArticleTagUserRelation().Create(&model.ArticleTagUserRelations{UserId: tag.UserId, TagId: tag.Id})
	return &tag, nil
}

func (a *ArticleTagService) DeleteTag(tagId, userId int) error {
	// 被引用则不能删除
	var count int64
	model.ArticleTagRelation().Where("tag_id = ?", tagId).Count(&count)
	if count > 0 {
		return errors.New("标签已被引用,不可删除")
	}
	mysql.GetInstance().Where("user_id = ? and tag_id = ?", userId, tagId).Delete(model.ArticleTagUserRelations{})
	mysql.GetInstance().Where("id = ? and user_id = ?", tagId, userId).Delete(model.ArticleTags{})
	return nil
}

func (a *ArticleTagService) GetTagArticleCount(userId int) []model.TagArticleCount {
	var tagsIds []int
	model.ArticleTagUserRelation().Where("user_id = ?", userId).Select("tag_id").Find(&tagsIds)
	if len(tagsIds) == 0 {
		return []model.TagArticleCount{}
	}
	var tagAcount []model.TagArticleCount
	db := mysql.GetInstance()
	db.Raw("select tag_id,count(article_id) as article_count from article_tag_relations WHERE tag_id in(?) GROUP BY tag_id", tagsIds).Scan(&tagAcount)
	if len(tagAcount) == 0 {
		return []model.TagArticleCount{}
	}
	setTagIds := mapset.NewSet[int]()
	for i := range tagAcount {
		setTagIds.Add(tagAcount[i].TagId)
	}
	for i := range tagsIds {
		id := tagsIds[i]
		if !setTagIds.Contains(id) {
			tagAcount = append(tagAcount, model.TagArticleCount{TagId: id, ArticleCount: 0})
		}
	}

	var articleTags []model.ArticleTags
	model.ArticleTag().Where("id in ?", tagsIds).Select("id", "tag_name").Find(&articleTags)

	var m = make(map[int]string)
	for i := range articleTags {
		m[articleTags[i].Id] = articleTags[i].TagName
	}
	for i := range tagAcount {
		tagAcount[i].TagName = m[tagAcount[i].TagId]
	}
	return tagAcount
}
