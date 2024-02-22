package services

import (
	"time"

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

func (*ArticleTagService) CreateTag(tag *model.ArticleTags) (result *model.ArticleTags, err error) {
	db := model.ArticleTag()
	db.Where(`tag_name = ?`, tag.TagName).First(result)
	if result == nil {
		result = &model.ArticleTags{
			TagName:     tag.TagName,
			Description: tag.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserId:      tag.UserId,
		}
		err = db.Create(result).Error
		if err != nil {
			result = nil
		}
	}
	model.ArticleTagUserRelation().Create(
		model.ArticleTagUserRelations{
			UserId: tag.UserId,
			TagId:  result.Id,
		})
	return tag, err
}

func (*ArticleTagService) DeleteTag(userId, tagId int) bool {
	var count int64 = 0
	// find user tag
	model.ArticleTagUserRelation().Where("tag_id = ?, user_id = ?", tagId, userId).Count(&count)
	if count == 0 {
		return false
	}
	// remove article tag ref and user ref
	var userArticleIds []int
	model.Article().Select("id").Where("user_id = ?", userId).Scan(&userArticleIds)
	model.ArticleTagRelation().Delete("tag_id = ? and article_id in ?", tagId, userArticleIds)
	model.ArticleTagUserRelation().Delete("user_id = ? and tag_id = ?", userId, tagId)
	return true
}
