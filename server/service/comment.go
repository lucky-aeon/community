package services

import "xhyovo.cn/community/server/model"

type CommentsService struct {
}

// 查询文章下的评论
func (*CommentsService) GetCommentsByArticleID(page, limit, businessId uint) []*model.Comments {

	var parentComments []*model.Comments
	var commentsMap map[uint][]*model.Comments
	comments := commentDao.GetCommentsByArticleID(page, limit, businessId)

	// 收集根评论
	for i := range comments {
		if comments[i].ParentId == 0 {
			parentComments = append(parentComments, comments[i])
		}
		commentsMap[comments[i].RootId] = append(commentsMap[comments[i].ParentId], comments[i])
	}
	for i := range parentComments {
		parentComments[i].ChildComments = commentsMap[parentComments[i].RootId]
	}

	return parentComments
}
