package dao

import "xhyovo.cn/community/server/model"

var InviteCode inviteCode

type inviteCode struct {
}

func (d *inviteCode) SaveCode(code uint16) {

	db.Create(&model.InviteCode{Code: code, State: false})

}

// 是否存在code
func (i *inviteCode) Exist(code uint16) bool {

	var count int64
	object := &model.InviteCode{}
	db.Where("code = ? and state = 0", code).Find(object).Count(&count)

	return count == 1
}

func (*inviteCode) Del(id int) {

	db.Where("id = ?", id).Delete(&model.InviteCode{})
}

func (i *inviteCode) SetState(id uint16) {

	db.Model(&model.InviteCode{}).Where("id = ?", id).Update("state", 1)

}
