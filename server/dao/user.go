package dao

import (
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

func (*UserDao) QueryUserSimple(user *model.Users) (result model.UserSimple, err error) {
	err = model.User().
		Joins("JOIN invite_codes ON invite_codes.code = users.invite_code").
		Joins("JOIN member_infos ON member_infos.id = invite_codes.member_id").
		Select("users.*, member_infos.name as u_role").Where(&user).Find(&result).Error
	return
}

func (*UserDao) CreateUser(account, name, pswd, ininviteCode string) int {

	// 自动接收邮箱
	user := model.Users{Account: account, Name: name, Password: pswd, InviteCode: ininviteCode, Subscribe: 2}
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
	model.User().Where("id in ?", ids).Select("id,name,account,`desc`,avatar").Find(&users)
	return users
}

func (d *UserDao) ExistById(id int) bool {
	var count int64
	model.User().Where("id = ?", id).Count(&count)
	return count == 1
}

func (d *UserDao) GetById(id int) model.Users {
	var user model.Users
	model.User().Where("id = ?", id).First(&user)
	user.Password = ""
	return user
}
