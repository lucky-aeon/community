package constant

const (
	Draft int = iota + 1
	Published
	Pending
	Resolved
	PrivateQuestion
)

var m map[int]string

func GetArticleMsg(id int) string {
	return m[id]
}
func init() {
	m = make(map[int]string)
	m[Draft] = "草稿"
	m[Published] = "发布"
	m[Pending] = "待解决"
	m[Resolved] = "已解决"
	m[PrivateQuestion] = "私密提问"
}
