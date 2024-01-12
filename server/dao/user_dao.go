package dao

import (
	"xhyovo.cn/community/server/model"
)

var UserDao userDao

type userDao struct {
}

func (*userDao) QueryUserByAccountAndPswd(account, pswd string) (*model.User, error) {
	return nil, nil
}

func (*userDao) ExistUserByAccount(account string) bool {
	return true
}

func (*userDao) CreateUser(account, name, pswd string, ininviteCode uint16) (int64, error) {
	return 1, nil
}
