package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var courseService services.CourseService

func InitCourseRouters(r *gin.Engine) {
	group := r.Group("/community/admin/courses")
	group.GET("/:id", GetCourseDetail)
	group.GET("", ListCourse)

	group.GET("/section", ListCourseSection)
	group.GET("/section/:id", GetCourseSectionDetail)

	group.POST("", PublishCourse)
	group.DELETE("/:id", DeleteCourse)
	group.POST("/section", PublishSection)
	group.DELETE("/section/:id", DeleteCourseSection)
}

func PublishCourse(ctx *gin.Context) {
	var course model.Courses
	if err := ctx.ShouldBindJSON(&course); err != nil {
		msg := utils.GetValidateErr(course, err)
		log.Warnf("发布课程时参数解析失败,err: %s", msg)
		result.Err(msg).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	course.UserId = userId
	courseService.Publish(course)
	result.OkWithMsg(nil, "发布成功").Json(ctx)
}

// 获取课程详细信息
func GetCourseDetail(ctx *gin.Context) {
	courseId, err := strconv.Atoi(ctx.Param("id"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d ,获取课程详细信息失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Ok(courseService.GetCourseDetail(courseId), "").Json(ctx)
}

// 获取课程列表
func ListCourse(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	courses, count := courseService.PageCourse(p, limit)
	result.Page(courses, count, nil).Json(ctx)
}

// 删除课程
func DeleteCourse(ctx *gin.Context) {
	courseId, err := strconv.Atoi(ctx.Param("id"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d ,删除课程失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	courseService.DeleteCourse(courseId)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 发布章节
func PublishSection(ctx *gin.Context) {
	sections := model.CoursesSections{}
	if err := ctx.ShouldBindJSON(sections); err != nil {
		msg := utils.GetValidateErr(sections, err)
		log.Warnf("发布章节时参数解析失败,err: %s", msg)
		result.Err(msg).Json(ctx)
		return
	}
	sections.UserId = middleware.GetUserId(ctx)
	courseService.PublishSection(sections)
	result.OkWithMsg(nil, "发布成功").Json(ctx)
}

// 获取课程详细信息
func GetCourseSectionDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d,获取课程详细信息失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	detail := courseService.GetCourseSectionDetail(id)
	result.Ok(detail, "").Json(ctx)
}

// 获取课程列表
func ListCourseSection(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	courses, count := courseService.PageCourseSection(p, limit)
	result.Page(courses, count, nil).Json(ctx)
	return
}

// 删除课程
func DeleteCourseSection(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d,删除课程失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	courseService.DeleteCourseSection(id)
}
