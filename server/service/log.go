package services

import (
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/server/model"
)

type LogServices struct {
}

func (*LogServices) GetPageOperLog(page, limit int, logSearch model.LogSearch) (logs []model.OperLogs, count int64) {
	db := model.OperLog()
	if logSearch.RequestMethod != "" {
		db.Where("request_method = ?", logSearch.RequestMethod)
	}
	if logSearch.RequestInfo != "" {
		db.Where("request_info like ?", "%s"+logSearch.RequestInfo+"%s")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%s"+logSearch.Ip+"%s")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at <= ? and ? >= created_at", logSearch.StartTime, logSearch.EndTime)
	}
	if logSearch.UserName != "" {
		var userS UserService
		ids := userS.SearchNameSelectId(logSearch.UserName)
		db.Where("user_id in ?", ids)
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	db.Count(&count)
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
	model.OperLog().Create(&log)
}

func (*LogServices) InsertLoginLog(log model.LoginLogs) {
	model.LoginLog().Create(&log)
}

func (s *LogServices) GetPageLoginPage(page, limit int, logSearch model.LogSearch) (logs []model.LoginLogs, count int64) {
	db := model.LoginLog()
	if logSearch.Account != "" {
		db.Where("account like ?", "%"+logSearch.Account+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	db.Count(&count)
	return
}
