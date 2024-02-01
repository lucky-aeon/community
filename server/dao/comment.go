package dao

import (
	"xhyovo.cn/community/server/model"
)

type CommentDao struct {
}

// 发布评论
func (a *CommentDao) AddComment(comment *model.Comments) {
	db := model.Comment()
	db.Create(&comment)
	// 如果为根评论,则需要获取id设置rootId
	if comment.ParentId == 0 {
		comment.RootId = comment.ID
	}
	db.Where("id = ?", comment.ID).Update("root_id", &comment.RootId)
}

// 删除评论以及下的子评论
func (a *CommentDao) Delete(id, userId int) {
	//model.Comment().Where("id = ?", id).Delete(&model.Comments{})
	sql := "select id,parent_id from comments where root_id = (select root_id from comments where id = ?) and user_id =?"
	db := model.Comment()
	var comments []*model.Comments
	db.Raw(sql, id, userId).Find(&comments)

	// key: parentId value: id
	m := make(map[int]int)

	for _, v := range comments {
		m[v.ParentId] = v.ID
	}
	var deleteIds []int
	deleteIds = append(deleteIds, id)
	for true {
		if v, ok := m[id]; ok {
			id = v
			deleteIds = append(deleteIds, id)
		} else {
			break
		}
	}

	db.Delete(&model.Comments{}, &deleteIds)
}

func (a *CommentDao) Create(comment *model.Comments) error {
	return model.Comment().Model(&model.Comments{}).Create(comment).Error
}

// 查询文章下的所有评论
func (a *CommentDao) GetAllCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {
	var comments []*model.Comments
	model.Comment().Where("business_id", businessId).Order("created_at").Select("id").Limit(limit).Offset((page - 1) * limit).Find(&comments)
	count := a.GetCommentsCountByArticleID(businessId)
	return comments, count
}

// 查询根评论下的子评论总数
func (a *CommentDao) GetCommentsCountByRootId(rootIds []int) map[int]int {
	sql := "SELECT root_id, COUNT(*) AS number FROM comments WHERE root_id IN (?)GROUP BY root_id;"
	var ChildCommentNumber []*model.ChildCommentNumber
	model.Comment().Raw(sql, rootIds).Scan(&ChildCommentNumber)
	m := make(map[int]int)
	for i := range ChildCommentNumber {
		commentNumber := ChildCommentNumber[i]
		// 把自身减去
		m[commentNumber.RootId] = commentNumber.Number - 1
	}

	return m
}

// 查询文章下的评论带分页并且只显示跟评论的前n条
func (a *CommentDao) GetCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {
	// 查询所有根评论,只想要根评论
	var parentIds []int
	var comments []*model.Comments
	model.Comment().Where("business_id", businessId).Order("created_at").Select("id").Group("root_id").Limit(limit).Offset(page - 1*limit).Find(&parentIds)

	if len(parentIds) == 0 {
		return comments, 0
	}

	// 根据根评论查
	sql := "select c.* from comments c where (select count(id) from comments where root_Id = c.root_id and id<=c.id ) <= 5 and  c.root_id in  ? order by root_id desc"
	model.Comment().Raw(sql, parentIds).Scan(&comments)

	count := a.GetCommentsCountByArticleID(businessId)
	return comments, count
}

// 根据根评论查询下的子评论
func (a *CommentDao) GetCommentsByCommentID(page, limit, rootId int) ([]*model.Comments, int64) {
	sql := "select * from comments where root_id = ? order by created_at desc limit ?,?"
	var comments []*model.Comments

	db := model.Comment()

	db.Raw(sql, rootId, (page-1)*limit, limit).Scan(&comments)
	count := a.GetRootCommentsCountByArticleID(rootId)
	return comments, count
}

// 查询跟评论下的评论总数
func (a *CommentDao) GetRootCommentsCountByArticleID(rootId int) int64 {
	sql := "select count(id) from comments where root_id =?"
	var count int64
	db := model.Comment()
	db.Raw(sql, rootId).Scan(&count)
	return count
}

// 获取文章评论总数
func (a *CommentDao) GetCommentsCountByArticleID(businessId int) int64 {
	var count int64
	model.Comment().Where("business_id", businessId).Count(&count)
	return count
}
