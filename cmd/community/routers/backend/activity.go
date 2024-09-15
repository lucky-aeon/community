package backend

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitActivityRouters(r *gin.Engine) {
	group := r.Group("/community/admin/activity")
	group.GET("", userDailyActive)
	group.GET("/courses/weekly-trend", courseWeeklyTrend)
	group.GET("/user/weekly-trend", userActivityLineData)
}

func courseWeeklyTrend(ctx *gin.Context) {
	var activityService services.ActivityService
	activity, err := activityService.GetCourseWeeklyTrend()
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.Ok(activity, "").Json(ctx)
}

func userDailyActive(ctx *gin.Context) {
	var activityService services.ActivityService
	activityUserDay, err := activityService.GetUserDailyActive()
	activityUserWeek, err := activityService.GetUserWeeklyActive()
	activityUserMonth, err := activityService.GeUserMonthlyActive()

	activityCourseDay, err := activityService.GetCourseDailyActivity()
	activityCourseWeek, err := activityService.GetCourseWeeklyActivity()
	activityCourseMonth, err := activityService.GetCourseMonthlyActivity()

	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	// 封装结果到 map
	response := map[string]interface{}{
		"userActivityDay":   activityUserDay.ActiveUsers,
		"userActivityWeek":  activityUserWeek.ActiveUsers,  // 假设返回的结构体字段为 ActiveUsers
		"userActivityMonth": activityUserMonth.ActiveUsers, // 假设返回的结构体字段为 ActiveUsers

		"courseActivityDay":   activityCourseDay,   // 假设返回的结构体字段为 ActiveUsers
		"courseActivityWeek":  activityCourseWeek,  // 假设返回的结构体字段为 ActiveUsers
		"courseActivityMonth": activityCourseMonth, // 假设返回的结构体字段为 ActiveUsers
	}
	result.Ok(response, "").Json(ctx)
}

func userActivityLineData(ctx *gin.Context) {
	var activityService services.ActivityService

	// 获取数据
	activity, err := activityService.GetUserActivityLineData()
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	// 定义两个数组，分别存储日期和活跃用户数
	var dates []string
	var activeUsers []int

	// 日期格式化布局
	layoutIn := "2006-01-02T15:04:05-07:00" // 输入的时间格式
	layoutOut := "2006-01-02"               // 输出的时间格式

	// 遍历数据，分别填充日期和活跃用户数
	for _, a := range activity {
		// 解析原始日期
		parsedDate, err := time.Parse(layoutIn, a.Date)
		if err != nil {
			log.Println("日期解析错误:", err)
			continue // 如果日期解析失败，跳过
		}
		// 格式化为 YYYY-MM-DD 格式
		formattedDate := parsedDate.Format(layoutOut)
		dates = append(dates, formattedDate)
		activeUsers = append(activeUsers, a.ActiveUsers) // 假设 a.ActiveUsers 是活跃用户数字段
	}

	// 将两个数组打包到 map 中返回
	response := map[string]interface{}{
		"dates":       dates,
		"activeUsers": activeUsers,
	}

	// 返回封装好的响应
	result.Ok(response, "").Json(ctx)
}
