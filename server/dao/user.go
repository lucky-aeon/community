package dao

import (
	"xhyovo.cn/community/server/model"
)

var UserDao userDao

type userDao struct {
}

func (*userDao) QueryUser(user *model.User) *model.User {

	db.Where(&user).Find(&user)
	return user
}

func (*userDao) CreateUser(account, name, pswd string, ininviteCode uint16) uint {
	user := model.User{Account: account, Name: name, Password: pswd, InviteCode: ininviteCode}
	db.Create(user)
	return user.ID
}

func (d *userDao) UpdateUser(id uint, username, pswd string) {
	db.Model(&model.User{}).Where("id = ?", id).Update("name", username).Update("password", pswd)
}
