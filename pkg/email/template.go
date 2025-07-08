package email

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
)

// ArticleData 文章邮件模板数据
type ArticleData struct {
	UserName       string
	UserAvatar     string
	ArticleTitle   string
	ArticleContent string
	ArticleURL     string
	PublishTime    string
}

// CommentReplyData 评论回复邮件模板数据
type CommentReplyData struct {
	UserName        string
	UserAvatar      string
	ReplyContent    string
	OriginalComment string
	ArticleTitle    string
	ArticleURL      string
	ReplyTime       string
}

// ArticleCommentData 文章评论邮件模板数据
type ArticleCommentData struct {
	UserName       string
	UserAvatar     string
	CommentContent string
	ArticleTitle   string
	ArticleURL     string
	CommentTime    string
}

// SectionCommentData 章节评论邮件模板数据
type SectionCommentData struct {
	UserName       string
	UserAvatar     string
	CommentContent string
	SectionTitle   string
	CourseTitle    string
	SectionURL     string
	CommentTime    string
}

// CourseCommentData 课程评论邮件模板数据
type CourseCommentData struct {
	UserName       string
	UserAvatar     string
	CommentContent string
	CourseTitle    string
	CourseURL      string
	CommentTime    string
}

// AdoptionData 采纳邮件模板数据
type AdoptionData struct {
	ArticleTitle   string
	CommentContent string
	ArticleURL     string
	AdoptionTime   string
}

// CourseUpdateData 课程更新邮件模板数据
type CourseUpdateData struct {
	CourseTitle  string
	SectionTitle string
	CourseURL    string
	UpdateTime   string
}

// SectionPublishData 章节发布邮件模板数据
type SectionPublishData struct {
	UserName       string
	UserAvatar     string
	CourseTitle    string
	SectionTitle   string
	SectionContent string
	SectionURL     string
	PublishTime    string
}

// GenerateArticlePublishHTML 生成文章发布邮件HTML模板
func GenerateArticlePublishHTML(data ArticleData) string {
	// 转义HTML内容防止XSS（但不转义 articleContent，因为它已经是安全的HTML）
	userName := html.EscapeString(data.UserName)
	articleTitle := html.EscapeString(data.ArticleTitle)
	articleContent := data.ArticleContent // 直接使用，不转义（假设已经是安全的HTML）
	articleURL := html.EscapeString(data.ArticleURL)
	publishTime := html.EscapeString(data.PublishTime)

	// 如果HTML内容过长，截取（基于去除标签后的文本长度）
	plainTextForLimit := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(articleContent, "")
	if len([]rune(plainTextForLimit)) > 200 {
		runes := []rune(plainTextForLimit)
		limitText := string(runes[:200]) + "..."
		articleContent = "<p>" + html.EscapeString(limitText) + "</p>"
	}

	// 不再处理头像

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>敲鸭社区 - 新文章通知</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
        }
        
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 24px;
            font-weight: 600;
            margin: 0;
        }
        
        .content {
            padding: 30px;
        }
        
        .author-info {
            display: flex;
            align-items: center;
            margin-bottom: 25px;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        
        .author-avatar {
            width: 50px;
            height: 50px;
            border-radius: 50%%;
            margin-right: 15px;
            border: 3px solid #fff;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        
        .author-details h3 {
            color: #333;
            font-size: 16px;
            margin-bottom: 4px;
        }
        
        .author-details p {
            color: #666;
            font-size: 14px;
            margin: 0;
        }
        
        .article-content {
            background-color: #fff;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 25px;
            margin-bottom: 25px;
        }
        
        .article-title {
            font-size: 22px;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 15px;
            line-height: 1.4;
        }
        
        .article-preview {
            color: #555;
            font-size: 15px;
            line-height: 1.7;
            margin-bottom: 15px;
        }
        
        .article-meta {
            color: #888;
            font-size: 13px;
            margin-bottom: 20px;
        }
        
        .cta-button {
            display: inline-block;
            padding: 12px 30px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white !important;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            font-size: 14px;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
        }
        
        .footer {
            background-color: #f8f9fa;
            padding: 20px;
            text-align: center;
            border-top: 1px solid #e9ecef;
        }
        
        .footer p {
            color: #6c757d;
            font-size: 12px;
            margin: 5px 0;
        }
        
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
        
        @media (max-width: 600px) {
            .email-container {
                margin: 0;
                box-shadow: none;
            }
            
            .content {
                padding: 20px;
            }
            
            .header {
                padding: 15px;
            }
            
            .header h1 {
                font-size: 20px;
            }
            
            .author-info {
                padding: 12px;
            }
            
            .article-content {
                padding: 18px;
            }
            
            .article-title {
                font-size: 18px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        
        <div class="content">
            <div class="author-info">
                <div class="author-details">
                    <h3>%s</h3>
                    <p>发布了新文章</p>
                </div>
            </div>
            
            <div class="article-content">
                <h2 class="article-title">%s</h2>
                
                <div class="article-preview">
                    %s
                </div>
                
                <div class="article-meta">
                    📅 发布于: %s
                </div>
                
                <a href="%s" class="cta-button">
                    📖 查看完整文章
                </a>
            </div>
        </div>
        
        <div class="footer">
            <p>感谢您关注敲鸭社区！</p>
            <p>
                <a href="https://code.xhyovo.cn">访问社区</a> | 
                <a href="#">邮件偏好设置</a>
            </p>
            <p>© %d 敲鸭社区 - 专注于技术分享与交流</p>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate,
		userName,
		articleTitle, articleContent, publishTime,
		articleURL, time.Now().Year())
}

// IsUserUpdateEvent 判断是否为用户更新事件
func IsUserUpdateEvent(template string) bool {
	return strings.Contains(template, "${user.name}") &&
		strings.Contains(template, "${article.title}") &&
		strings.Contains(template, "发布了最新文章")
}

// IsCommentReplyEvent 判断是否为评论回复事件
func IsCommentReplyEvent(template string) bool {
	return strings.Contains(template, "回复了你的评论") &&
		strings.Contains(template, "${comment.content}")
}

// IsArticleCommentEvent 判断是否为文章评论事件
func IsArticleCommentEvent(template string) bool {
	return strings.Contains(template, "有最新评论了") &&
		strings.Contains(template, "${comment.content}")
}

// IsAdoptionEvent 判断是否为采纳事件
func IsAdoptionEvent(template string) bool {
	return strings.Contains(template, "被采纳") &&
		strings.Contains(template, "${comment.content}")
}

// IsCourseUpdateEvent 判断是否为课程更新事件
func IsCourseUpdateEvent(template string) bool {
	return strings.Contains(template, "更新了章节") &&
		strings.Contains(template, "${courses_section.title}")
}

// GenerateCommentReplyHTML 生成评论回复邮件HTML模板
func GenerateCommentReplyHTML(data CommentReplyData) string {
	userName := html.EscapeString(data.UserName)
	replyContent := data.ReplyContent       // 回复内容（HTML格式）
	originalComment := data.OriginalComment // 被回复的评论内容（HTML格式）
	articleTitle := html.EscapeString(data.ArticleTitle)
	articleURL := html.EscapeString(data.ArticleURL)
	replyTime := html.EscapeString(data.ReplyTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 评论回复通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .reply-info { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="reply-info">
                <h3>💬 %s 回复了你在《%s》中的评论</h3>
                <p><strong>回复时间：</strong>%s</p>
            </div>
            
            <!-- 显示被回复的评论 -->
            <div style="background: #e9ecef; padding: 15px; border-radius: 8px; margin: 20px 0; border-left: 3px solid #667eea;">
                <h4 style="margin-top: 0; color: #495057;">你的评论：</h4>
                <div style="color: #6c757d;">%s</div>
            </div>
            
            <!-- 显示回复内容 -->
            <div style="background: #fff3cd; padding: 15px; border-radius: 8px; margin: 20px 0; border-left: 3px solid #ffc107;">
                <h4 style="margin-top: 0; color: #856404;">%s 的回复：</h4>
                <div style="color: #856404;">%s</div>
            </div>
            
            <div style="text-align: center;">
                <a href="%s" class="cta-button">💬 查看完整对话</a>
            </div>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, userName, articleTitle, replyTime, originalComment, userName, replyContent, articleURL)
}

// GenerateAdoptionHTML 生成采纳邮件HTML模板
func GenerateAdoptionHTML(data AdoptionData) string {
	articleTitle := html.EscapeString(data.ArticleTitle)
	commentContent := data.CommentContent // 直接使用，不转义（假设已经是安全的HTML）
	articleURL := html.EscapeString(data.ArticleURL)
	adoptionTime := html.EscapeString(data.AdoptionTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 评论采纳通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #28a745 0%%, #20c997 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .adoption-badge { background: #d4edda; color: #155724; padding: 15px; border-radius: 8px; text-align: center; margin-bottom: 20px; }
        .comment-content { background: #f8f9fa; padding: 15px; border-left: 3px solid #28a745; margin: 15px 0; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #28a745 0%%, #20c997 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="adoption-badge">
                <h2>🎉 恭喜！你的评论被采纳了</h2>
            </div>
            <p><strong>文章：</strong>%s</p>
            <p><strong>采纳时间：</strong>%s</p>
            <div class="comment-content">
                <p><strong>被采纳的评论：</strong></p>
                <p>%s</p>
            </div>
            <p>感谢你的精彩回答！这对社区其他成员很有帮助。</p>
            <a href="%s" class="cta-button">🏆 查看详情</a>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, articleTitle, adoptionTime, commentContent, articleURL)
}

// GenerateCourseUpdateHTML 生成课程更新邮件HTML模板
func GenerateCourseUpdateHTML(data CourseUpdateData) string {
	courseTitle := html.EscapeString(data.CourseTitle)
	sectionTitle := html.EscapeString(data.SectionTitle)
	courseURL := html.EscapeString(data.CourseURL)
	updateTime := html.EscapeString(data.UpdateTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 课程更新通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #fd7e14 0%%, #ffc107 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .course-info { background: #fff3cd; border: 1px solid #ffeaa7; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .new-section { background: #f8f9fa; padding: 15px; border-left: 3px solid #fd7e14; margin: 15px 0; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #fd7e14 0%%, #ffc107 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="course-info">
                <h3>📚 课程有新内容啦！</h3>
                <p><strong>课程：</strong>%s</p>
                <p><strong>更新时间：</strong>%s</p>
            </div>
            <div class="new-section">
                <p><strong>新增章节：</strong></p>
                <h4>%s</h4>
            </div>
            <p>快来学习最新的内容吧！</p>
            <a href="%s" class="cta-button">📖 继续学习</a>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, courseTitle, updateTime, sectionTitle, courseURL)
}

// GenerateArticleCommentHTML 生成文章评论邮件HTML模板
func GenerateArticleCommentHTML(data ArticleCommentData) string {
	userName := html.EscapeString(data.UserName)
	commentContent := data.CommentContent // 直接使用，不转义（假设已经是安全的HTML）
	articleTitle := html.EscapeString(data.ArticleTitle)
	articleURL := html.EscapeString(data.ArticleURL)
	commentTime := html.EscapeString(data.CommentTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 文章评论通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .comment-info { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .comment-content { background: #e9ecef; padding: 15px; border-left: 3px solid #667eea; margin: 15px 0; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="comment-info">
                <h3>💬 %s 评论了你的文章</h3>
                <p><strong>文章：</strong>%s</p>
                <p><strong>评论时间：</strong>%s</p>
            </div>
            <div class="comment-content">
                <p><strong>%s 的评论：</strong></p>
                <p>%s</p>
            </div>
            <a href="%s" class="cta-button">💬 查看完整评论</a>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, userName, articleTitle, commentTime, userName, commentContent, articleURL)
}

// GenerateSectionCommentHTML 生成章节评论邮件HTML模板
func GenerateSectionCommentHTML(data SectionCommentData) string {
	userName := html.EscapeString(data.UserName)
	commentContent := data.CommentContent // 直接使用，不转义（假设已经是安全的HTML）
	sectionTitle := html.EscapeString(data.SectionTitle)
	courseTitle := html.EscapeString(data.CourseTitle)
	sectionURL := html.EscapeString(data.SectionURL)
	commentTime := html.EscapeString(data.CommentTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 章节评论通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .comment-info { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .course-info { background: #e3f2fd; padding: 10px; border-left: 3px solid #2196f3; margin: 10px 0; }
        .comment-content { background: #e9ecef; padding: 15px; border-left: 3px solid #667eea; margin: 15px 0; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="comment-info">
                <h3>💬 %s 评论了你的章节</h3>
                <p><strong>章节：</strong>%s</p>
                <p><strong>评论时间：</strong>%s</p>
            </div>
            <div class="course-info">
                <p><strong>所属课程：</strong>%s</p>
            </div>
            <div class="comment-content">
                <p><strong>%s 的评论：</strong></p>
                <p>%s</p>
            </div>
            <a href="%s" class="cta-button">💬 查看完整评论</a>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, userName, sectionTitle, commentTime, courseTitle, userName, commentContent, sectionURL)
}

// GenerateCourseCommentHTML 生成课程评论邮件HTML模板
func GenerateCourseCommentHTML(data CourseCommentData) string {
	userName := html.EscapeString(data.UserName)
	commentContent := data.CommentContent // 直接使用，不转义（假设已经是安全的HTML）
	courseTitle := html.EscapeString(data.CourseTitle)
	courseURL := html.EscapeString(data.CourseURL)
	commentTime := html.EscapeString(data.CommentTime)

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>敲鸭社区 - 课程评论通知</title>
    <style>
        .email-container { max-width: 600px; margin: 0 auto; background: #fff; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; }
        .comment-info { background: #f8f9fa; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
        .comment-content { background: #e9ecef; padding: 15px; border-left: 3px solid #667eea; margin: 15px 0; }
        .cta-button { display: inline-block; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-decoration: none; border-radius: 25px; }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        <div class="content">
            <div class="comment-info">
                <h3>💬 %s 评论了你的课程</h3>
                <p><strong>课程：</strong>%s</p>
                <p><strong>评论时间：</strong>%s</p>
            </div>
            <div class="comment-content">
                <p><strong>%s 的评论：</strong></p>
                <p>%s</p>
            </div>
            <a href="%s" class="cta-button">💬 查看完整评论</a>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate, userName, courseTitle, commentTime, userName, commentContent, courseURL)
}

// GenerateSectionPublishHTML 生成章节发布邮件HTML模板
func GenerateSectionPublishHTML(data SectionPublishData) string {
	// 转义HTML内容防止XSS（但不转义 sectionContent，因为它已经是安全的HTML）
	userName := html.EscapeString(data.UserName)
	courseTitle := html.EscapeString(data.CourseTitle)
	sectionTitle := html.EscapeString(data.SectionTitle)
	sectionContent := data.SectionContent // 直接使用，不转义（假设已经是安全的HTML）
	sectionURL := html.EscapeString(data.SectionURL)
	publishTime := html.EscapeString(data.PublishTime)

	// 如果HTML内容过长，截取（基于去除标签后的文本长度）
	plainTextForLimit := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(sectionContent, "")
	if len([]rune(plainTextForLimit)) > 200 {
		runes := []rune(plainTextForLimit)
		limitText := string(runes[:200]) + "..."
		sectionContent = "<p>" + html.EscapeString(limitText) + "</p>"
	}

	// 不再处理头像
	userAvatar := data.UserAvatar
	if userAvatar == "" {
		userAvatar = "https://via.placeholder.com/50x50/667eea/ffffff?text=" + string([]rune(userName)[:1])
	}

	htmlTemplate := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>敲鸭社区 - 课程更新通知</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f5f5;
        }
        
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 20px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 24px;
            font-weight: 600;
            margin: 0;
        }
        
        .content {
            padding: 30px;
        }
        
        .author-info {
            display: flex;
            align-items: center;
            margin-bottom: 25px;
            padding: 15px;
            background-color: #f8f9fa;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        
        .author-avatar {
            width: 50px;
            height: 50px;
            border-radius: 50%%;
            margin-right: 15px;
            border: 3px solid #fff;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        }
        
        .author-details h3 {
            color: #333;
            font-size: 16px;
            margin-bottom: 4px;
        }
        
        .author-details p {
            color: #666;
            font-size: 14px;
            margin: 0;
        }
        
        .section-content {
            background-color: #fff;
            border: 1px solid #e9ecef;
            border-radius: 8px;
            padding: 25px;
            margin-bottom: 25px;
        }
        
        .section-title {
            font-size: 22px;
            font-weight: 600;
            color: #2c3e50;
            margin-bottom: 15px;
            line-height: 1.4;
        }
        
        .course-info {
            background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%);
            color: white;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 15px;
        }
        
        .course-info h4 {
            margin: 0 0 5px 0;
            font-size: 14px;
            opacity: 0.9;
        }
        
        .course-info h3 {
            margin: 0;
            font-size: 18px;
        }
        
        .section-preview {
            color: #555;
            font-size: 15px;
            line-height: 1.7;
            margin-bottom: 15px;
        }
        
        .section-meta {
            color: #888;
            font-size: 13px;
            margin-bottom: 20px;
        }
        
        .cta-button {
            display: inline-block;
            padding: 12px 30px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white !important;
            text-decoration: none;
            border-radius: 25px;
            font-weight: 600;
            font-size: 14px;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
        }
        
        .footer {
            background-color: #f8f9fa;
            padding: 20px;
            text-align: center;
            border-top: 1px solid #e9ecef;
        }
        
        .footer p {
            color: #6c757d;
            font-size: 12px;
            margin: 5px 0;
        }
        
        .footer a {
            color: #667eea;
            text-decoration: none;
        }
        
        @media (max-width: 600px) {
            .email-container {
                margin: 0;
                box-shadow: none;
            }
            
            .content {
                padding: 20px;
            }
            
            .header {
                padding: 15px;
            }
            
            .header h1 {
                font-size: 20px;
            }
            
            .author-info {
                padding: 12px;
            }
            
            .section-content {
                padding: 18px;
            }
            
            .section-title {
                font-size: 18px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="header">
            <h1>敲鸭社区</h1>
        </div>
        
        <div class="content">
            <div class="author-info">
                <div class="author-details">
                    <h3>%s</h3>
                    <p>发布了新章节</p>
                </div>
            </div>
            
            <div class="section-content">
                <div class="course-info">
                    <h4>课程</h4>
                    <h3>%s</h3>
                </div>
                
                <h2 class="section-title">新章节：%s</h2>
                
                <div class="section-preview">
                    %s
                </div>
                
                <div class="section-meta">
                    📅 发布于: %s
                </div>
                
                <a href="%s" class="cta-button">
                    📖 继续学习
                </a>
            </div>
        </div>
        
        <div class="footer">
            <p>感谢您关注敲鸭社区！</p>
            <p>
                <a href="https://code.xhyovo.cn">访问社区</a> | 
                <a href="#">邮件偏好设置</a>
            </p>
            <p>© %d 敲鸭社区 - 专注于技术分享与交流</p>
        </div>
    </div>
</body>
</html>`

	return fmt.Sprintf(htmlTemplate,
		userName, courseTitle, sectionTitle, sectionContent, publishTime,
		sectionURL, time.Now().Year())
}
