package services

import (
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

// get user information
func GetUserById(id uint) *model.User {

	return dao.UserDao.QueryUser(&model.User{ID: id})
}

// update user information
func UpdateUser(id uint, username, pswd string) {

	dao.UserDao.UpdateUser(id, username, pswd)

}
