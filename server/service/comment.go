package services

import (
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type CommentsService struct {
	ctx *gin.Context
}

func NewCommentService(ctx *gin.Context) *CommentsService {
	return &CommentsService{ctx: ctx}
}

// 发布评论
func (a *CommentsService) Comment(comment *model.Comments) error {

	articleId := comment.BusinessId
	if f := articleDao.ExistById(articleId); !f {
		return errors.New("文章不存在")
	}

	parentId := comment.ParentId
	// 父评论是否存在
	if parentId != 0 {
		parentComment := commentDao.GetByParentId(parentId)
		if parentComment.ID == 0 {
			parentComment = commentDao.GetByRootId(parentId)
		}
		comment.ToUserId = parentComment.FromUserId
		comment.RootId = parentComment.RootId
		comment.BusinessId = parentComment.BusinessId
		comment.RootId = parentComment.RootId
	}

	commentDao.AddComment(comment)
	var subscriptionService SubscriptionService
	var b BusinessId
	b.CommentId = comment.BusinessId
	b.UserId = comment.FromUserId
	b.ArticleId = comment.BusinessId
	b.CurrentBusinessId = comment.BusinessId
	subscriptionService.ConstantAtSend(event.CommentAt, comment.FromUserId, comment.Content, b)
	subscriptionService.Do(event.CommentUpdateEvent, b)
	return nil
}

// 删除评论
func (a *CommentsService) DeleteComment(id, userId int) bool {

	return commentDao.Delete(id, userId) == 1
}

// 查询文章下的评论
func (*CommentsService) GetCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {

	var parentComments []*model.Comments
	childCommentsMap := make(map[int][]*model.Comments)
	comments, count := commentDao.GetCommentsByArticleID(page, limit, businessId)
	if count == 0 {
		return parentComments, 0
	}
	parentIds := make([]int, len(comments))
	userIds := make([]int, len(comments))
	// 收集根评论
	for i := range comments {
		comment := comments[i]
		if comment.ParentId == 0 {
			parentComments = append(parentComments, comment)
			parentIds = append(parentIds, comment.ID)
		} else {
			childCommentsMap[comment.RootId] = append(childCommentsMap[comment.RootId], comment)
		}
		userIds = append(userIds, comment.FromUserId)
	}

	setCommentUserInfoAndArticleTitle(comments)

	ChildCommentNumberMap := commentDao.GetCommentsCountByRootId(parentIds)
	for i := range parentComments {
		parentComments[i].ChildComments = childCommentsMap[parentComments[i].RootId]
		parentComments[i].ChildCommentNumber = ChildCommentNumberMap[parentComments[i].RootId]
	}

	return parentComments, count
}

// 设置用户的昵称和头像
func setCommentUserInfoAndArticleTitle(comments []*model.Comments) {
	userIds := mapset.NewSetWithSize[int](len(comments))
	articleIds := mapset.NewSetWithSize[int](len(comments))
	for i := range comments {
		comment := comments[i]
		articleIds.Add(comment.BusinessId)
		userIds.Add(comment.FromUserId)
		userIds.Add(comment.ToUserId)
	}

	if userIds.IsEmpty() {
		return
	}
	var u UserService
	userNameMap := u.ListByIdsToMap(userIds.ToSlice())

	var a ArticleService
	articleTitleMap := a.ListByIdsSelectIdTitleMap(articleIds.ToSlice())

	for i := range comments {
		comment := comments[i]
		comment.ArticleTitle = articleTitleMap[comment.BusinessId]
		comment.FromUserName = userNameMap[comment.FromUserId].Name
		comment.FromUserAvatar = userNameMap[comment.FromUserId].Avatar
		if comment.ParentId != 0 {
			comment.ToUserName = userNameMap[comment.ToUserId].Name
		}
	}
}

// 查询文章下的所有评论(可指定)
func (*CommentsService) GetAllCommentsByArticleID(page, limit, userId, businessId int) ([]*model.Comments, int64) {
	comments, count := commentDao.GetAllCommentsByArticleID(page, userId, limit, businessId)
	if count == 0 {
		return comments, count
	}
	setCommentUserInfoAndArticleTitle(comments)
	return comments, count
}

// 查询根评论下的子评论
func (*CommentsService) GetCommentsByRootID(page, limit, rootId int) ([]*model.Comments, int64) {

	comments, count := commentDao.GetCommentsByCommentID(page, limit, rootId)
	// 如果根评论为空,说明是查询指定根评论下的子评论
	if count == 0 {
		return comments, count
	}

	setCommentUserInfoAndArticleTitle(comments)
	return comments, count
}

func (a *CommentsService) PageComment(p, limit int) (comments []*model.Comments, count int64) {
	model.Comment().Limit(limit).Offset((p - 1) * limit).Find(&comments)

	model.Comment().Count(&count)
	setCommentUserInfoAndArticleTitle(comments)
	return comments, count
}
