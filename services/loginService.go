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

func Register(account, pswd, name string, inviteCode int) error {

	// query code
	code := dao.InviteCode.QuerySingle(&model.InviteCode{Code: inviteCode})
	if code == nil {
		errors.New("验证码不存在")
	}

	// query account
	user := dao.UserDao.QuerySingle(&model.User{Account: account})
	if user != nil {
		errors.New("账户已存在,换一个吧")
	}

	return nil
}
