package model

type CourseActivity struct {
	CourseTitle      string `json:"courseTitle"`
	ActivityDate     string `json:"activityDate,omitempty"` // 活跃度趋势会使用这个字段
	DailyActiveUsers int    `json:"DailyActiveUsers"`       // 用于日活、周活、月活
}

type CourseWeekActivity struct {
	CourseTitle       string `json:"courseTitle"`
	WeeklyActiveUsers int    `json:"weeklyActiveUsers"`
}

type CourseMonthActivity struct {
	CourseTitle        string `json:"courseTitle"`
	MonthlyActiveUsers int    `json:"monthlyActiveUsers"`
}

// 用于返回日活、周活、月活的数据
type ActivityData struct {
	ActiveUsers int `json:"activeUsers"`
}

// 用于折线图的数据
type ActivityLine struct {
	Date        string `json:"date"`
	ActiveUsers int    `json:"activeUsers"`
}
