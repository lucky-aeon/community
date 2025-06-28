package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strings"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"xhyovo.cn/community/pkg/log"
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
		model.Course().Where("id = ?", course.ID).Select("*").Updates(&course)
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

	// 计算章节内容的阅读时间
	section.ReadingTime = c.calculateReadingTime(section.Content)

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

// calculateReadingTime 计算章节的阅读时间（分钟）
func (c *CourseService) calculateReadingTime(content string) int {
	if content == "" {
		return 0
	}

	// 统计计数
	var textWordCount int     // 文本字数
	var imageCount int        // 图片数量
	var totalVideoSeconds int // 视频总时长（秒）

	// 1. 计算视频时间
	// 匹配视频标记 !video[视频](/api/community/file/singUrl?fileKey=13/1741427148178)
	videoRegex := regexp.MustCompile(`!video\[.*?\]\(/api/community/file/singUrl\?fileKey=([^)]+)\)`)
	videoMatches := videoRegex.FindAllStringSubmatch(content, -1)

	for _, match := range videoMatches {
		if len(match) >= 2 {
			fileKey := match[1]
			// 对fileKey进行URL解码
			decodedFileKey, err := url.QueryUnescape(fileKey)
			if err != nil {
				log.Error("解码视频fileKey失败: " + err.Error() + ", 原始fileKey: " + fileKey)
				// 解码失败时继续使用原始fileKey
				videoSeconds := c.getVideoSeconds(fileKey)
				totalVideoSeconds += videoSeconds
			} else {
				// 解码成功，使用解码后的fileKey
				videoSeconds := c.getVideoSeconds(decodedFileKey)
				totalVideoSeconds += videoSeconds
			}
		}
	}

	// 2. 计算图片数量
	// 匹配图片标记 ![image#S #R #100% #100%](/api/community/file/singUrl?fileKey=13%2F1741427453119)
	imageRegex := regexp.MustCompile(`!\[image.*?\]\(/api/community/file/singUrl\?fileKey=([^)]+)\)`)
	imageMatches := imageRegex.FindAllString(content, -1)
	imageCount = len(imageMatches)

	// 3. 去除标记后计算文本字数
	// 去除视频和图片标记
	cleanedContent := videoRegex.ReplaceAllString(content, "")
	cleanedContent = imageRegex.ReplaceAllString(cleanedContent, "")

	// 去除Markdown标记
	markdownRegex := regexp.MustCompile(`[#*_~\[\](){}|>]+`)
	cleanedContent = markdownRegex.ReplaceAllString(cleanedContent, "")

	// 统计字数（每个中文字符或单词作为一个字）
	textWordCount = len([]rune(cleanedContent))

	// 计算阅读时间
	// 1. 文本阅读速度：平均每分钟阅读300字
	textMinutes := textWordCount / 300
	if textWordCount%300 > 0 {
		textMinutes++
	}

	// 2. 图片观看时间：每张图片平均10秒
	imageSeconds := imageCount * 10

	// 3. 视频观看时间：已经计算了总秒数

	// 总阅读时间(分钟) = 文本阅读时间 + 图片时间 + 视频时间
	totalMinutes := textMinutes + (imageSeconds+totalVideoSeconds)/60

	// 确保至少为1分钟
	if totalMinutes < 1 {
		totalMinutes = 1
	}

	return totalMinutes
}

// getVideoSeconds 获取视频时长（秒）
func (c *CourseService) getVideoSeconds(fileKey string) int {
	// 对fileKey进行URL解码，处理%2F等转义字符
	decodedFileKey, err := url.QueryUnescape(fileKey)
	if err != nil {
		log.Error("解码fileKey失败: " + err.Error() + ", 原始fileKey: " + fileKey)
		// 解码失败时继续使用原始fileKey
	} else {
		// 解码成功，使用解码后的fileKey
		fileKey = decodedFileKey
	}

	// 调用阿里云OSS获取视频信息
	provider, err := alioss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Error("获取阿里云凭证失败: " + err.Error())
		return 60 // 默认1分钟
	}

	client, err := alioss.New("https://oss-cn-beijing.aliyuncs.com", "", "",
		alioss.SetCredentialsProvider(&provider),
		alioss.AuthVersion(alioss.AuthV4),
		alioss.Region("cn-beijing"))
	if err != nil {
		log.Error("创建阿里云客户端失败: " + err.Error())
		return 60
	}

	bucketName := "luckly-community"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Error("获取Bucket失败: " + err.Error())
		return 60
	}

	// 获取视频信息
	body, err := bucket.GetObject(fileKey, alioss.Process("video/info"))
	if err != nil {
		log.Error("获取视频信息失败: " + err.Error() + ", fileKey: " + fileKey)
		return 60
	}
	defer body.Close()

	// 读取响应
	data, err := io.ReadAll(body)
	if err != nil {
		log.Error("读取视频信息失败: " + err.Error())
		return 60
	}

	// 解析JSON
	// 根据提供的实际JSON结构定义
	type OSSVideoInfo struct {
		Duration     float64 `json:"Duration"`   // 视频时长（秒）
		Bitrate      int     `json:"Bitrate"`    // 比特率
		FormatName   string  `json:"FormatName"` // 格式名称
		Size         int64   `json:"Size"`       // 大小
		VideoStreams []struct {
			Duration float64 `json:"Duration"` // 视频流时长
		} `json:"VideoStreams"`
		AudioStreams []struct {
			Duration float64 `json:"Duration"` // 音频流时长
		} `json:"AudioStreams"`
	}

	var videoInfo OSSVideoInfo
	if err := json.Unmarshal(data, &videoInfo); err != nil {
		log.Error("解析视频信息失败: " + err.Error())
		return 60
	}

	// 使用Duration字段
	if videoInfo.Duration > 0 {
		return int(videoInfo.Duration)
	}

	// 如果顶层Duration为0，尝试使用视频流的Duration
	if len(videoInfo.VideoStreams) > 0 && videoInfo.VideoStreams[0].Duration > 0 {
		return int(videoInfo.VideoStreams[0].Duration)
	}

	// 如果视频流Duration也为0，尝试使用音频流的Duration
	if len(videoInfo.AudioStreams) > 0 && videoInfo.AudioStreams[0].Duration > 0 {
		return int(videoInfo.AudioStreams[0].Duration)
	}

	log.Error("视频信息中未找到有效的Duration字段")
	return 60
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
	// 获取章节列表，包含阅读时间
	model.CoursesSection().
		Where("course_id = ? ", courseId).
		Order("sort").
		Select("id", "title", "reading_time").
		Find(&courses)

	// 批量查询所有章节的评论数量，避免N+1查询问题
	if len(courses) > 0 {
		sectionIds := make([]int, len(courses))
		for i, section := range courses {
			sectionIds[i] = section.ID
		}

		// 批量查询评论数量
		var commentCounts []struct {
			BusinessId int   `json:"business_id"`
			Count      int64 `json:"count"`
		}

		model.Comment().
			Select("business_id, COUNT(*) as count").
			Where("business_id IN ? AND tenant_id = ? AND deleted_at IS NULL", sectionIds, 1).
			Group("business_id").
			Scan(&commentCounts)

		// 建立评论数量映射
		commentCountMap := make(map[int]int64)
		for _, cc := range commentCounts {
			commentCountMap[cc.BusinessId] = cc.Count
		}

		// 设置每个章节的评论数量
		for i := range courses {
			courses[i].CommentCount = commentCountMap[courses[i].ID]
		}
	}

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
		Select("id", "title", "course_id", "sort", "reading_time").
		Order("sort asc").
		Find(&allSections)

	// 4. 构建映射关系，将章节分配给对应的课程，只保留章节的必要信息
	for _, section := range allSections {
		// 清除章节中的详细内容
		simplifiedSection := model.CoursesSections{
			ID:          section.ID,
			Title:       section.Title,
			CourseId:    section.CourseId,
			Sort:        section.Sort,
			ReadingTime: section.ReadingTime,
		}

		if course, ok := courseMap[section.CourseId]; ok {
			course.Sections = append(course.Sections, simplifiedSection)
		}
	}

	return courses, nil
}

// UpdateAllSectionsReadingTime 更新所有课程章节的阅读时间
func (c *CourseService) UpdateAllSectionsReadingTime() (int, error) {
	// 获取所有课程章节
	var sections []model.CoursesSections
	result := model.CoursesSection().Find(&sections)
	if result.Error != nil {
		log.Error("获取课程章节失败: " + result.Error.Error())
		return 0, result.Error
	}

	// 记录更新数量
	updatedCount := 0

	// 遍历所有章节并计算阅读时间
	for i := range sections {
		// 跳过已有阅读时间的章节（可选，取消注释以跳过）
		// if sections[i].ReadingTime > 0 {
		//     continue
		// }

		// 计算阅读时间
		readingTime := c.calculateReadingTime(sections[i].Content)

		// 如果阅读时间有变化，更新数据库
		if sections[i].ReadingTime != readingTime {
			sections[i].ReadingTime = readingTime
			result := model.CoursesSection().Where("id = ?", sections[i].ID).
				Update("reading_time", readingTime)

			if result.Error != nil {
				log.Error(fmt.Sprintf("更新章节ID=%d的阅读时间失败: %s",
					sections[i].ID, result.Error.Error()))
				continue
			}

			if result.RowsAffected > 0 {
				updatedCount++
			}
		}
	}

	log.Info(fmt.Sprintf("成功更新了%d个章节的阅读时间", updatedCount))
	return updatedCount, nil
}

// UpdateCourseSectionsReadingTime 更新指定课程的所有章节阅读时间
func (c *CourseService) UpdateCourseSectionsReadingTime(courseId int) (int, error) {
	// 获取指定课程的所有章节
	var sections []model.CoursesSections
	result := model.CoursesSection().Where("course_id = ?", courseId).Find(&sections)
	if result.Error != nil {
		log.Error(fmt.Sprintf("获取课程ID=%d的章节失败: %s",
			courseId, result.Error.Error()))
		return 0, result.Error
	}

	// 记录更新数量
	updatedCount := 0

	// 遍历所有章节并计算阅读时间
	for i := range sections {
		// 计算阅读时间
		readingTime := c.calculateReadingTime(sections[i].Content)

		// 如果阅读时间有变化，更新数据库
		if sections[i].ReadingTime != readingTime {
			sections[i].ReadingTime = readingTime
			result := model.CoursesSection().Where("id = ?", sections[i].ID).
				Update("reading_time", readingTime)

			if result.Error != nil {
				log.Error(fmt.Sprintf("更新章节ID=%d的阅读时间失败: %s",
					sections[i].ID, result.Error.Error()))
				continue
			}

			if result.RowsAffected > 0 {
				updatedCount++
			}
		}
	}

	log.Info(fmt.Sprintf("成功更新了课程ID=%d的%d个章节的阅读时间",
		courseId, updatedCount))
	return updatedCount, nil
}
