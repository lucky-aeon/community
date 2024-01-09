package services

import (
	"errors"

	"xhyovo.cn/community/dao"
	"xhyovo.cn/community/model"
)

func Login(account, pswd string) model.User {

	user := &model.User{Account: account, Password: pswd}

	user = dao.UserDao.QuerySingle(user)
	if user == nil {
		errors.New("登陆失败！账号或密码错误")
	}

	return *user
}

func Register(account, pswd, name string, inviteCode int) {

	// query code

}
