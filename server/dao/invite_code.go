package dao

import "xhyovo.cn/community/server/model"

var InviteCode inviteCode

type inviteCode struct {
}

func (d *inviteCode) SaveCode(code int) error {

	return nil
}

// 是否存在code
func (i *inviteCode) Exist(code uint16) bool {

	var count int64
	object := &model.InviteCode{}
	db.Where("code = ?", code).Find(object).Count(&count)

	return count == 1
}

func (*inviteCode) Del(code int) error {

	return nil
}

func (i *inviteCode) SetState(code uint16) error {

	return nil
}
