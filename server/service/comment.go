package services

import "xhyovo.cn/community/server/model"

type CommentsService struct {
}

// 发布评论
func (a *CommentsService) Comment(comment *model.Comments) {
	// todo 判断文章是否存在
	commentDao.AddComment(comment)
}

// 删除评论
func (a *CommentsService) DeleteComment(id, userId int) {
	commentDao.Delete(id, userId)
}

// 查询文章下的评论
func (*CommentsService) GetCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {

	var parentComments []*model.Comments
	commentsMap := make(map[int][]*model.Comments)
	comments, count := commentDao.GetCommentsByArticleID(page, limit, businessId)

	var parentIds []int
	// 收集根评论
	for i := range comments {
		if comments[i].ParentId == 0 {
			parentComments = append(parentComments, comments[i])
			parentIds = append(parentIds, comments[i].ID)
		} else {
			commentsMap[comments[i].RootId] = append(commentsMap[comments[i].RootId], comments[i])
		}
	}
	ChildCommentNumberMap := commentDao.GetCommentsCountByRootId(parentIds)
	for i := range parentComments {
		parentComments[i].ChildComments = commentsMap[parentComments[i].RootId]
		parentComments[i].ChildCommentNumber = ChildCommentNumberMap[parentComments[i].RootId]
	}

	return parentComments, count
}

// 查询文章下的所有评论
func (*CommentsService) GetAllCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {
	return commentDao.GetAllCommentsByArticleID(page, limit, businessId)
}

// 查询指定评论下的评论
func (*CommentsService) GetCommentsByRootID(page, limit, rootId int) ([]*model.Comments, int64) {

	var parentComments []*model.Comments
	commentsMap := make(map[int][]*model.Comments)
	comments, count := commentDao.GetCommentsByCommentID(page, limit, rootId)

	// 收集根评论
	for i := range comments {
		if comments[i].ParentId == 0 {
			parentComments = append(parentComments, comments[i])
		} else {
			commentsMap[comments[i].RootId] = append(commentsMap[comments[i].RootId], comments[i])
		}

	}
	// 如果根评论为空,说明是查询指定根评论下的子评论
	if len(parentComments) == 0 {
		return comments, count
	}
	for i := range parentComments {
		parentComments[i].ChildComments = commentsMap[parentComments[i].RootId]
	}

	return parentComments, count
}
