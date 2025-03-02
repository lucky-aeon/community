package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/oss"
	"xhyovo.cn/community/pkg/postgre"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service/llm"
)

func main() {
	log.Init()
	// 设置程序使用中国时区
	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	time.Local = chinaLoc
	if err != nil {
		log.Errorf("Error loading China location:", err)
		return
	}
	r := gin.Default()
	r.SetFuncMap(utils.GlobalFunc())
	config.Init()
	appConfig := config.GetInstance()
	db := appConfig.DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)
	pgDbConfig := appConfig.PGDbConfig
	postgre.Init(pgDbConfig.Username, pgDbConfig.Password, pgDbConfig.Address, pgDbConfig.Database)
	ossConfig := appConfig.OssConfig
	oss.Init(ossConfig.Endpoint, ossConfig.AccessKey, ossConfig.SecretKey, ossConfig.Bucket)
	emailConfig := appConfig.EmailConfig
	email.Init(emailConfig.Username, emailConfig.Password, emailConfig.Host)

	routers.InitFrontedRouter(r)
	cache.Init()
	log.Info("start web :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Errorln(err)
	}

}
func GetPwd(pwd string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return hash, err
}

func knowledgeScript() {
	// 定义四个切片分别存储每个类型的执行失败ID
	var failedArticleIDs []int
	var failedCommentIDs []int
	var failedCourseIDs []int
	var failedSectionIDs []int

	var kService services.KnowledgeBaseService
	// 脚本，获取所有文章、评论、课程、章节
	var articles []model.Articles
	model.Article().Where("id in ?", []int{225}).Find(&articles)

	for i := range articles {
		article := articles[i]
		err := kService.AddKnowledge(article.Content, "https://code.xhyovo.cn/article/view?articleId="+strconv.Itoa(article.ID), "", constant.InternalArticle, article.ID)
		if err != nil {
			failedArticleIDs = append(failedArticleIDs, article.ID)
			log.Warnf("文章添加知识库失败:id: %d,err: %v", article.ID, err)
		}
	}

	var comments []model.Comments
	model.Comment().Find(&comments)

	for i := range comments {
		comment := comments[i]
		var link = ""
		var remake = ""
		if comment.TenantId == 0 {
			link = "https://code.xhyovo.cn/article/view?articleId=" + strconv.Itoa(comment.BusinessUserId)
		} else if comment.TenantId == 1 {
			link = ""
			remake = "该知识来源于课程评论,不支持跳转到原文"
		} else if comment.TenantId == 2 {
			link = "https://code.xhyovo.cn/article/view?sectionId=" + strconv.Itoa(comment.BusinessUserId)
		} else if comment.TenantId == 3 {
			link = ""
			remake = "该知识来源于会议评论,不支持跳转到原文"
		}
		err := kService.AddKnowledge(comment.Content, link, remake, constant.InternalComment, comment.BusinessId)
		if err != nil {
			failedCommentIDs = append(failedCommentIDs, comment.ID)
			log.Warnf("评论添加知识库失败:id: %d,err: %v", comment.BusinessId, err)
		}
	}

	var courses []model.Courses
	model.Course().Find(&courses)

	for i := range courses {
		course := courses[i]
		err := kService.AddKnowledge(course.Desc, "", "该知识来源于课程,不支持跳转到原文", constant.InternalCourse, course.ID)
		if err != nil {
			log.Warnf("课程添加知识库失败:id: %d,err: %v", course.ID, err)
		}
	}

	var courseSections []model.CoursesSections
	model.CoursesSection().Find(&courseSections)

	for i := range courseSections {
		sections := courseSections[i]
		err := kService.AddKnowledge(sections.Content, "https://code.xhyovo.cn/article/view?sectionId="+strconv.Itoa(sections.ID), "", constant.InternalChapter, sections.ID)
		if err != nil {
			failedSectionIDs = append(failedSectionIDs, sections.ID)
			log.Warnf("章节添加知识库失败:id: %d,err: %v", sections.ID, err)
		}
	}
	// 打印所有执行失败的ID
	fmt.Println("执行完成")
	if len(failedArticleIDs) > 0 {
		fmt.Printf("文章执行失败的ID: %v\n", failedArticleIDs)
	}
	if len(failedCommentIDs) > 0 {
		fmt.Printf("评论执行失败的ID: %v\n", failedCommentIDs)
	}
	if len(failedCourseIDs) > 0 {
		fmt.Printf("课程执行失败的ID: %v\n", failedCourseIDs)
	}
	if len(failedSectionIDs) > 0 {
		fmt.Printf("章节执行失败的ID: %v\n", failedSectionIDs)
	}

}
