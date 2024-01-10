package dao

import (
	"xhyovo.cn/community/model"
)

var InviteCode inviteCode

type inviteCode struct {
}

func (*inviteCode) QuerySingle(code *model.InviteCode) *model.InviteCode {

	o := new(model.InviteCode)
	db.Where(&code).Find(&o)

	return o
}

func (*inviteCode) Del(code int) {
	db.Delete(&model.InviteCode{}, code)
}
