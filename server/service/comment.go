package services

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
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

	parentId := comment.ParentId
	var subscriptionService SubscriptionService
	var b SubscribeData

	// 父评论是否存在

	if parentId != 0 {
		parentComment := commentDao.GetByParentId(parentId)
		comment.ToUserId = parentComment.FromUserId
		comment.RootId = parentComment.RootId
		comment.BusinessId = parentComment.BusinessId
		comment.RootId = parentComment.RootId
	}
	commentDao.AddComment(comment)
	b.UserId = comment.FromUserId
	b.ArticleId = comment.BusinessId
	b.CurrentBusinessId = comment.BusinessId
	b.SubscribeId = comment.BusinessId
	b.SectionId = comment.BusinessId
	b.CourseId = comment.BusinessId
	b.CommentId = comment.ID
	eventId := event.CommentUpdateEvent
	// 延迟发送评论事件
	if parentId != 0 {
		subscriptionService.Send(event.ReplyComment, constant.NOTICE, comment.FromUserId, comment.ToUserId, b)
	}
	userId := 0
	if comment.TenantId == 0 {
		var articles ArticleService
		userId = articles.GetById(comment.BusinessId).UserId
	}
	if comment.TenantId == 1 {
		var courS CourseService
		userId = courS.GetCourseSectionDetail(comment.BusinessId).UserId
		eventId = event.SectionComment
	} else if comment.TenantId == 2 {
		var courS CourseService
		userId = courS.GetCourseDetail(comment.BusinessId).UserId
		eventId = event.CourseComment
	}
	subscriptionService.ConstantAtSend(event.CommentAt, comment.FromUserId, comment.Content, b)
	subscriptionService.Do(eventId, b)
	// 文章发布者收到消息
	subscriptionService.Send(eventId, constant.NOTICE, comment.FromUserId, userId, b)
	jsonBody, _ := json.Marshal(comment)
	log.Infof("用户id: %d,发布评论: %s", comment.FromUserId, jsonBody)
	return nil
}

// 删除评论
func (a *CommentsService) DeleteComment(id, userId int) bool {

	log.Infof("用户id: %d,删除评论: %d", userId, id)
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
			comment.ToUserAvatar = userNameMap[comment.ToUserId].Avatar
		}
	}
}

// 查询用户的评论(管理端)
func (*CommentsService) GetAllCommentsByArticleID(page, limit, userId, businessId, tenantId int) ([]*model.Comments, int64) {
	comments, count := commentDao.GetAllCommentsByArticleID(page, limit, userId, businessId, tenantId)
	if count == 0 {
		return comments, count
	}
	setCommentUserInfoAndArticleTitle(comments)
	return comments, count
}

// 查询根评论下的子评论
func (*CommentsService) GetCommentsByRootID(page, limit, rootId int) (comments []*model.Comments, count int64) {

	model.Comment().Where("root_id = ? and id <> root_id", rootId).Count(&count)
	if count == 0 {
		return
	}
	comments = commentDao.GetCommentsByCommentID(page, limit, rootId)

	setCommentUserInfoAndArticleTitle(comments)
	return
}

func (a *CommentsService) PageComment(p, limit int) (comments []*model.Comments, count int64) {
	db := model.Comment()
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((p - 1) * limit).Find(&comments)

	setCommentUserInfoAndArticleTitle(comments)
	return comments, count
}

func (a *CommentsService) Exist(commentId int) bool {
	var count int64
	model.Comment().Where("id = ?", commentId).Count(&count)
	return count == 1
}

func (a *CommentsService) GetById(id int) (comment model.Comments) {
	model.Comment().Where("id = ?", id).Find(&comment)
	return
}

func (a *CommentsService) ListAdoptionsByArticleId(articleId, page, limit int) (comments []*model.Comments, count int64) {
	db := model.Comment().
		Joins("JOIN qa_adoptions ON qa_adoptions.comment_id = comments.id").
		Order("qa_adoptions.created_at DESC"). // 按照采纳时间降序排列
		Where("qa_adoptions.article_id = ?", articleId)
	db.Count(&count)
	db.Limit(limit).Offset((page - 1) * limit)
	db.Find(&comments)
	for i := range comments {
		comments[i].AdoptionState = true
	}
	setCommentUserInfoAndArticleTitle(comments)
	return
}

func (a *CommentsService) ListCommentsByArticleIdNoTree(businessId, tenantId int) (comments []*model.Comments) {

	model.Comment().Where("business_id = ? and tenant_id = ? ", businessId, tenantId).Order("created_at desc").Find(&comments)
	setCommentUserInfoAndArticleTitle(comments)
	return comments
}
