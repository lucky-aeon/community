package email

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"

	"github.com/yuin/goldmark"
)

// 邮件配置 - 实际使用时应从环境变量获取
const (
	realTestEmailAddr = ""
	realTestSMTPHost  = ""
	realTestSMTPUser  = ""
	realTestSMTPPass  = ""
)

// 测试用的真实数据ID - 用户可以修改这些ID值
var testDataIds = struct {
	ArticleId      int // 文章ID
	CommentId      int // 评论ID
	CourseId       int // 课程ID
	SectionId      int // 课程章节ID
	AdoptionId     int // 采纳ID
	ReplyCommentId int // 回复评论ID
}{
	ArticleId:      272,  // 请填写真实的文章ID
	CommentId:      1157, // 请填写真实的评论ID
	CourseId:       21,   // 请填写真实的课程ID
	SectionId:      104,  // 请填写真实的课程章节ID
	AdoptionId:     1,    // 请填写真实的采纳ID
	ReplyCommentId: 1175, // 请填写真实的回复评论ID
}

// 初始化测试环境
func initTestEnv() {
	// 初始化日志系统
	log.Init()

	// 初始化配置
	config.Init()

	// 初始化数据库连接
	mysql.Init("community", "liuzg0815", "124.220.234.136:3306", "community")

	// 初始化邮件服务
	Init(realTestSMTPUser, realTestSMTPPass, realTestSMTPHost)
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

// TestSendArticlePublishEmailWithRealData 测试使用真实数据发送文章发布邮件
func TestSendArticlePublishEmailWithRealData(t *testing.T) {
	initTestEnv()

	// 从数据库获取真实的文章数据
	articleDao := &dao.Article{}
	userDao := &dao.UserDao{}

	// 获取文章信息
	article := articleDao.GetById(testDataIds.ArticleId)
	if article.ID == 0 {
		t.Skipf("文章ID %d 不存在，跳过测试", testDataIds.ArticleId)
		return
	}

	// 获取用户信息
	user := userDao.GetById(article.UserId)
	if user.ID == 0 {
		t.Skipf("用户ID %d 不存在，跳过测试", article.UserId)
		return
	}

	// 处理文章内容：将 Markdown 转换为 HTML
	processedHTML := markdownToHTML(article.Content)

	// 截取内容用于邮件预览（去除HTML标签后截取文本，但邮件中显示HTML）
	plainTextForLimit := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(processedHTML, "")
	if len([]rune(plainTextForLimit)) > 300 {
		// 如果纯文本超过300字符，需要截取HTML
		runes := []rune(plainTextForLimit)
		limitText := string(runes[:300]) + "..."
		// 简单处理：如果内容过长，在邮件中显示截取的纯文本加省略号
		processedHTML = "<p>" + limitText + "</p>"
	}

	// 构建邮件数据
	testData := ArticleData{
		UserName:       user.Name,
		UserAvatar:     "", // 不设置头像，避免防盗链问题
		ArticleTitle:   article.Title,
		ArticleContent: processedHTML, // 使用 HTML 格式内容
		ArticleURL:     fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID),
		PublishTime:    time.Time(article.CreatedAt).Format("2006年01月02日 15:04"),
	}

	// 生成HTML邮件内容
	htmlContent := GenerateArticlePublishHTML(testData)

	// 发送邮件
	err := Send([]string{realTestEmailAddr}, htmlContent, "敲鸭社区 - 文章发布通知 (真实数据测试)")
	if err != nil {
		t.Fatalf("发送文章发布邮件失败: %v", err)
	}

	t.Logf("文章发布邮件已发送到: %s", realTestEmailAddr)
	t.Logf("文章标题: %s", article.Title)
	t.Logf("作者: %s", user.Name)
	t.Log("请检查邮箱查看真实数据邮件效果")
}

// TestSendCommentReplyEmailWithRealData 测试使用真实数据发送评论回复邮件
func TestSendCommentReplyEmailWithRealData(t *testing.T) {
	initTestEnv()

	// 从数据库获取真实的评论数据
	commentDao := &dao.CommentDao{}
	userDao := &dao.UserDao{}
	articleDao := &dao.Article{}

	// 获取回复评论信息
	replyComment := model.Comments{}
	model.Comment().Where("id = ?", testDataIds.ReplyCommentId).First(&replyComment)
	if replyComment.ID == 0 {
		t.Skipf("回复评论ID %d 不存在，跳过测试", testDataIds.ReplyCommentId)
		return
	}

	// 获取被回复的原始评论
	originalComment := commentDao.GetByParentId(replyComment.ParentId)
	if originalComment.ID == 0 {
		t.Skipf("原始评论ID %d 不存在，跳过测试", replyComment.ParentId)
		return
	}

	// 获取回复者信息
	replyUser := userDao.GetById(replyComment.FromUserId)
	if replyUser.ID == 0 {
		t.Skipf("回复用户ID %d 不存在，跳过测试", replyComment.FromUserId)
		return
	}

	// 获取文章信息
	article := articleDao.GetById(replyComment.BusinessId)
	if article.ID == 0 {
		t.Skipf("文章ID %d 不存在，跳过测试", replyComment.BusinessId)
		return
	}

	// 处理评论内容：将 Markdown 转换为 HTML
	replyContentHTML := markdownToHTML(replyComment.Content)
	originalContentHTML := markdownToHTML(originalComment.Content)

	// 构建邮件数据
	testData := CommentReplyData{
		UserName:        replyUser.Name,
		UserAvatar:      "",                  // 不设置头像，避免防盗链问题
		ReplyContent:    replyContentHTML,    // 使用 HTML 格式
		OriginalComment: originalContentHTML, // 使用 HTML 格式
		ArticleTitle:    article.Title,
		ArticleURL:      fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID),
		ReplyTime:       time.Time(replyComment.CreatedAt).Format("2006年01月02日 15:04"),
	}

	// 生成HTML邮件内容
	htmlContent := GenerateCommentReplyHTML(testData)

	// 发送邮件
	err := Send([]string{realTestEmailAddr}, htmlContent, "敲鸭社区 - 评论回复通知 (真实数据测试)")
	if err != nil {
		t.Fatalf("发送评论回复邮件失败: %v", err)
	}

	t.Logf("评论回复邮件已发送到: %s", realTestEmailAddr)
	t.Logf("回复者: %s", replyUser.Name)
	t.Logf("文章标题: %s", article.Title)
	t.Log("请检查邮箱查看真实数据邮件效果")
}

// TestSendAdoptionEmailWithRealData 测试使用真实数据发送采纳邮件
func TestSendAdoptionEmailWithRealData(t *testing.T) {
	initTestEnv()

	// 从数据库获取真实的采纳数据
	adoption := model.QaAdoptions{}
	model.QaAdoption().Where("id = ?", testDataIds.AdoptionId).First(&adoption)
	if adoption.ID == 0 {
		t.Skipf("采纳ID %d 不存在，跳过测试", testDataIds.AdoptionId)
		return
	}

	// 获取被采纳的评论
	comment := model.Comments{}
	model.Comment().Where("id = ?", adoption.CommentId).First(&comment)
	if comment.ID == 0 {
		t.Skipf("评论ID %d 不存在，跳过测试", adoption.CommentId)
		return
	}

	// 获取文章信息
	articleDao := &dao.Article{}
	article := articleDao.GetById(adoption.ArticleId)
	if article.ID == 0 {
		t.Skipf("文章ID %d 不存在，跳过测试", adoption.ArticleId)
		return
	}

	// 处理评论内容：将 Markdown 转换为 HTML
	commentContentHTML := markdownToHTML(comment.Content)

	// 构建邮件数据
	testData := AdoptionData{
		ArticleTitle:   article.Title,
		CommentContent: commentContentHTML, // 使用 HTML 格式
		ArticleURL:     fmt.Sprintf("https://code.xhyovo.cn/article/view?articleId=%d", article.ID),
		AdoptionTime:   time.Time(adoption.CreatedAt).Format("2006年01月02日 15:04"),
	}

	// 生成HTML邮件内容
	htmlContent := GenerateAdoptionHTML(testData)

	// 发送邮件
	err := Send([]string{realTestEmailAddr}, htmlContent, "敲鸭社区 - 评论采纳通知 (真实数据测试)")
	if err != nil {
		t.Fatalf("发送采纳邮件失败: %v", err)
	}

	t.Logf("采纳邮件已发送到: %s", realTestEmailAddr)
	t.Logf("文章标题: %s", article.Title)
	t.Logf("采纳时间: %s", time.Time(adoption.CreatedAt).Format("2006年01月02日 15:04"))
	t.Log("请检查邮箱查看真实数据邮件效果")
}

// TestSendCourseUpdateEmailWithRealData 测试使用真实数据发送课程更新邮件
func TestSendCourseUpdateEmailWithRealData(t *testing.T) {
	initTestEnv()

	// 从数据库获取真实的课程和章节数据
	course := model.Courses{}
	model.Course().Where("id = ?", testDataIds.CourseId).First(&course)
	if course.ID == 0 {
		t.Skipf("课程ID %d 不存在，跳过测试", testDataIds.CourseId)
		return
	}

	// 获取课程章节信息
	section := model.CoursesSections{}
	model.CoursesSection().Where("id = ?", testDataIds.SectionId).First(&section)
	if section.ID == 0 {
		t.Skipf("课程章节ID %d 不存在，跳过测试", testDataIds.SectionId)
		return
	}

	// 构建邮件数据
	testData := CourseUpdateData{
		CourseTitle:  course.Title,
		SectionTitle: section.Title,
		CourseURL:    fmt.Sprintf("https://code.xhyovo.cn/article/view?sectionId=%d", section.ID),
		UpdateTime:   time.Time(section.CreatedAt).Format("2006年01月02日 15:04"),
	}

	// 生成HTML邮件内容
	htmlContent := GenerateCourseUpdateHTML(testData)

	// 发送邮件
	err := Send([]string{realTestEmailAddr}, htmlContent, "敲鸭社区 - 课程更新通知 (真实数据测试)")
	if err != nil {
		t.Fatalf("发送课程更新邮件失败: %v", err)
	}

	t.Logf("课程更新邮件已发送到: %s", realTestEmailAddr)
	t.Logf("课程标题: %s", course.Title)
	t.Logf("章节标题: %s", section.Title)
	t.Log("请检查邮箱查看真实数据邮件效果")
}

// TestSendAllEventTypesWithRealData 测试使用真实数据发送所有类型的邮件
func TestSendAllEventTypesWithRealData(t *testing.T) {
	initTestEnv()

	t.Log("开始使用真实数据批量测试所有事件类型的邮件发送...")

	// 1. 文章发布邮件
	t.Run("ArticlePublish", func(t *testing.T) {
		TestSendArticlePublishEmailWithRealData(t)
	})

	// 2. 评论回复邮件
	t.Run("CommentReply", func(t *testing.T) {
		TestSendCommentReplyEmailWithRealData(t)
	})

	// 3. 采纳邮件
	t.Run("Adoption", func(t *testing.T) {
		TestSendAdoptionEmailWithRealData(t)
	})

	// 4. 课程更新邮件
	t.Run("CourseUpdate", func(t *testing.T) {
		TestSendCourseUpdateEmailWithRealData(t)
	})

	t.Log("真实数据批量邮件测试完成！请检查邮箱查看所有邮件效果")
}

// TestDataValidation 测试数据验证 - 检查配置的ID是否在数据库中存在
func TestDataValidation(t *testing.T) {
	initTestEnv()

	t.Log("开始验证测试数据ID的有效性...")

	// 验证文章ID
	articleDao := &dao.Article{}
	if !articleDao.ExistById(testDataIds.ArticleId) {
		t.Errorf("文章ID %d 不存在，请修改 testDataIds.ArticleId", testDataIds.ArticleId)
	} else {
		t.Logf("✓ 文章ID %d 存在", testDataIds.ArticleId)
	}

	// 验证评论ID
	var commentCount int64
	model.Comment().Where("id = ?", testDataIds.CommentId).Count(&commentCount)
	if commentCount == 0 {
		t.Errorf("评论ID %d 不存在，请修改 testDataIds.CommentId", testDataIds.CommentId)
	} else {
		t.Logf("✓ 评论ID %d 存在", testDataIds.CommentId)
	}

	// 验证课程ID
	var courseCount int64
	model.Course().Where("id = ?", testDataIds.CourseId).Count(&courseCount)
	if courseCount == 0 {
		t.Errorf("课程ID %d 不存在，请修改 testDataIds.CourseId", testDataIds.CourseId)
	} else {
		t.Logf("✓ 课程ID %d 存在", testDataIds.CourseId)
	}

	// 验证课程章节ID
	var sectionCount int64
	model.CoursesSection().Where("id = ?", testDataIds.SectionId).Count(&sectionCount)
	if sectionCount == 0 {
		t.Errorf("课程章节ID %d 不存在，请修改 testDataIds.SectionId", testDataIds.SectionId)
	} else {
		t.Logf("✓ 课程章节ID %d 存在", testDataIds.SectionId)
	}

	// 验证采纳ID
	var adoptionCount int64
	model.QaAdoption().Where("id = ?", testDataIds.AdoptionId).Count(&adoptionCount)
	if adoptionCount == 0 {
		t.Errorf("采纳ID %d 不存在，请修改 testDataIds.AdoptionId", testDataIds.AdoptionId)
	} else {
		t.Logf("✓ 采纳ID %d 存在", testDataIds.AdoptionId)
	}

	// 验证回复评论ID
	var replyCommentCount int64
	model.Comment().Where("id = ?", testDataIds.ReplyCommentId).Count(&replyCommentCount)
	if replyCommentCount == 0 {
		t.Errorf("回复评论ID %d 不存在，请修改 testDataIds.ReplyCommentId", testDataIds.ReplyCommentId)
	} else {
		t.Logf("✓ 回复评论ID %d 存在", testDataIds.ReplyCommentId)
	}

	t.Log("数据验证完成！")
}
