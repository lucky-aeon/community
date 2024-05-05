package services

import (
	"errors"
	"sort"
	"strings"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type CourseService struct {
}

// 发布课程
func (*CourseService) Publish(course model.Courses) {
	course.Technology = strings.Join(course.TechnologyS, ",")
	if course.ID == 0 {
		model.Course().Save(&course)
	} else {
		model.Course().Where("id = ?", course.ID).Updates(&course)
	}
}

// 获取课程详细信息
func (*CourseService) GetCourseDetail(id int) *model.Courses {
	var course *model.Courses
	model.Course().Where("id = ?", id).Find(&course)
	course.TechnologyS = strings.Split(course.Technology, ",")
	return course
}

// 获取课程列表
func (*CourseService) PageCourse(page, limit int) (courses []model.Courses, count int64) {
	model.Course().Offset((page - 1) * limit).Limit(limit).Order("created_at desc").Find(&courses)
	model.Course().Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourse(id int) {
	model.Course().Delete("id = ?", id)
	model.CoursesSection().Where("course_id = ?", id).Delete(&model.CoursesSections{})
}

// 发布章节
func (c *CourseService) PublishSection(section model.CoursesSections) error {
	if c.GetCourseDetail(section.CourseId).ID == 0 {
		return errors.New("对应课程不存在")
	}
	if section.ID == 0 {

		model.CoursesSection().Save(&section)
		var b SubscribeData
		var subscriptionService SubscriptionService
		b.UserId = section.UserId
		b.CurrentBusinessId = section.CourseId
		b.SubscribeId = section.CourseId
		b.CourseId = section.CourseId
		b.SectionId = section.ID
		subscriptionService.Do(event.CourseUpdate, b)
	} else {
		model.CoursesSection().Where("id = ?", section.ID).Updates(&section)
	}
	return nil
}

// 获取章节详细信息
func (*CourseService) GetCourseSectionDetail(id int) *model.CoursesSections {
	var sections *model.CoursesSections
	model.CoursesSection().Where("id = ?", id).Find(&sections)
	var userS UserService
	sections.UserSimple = userS.GetUserSimpleById(sections.UserId)
	return sections
}

// 获取课程列表
func (*CourseService) PageCourseSection(page, limit, courseId int) (courses []model.CoursesSections, count int64) {
	model.CoursesSection().Limit(limit).Offset((page-1)*limit).Where("course_id = ? ", courseId).Select("id", "title").Find(&courses)
	model.CoursesSection().Where("course_id = ? ", courseId).Count(&count)
	return
}

// 删除课程
func (*CourseService) DeleteCourseSection(id int) {
	model.CoursesSection().Delete("id = ?", id)
	// 对应评论一并删除 todo
}

func (c *CourseService) ListByIdsSelectIdTitleMap(ids []int) (m map[int]string) {
	rows, err := model.Course().Where("id in ?", ids).Select("id", "title").Rows()
	defer rows.Close()
	if err != nil {
		// 处理错误
		return nil
	}
	m = make(map[int]string)
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			// 处理错误
		}
		m[id] = title
	}
	return
}

func (c *CourseService) ListCourseTree() (courses []model.Courses) {
	// 获取所有章节
	model.Course().Find(&courses)
	// 获取所有课程
	var sections []model.CoursesSections
	model.CoursesSection().Find(&sections)
	// 构建课程树
	c.buildCourseTree(courses, sections)
	return courses
}

func (c *CourseService) buildCourseTree(courses []model.Courses, sections []model.CoursesSections) {
	// 章节转为 map
	var sectionMap = make(map[int][]model.CoursesSections)
	for i := range sections {
		sectionMap[sections[i].CourseId] = append(sectionMap[sections[i].CourseId], sections[i])
	}
	// 对课程中的章节排序
	for i := range sectionMap {

		sort.Slice(sectionMap[i], func(j, k int) bool {
			return sectionMap[i][j].Sort < sectionMap[i][k].Sort
		})
	}
	// 遍历课程，对应的章节数组放入课程
	for i := range courses {
		courses[i].Sections = sectionMap[courses[i].ID]
	}
}

func (c *CourseService) ListCourseTitle() []model.Courses {
	var courses []model.Courses
	model.Course().Select("id", "title").Find(&courses)
	return courses
}
