package services

import (
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/pkg/time"
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
	ID        int            `json:"id"`
	Question  string         `json:"question"`
	UserName  string         `json:"userName"`
	CreatedAt time.LocalTime `json:"createdAt"`
}, int64) {
	var count int64
	var logs []struct {
		ID        int            `json:"id"`
		Question  string         `json:"question"`
		UserName  string         `json:"userName"`
		CreatedAt time.LocalTime `json:"createdAt"`
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

	// 3. 查询这些章节对应的课程ID
	var sections []model.CoursesSections
	model.CoursesSection().
		Select("id, course_id").
		Where("id IN ?", sectionIds).
		Find(&sections)

	// 4. 建立章节到课程的映射
	sectionToCourse := make(map[int]int)
	for _, section := range sections {
		sectionToCourse[section.ID] = section.CourseId
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
