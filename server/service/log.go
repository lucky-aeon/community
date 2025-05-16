package services

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	ltime "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

type LogServices struct {
}

func (*LogServices) GetPageOperLog(page, limit int, logSearch model.LogSearch, flag bool) (logs []model.OperLogs, count int64) {
	db := model.OperLog()
	if logSearch.RequestMethod != "" {
		db.Where("request_method = ?", logSearch.RequestMethod)
	}
	if logSearch.RequestInfo != "" {
		db.Where("request_info like ?", "%"+logSearch.RequestInfo+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at >= ? and created_at <= ?", logSearch.StartTime, logSearch.EndTime)
	}
	if logSearch.UserName != "" {
		var userS UserService
		ids := userS.SearchNameSelectId(logSearch.UserName)
		db.Where("user_id in ?", ids)
	}
	if flag {
		db.Where("request_info != '/community/file/singUrl'")
	} else {
		db.Where("request_info = '/community/file/singUrl'")
	}
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)

	set := mapset.NewSet[int]()

	for _, lo := range logs {
		set.Add(lo.UserId)
	}
	var u UserService
	userMap := u.ListByIdsToMap(set.ToSlice())
	for i := range logs {
		logs[i].UserName = userMap[logs[i].UserId].Name
	}
	return
}

func (*LogServices) InsertOperLog(log model.OperLogs) {
	go func(log model.OperLogs) {
		model.OperLog().Create(&log)
	}(log)
}

func (*LogServices) InsertLoginLog(log model.LoginLogs) {
	go func(log model.LoginLogs) {
		model.LoginLog().Create(&log)
	}(log)
}

func (s *LogServices) GetPageLoginPage(page, limit int, logSearch model.LogSearch) (logs []model.LoginLogs, count int64) {
	db := model.LoginLog()
	if logSearch.Account != "" {
		db.Where("account like ?", "%"+logSearch.Account+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at >= ? and created_at <= ?", logSearch.StartTime, logSearch.EndTime)
	}

	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	return
}

func (s *LogServices) GetPageQuestionLogs(page int, limit int) ([]struct {
	ID        int             `json:"id"`
	Question  string          `json:"question"`
	UserName  string          `json:"userName"`
	CreatedAt ltime.LocalTime `json:"createdAt"`
}, int64) {
	var count int64
	var logs []struct {
		ID        int             `json:"id"`
		Question  string          `json:"question"`
		UserName  string          `json:"userName"`
		CreatedAt ltime.LocalTime `json:"createdAt"`
	}

	db := model.OperLog()

	// 查找 request_info 以 "/community/knowledge/query" 开头的记录，并提取 question 参数、关联的 user name 和提问时间
	db.Select("oper_logs.id, SUBSTRING_INDEX(SUBSTRING(oper_logs.request_info, LOCATE('?question=', oper_logs.request_info) + 10), '&', 1) AS question, users.name AS user_name, oper_logs.created_at").
		Joins("JOIN users ON oper_logs.user_id = users.id").
		Where("oper_logs.request_info LIKE ?", "/community/knowledge/query%").
		Count(&count)

	if count == 0 {
		return nil, 0
	}

	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)

	return logs, count
}

// GetArticleReadCount 获取文章被阅读数量（去重）
func (*LogServices) GetArticleReadCount(articleId int) int64 {
	var count int64

	// 构建查询路径
	requestPath := "/community/articles/" + strconv.Itoa(articleId)

	// 使用 DISTINCT 对 user_id 去重计数
	model.OperLog().
		Where("request_info = ?", requestPath).
		Distinct("user_id").
		Count(&count)

	return count
}

// GetArticleReadCountBatch 批量获取多篇文章的阅读数量
func (*LogServices) GetArticleReadCountBatch(articleIds []int) map[int]int64 {
	result := make(map[int]int64)
	if len(articleIds) == 0 {
		return result
	}

	// 构建查询条件
	var conditions []string
	for _, id := range articleIds {
		conditions = append(conditions, "/community/articles/"+strconv.Itoa(id))
	}

	// 查询结果结构
	type ReadCount struct {
		RequestInfo string
		Count       int64
	}
	var counts []ReadCount

	// 使用 GROUP BY 进行批量查询
	model.OperLog().
		Select("request_info, COUNT(DISTINCT user_id) as count").
		Where("request_info IN ?", conditions).
		Group("request_info").
		Find(&counts)

	// 将结果转换为 map
	for _, c := range counts {
		// 从 request_info 中提取文章 ID
		idStr := strings.TrimPrefix(c.RequestInfo, "/community/articles/")
		if id, err := strconv.Atoi(idStr); err == nil {
			result[id] = c.Count
		}
	}

	return result
}

// GetCoursesLearnCount 获取所有课程的学习人数数量
func (*LogServices) GetCoursesLearnCount() map[int]int64 {
	result := make(map[int]int64)

	// 查询结果结构
	type LearnCount struct {
		SectionId int
		Count     int64
	}
	var counts []LearnCount

	// 1. 首先获取所有访问章节的记录和去重后的访问人数
	model.OperLog().
		Select("CAST(SUBSTRING_INDEX(request_info, '/', -1) AS UNSIGNED) as section_id, COUNT(DISTINCT user_id) as count").
		Where("request_info LIKE ?", "/community/courses/section/%").
		Group("section_id").
		Find(&counts)

	if len(counts) == 0 {
		return result
	}

	// 2. 获取所有涉及的章节ID
	var sectionIds []int
	for _, c := range counts {
		sectionIds = append(sectionIds, c.SectionId)
	}

	// 3. 查询这些章节的信息和所属课程
	var sections []model.CoursesSections
	model.CoursesSection().
		Select("id, title, course_id, created_at, sort").
		Where("id IN ?", sectionIds).
		Find(&sections)

	// 4. 建立章节到课程的映射
	sectionToCourse := make(map[int]int)                // 章节ID -> 课程ID
	sectionToTitle := make(map[int]string)              // 章节ID -> 章节标题
	courseToSections := make(map[int][]int)             // 课程ID -> 章节ID列表
	sectionToCreatedAt := make(map[int]ltime.LocalTime) // 章节ID -> 创建时间
	sectionToSort := make(map[int]int)                  // 章节ID -> 排序值

	for _, section := range sections {
		sectionToCourse[section.ID] = section.CourseId
		sectionToTitle[section.ID] = section.Title
		courseToSections[section.CourseId] = append(courseToSections[section.CourseId], section.ID)
		sectionToCreatedAt[section.ID] = section.CreatedAt
		sectionToSort[section.ID] = section.Sort
	}

	// 5. 统计每个课程的总学习人数（对同一用户访问同一课程的不同章节去重）
	courseUserMap := make(map[int]mapset.Set[int])

	// 再次查询日志获取用户ID信息
	var logs []model.OperLogs
	model.OperLog().
		Select("user_id, request_info").
		Where("request_info LIKE ?", "/community/courses/section/%").
		Find(&logs)

	// 处理每条日志
	for _, log := range logs {
		// 从 request_info 提取章节ID
		sectionStr := strings.TrimPrefix(log.RequestInfo, "/community/courses/section/")
		sectionId, err := strconv.Atoi(sectionStr)
		if err != nil {
			continue
		}

		// 获取课程ID
		courseId, exists := sectionToCourse[sectionId]
		if !exists {
			continue
		}

		// 初始化课程的用户集合
		if courseUserMap[courseId] == nil {
			courseUserMap[courseId] = mapset.NewSet[int]()
		}

		// 添加用户ID到集合
		courseUserMap[courseId].Add(log.UserId)
	}

	// 6. 生成最终结果
	for courseId, userSet := range courseUserMap {
		result[courseId] = int64(userSet.Cardinality())
	}

	return result
}

// CourseStats 课程统计信息结构
type CourseStats struct {
	CourseId       int             `json:"courseId"`
	CourseTitle    string          `json:"courseTitle"`
	UniqueUsers    int64           `json:"uniqueUsers"`    // 去重用户数
	TotalVisits    int64           `json:"totalVisits"`    // 总访问量
	SectionStats   []SectionStats  `json:"sectionStats"`   // 各章节统计
	CreatedAt      ltime.LocalTime `json:"createdAt"`      // 课程创建时间
	TotalSections  int             `json:"totalSections"`  // 总章节数
	CompletionRate float64         `json:"completionRate"` // 完成度，计算方式: 访问章节数/总章节数
}

// SectionStats 章节统计信息结构
type SectionStats struct {
	SectionId    int             `json:"sectionId"`
	SectionTitle string          `json:"sectionTitle"`
	UniqueUsers  int64           `json:"uniqueUsers"` // 去重用户数
	TotalVisits  int64           `json:"totalVisits"` // 总访问量
	CreatedAt    ltime.LocalTime `json:"createdAt"`   // 章节创建时间
	Sort         int             `json:"sort"`        // 章节排序值
}

// CoursesStatisticsResult 课程统计结果
type CoursesStatisticsResult struct {
	Courses     []CourseStats         `json:"courses"`     // 课程统计
	TotalUsers  int64                 `json:"totalUsers"`  // 系统总用户数
	ActiveUsers int64                 `json:"activeUsers"` // 有学习行为的活跃用户数
	TimeRanges  []TimeRangeStatistics `json:"timeRanges"`  // 时间段统计，默认按月统计近6个月
}

// TimeRangeStatistics 时间段统计
type TimeRangeStatistics struct {
	Period      string `json:"period"`      // 时间段，如 "2023-01"
	UniqueUsers int64  `json:"uniqueUsers"` // 该时间段内的去重用户数
	TotalVisits int64  `json:"totalVisits"` // 该时间段内的总访问量
}

// GetCoursesStatistics 获取课程和章节的访问统计数据
func (*LogServices) GetCoursesStatistics() CoursesStatisticsResult {
	result := CoursesStatisticsResult{
		Courses:    []CourseStats{},
		TimeRanges: []TimeRangeStatistics{},
	}

	// 1. 获取所有课程章节访问日志
	var logs []model.OperLogs
	model.OperLog().
		Select("id, user_id, request_info, created_at").
		Where("request_info LIKE ?", "/community/courses/section/%").
		Find(&logs)

	// 2. 提取所有涉及的章节ID
	sectionIds := make(map[int]bool)
	sectionLogs := make(map[int][]model.OperLogs) // 每个章节的访问日志

	// 收集所有有学习行为的用户ID
	allLearningUsers := mapset.NewSet[int]()

	// 按月份统计访问数据
	monthlyStats := make(map[string]map[string]interface{})

	// 获取最近6个月的月份列表
	now := time.Now()
	monthPeriods := make([]string, 6)
	for i := 0; i < 6; i++ {
		month := now.AddDate(0, -i, 0)
		monthPeriods[i] = month.Format("2006-01")
	}

	// 初始化月度统计数据
	for _, period := range monthPeriods {
		monthlyStats[period] = map[string]interface{}{
			"uniqueUsers": mapset.NewSet[int](),
			"totalVisits": int64(0),
		}
	}

	for _, log := range logs {
		// 从 request_info 提取章节ID
		sectionStr := strings.TrimPrefix(log.RequestInfo, "/community/courses/section/")
		sectionId, err := strconv.Atoi(sectionStr)
		if err != nil {
			continue
		}

		sectionIds[sectionId] = true

		// 按章节分组日志
		sectionLogs[sectionId] = append(sectionLogs[sectionId], log)

		// 收集所有学习用户
		allLearningUsers.Add(log.UserId)

		// 按月统计
		month := time.Time(log.CreatedAt).Format("2006-01")
		// 只统计最近6个月
		if stats, exists := monthlyStats[month]; exists {
			users := stats["uniqueUsers"].(mapset.Set[int])
			users.Add(log.UserId)
			stats["totalVisits"] = stats["totalVisits"].(int64) + 1
		}
	}

	// 3. 查询这些章节的信息和所属课程
	var sections []model.CoursesSections
	if len(sectionIds) > 0 {
		sectionIdList := keys(sectionIds)
		model.CoursesSection().
			Select("id, title, course_id, created_at, sort").
			Where("id IN ?", sectionIdList).
			Find(&sections)
	}

	// 获取系统总用户数，即使没有章节数据
	var totalUserCount int64
	model.User().Count(&totalUserCount)

	// 如果没有章节访问数据，至少返回部分统计信息
	if len(logs) == 0 || len(sections) == 0 {
		// 获取所有课程
		var allCourses []model.Courses
		model.Course().
			Select("id, title, created_at").
			Order("created_at DESC").
			Find(&allCourses)

		// 获取每个课程的章节数量
		for _, course := range allCourses {
			var sectionCount int64
			model.CoursesSection().
				Where("course_id = ?", course.ID).
				Count(&sectionCount)

			result.Courses = append(result.Courses, CourseStats{
				CourseId:       course.ID,
				CourseTitle:    course.Title,
				UniqueUsers:    0,
				TotalVisits:    0,
				SectionStats:   []SectionStats{},
				CreatedAt:      course.CreatedAt,
				TotalSections:  int(sectionCount),
				CompletionRate: 0,
			})
		}

		// 返回至少包含课程基本信息的结果
		result.TotalUsers = totalUserCount
		result.ActiveUsers = int64(allLearningUsers.Cardinality())

		// 构建时间段统计数据
		timeRangesResult := make([]TimeRangeStatistics, 0, len(monthPeriods))
		for _, period := range monthPeriods {
			if stats, exists := monthlyStats[period]; exists {
				users := stats["uniqueUsers"].(mapset.Set[int])
				visits := stats["totalVisits"].(int64)

				timeRangesResult = append(timeRangesResult, TimeRangeStatistics{
					Period:      period,
					UniqueUsers: int64(users.Cardinality()),
					TotalVisits: visits,
				})
			}
		}

		// 按时间倒序排列
		sort.Slice(timeRangesResult, func(i, j int) bool {
			return timeRangesResult[i].Period > timeRangesResult[j].Period
		})

		result.TimeRanges = timeRangesResult

		// 按创建时间对课程排序
		sort.Slice(result.Courses, func(i, j int) bool {
			timeI := time.Time(result.Courses[i].CreatedAt)
			timeJ := time.Time(result.Courses[j].CreatedAt)
			return timeI.After(timeJ)
		})

		return result
	}

	// 4. 获取所有章节相关的课程数据
	courseIds := make(map[int]bool)
	for _, section := range sections {
		courseIds[section.CourseId] = true
	}

	// 5. 查询这些课程的信息
	var courses []model.Courses
	if len(courseIds) > 0 {
		courseIdList := keys(courseIds)
		model.Course().
			Select("id, title, created_at").
			Where("id IN ?", courseIdList).
			Find(&courses)
	}

	// 如果没有找到课程信息，获取所有课程
	if len(courses) == 0 {
		model.Course().
			Select("id, title, created_at").
			Find(&courses)
	}

	// 6. 创建映射
	courseToTitle := make(map[int]string)
	courseToCreatedAt := make(map[int]ltime.LocalTime)
	for _, course := range courses {
		courseToTitle[course.ID] = course.Title
		courseToCreatedAt[course.ID] = course.CreatedAt
	}

	// 7. 获取每个课程的总章节数
	courseToTotalSections := make(map[int]int)
	for _, course := range courses {
		var count int64
		model.CoursesSection().
			Where("course_id = ?", course.ID).
			Count(&count)
		courseToTotalSections[course.ID] = int(count)
	}

	// 8. 构建章节到课程的映射
	sectionToCourse := make(map[int]int)                // 章节ID -> 课程ID
	sectionToTitle := make(map[int]string)              // 章节ID -> 章节标题
	sectionToCreatedAt := make(map[int]ltime.LocalTime) // 章节ID -> 创建时间
	sectionToSort := make(map[int]int)                  // 章节ID -> 排序值

	// 收集所有课程的所有章节信息，不仅仅是有访问记录的章节
	allCourseSections := make(map[int][]model.CoursesSections) // 课程ID -> 章节列表

	for _, section := range sections {
		sectionToCourse[section.ID] = section.CourseId
		sectionToTitle[section.ID] = section.Title
		sectionToCreatedAt[section.ID] = section.CreatedAt
		sectionToSort[section.ID] = section.Sort

		// 按课程分组章节
		if _, exists := allCourseSections[section.CourseId]; !exists {
			allCourseSections[section.CourseId] = []model.CoursesSections{}
		}
		allCourseSections[section.CourseId] = append(allCourseSections[section.CourseId], section)
	}

	// 9. 为每个课程创建统计对象
	courseStats := make(map[int]*CourseStats)
	for _, course := range courses {
		courseStats[course.ID] = &CourseStats{
			CourseId:      course.ID,
			CourseTitle:   courseToTitle[course.ID],
			UniqueUsers:   0,
			TotalVisits:   0,
			SectionStats:  []SectionStats{},
			CreatedAt:     courseToCreatedAt[course.ID],
			TotalSections: courseToTotalSections[course.ID],
		}
	}

	// 10. 获取每个课程的章节统计数据
	for courseId, courseSections := range allCourseSections {
		courseUserSet := mapset.NewSet[int]()
		visitedSections := mapset.NewSet[int]()
		var courseTotalVisits int64 = 0

		// 处理所有章节
		for _, section := range courseSections {
			// 默认统计值
			uniqueUsers := int64(0)
			totalVisits := int64(0)

			// 如果有访问记录，更新统计值
			if logs, hasLogs := sectionLogs[section.ID]; hasLogs {
				// 标记章节为已访问
				visitedSections.Add(section.ID)

				// 计算章节访问量
				totalVisits = int64(len(logs))

				// 计算章节去重用户量
				userSet := mapset.NewSet[int]()
				for _, log := range logs {
					userSet.Add(log.UserId)
					courseUserSet.Add(log.UserId) // 添加到课程用户集合
				}
				uniqueUsers = int64(userSet.Cardinality())

				// 累加到课程总访问量
				courseTotalVisits += totalVisits
			}

			// 创建章节统计数据
			sectionStat := SectionStats{
				SectionId:    section.ID,
				SectionTitle: section.Title,
				UniqueUsers:  uniqueUsers,
				TotalVisits:  totalVisits,
				CreatedAt:    section.CreatedAt,
				Sort:         section.Sort,
			}

			// 添加到课程统计数据
			if stats, exists := courseStats[courseId]; exists {
				stats.SectionStats = append(stats.SectionStats, sectionStat)
			}
		}

		// 更新课程统计数据
		if stats, exists := courseStats[courseId]; exists {
			stats.UniqueUsers = int64(courseUserSet.Cardinality())
			stats.TotalVisits = courseTotalVisits

			// 计算完成度
			if stats.TotalSections > 0 {
				stats.CompletionRate = float64(visitedSections.Cardinality()) / float64(stats.TotalSections)
			}

			// 对章节按sort字段排序
			if len(stats.SectionStats) > 0 {
				sort.Slice(stats.SectionStats, func(i, j int) bool {
					return stats.SectionStats[i].Sort < stats.SectionStats[j].Sort
				})
			}
		}
	}

	// 11. 转换为数组结果
	courseResult := make([]CourseStats, 0, len(courseStats))
	for _, stats := range courseStats {
		courseResult = append(courseResult, *stats)
	}

	// 12. 根据课程创建时间倒序排序
	sort.Slice(courseResult, func(i, j int) bool {
		timeI := time.Time(courseResult[i].CreatedAt)
		timeJ := time.Time(courseResult[j].CreatedAt)
		return timeI.After(timeJ)
	})

	// 13. 构建时间段统计数据
	timeRangesResult := make([]TimeRangeStatistics, 0, len(monthPeriods))
	for _, period := range monthPeriods {
		if stats, exists := monthlyStats[period]; exists {
			users := stats["uniqueUsers"].(mapset.Set[int])
			visits := stats["totalVisits"].(int64)

			timeRangesResult = append(timeRangesResult, TimeRangeStatistics{
				Period:      period,
				UniqueUsers: int64(users.Cardinality()),
				TotalVisits: visits,
			})
		}
	}

	// 按时间倒序排列
	sort.Slice(timeRangesResult, func(i, j int) bool {
		return timeRangesResult[i].Period > timeRangesResult[j].Period
	})

	// 14. 设置最终结果
	result.Courses = courseResult
	result.TotalUsers = totalUserCount
	result.ActiveUsers = int64(allLearningUsers.Cardinality())
	result.TimeRanges = timeRangesResult

	return result
}

// 辅助函数：从map获取所有键
func keys(m map[int]bool) []int {
	result := make([]int, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}

// TimeSeriesStatistics 时间序列统计
type TimeSeriesStatistics struct {
	TimePoint   string `json:"timePoint"`   // 时间点，格式取决于粒度(2023-10-01, 2023-W40, 2023-10)
	UniqueUsers int64  `json:"uniqueUsers"` // 该时间点的独立用户数
	TotalVisits int64  `json:"totalVisits"` // 该时间点的总访问量
}

// TimeSeriesSummary 时间序列统计摘要
type TimeSeriesSummary struct {
	TotalUniqueUsers int64   `json:"totalUniqueUsers"` // 时间范围内总独立用户数
	TotalVisits      int64   `json:"totalVisits"`      // 时间范围内总访问量
	AvgDailyUsers    float64 `json:"avgDailyUsers"`    // 平均每天用户数
	PeakTimePoint    string  `json:"peakTimePoint"`    // 访问峰值时间点
	PeakVisits       int64   `json:"peakVisits"`       // 峰值访问量
}

// TimeSeriesResponse 时间序列API响应
type TimeSeriesResponse struct {
	Trend   []TimeSeriesStatistics `json:"trend"`   // 时间序列趋势数据
	Summary TimeSeriesSummary      `json:"summary"` // 时间范围统计摘要
}

// GetCoursesTimeSeries 获取课程访问时间序列数据
func (*LogServices) GetCoursesTimeSeries(courseId int, startDate, endDate string, granularity string) TimeSeriesResponse {
	result := TimeSeriesResponse{
		Trend:   []TimeSeriesStatistics{},
		Summary: TimeSeriesSummary{},
	}

	// 1. 参数验证
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		// 默认为90天前
		startTime = time.Now().AddDate(0, 0, -90)
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		// 默认为今天
		endTime = time.Now()
	}

	// 确保结束日期不早于开始日期
	if endTime.Before(startTime) {
		endTime = startTime.AddDate(0, 0, 90) // 默认90天范围
	}

	// 2. 验证粒度
	switch granularity {
	case "week", "month", "quarter":
		// 支持的粒度
	default:
		// 使用默认的"day"粒度
		granularity = "day"
	}

	// 3. 构建查询
	var logs []model.OperLogs
	query := model.OperLog().
		Select("id, user_id, request_info, created_at").
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// 如果指定了课程ID
	if courseId > 0 {
		// 先获取课程的所有章节ID
		var sectionIds []int
		model.CoursesSection().
			Where("course_id = ?", courseId).
			Pluck("id", &sectionIds)

		if len(sectionIds) > 0 {
			// 构建条件数组
			conditions := make([]string, 0, len(sectionIds))
			for _, id := range sectionIds {
				conditions = append(conditions, "/community/courses/section/"+strconv.Itoa(id))
			}
			query = query.Where("request_info IN ?", conditions)
		} else {
			// 没有找到章节，返回空结果
			return result
		}
	} else {
		// 所有课程
		query = query.Where("request_info LIKE ?", "/community/courses/section/%")
	}

	// 执行查询
	query.Find(&logs)

	if len(logs) == 0 {
		return result
	}

	// 4. 按时间点分组统计
	timePointStats := make(map[string]map[string]interface{})
	allUsers := mapset.NewSet[int]()

	var peakTimePoint string
	var peakVisits int64 = 0

	for _, log := range logs {
		// 获取时间点字符串
		timeStr := formatTimePoint(time.Time(log.CreatedAt), granularity)

		// 初始化时间点统计
		if _, exists := timePointStats[timeStr]; !exists {
			timePointStats[timeStr] = map[string]interface{}{
				"uniqueUsers": mapset.NewSet[int](),
				"totalVisits": int64(0),
			}
		}

		// 更新统计
		stats := timePointStats[timeStr]
		userSet := stats["uniqueUsers"].(mapset.Set[int])
		userSet.Add(log.UserId)
		stats["totalVisits"] = stats["totalVisits"].(int64) + 1

		// 全局用户统计
		allUsers.Add(log.UserId)

		// 更新峰值
		if stats["totalVisits"].(int64) > peakVisits {
			peakVisits = stats["totalVisits"].(int64)
			peakTimePoint = timeStr
		}
	}

	// 5. 构建趋势数据
	// 获取时间范围内的所有时间点
	timePoints := generateTimePoints(startTime, endTime, granularity)

	for _, timePoint := range timePoints {
		var uniqueUsers int64 = 0
		var totalVisits int64 = 0

		if stats, exists := timePointStats[timePoint]; exists {
			userSet := stats["uniqueUsers"].(mapset.Set[int])
			uniqueUsers = int64(userSet.Cardinality())
			totalVisits = stats["totalVisits"].(int64)
		}

		result.Trend = append(result.Trend, TimeSeriesStatistics{
			TimePoint:   timePoint,
			UniqueUsers: uniqueUsers,
			TotalVisits: totalVisits,
		})
	}

	// 6. 计算摘要统计
	totalDays := endTime.Sub(startTime).Hours() / 24

	result.Summary = TimeSeriesSummary{
		TotalUniqueUsers: int64(allUsers.Cardinality()),
		TotalVisits:      getTotalVisits(timePointStats),
		AvgDailyUsers:    float64(allUsers.Cardinality()) / totalDays,
		PeakTimePoint:    peakTimePoint,
		PeakVisits:       peakVisits,
	}

	return result
}

// formatTimePoint 格式化时间点字符串
func formatTimePoint(t time.Time, granularity string) string {
	switch granularity {
	case "week":
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		return t.Format("2006-01")
	case "quarter":
		quarter := (t.Month()-1)/3 + 1
		return fmt.Sprintf("%d-Q%d", t.Year(), quarter)
	default: // day
		return t.Format("2006-01-02")
	}
}

// generateTimePoints 生成时间范围内所有的时间点
func generateTimePoints(start, end time.Time, granularity string) []string {
	result := []string{}
	current := start

	// 根据粒度确定步长
	var stepFunc func(time.Time) time.Time

	switch granularity {
	case "week":
		stepFunc = func(t time.Time) time.Time {
			return t.AddDate(0, 0, 7)
		}
	case "month":
		stepFunc = func(t time.Time) time.Time {
			return t.AddDate(0, 1, 0)
		}
	case "quarter":
		stepFunc = func(t time.Time) time.Time {
			return t.AddDate(0, 3, 0)
		}
	default: // day
		stepFunc = func(t time.Time) time.Time {
			return t.AddDate(0, 0, 1)
		}
	}

	// 生成所有时间点
	for current.Before(end) || current.Equal(end) {
		result = append(result, formatTimePoint(current, granularity))
		current = stepFunc(current)
	}

	return result
}

// getTotalVisits 获取总访问量
func getTotalVisits(timePointStats map[string]map[string]interface{}) int64 {
	var total int64 = 0
	for _, stats := range timePointStats {
		total += stats["totalVisits"].(int64)
	}
	return total
}
