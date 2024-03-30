package constant

import "xhyovo.cn/community/server/response"

const (
	Draft int = iota + 1
	Published
	Pending
	Resolved
	PrivateQuestion
	QADraft
	Top
)

var name map[int]string

var msg map[int]string

func GetArticleName(id int) string {
	return name[id]
}
func GetArticleMsg(id int) string {
	return msg[id]
}
func init() {
	name = make(map[int]string)
	name[Draft] = "草稿"
	name[Published] = "发布"
	name[Pending] = "待解决"
	name[Resolved] = "已解决"
	name[PrivateQuestion] = "私密提问"
	name[Top] = "置顶"
	msg = make(map[int]string)
	msg[Draft] = "保存草稿"
	msg[Published] = "发布成功"
	msg[Pending] = "发布待解决"
	msg[Resolved] = "修改为已解决"
	msg[PrivateQuestion] = "发布为私密提问"
	msg[QADraft] = "保存草稿"
	msg[Top] = "置顶"
}

func ListState() []response.ArticleState {
	var states = make([]response.ArticleState, 0, len(name))
	for k, v := range name {
		states = append(states, response.ArticleState{Id: k, Name: v})

	}
	return states
}
