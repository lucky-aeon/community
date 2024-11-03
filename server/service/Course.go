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
		var subscriptionService SubscriptionService
		var b SubscribeData
		b.UserId = course.UserId
		b.CurrentBusinessId = course.ID
		b.CourseId = course.ID
		b.SubscribeId = course.UserId
		var messageTemp = "你关注的用户 ${user.name} 发布了最新课程: ${course.title}"
		subscriptionService.DoWithMessageTempl(event.UserFollowingEvent, b, messageTemp)
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
		b.CurrentBusinessId = section.ID
		b.SubscribeId = section.CourseId
		b.CourseId = section.CourseId
		b.SectionId = section.ID
		subscriptionService.Do(event.CourseUpdate, b)

		var b2 SubscribeData
		b2.UserId = section.UserId
		b2.CurrentBusinessId = section.ID
		b2.CourseId = section.CourseId
		b2.SubscribeId = section.UserId
		b2.SectionId = section.ID
		var messageTemp = "你关注的用户 ${user.name} 在课程: ${course.title} 发布了最新章节: ${courses_section.title}"
		subscriptionService.DoWithMessageTempl(event.UserFollowingEvent, b2, messageTemp)
	} else {
		model.CoursesSection().Where("id = ?", section.ID).Updates(&section)
	}
	return nil
}

// 获取章节详细信息
func (*CourseService) GetCourseSectionDetail(id int) *model.CoursesSections {
	var sections *model.CoursesSections
	var courseId int
	var courses []model.CoursesSections
	model.CoursesSection().Where("id = ?", id).Select("course_id").Find(&courseId)
	model.CoursesSection().Where("course_id = ? ", courseId).Order("sort").Find(&courses)
	// 遍历course 找到id 和 id相同的
	for i := range courses {
		if courses[i].ID == id {
			sections = &courses[i]
			// 如果当前下标还有上一个则设置PreId
			if i > 0 {
				sections.PreId = courses[i-1].ID
			}
			// 如果还有下一个下标
			if i < len(courses)-1 {
				sections.NextId = courses[i+1].ID
			}
			break
		}
	}
	var userS UserService
	sections.UserSimple = userS.GetUserSimpleById(sections.UserId)
	return sections
}

// 获取课程列表
func (*CourseService) PageCourseSection(page, limit, courseId int) (courses []model.CoursesSections, count int64) {
	model.CoursesSection().Where("course_id = ? ", courseId).Order("sort").Select("id", "title").Find(&courses)
	count = int64(len(courses))
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

func (c *CourseService) ListSectionByIds(ids []int) (m map[int]string) {
	rows, err := model.CoursesSection().Where("id in ?", ids).Select("id", "title").Rows()
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

func (c *CourseService) GetNewestCourseSection() []model.CoursesSections {

	var sections []model.CoursesSections

	// 使用联表查询，确保可以获取课程的名称
	model.CoursesSection().Select("courses_sections.title, courses_sections.id,courses.title AS courseTitle").
		Joins("LEFT JOIN courses ON courses.id = courses_sections.course_id").
		Order("courses_sections.created_at desc").
		Limit(8).
		Find(&sections)

	return sections
}
