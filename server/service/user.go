package services

import (
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

var UserDao dao.UserDao

type UserService struct {
}

// get user information
func (*UserService) GetUserById(id uint) *model.User {

	return UserDao.QueryUser(&model.User{ID: id})
}

// update user information
func (*UserService) UpdateUser(id uint, username, pswd string) {

	UserDao.UpdateUser(id, username, pswd)

}
