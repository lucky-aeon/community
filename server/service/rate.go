package services

import (
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type RateService struct {
}

func (RateService) Comment(notes model.Rates) {
	mysql.GetInstance().Save(&notes)
}

func (RateService) Delete(id, userId int) {
	model.Rate().Where("id = ? and user_id = ?", id, userId).Delete(&model.Rates{})
}

func (RateService) GetById(userId int) (note model.Rates) {
	model.Rate().Where("user_id = ?", userId).Find(&note)
	return
}

func (RateService) Page(page, limit int) (state bool, notes []model.Rates) {
	model.Rate().Offset((page - 1) * limit).Limit(limit).Order("created_at desc").Find(&notes)
	if len(notes) == 0 {
		notes = make([]model.Rates, 0)
		return
	}
	state = false
	if len(notes) < limit {
		state = true
	}
	// 获取用户id
	userIds := make([]int, 0)
	for _, v := range notes {
		userIds = append(userIds, v.UserId)
	}
	// 根据用户id获取用户昵称
	users := make([]model.Users, 0)
	model.User().Where("id in (?)", userIds).Select("id", "name").Find(&users)

	for i, v := range notes {
		for _, u := range users {
			if v.UserId == u.ID {
				notes[i].Nickname = u.Name
			}
		}
	}
	return
}
