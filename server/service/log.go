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
		db.Where("request_info like ?", "%"+logSearch.RequestInfo+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at <= ? and ? >= created_at", logSearch.StartTime, logSearch.EndTime)
	}
	if logSearch.UserName != "" {
		var userS UserService
		ids := userS.SearchNameSelectId(logSearch.UserName)
		db.Where("user_id in ?", ids)
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
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	return
}
