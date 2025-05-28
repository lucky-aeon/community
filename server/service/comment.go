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
	// 发送回复评论事件
	if parentId != 0 {
		subscriptionService.Send(event.ReplyComment, constant.NOTICE, comment.FromUserId, comment.ToUserId, b)
	}
	userId := 0
	// 文章评论
	if comment.TenantId == 0 {
		var articles ArticleService
		userId = articles.GetById(comment.BusinessId).UserId
	}
	// 章节评论
	if comment.TenantId == 1 {
		var courS CourseService
		userId = courS.GetCourseSectionDetail(comment.BusinessId).UserId
		eventId = event.SectionComment
	} else if comment.TenantId == 2 {
		// 课程评论
		var courS CourseService
		userId = courS.GetCourseDetail(comment.BusinessId).UserId
		eventId = event.CourseComment
	} else if comment.TenantId == 3 {
		// 分享会评论
		var meetingS MeetingService
		userId = meetingS.GetById(comment.BusinessId).InitiatorId
		eventId = event.Meeting
	} else if comment.TenantId == 4 {
		// AI日报没有具体的用户所有者，设置userId为0
		userId = 0
		eventId = event.CommentUpdateEvent
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

// ListLatestComments 获取最新的10条评论
func (a *CommentsService) ListLatestComments() ([]*model.Comments, int64) {
	var comments []*model.Comments
	var count int64

	// 查询最新的10条评论
	db := model.Comment().Order("created_at desc").Limit(20)
	db.Count(&count)
	if count == 0 {
		return comments, 0
	}
	db.Find(&comments)

	// 设置评论信息
	setLatestCommentsInfo(comments)

	return comments, count
}

// setLatestCommentsInfo 设置最新评论的相关信息，包括业务标题和用户信息
func setLatestCommentsInfo(comments []*model.Comments) {
	// 收集用户ID和各类业务ID
	userIds := mapset.NewSetWithSize[int](len(comments) * 2) // 评论者和回复者
	articleIds := mapset.NewSetWithSize[int](len(comments))
	sectionIds := mapset.NewSetWithSize[int](len(comments))
	courseIds := mapset.NewSetWithSize[int](len(comments))
	meetingIds := mapset.NewSetWithSize[int](len(comments))
	aiNewsIds := mapset.NewSetWithSize[int](len(comments))

	for i := range comments {
		comment := comments[i]
		userIds.Add(comment.FromUserId)
		userIds.Add(comment.ToUserId)

		// 根据TenantId分类收集业务ID
		switch comment.TenantId {
		case 0: // 文章评论
			articleIds.Add(comment.BusinessId)
		case 1: // 章节评论
			sectionIds.Add(comment.BusinessId)
		case 2: // 课程评论
			courseIds.Add(comment.BusinessId)
		case 3: // 分享会评论
			meetingIds.Add(comment.BusinessId)
		case 4: // AI日报评论
			aiNewsIds.Add(comment.BusinessId)
		}
	}

	// 获取用户信息
	var userService UserService
	userNameMap := userService.ListByIdsToMap(userIds.ToSlice())

	// 获取各类业务标题
	var articleTitleMap map[int]string
	var sectionTitleMap map[int]string
	var courseTitleMap map[int]string
	var meetingTitleMap map[int]string
	var aiNewsTitleMap map[int]string

	if !articleIds.IsEmpty() {
		var articleService ArticleService
		articleTitleMap = articleService.ListByIdsSelectIdTitleMap(articleIds.ToSlice())
	}

	if !sectionIds.IsEmpty() {
		var courseService CourseService
		sectionTitleMap = courseService.ListSectionByIds(sectionIds.ToSlice())
	}

	if !courseIds.IsEmpty() {
		var courseService CourseService
		courseTitleMap = courseService.ListByIdsSelectIdTitleMap(courseIds.ToSlice())
	}

	if !meetingIds.IsEmpty() {
		// 分享会标题需要单独处理
		meetingTitleMap = make(map[int]string)
		var meetings []model.Meetings
		model.Meeting().Where("id in ?", meetingIds.ToSlice()).Select("id, title").Find(&meetings)
		for _, meeting := range meetings {
			meetingTitleMap[meeting.Id] = meeting.Title
		}
	}

	if !aiNewsIds.IsEmpty() {
		// AI日报标题处理
		aiNewsTitleMap = make(map[int]string)
		var aiNewsService AiNewsService
		for _, id := range aiNewsIds.ToSlice() {
			article, err := aiNewsService.GetNewsById(id)
			if err == nil {
				aiNewsTitleMap[id] = article.Title
			}
		}
	}

	// 设置评论相关信息
	for i := range comments {
		comment := comments[i]
		// 设置用户信息
		comment.FromUserName = userNameMap[comment.FromUserId].Name
		comment.FromUserAvatar = userNameMap[comment.FromUserId].Avatar
		if comment.ParentId != 0 && comment.ToUserId > 0 {
			comment.ToUserName = userNameMap[comment.ToUserId].Name
			comment.ToUserAvatar = userNameMap[comment.ToUserId].Avatar
		}

		// 设置业务标题
		switch comment.TenantId {
		case 0: // 文章评论
			comment.ArticleTitle = articleTitleMap[comment.BusinessId]
		case 1: // 章节评论
			comment.ArticleTitle = sectionTitleMap[comment.BusinessId]
		case 2: // 课程评论
			comment.ArticleTitle = courseTitleMap[comment.BusinessId]
		case 3: // 分享会评论
			comment.ArticleTitle = meetingTitleMap[comment.BusinessId]
		case 4: // AI日报评论
			comment.ArticleTitle = aiNewsTitleMap[comment.BusinessId]
		}
	}
}
