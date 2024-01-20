package dao

import (
	"xhyovo.cn/community/server/model"
)

type UserDao struct {
}

func (*UserDao) QueryUser(user *model.User) *model.User {

	db.Where(&user).Find(&user)
	return user
}

func (*UserDao) CreateUser(account, name, pswd string, ininviteCode uint16) uint {

	user := model.User{Account: account, Name: name, Password: pswd, InviteCode: ininviteCode}
	db.Create(user)
	return user.ID
}

func (d *UserDao) UpdateUser(id uint, username, pswd string) {
	db.Model(&model.User{}).Where("id = ?", id).Update("name", username).Update("password", pswd)
}
