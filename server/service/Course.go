package services

import "xhyovo.cn/community/server/model"

type CourseService struct {
}

// 发布课程
func (*CourseService) Publish(course model.Courses) {

	if course.ID == 0 {
		model.Course().Create(&course)
	} else {
		model.Course().Save(&course)
	}

}

// 获取课程详细信息
func (*CourseService) GetCourseDetail(id int) model.Courses {
	var course model.Courses
	model.Course().Where("id = ?", id).First(&course)
	return course
}

// 获取课程列表
func (*CourseService) PageCourse(page, limit int) (courses []model.Courses, count int64) {
	model.Course().Offset(limit).Limit((page - 1) * limit).Order("created_at desc").Find(&courses)
	model.Course().Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourse(id int) {
	model.Course().Delete("id = ?", id)
	model.CoursesSection().Delete("course_id = ?", id)
}

// 发布章节
func (*CourseService) PublishSection(course model.CoursesSections) {

	if course.ID == 0 {
		model.CoursesSection().Create(&course)
	} else {
		model.CoursesSection().Save(&course)
	}
}

// 获取课程详细信息
func (*CourseService) GetCourseSectionDetail(id int) model.CoursesSections {
	var sections model.CoursesSections
	model.CoursesSection().Where("id = ?", id).First(&sections)
	return sections
}

// 获取课程列表
func (*CourseService) PageCourseSection(page, limit int) (courses []model.CoursesSections, count int64) {
	model.CoursesSection().Limit((page-1)*limit).Offset(limit).Select("id", "title").Order("created_at desc").Find(&courses)
	model.CoursesSection().Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourseSection(id int) {
	model.CoursesSection().Delete("id = ?", id)
	// 对应评论一并删除 todo
}
