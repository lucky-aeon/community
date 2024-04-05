package services

import (
	"strconv"
	"strings"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type Draft struct {
}

func (*Draft) Get(userId int) *model.Drafts {
	draft := new(model.Drafts)
	model.Draft().Where("user_id = ?", userId).Find(&draft)
	if len(draft.LabelIds) > 0 {
		split := strings.Split(draft.LabelIds, ",")
		var ids = make([]int, 0, len(split))
		for i := range split {
			id, _ := strconv.Atoi(split[i])
			ids = append(ids, id)
		}
	}

	return draft
}

func (*Draft) InitDraft(userId int) {
	go func() {
		mysql.GetInstance().Create([]model.Drafts{model.Drafts{UserId: userId, State: 1}, model.Drafts{UserId: userId, State: 2}})

	}()
}

func (*Draft) Save(draft model.Drafts) {
	if len(draft.Labels) > 0 {
		var ls = make([]string, 0, len(draft.Labels))
		for i := range draft.Labels {
			ls = append(ls, strconv.Itoa(draft.Labels[i]))
		}
		draft.LabelIds = strings.Join(ls, ",")
	}
	model.Draft().Where("user_id = ? and state = ?", draft.UserId, 2).Updates(&draft)
}

func (*Draft) DelDraft(userId int) {

	model.Draft().Where("user_id = ? and state = ?", userId, 1).Updates(map[string]interface{}{"content": "", "type": 0, "label_ids": ""})
}
