package services

import (
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type ActivityService struct {
}

func (s *ActivityService) GetCourseDailyActivity() ([]model.CourseActivity, error) {
	var a dao.ActivityDAO
	return a.GetCourseDailyActivity()
}

func (s *ActivityService) GetCourseWeeklyActivity() ([]model.CourseWeekActivity, error) {
	var a dao.ActivityDAO
	return a.GetCourseWeeklyActivity()
}

func (s *ActivityService) GetCourseMonthlyActivity() ([]model.CourseMonthActivity, error) {
	var a dao.ActivityDAO
	return a.GetCourseMonthlyActivity()
}

func (s *ActivityService) GetCourseWeeklyTrend() ([]model.CourseActivity, error) {
	var a dao.ActivityDAO
	return a.GetCourseWeeklyTrend()
}

func (s *ActivityService) GetUserDailyActive() (model.ActivityData, error) {
	var a dao.ActivityDAO
	return a.GetUserDailyActiveUsers()
}

func (s *ActivityService) GetUserWeeklyActive() (model.ActivityData, error) {
	var a dao.ActivityDAO
	return a.GetUserWeeklyActiveUsers()
}

func (s *ActivityService) GeUserMonthlyActive() (model.ActivityData, error) {
	var a dao.ActivityDAO
	return a.GetUserMonthlyActiveUsers()
}

func (s *ActivityService) GetUserActivityLineData() ([]model.ActivityLine, error) {
	var a dao.ActivityDAO
	return a.GetUserActivityLineData()
}
