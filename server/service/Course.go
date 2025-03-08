package services

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type CourseService struct {
}

// CourseResource 课程配套学习资源
type CourseResource struct {
	Title       string      `json:"title"`       // 资源标题
	Description string      `json:"description"` // 资源描述
	Icon        interface{} `json:"icon"`        // 资源图标，可以是字符串或数字
}

// 发布课程
func (*CourseService) Publish(course model.Courses) {
	course.Technology = strings.Join(course.TechnologyS, ",")

	// 处理配套学习资源
	if course.Resources != nil {
		// 首先尝试将其作为单个CourseResource对象处理
		if singleResource, ok := course.Resources.(map[string]interface{}); ok {
			// 将单个对象封装为数组
			resourceArr := []map[string]interface{}{singleResource}
			resourcesJSON, err := json.Marshal(resourceArr)
			if err == nil {
				course.ResourcesJSON = string(resourcesJSON)
			} else {
				course.ResourcesJSON = "[]"
			}
		} else if resources, ok := course.Resources.([]CourseResource); ok && len(resources) > 0 {
			// 处理CourseResource数组
			resourcesJSON, err := json.Marshal(resources)
			if err == nil {
				course.ResourcesJSON = string(resourcesJSON)
			} else {
				course.ResourcesJSON = "[]"
			}
		} else if resources, ok := course.Resources.([]interface{}); ok && len(resources) > 0 {
			// 处理interface{}数组
			resourcesJSON, err := json.Marshal(resources)
			if err == nil {
				course.ResourcesJSON = string(resourcesJSON)
			} else {
				course.ResourcesJSON = "[]"
			}
		} else {
			// 如果Resources为空或转换失败，设置为空数组
			course.ResourcesJSON = "[]"
		}
	} else {
		// 确保ResourcesJSON至少是一个空数组
		course.ResourcesJSON = "[]"
	}

	if course.ID == 0 {
		model.Course().Create(&course)
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
	var c model.Courses
	model.Course().First(&c, id)
	c.TechnologyS = strings.Split(c.Technology, ",")

	// 处理课程学习资源
	if c.ResourcesJSON != "" && c.ResourcesJSON != "null" {
		var resources []CourseResource
		err := json.Unmarshal([]byte(c.ResourcesJSON), &resources)
		if err == nil {
			c.Resources = resources
		} else {
			// 如果解析失败，设置为空数组
			c.Resources = []CourseResource{}
		}
	} else {
		// 如果ResourcesJSON为空或为"null"，设置为空数组
		c.Resources = []CourseResource{}
	}

	model.CoursesSection().Where("course_id = ?", id).Find(&c.Sections)
	// 暂时不处理Views字段计数
	return &c
}

// 获取课程列表
func (*CourseService) PageCourse(page, limit int) (courses []model.Courses, count int64) {
	model.Course().Offset((page - 1) * limit).Limit(limit).Order("created_at desc").Find(&courses)
	model.Course().Count(&count)

	if count == 0 {
		return
	}

	// 获取所有章节的访问量
	//var logService LogServices
	//viewsMap := logService.GetCoursesLearnCount()
	//
	//// 将访问量添加到章节信息中
	//for i := range courses {
	//	if views, exists := viewsMap[courses[i].ID]; exists {
	//		courses[i].Views = views
	//	}
	//}
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

		model.CoursesSection().Create(&section)
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
	// 获取章节列表
	model.CoursesSection().
		Where("course_id = ? ", courseId).
		Order("sort").
		Select("id", "title").
		Find(&courses)
	count = int64(len(courses))
	return
}

// 删除课程章节
func (*CourseService) DeleteCourseSection(id int) {
	model.CoursesSection().Delete("id = ?", id)
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

// GetHomePageCourses 获取首页课程列表
func (c *CourseService) GetHomePageCourses(limit int) ([]model.Courses, error) {
	var courses []model.Courses

	// 获取推荐课程列表
	model.Course().
		Order("created_at desc").
		Limit(limit).
		Find(&courses)

	// 处理技术栈标签和学习资源
	for i := range courses {
		courses[i].TechnologyS = strings.Split(courses[i].Technology, ",")

		// 处理课程学习资源
		if courses[i].ResourcesJSON != "" && courses[i].ResourcesJSON != "null" {
			var resources []CourseResource
			err := json.Unmarshal([]byte(courses[i].ResourcesJSON), &resources)
			if err == nil {
				courses[i].Resources = resources
			} else {
				// 如果解析失败，设置为空数组
				courses[i].Resources = []CourseResource{}
			}
		} else {
			// 如果ResourcesJSON为空或为"null"，设置为空数组
			courses[i].Resources = []CourseResource{}
		}

		// 清空URL字段（机密信息）
		courses[i].Url = ""
	}

	return courses, nil
}

// GetDefaultCourseResources 获取默认的课程学习资源
func (c *CourseService) GetDefaultCourseResources(courseId int) []CourseResource {
	// 示例数据
	resources := []CourseResource{
		{
			Title:       "课程讲义",
			Description: "PDF格式课程讲义，随课程更新",
			Icon:        "document-icon",
		},
		{
			Title:       "示例代码",
			Description: "GitHub代码仓库，包含所有示例",
			Icon:        "github-icon",
		},
		{
			Title:       "学习社区",
			Description: "加入微信交流群，与同学共同学习",
			Icon:        "community-icon",
		},
	}

	return resources
}

// GetAllCoursesWithDetails 获取所有课程的详细信息（包括章节列表和学习资源）
func (c *CourseService) GetAllCoursesWithDetails() ([]model.Courses, error) {
	var courses []model.Courses

	// 1. 获取所有课程
	model.Course().
		Order("created_at desc").
		Find(&courses)

	if len(courses) == 0 {
		return courses, nil
	}

	// 2. 提取所有课程ID
	courseIds := make([]int, len(courses))
	courseMap := make(map[int]*model.Courses)
	for i, course := range courses {
		courseIds[i] = course.ID
		courseMap[course.ID] = &courses[i]
		// 处理技术栈标签
		courses[i].TechnologyS = strings.Split(courses[i].Technology, ",")

		// 处理课程学习资源
		if courses[i].ResourcesJSON != "" && courses[i].ResourcesJSON != "null" {
			var resources []CourseResource
			err := json.Unmarshal([]byte(courses[i].ResourcesJSON), &resources)
			if err == nil {
				courses[i].Resources = resources
			} else {
				// 如果解析失败，设置为空数组
				courses[i].Resources = []CourseResource{}
			}
		} else {
			// 如果ResourcesJSON为空或为"null"，设置为空数组
			courses[i].Resources = []CourseResource{}
		}

		// 清空URL字段（机密信息）
		courses[i].Url = ""
	}

	// 3. 一次性查询所有课程的章节，只查询ID和标题字段
	var allSections []model.CoursesSections
	model.CoursesSection().
		Where("course_id IN ?", courseIds).
		Select("id", "title", "course_id", "sort").
		Order("sort asc").
		Find(&allSections)

	// 4. 构建映射关系，将章节分配给对应的课程，只保留章节的必要信息
	for _, section := range allSections {
		// 清除章节中的详细内容
		simplifiedSection := model.CoursesSections{
			ID:       section.ID,
			Title:    section.Title,
			CourseId: section.CourseId,
			Sort:     section.Sort,
		}

		if course, ok := courseMap[section.CourseId]; ok {
			course.Sections = append(course.Sections, simplifiedSection)
		}
	}

	return courses, nil
}
