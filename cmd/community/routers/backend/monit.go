package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitMonitRouters(r *gin.Engine) {
	group := r.Group("/community/admin/monit")
	group.GET("", listMonitUser)
	group.GET("/ip/:userId", getMonitUserIpDetails)
	group.GET("/section/:userId", getMonitUserSectionDetails)
}

// 监控用户 监控：ip 登陆大于 5 次不同 并且不是 127.0.0.1
func listMonitUser(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	type Result struct {
		UserID int    `json:"userID"`
		Count  int    `json:"count"`
		Name   string `gorm:"-" json:"name"`
	}

	var results []Result
	var count int64
	tx := mysql.GetInstance().Table("oper_logs").
		Select("user_id, COUNT(DISTINCT ip) as count").
		Where("ip <> ?", "127.0.0.1").
		Where("user_id <> ?", 13).
		Group("user_id").
		Having("COUNT(DISTINCT ip) > ?", 5)
	tx.Count(&count)
	tx.Offset((p - 1) * limit).
		Limit(limit).
		Order("count DESC")
	tx.Scan(&results)

	var userIds []int
	for _, item := range results {
		userIds = append(userIds, item.UserID)
	}

	var userService services.UserService
	nameMap := userService.ListByIdsToMap(userIds)

	for i, _ := range results {
		results[i].Name = nameMap[results[i].UserID].Name
	}
	result.Page(results, count, nil).Json(ctx)
	return
}

// 获取监控用户的操作
func getMonitUserIpDetails(ctx *gin.Context) {

	userId := ctx.Param("userId")
	p, limit := page.GetPage(ctx)

	type ipMonit struct {
		IP        string `json:"ip"`
		Platform  string `json:"platform"`
		UserAgent string `json:"userAgent"`
		Count     int    `json:"count"`
	}

	var ipResult []ipMonit
	var count int64

	// 查询去重后的IP，不包括 127.0.0.1
	tx := mysql.GetInstance().Table("oper_logs").
		Select("ip, platform, user_agent, COUNT(ip) as count").
		Where("ip <> ?", "127.0.0.1").
		Where("user_id = ?", userId).
		Group("ip, platform, user_agent")
	tx.Count(&count)
	tx.Limit(limit).
		Offset((p - 1) * limit).
		Order("count DESC")
	tx.Scan(&ipResult)

	result.Page(ipResult, count, nil).Json(ctx)

}

func getMonitUserSectionDetails(ctx *gin.Context) {
	userId := ctx.Param("userId")
	p, limit := page.GetPage(ctx)
	// 收集章节观看次数
	type sectionMonit struct {
		UserID       int    `json:"userID"`
		RequestInfo  string `json:"requestInfo"`
		Count        int    `json:"count"`
		SectionTitle string `json:"sectionTitle" gorm:"-"`
		SectionId    int    `json:"sectionId" gorm:"-"`
	}

	var sectionResult []sectionMonit
	var count int64
	tx := mysql.GetInstance().Table("oper_logs").
		Select("user_id, request_info, COUNT(*) as count").
		Where("request_info REGEXP ?", `^/community/courses/section/[0-9]+$`).
		Where("user_id = ?", userId).
		Group("user_id, request_info")

	tx.Count(&count).
		Limit(limit).
		Offset((p - 1) * limit).
		Order("count DESC")
	tx.Scan(&sectionResult)

	var sectionIds []int
	for i, item := range sectionResult {
		id, err := extractIDFromURL(item.RequestInfo)
		if err != nil {
			continue
		}
		idValue, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		sectionIds = append(sectionIds, idValue)
		sectionResult[i].SectionId = idValue
	}

	var sectionService services.CourseService
	selectIdTitleMap := sectionService.ListSectionByIds(sectionIds)
	for i := range sectionResult {
		sectionResult[i].SectionTitle = selectIdTitleMap[sectionResult[i].SectionId]
	}

	result.Page(sectionResult, count, nil).Json(ctx)
}

func extractIDFromURL(url string) (string, error) {
	re := regexp.MustCompile(`/community/courses/section/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("no ID found in URL")
	}
	return matches[1], nil
}
