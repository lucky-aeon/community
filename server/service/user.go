package services

import (
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserService struct {
}

// get user information
func (*UserService) GetUserById(id uint) *model.Users {

	user := UserDao.QueryUser(&model.Users{ID: id})
	user.Avatar = utils.BuildFileUrl(user.Avatar)
	return user
}

// update user information
func (*UserService) UpdateUser(user *model.Users) {

	UserDao.UpdateUser(user)

}
