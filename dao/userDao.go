package dao

import "xhyovo.cn/community/model"

var UserDao userDao

type userDao struct {
}

func (*userDao) QuerySingle(user *model.User) *model.User {

	userObject := new(model.User)
	// db.Where(user).Find(&userObject)

	return userObject
}

func (*userDao) QueryList(user *model.User) []*model.User {

	var users []*model.User
	// db.Where(user).Find(&users)

	return users
}

func (*userDao) Create(user *model.User) int64 {

	res := db.Create(user)
	return res.RowsAffected
}
