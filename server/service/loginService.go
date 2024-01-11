package services

import (
	"errors"

	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

func Login(account, pswd string) (model.User, error) {

	user := &model.User{Account: account, Password: pswd}

	user = dao.UserDao.QuerySingle(user)
	if user == nil {
		return *user, errors.New("登陆失败！账号或密码错误")

	}

	return *user, nil
}

func Register(account, pswd, name string, inviteCode int) error {

	// query code
	if !dao.InviteCode.Exist(inviteCode) {
		return errors.New("验证码不存在")
	}

	// query account
	user := dao.UserDao.QuerySingle(&model.User{Account: account})
	if user != nil {
		return errors.New("账户已存在,换一个吧")
	}

	return nil
}
