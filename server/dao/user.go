package dao

import (
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type UserDao struct {
}

// set user info by article user id
func (*UserDao) SetUserInfo(articles []*model.Articles) {
	var userIds []int

	for _, article := range articles {
		userIds = append(userIds, article.UserId)
	}
	users := []model.Users{}
	model.User().Find(&users, userIds)
	var userMap = make(map[int]model.Users)
	for _, user := range users {
		user.Avatar = utils.BuildFileUrl(user.Avatar)
		userMap[user.ID] = user
	}
	// articles set userinfo
	for _, article := range articles {
		article.Users = userMap[article.UserId]
	}

}

func (*UserDao) QueryUsersByUserIds(ids []int) []model.Users {
	users := []model.Users{}
	model.User().Find(&users, ids)
	return users
}

func (*UserDao) QueryUser(user *model.Users) *model.Users {

	model.User().Where(&user).Find(&user)
	return user
}

func (*UserDao) CreateUser(account, name, pswd string, ininviteCode int) int {

	user := model.Users{Account: account, Name: name, Password: pswd, InviteCode: ininviteCode}
	model.User().Create(&user)
	return user.ID
}

func (d *UserDao) UpdateUser(user *model.Users) {

	model.User().Where("id = ?", user.ID).Updates(&user)
}

func (d *UserDao) ListByIds(id ...int) []string {
	var email []string
	model.User().Where("id in ?", id).Select("account").Find(&email)
	return email

}

func (d *UserDao) ListByIdsSelectIdName(ids []int) []model.Users {
	var users []model.Users
	model.User().Where("id in ?", ids).Select("id,name").Find(&users)
	return users
}
