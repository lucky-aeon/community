package services

import "xhyovo.cn/community/server/dao"

var (
	messageDao dao.MessageDao
	articleDao dao.Article
	fileDao    dao.File
	codeDao    dao.InviteCode
	typeDao    dao.Type
	userDao    dao.UserDao
	commentDao dao.CommentDao
)
