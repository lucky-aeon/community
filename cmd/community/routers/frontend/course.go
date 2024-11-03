package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

var courseService services.CourseService

func InitCourseRouters(r *gin.Engine) {
	group := r.Group("/community/courses")
	group.GET("", ListCourse)
	group.GET("/section", ListCourseSection)
	group.GET("/section/newest", ListNewest)
	group.Use(middleware.OperLogger())
	group.GET("/:id", GetCourseDetail)
	group.GET("/section/:id", GetCourseSectionDetail)

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

// 获取课程详细信息
func GetCourseSectionDetail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
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
	courseId, err := strconv.Atoi(ctx.Query("courseId"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d,获取课程列表信息失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	courses, count := courseService.PageCourseSection(p, limit, courseId)
	result.Page(courses, count, nil).Json(ctx)
	return
}

// 获取最新课程章节，只获取8个
func ListNewest(ctx *gin.Context) {

	courses := courseService.GetNewestCourseSection()

	result.Ok(courses, "").Json(ctx)
}
