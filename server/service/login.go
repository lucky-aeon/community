package services

import (
	"errors"

	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

func Login(account, pswd string) (*model.User, error) {

	user := dao.UserDao.QueryUser(&model.User{Account: account, Password: pswd})
	if user.ID == 0 {
		return nil, errors.New("登录失败！账号或密码错误")
	}

	return user, nil
}

func Register(account, pswd, name string, inviteCode uint16) error {

	if err := utils.NotBlank(account, pswd, name, inviteCode); err != nil {
		return err
	}

	// query code
	if !dao.InviteCode.Exist(inviteCode) {
		return errors.New("验证码不存在")
	}

	// 查询账户
	user := dao.UserDao.QueryUser(&model.User{Account: account})
	if user.ID == 1 {
		return errors.New("账户已存在,换一个吧")
	}

	// 保存用户
	dao.UserDao.CreateUser(account, name, pswd, inviteCode)
	// 修改code状态
	SetState(inviteCode)

	return nil
}
