package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var courseService services.CourseService

func InitCourseRouters(r *gin.Engine) {
	group := r.Group("/community/admin/courses")
	group.GET("/tree", ListCourseTree)
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
	userId := middleware.GetUserId(ctx)

	if err := ctx.ShouldBindJSON(&sections); err != nil {
		msg := utils.GetValidateErr(sections, err)
		log.Warnf("用户: % d ,发布章节时参数解析失败,err: %s", userId, msg)
		result.Err(msg).Json(ctx)
		return
	}
	sections.UserId = userId
	if err := courseService.PublishSection(sections); err != nil {
		log.Warnf("用户: %d 发布章节时对应文章不存在,课程 id : %s", userId, sections.CourseId)
		result.Err("对应课程不存在").Json(ctx)
		return
	}

	result.OkWithMsg(nil, "发布成功").Json(ctx)
}

// 删除章节
func DeleteCourseSection(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	userId := middleware.GetUserId(ctx)
	if err != nil {
		log.Warnf("用户: %d,删除课程失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	courseService.DeleteCourseSection(id)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 获取所有的课程以及章节为树形
func ListCourseTree(ctx *gin.Context) {

	result.Ok(courseService.ListCourseTree(), "").Json(ctx)
}

func ListCourseTitle() {

}
