package dao

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/server/model"
)

type Comment struct {
}

func (a *Comment) Delete(articleId, userId uint) error {
	return db.Model(&model.Comment{}).Delete(&model.Comment{
		Model: gorm.Model{
			ID: articleId,
		},
		UserId: userId,
	}).Error
}

func (a *Comment) Create(comment *model.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return db.Model(&model.Comment{}).Create(comment).Error
}

// 查询文章下的评论
func (a *Comment) GetCommentsByArticleID(businessId uint) ([]model.Comment, error) {
	var comments = make([]model.Comment, 0)
	commentDb := db.Model(&model.Comment{}).Where("business_id = ?", businessId).Find(&comments)
	return comments, commentDb.Error
}
