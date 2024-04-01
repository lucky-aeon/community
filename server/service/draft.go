package services

import (
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type Draft struct {
}

func (*Draft) Get(userId int) *model.Drafts {

	draft := new(model.Drafts)
	model.Draft().Where("user_id = ?", userId).Find(&draft)
	return draft
}

func (*Draft) Save(draft model.Drafts) {
	if draft.ID > 0 {
		model.Draft().Where("user_id = ?", draft.UserId).Update("content", draft.Content)
	} else {
		mysql.GetInstance().Save(&draft)
	}
}

func (*Draft) Del(userId int) {
	model.Draft().Where("user_id = ?", userId).Update("content", "")
}
