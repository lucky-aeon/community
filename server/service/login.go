package services

import (
	"errors"

	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

func Login(account, pswd string) (*model.User, error) {

	user, err := dao.UserDao.QueryUserByAccountAndPswd(account, pswd)
	if err != nil {
		return user, errors.New("登录失败！账号或密码错误")

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

	// query account
	if dao.UserDao.ExistUserByAccount(account) {
		return errors.New("账户已存在,换一个吧")
	}

	// save
	_, err := dao.UserDao.CreateUser(account, name, pswd, inviteCode)
	if err != nil {
		return err
	}

	// del code
	if err = SetState(inviteCode); err != nil {
		return err
	}

	return nil
}
