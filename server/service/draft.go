package services

import (
	"strconv"
	"strings"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type Draft struct {
}

// 获取临时文本分为： 编辑 / 发布
// 编辑：携带文章id
// 发布：不知文章id，使用 state = 2来判断
func (*Draft) Get(userId, articleId int) *model.Drafts {
	draft := new(model.Drafts)
	tx := model.Draft().Where("user_id = ?", userId)
	// 获取编辑文章的临时文本
	if articleId > 0 {
		tx.Where("article_id = ?", articleId)
	} else {
		tx.Where("state = 2")
	}
	tx.Find(&draft)
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

func (*Draft) Save(draft model.Drafts) {
	if len(draft.Labels) > 0 {
		var ls = make([]string, 0, len(draft.Labels))
		for i := range draft.Labels {
			ls = append(ls, strconv.Itoa(draft.Labels[i]))
		}
		draft.LabelIds = strings.Join(ls, ",")
	}
	if draft.ID > 0 {
		model.Draft().Where("user_id = ? and article_id = ?", draft.UserId, draft.ArticleId).Updates(&draft)
	} else {
		if draft.ArticleId < 1 {
			draft.State = 2
		}
		mysql.GetInstance().Save(&draft)
	}
}

func (*Draft) Del(userId, articleId int) {
	tx := model.Draft().Where("user_id = ?", userId)

	if articleId < 1 {
		tx.Where("state = 2")
	} else {
		tx.Where("article = ?", articleId)
	}
	tx.Update("content", "")
}
