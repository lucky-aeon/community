package dao

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type Comment struct {
}

func (a *Comment) Delete(articleId, userId uint) error {
	return model.Comment().Model(&model.Comments{}).Delete(&model.Comments{
		Model: gorm.Model{
			ID: articleId,
		},
		UserId: userId,
	}).Error
}

func (a *Comment) Create(comment *model.Comments) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return model.Comment().Model(&model.Comments{}).Create(comment).Error
}

// 查询文章下的评论
func (a *Comment) GetCommentsByArticleID(businessId uint) ([]model.Comments, error) {
	var comments = make([]model.Comments, 0)
	commentDb := model.Comment().Model(&model.Comments{}).Where("business_id = ?", businessId).Find(&comments)
	return comments, commentDb.Error
}
