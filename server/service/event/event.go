package event

const (
	CommentUpdateEvent = iota + 1 // 文章下评论更新事件
	UserFollowingEvent            // 用户关注的人发布文章事件
	ArticleAt                     // 文章中@
	CommentAt                     // 评论中@
	ReplyComment                  // 评论回复
)

var events []*event

var eventMap = make(map[int]string)

func init() {
	events = append(events, nil)
	events = append(events, &event{Id: CommentUpdateEvent, Msg: "文章评论"})
	events = append(events, &event{Id: UserFollowingEvent, Msg: "用户文章更新"})
	events = append(events, &event{Id: ArticleAt, Msg: "文章 @"})
	events = append(events, &event{Id: CommentAt, Msg: "评论 @"})
	events = append(events, &event{Id: ReplyComment, Msg: "评论回复"})

	eventMap[CommentUpdateEvent] = "文章评论"
	eventMap[UserFollowingEvent] = "用户文章更新"
	eventMap[ArticleAt] = "文章 @"
	eventMap[CommentAt] = "评论 @"
	eventMap[ReplyComment] = "评论回复"
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
func Map() map[int]string {
	return eventMap
}
