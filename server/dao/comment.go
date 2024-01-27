package dao

import (
	"time"

	"xhyovo.cn/community/server/model"
)

type CommentDao struct {
}

func (a *CommentDao) Delete(articleId, userId uint) error {
	return model.Comment().Model(&model.Comments{}).Delete(&model.Comments{
		ID:     articleId,
		UserId: userId,
	}).Error
}

func (a *CommentDao) Create(comment *model.Comments) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return model.Comment().Model(&model.Comments{}).Create(comment).Error
}

// 查询文章下的评论
func (a *CommentDao) GetCommentsByArticleID(page, limit, businessId uint) []*model.Comments {
	// 查询所有根评论,只想要根评论
	var parentIds []int
	var comments []*model.Comments
	model.Comment().Where("business_id", businessId).Order("created_at").Select("id").Limit(int(limit)).Offset(int(page-1) * int(limit)).Find(parentIds)

	if len(parentIds) == 0 {
		return comments
	}

	// 根据根评论查
	sql := "select c.* from comments c where (select count(id) from comments where root_Id = c.root_id and id<=c.id  ) <= 3 and  c.root_id in = ? order by root_id desc"
	model.Comment().Raw(sql, parentIds).Scan(&comments)

	return comments
}
