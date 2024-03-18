package event

const (
	CommentUpdateEvent = iota + 1 // 文章下评论更新事件
	UserFollowingEvent            // 用户关注的人发布文章事件
	ArticleAt                     // 文章中@
	CommentAt                     // 评论中@
)

var events []*event

func init() {
	events = append(events, nil)
	events = append(events, &event{Id: CommentUpdateEvent, Msg: "文章评论"})
	events = append(events, &event{Id: UserFollowingEvent, Msg: "用户"})
}

// 事件
type event struct {
	Id  int    `json:"id"`
	Msg string `json:"msg"`
}

func GetMsg(eventId int) string {
	v := events[eventId]
	return v.Msg
}

func List() []*event {
	return events
}
