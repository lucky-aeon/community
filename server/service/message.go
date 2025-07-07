package services

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/mysql"
	localTime "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"

	"github.com/yuin/goldmark"
)

var messageTemplateVar = make(map[string]map[string]string)

func init() {
	messageTemplateVar["user"] = map[string]string{
		"user.id":      "id",
		"user.name":    "name",
		"user.account": "account",
		"user.avatar":  "avatar",
	}
	messageTemplateVar["article"] = map[string]string{
		"article.id":      "id",
		"article.title":   "title",
		"article.content": "content",
		"article.userId":  "user_id",
	}
	messageTemplateVar["comment"] = map[string]string{
		"comment.id":           "id",
		"comment.content":      "content",
		"comment.FromUserName": "from_user_name",
		"comment.ToUserName":   "to_user_name",
		"comment.ArticleTitle": "article_title",
	}
	messageTemplateVar["course"] = map[string]string{
		"course.id":    "id",
		"course.title": "title",
	}
	messageTemplateVar["courses_section"] = map[string]string{
		"courses_section.id":    "id",
		"courses_section.title": "title",
	}
}

// 发送消息中消息模板需要用到的业务id
type SubscribeData struct {
	ArticleId         int
	UserId            int
	CommentId         int
	CurrentBusinessId int // 当前主业务id
	SubscribeId       int // 订阅业务的id(在消息中可以点击跳转的)
	SectionId         int // 章节id
	CourseId          int // 课程id
}

type MessageService struct {
}

func (*MessageService) ListMessageTemplate(page, limit int) ([]*model.MessageTemplates, int64) {
	var count int64
	model.MessageTemplate().Count(&count)
	templates := messageDao.ListMessageTemplate(page, limit)
	eventMap := event.Map()
	for i := range templates {
		templates[i].EventName = eventMap[templates[i].EventId]
	}

	return templates, count
}

func (*MessageService) SaveMessageTemplate(template model.MessageTemplates) error {
	if err := messageDao.SaveMessageTemplate(template); err != nil {
		return errors.New("创建消息模板对应的事件已经存在")
	}
	return nil
}

func (*MessageService) DeleteMessageTemplate(id int) {
	messageDao.DeleteMessageTemplate(id)
}

func (*MessageService) AddMessageLogs(from, types, eventId, businessId int, to []int, content string) {
	var messageLogs []*model.MessageLogs
	for i := range to {
		log := &model.MessageLogs{
			From:      from,
			To:        i,
			Content:   content,
			Type:      types,
			ArticleId: businessId,
			EventId:   eventId,
		}
		messageLogs = append(messageLogs, log)
	}
	messageDao.SaveMessageLogs(messageLogs)
}

func (*MessageService) DeleteMessageLogs(id []int) {
	messageDao.DeleteMessageLogs(id)
}

func (m *MessageService) SendMessage(from, to, types, eventId, businessId int, content string) {

	messageDao.SaveMessage(from, types, eventId, businessId, []int{to}, content)
	// 添加记录
	m.AddMessageLogs(from, types, eventId, businessId, []int{to}, content)
}

func (m *MessageService) SendMessages(from, types, eventId, businessId int, to []int, content string) {

	messageDao.SaveMessage(from, types, eventId, businessId, to, content)
	// 添加记录
	m.AddMessageLogs(from, types, eventId, businessId, to, content)
}

func (*MessageService) ReadMessage2(typee, eventId, businessId, userId int) int64 {
	return messageDao.ReadMessage2(typee, eventId, businessId, userId)
}

func (*MessageService) ReadMessage(id []int, userId int) int64 {
	return messageDao.ReadMessage(id, userId)
}

func (m *MessageService) PageMessage(page, limit, userId, types, state int) (msgs []*model.MessageStates, count int64) {
	msgs = messageDao.ListMessage(page, limit, userId, types, state)
	count = messageDao.CountMessage(userId, types, state)
	return msgs, count
}

// 人：你订阅的 xxx 用户发布了文章，文章标题：xxx
// 文章：你订阅的 xxx 文章，被xxx评论了，评论内容：xxx
func (m *MessageService) GetMsg(template string, b SubscribeData) string {
	return m.GetMsgWithEventId(template, b, 0) // 兼容旧调用，eventId=0表示使用字符串判断
}

// GetMsgWithEventId 根据事件ID生成消息内容
func (m *MessageService) GetMsgWithEventId(template string, b SubscribeData, eventId int) string {
	// 如果有 eventId，优先使用事件ID判断
	if eventId > 0 {
		switch eventId {
		case event.CommentUpdateEvent: // 文章下评论更新事件
			return m.generateCommentUpdateHTML(b)
		case event.UserFollowingEvent: // 用户关注的人事件（文章发布）
			return m.generateArticlePublishHTML(b)
		case event.ArticleAt: // 文章中@
			return m.generateArticleAtHTML(b)
		case event.CommentAt: // 评论中@
			return m.generateCommentAtHTML(b)
		case event.ReplyComment: // 评论回复
			return m.generateCommentReplyHTML(b)
		case event.Adoption: // 采纳
			return m.generateAdoptionHTML(b)
		case event.SectionComment: // 章节评论
			return m.generateSectionCommentHTML(b)
		case event.CourseComment: // 课程回复
			return m.generateCourseCommentHTML(b)
		case event.CourseUpdate: // 课程更新
			return m.generateCourseUpdateHTML(b)
		case event.Meeting: // 会议
			return m.generateMeetingHTML(b)
		}
	}

	// 兼容原有字符串判断逻辑（用于旧代码或者eventId=0的情况）
	if email.IsUserUpdateEvent(template) {
		return m.generateArticlePublishHTML(b)
	}

	if email.IsCommentReplyEvent(template) {
		return m.generateCommentReplyHTML(b)
	}

	if email.IsAdoptionEvent(template) {
		return m.generateAdoptionHTML(b)
	}

	if email.IsCourseUpdateEvent(template) {
		return m.generateCourseUpdateHTML(b)
	}

	// 其他事件保持原有逻辑
	BusinessIdMap := businessIdToMap(b)
	for s, v := range messageTemplateVar {
		// 拼接 ${ + s + "." 如果存在则找
		str := fmt.Sprintf("${%s.", s)
		if strings.Contains(template, str) {
			var objet map[string]interface{}
			mysql.GetInstance().Table(s+"s").Where("id = ?", BusinessIdMap[s]).Find(&objet)
			// 遍历 v 从key找template
			for s2 := range v {
				varTemlp := fmt.Sprintf("${%s}", s2)
				if strings.Contains(template, varTemlp) {
					i := objet[v[s2]]
					// 检查值是否为nil
					var valueStr string
					if i != nil {
						valueStr = fmt.Sprintf("%s", i)
					} else {
						valueStr = "" // 如果是nil，替换为空字符串
					}
					template = strings.ReplaceAll(template, varTemlp, valueStr)
				}
			}
		}
	}
	return template
}

func businessIdToMap(b SubscribeData) map[string]int {
	m := map[string]int{
		"user":            b.UserId,
		"article":         b.ArticleId,
		"comment":         b.CommentId,
		"course":          b.CourseId,
		"courses_section": b.SectionId,
	}
	return m
}

func (*MessageService) GetMessageTemplateVar() map[string]map[string]string {
	return messageTemplateVar
}

func (m *MessageService) ClearUnReadMessage(msgType, userId int) {
	model.MessageState().Where("`to` = ? and type = ?", userId, msgType).Update("state", 0)
}

func (m *MessageService) GetUnReadMessageCountByUserId(userId int) (count int64) {
	model.MessageState().Where("`to` = ? and state = 1", userId).Count(&count)
	return
}

// generateArticlePublishHTML 生成文章/章节发布HTML邮件内容
func (m *MessageService) generateArticlePublishHTML(b SubscribeData) string {
	// 智能判断是文章发布还是章节发布
	if b.SectionId > 0 && b.CourseId > 0 && b.ArticleId == 0 {
		// 这是章节发布事件，获取课程和章节信息
		return m.generateSectionPublishHTML(b)
	}

	// 这是文章发布事件，获取文章信息
	var article model.Articles
	err := mysql.GetInstance().Table("articles").Where("id = ?", b.ArticleId).First(&article).Error
	if err != nil {
		// 如果查询失败，使用默认数据
		article.ID = b.ArticleId
		article.Title = "文章标题获取失败"
		article.Content = "文章内容获取失败，请点击查看详情"
	}

	// 获取用户详细信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		// 如果查询失败，使用默认数据
		user.Name = "用户"
		user.Avatar = ""
	}
	// 生成文章链接
	articleURL := fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID)

	// 格式化发布时间
	publishTime := "刚刚"
	createdTime := time.Time(article.CreatedAt)
	if !createdTime.IsZero() {
		publishTime = createdTime.Format("2006年01月02日 15:04")
	}

	// 准备模板数据
	templateData := email.ArticleData{
		UserName:       user.Name,
		UserAvatar:     "", // 不设置头像，避免防盗链问题
		ArticleTitle:   article.Title,
		ArticleContent: getArticlePreview(article),
		ArticleURL:     articleURL,
		PublishTime:    publishTime,
	}

	return email.GenerateArticlePublishHTML(templateData)
}

// generateSectionPublishHTML 生成章节发布HTML邮件内容
func (m *MessageService) generateSectionPublishHTML(b SubscribeData) string {
	// 获取课程信息
	var course model.Courses
	err := mysql.GetInstance().Table("courses").Where("id = ?", b.CourseId).First(&course).Error
	if err != nil {
		course.ID = b.CourseId
		course.Title = "课程标题获取失败"
		course.Desc = "课程描述获取失败，请点击查看详情"
	}

	// 获取章节信息
	var section model.CoursesSections
	err = mysql.GetInstance().Table("courses_sections").Where("id = ?", b.SectionId).First(&section).Error
	if err != nil {
		section.ID = b.SectionId
		section.Title = "章节标题获取失败"
		section.Content = "章节内容获取失败，请点击查看详情"
	}

	// 获取用户详细信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		user.Name = "用户"
		user.Avatar = ""
	}

	// 生成章节链接
	sectionURL := fmt.Sprintf("https://code.xhyovo.cn/article/view?sectionId=%d", section.ID)

	// 格式化发布时间
	publishTime := "刚刚"
	createdTime := time.Time(section.CreatedAt)
	if !createdTime.IsZero() {
		publishTime = createdTime.Format("2006年01月02日 15:04")
	}

	// 准备章节发布邮件模板数据
	templateData := email.SectionPublishData{
		UserName:       user.Name,
		UserAvatar:     "", // 不设置头像，避免防盗链问题
		CourseTitle:    course.Title,
		SectionTitle:   section.Title,
		SectionContent: getSectionPreview(section),
		SectionURL:     sectionURL,
		PublishTime:    publishTime,
	}

	return email.GenerateSectionPublishHTML(templateData)
}

// markdownToHTML 使用 goldmark 库将 Markdown 转换为 HTML
func markdownToHTML(markdown string) string {
	if markdown == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := goldmark.New().Convert([]byte(markdown), &buf); err != nil {
		// 如果转换失败，返回原文本
		return markdown
	}

	return buf.String()
}

// getArticlePreview 获取文章预览内容
func getArticlePreview(article model.Articles) string {
	content := article.Content
	if content == "" {
		content = article.Abstract
	}
	if content == "" {
		content = "点击查看文章详情..."
		return content
	}

	// 将 Markdown 转换为 HTML
	htmlContent := markdownToHTML(content)

	// 截取内容用于邮件预览（去除HTML标签后截取文本，但邮件中显示HTML）
	plainTextForLimit := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(htmlContent, "")
	if len([]rune(plainTextForLimit)) > 300 {
		// 如果纯文本超过300字符，需要截取HTML
		runes := []rune(plainTextForLimit)
		limitText := string(runes[:300]) + "..."
		// 简单处理：如果内容过长，在邮件中显示截取的纯文本加省略号
		htmlContent = "<p>" + limitText + "</p>"
	}

	return htmlContent
}

// getSectionPreview 获取章节预览内容
func getSectionPreview(section model.CoursesSections) string {
	content := section.Content
	if content == "" {
		content = "点击查看章节详情..."
		return content
	}

	// 将 Markdown 转换为 HTML
	htmlContent := markdownToHTML(content)

	// 截取内容用于邮件预览（去除HTML标签后截取文本，但邮件中显示HTML）
	plainTextForLimit := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(htmlContent, "")
	if len([]rune(plainTextForLimit)) > 200 {
		// 如果纯文本超过200字符，需要截取HTML
		runes := []rune(plainTextForLimit)
		limitText := string(runes[:200]) + "..."
		// 简单处理：如果内容过长，在邮件中显示截取的纯文本加省略号
		htmlContent = "<p>" + limitText + "</p>"
	}

	return htmlContent
}

// generateCommentReplyHTML 生成评论回复HTML邮件内容
func (m *MessageService) generateCommentReplyHTML(b SubscribeData) string {
	// 获取评论信息
	var comment model.Comments
	err := mysql.GetInstance().Table("comments").Where("id = ?", b.CommentId).First(&comment).Error
	if err != nil {
		comment.Content = "评论内容获取失败"
	}

	// 获取回复用户信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		user.Name = "用户"
	}

	// 获取被回复的原始评论内容
	originalComment := "你的评论"
	if comment.ParentId > 0 {
		var parentComment model.Comments
		err = mysql.GetInstance().Table("comments").Where("id = ?", comment.ParentId).First(&parentComment).Error
		if err == nil {
			originalComment = markdownToHTML(parentComment.Content)
		}
	}

	// 获取文章信息
	var article model.Articles
	articleURL := ""
	if b.ArticleId > 0 {
		mysql.GetInstance().Table("articles").Where("id = ?", b.ArticleId).First(&article)
		articleURL = fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID)
	}

	// 格式化时间
	replyTime := "刚刚"
	if comment.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(comment.CreatedAt)
		if !createdTime.IsZero() {
			replyTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.CommentReplyData{
		UserName:        user.Name,
		UserAvatar:      "",                              // 不设置头像，避免防盗链问题
		ReplyContent:    markdownToHTML(comment.Content), // 转换 Markdown 为 HTML
		OriginalComment: originalComment,                 // 获取实际的被回复评论内容
		ArticleTitle:    article.Title,
		ArticleURL:      articleURL,
		ReplyTime:       replyTime,
	}

	return email.GenerateCommentReplyHTML(templateData)
}

// generateAdoptionHTML 生成采纳HTML邮件内容
func (m *MessageService) generateAdoptionHTML(b SubscribeData) string {
	// 获取评论信息
	var comment model.Comments
	err := mysql.GetInstance().Table("comments").Where("id = ?", b.CommentId).First(&comment).Error
	if err != nil {
		comment.Content = "评论内容获取失败"
	}

	// 获取文章信息
	var article model.Articles
	err = mysql.GetInstance().Table("articles").Where("id = ?", b.ArticleId).First(&article).Error
	if err != nil {
		article.Title = "文章标题获取失败"
	}

	// 生成文章链接
	articleURL := fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID)

	// 格式化采纳时间
	adoptionTime := "刚刚"
	if comment.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(comment.CreatedAt)
		if !createdTime.IsZero() {
			adoptionTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.AdoptionData{
		ArticleTitle:   article.Title,
		CommentContent: markdownToHTML(comment.Content), // 转换 Markdown 为 HTML
		ArticleURL:     articleURL,
		AdoptionTime:   adoptionTime,
	}

	return email.GenerateAdoptionHTML(templateData)
}

// generateCourseUpdateHTML 生成课程更新HTML邮件内容
func (m *MessageService) generateCourseUpdateHTML(b SubscribeData) string {
	// 获取课程信息
	var course model.Courses
	err := mysql.GetInstance().Table("courses").Where("id = ?", b.CourseId).First(&course).Error
	if err != nil {
		course.Title = "课程标题获取失败"
	}

	// 获取章节信息
	var section model.CoursesSections
	err = mysql.GetInstance().Table("courses_sections").Where("id = ?", b.SectionId).First(&section).Error
	if err != nil {
		section.Title = "章节标题获取失败"
	}

	// 生成章节链接
	courseURL := fmt.Sprintf("https://code.xhyovo.cn/article/view?sectionId=%d", section.ID)

	// 格式化更新时间
	updateTime := "刚刚"
	if section.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(section.CreatedAt)
		if !createdTime.IsZero() {
			updateTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.CourseUpdateData{
		CourseTitle:  course.Title,
		SectionTitle: section.Title,
		CourseURL:    courseURL,
		UpdateTime:   updateTime,
	}

	return email.GenerateCourseUpdateHTML(templateData)
}

// generateCommentUpdateHTML 生成文章评论更新HTML邮件内容
func (m *MessageService) generateCommentUpdateHTML(b SubscribeData) string {
	// 获取评论信息
	var comment model.Comments
	err := mysql.GetInstance().Table("comments").Where("id = ?", b.CommentId).First(&comment).Error
	if err != nil {
		comment.Content = "评论内容获取失败"
	}

	// 获取评论用户信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		user.Name = "用户"
	}

	// 获取文章信息
	var article model.Articles
	articleURL := ""
	if b.ArticleId > 0 {
		mysql.GetInstance().Table("articles").Where("id = ?", b.ArticleId).First(&article)
		articleURL = fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID)
	}

	// 格式化时间
	commentTime := "刚刚"
	if comment.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(comment.CreatedAt)
		if !createdTime.IsZero() {
			commentTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.ArticleCommentData{
		UserName:       user.Name,
		UserAvatar:     "",                              // 不设置头像，避免防盗链问题
		CommentContent: markdownToHTML(comment.Content), // 转换 Markdown 为 HTML
		ArticleTitle:   article.Title,
		ArticleURL:     articleURL,
		CommentTime:    commentTime,
	}

	return email.GenerateArticleCommentHTML(templateData)
}

// generateArticleAtHTML 生成文章@HTML邮件内容
func (m *MessageService) generateArticleAtHTML(b SubscribeData) string {
	// 复用文章发布模板，@事件本质上是提醒用户关注某篇文章
	return m.generateArticlePublishHTML(b)
}

// generateCommentAtHTML 生成评论@HTML邮件内容
func (m *MessageService) generateCommentAtHTML(b SubscribeData) string {
	// 复用评论回复模板，@事件本质上是评论中提到了用户
	return m.generateCommentReplyHTML(b)
}

// generateSectionCommentHTML 生成章节评论HTML邮件内容
func (m *MessageService) generateSectionCommentHTML(b SubscribeData) string {
	// 获取评论信息
	var comment model.Comments
	err := mysql.GetInstance().Table("comments").Where("id = ?", b.CommentId).First(&comment).Error
	if err != nil {
		comment.Content = "评论内容获取失败"
	}

	// 获取评论用户信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		user.Name = "用户"
	}

	// 获取章节信息
	var section model.CoursesSections
	var course model.Courses
	sectionURL := ""
	if b.SectionId > 0 {
		mysql.GetInstance().Table("courses_sections").Where("id = ?", b.SectionId).First(&section)
		// 获取所属课程信息
		mysql.GetInstance().Table("courses").Where("id = ?", section.CourseId).First(&course)
		sectionURL = fmt.Sprintf("https://code.xhyovo.cn/article/view?sectionId=%d", section.ID)
	}

	// 格式化时间
	commentTime := "刚刚"
	if comment.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(comment.CreatedAt)
		if !createdTime.IsZero() {
			commentTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.SectionCommentData{
		UserName:       user.Name,
		UserAvatar:     "",                              // 不设置头像，避免防盗链问题
		CommentContent: markdownToHTML(comment.Content), // 转换 Markdown 为 HTML
		SectionTitle:   section.Title,
		CourseTitle:    course.Title,
		SectionURL:     sectionURL,
		CommentTime:    commentTime,
	}

	return email.GenerateSectionCommentHTML(templateData)
}

// generateCourseCommentHTML 生成课程评论HTML邮件内容
func (m *MessageService) generateCourseCommentHTML(b SubscribeData) string {
	// 获取评论信息
	var comment model.Comments
	err := mysql.GetInstance().Table("comments").Where("id = ?", b.CommentId).First(&comment).Error
	if err != nil {
		comment.Content = "评论内容获取失败"
	}

	// 获取评论用户信息
	var user model.Users
	err = mysql.GetInstance().Table("users").Where("id = ?", b.UserId).First(&user).Error
	if err != nil {
		user.Name = "用户"
	}

	// 获取课程信息
	var course model.Courses
	courseURL := ""
	if b.CourseId > 0 {
		mysql.GetInstance().Table("courses").Where("id = ?", b.CourseId).First(&course)
		courseURL = fmt.Sprintf("https://code.xhyovo.cn/course/view?courseId=%d", course.ID)
	}

	// 格式化时间
	commentTime := "刚刚"
	if comment.CreatedAt != (localTime.LocalTime{}) {
		createdTime := time.Time(comment.CreatedAt)
		if !createdTime.IsZero() {
			commentTime = createdTime.Format("2006年01月02日 15:04")
		}
	}

	templateData := email.CourseCommentData{
		UserName:       user.Name,
		UserAvatar:     "",                              // 不设置头像，避免防盗链问题
		CommentContent: markdownToHTML(comment.Content), // 转换 Markdown 为 HTML
		CourseTitle:    course.Title,
		CourseURL:      courseURL,
		CommentTime:    commentTime,
	}

	return email.GenerateCourseCommentHTML(templateData)
}

// generateMeetingHTML 生成会议HTML邮件内容
func (m *MessageService) generateMeetingHTML(b SubscribeData) string {
	// 获取会议信息（假设有会议表）
	// 这里需要根据实际的会议数据结构来实现
	// 暂时返回一个简单的会议通知

	templateData := email.ArticleData{
		UserName:       "系统",
		UserAvatar:     "", // 不设置头像，避免防盗链问题
		ArticleTitle:   "敲鸭社区 - 会议通知",
		ArticleContent: "<p>您有新的会议通知，请及时查看。</p>",
		ArticleURL:     "https://code.xhyovo.cn/meeting",
		PublishTime:    time.Now().Format("2006年01月02日 15:04"),
	}

	return email.GenerateArticlePublishHTML(templateData)
}
