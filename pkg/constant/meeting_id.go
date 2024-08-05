package constant

/*
*
因订阅需要依赖
SubscriberId int `json:"userId"`                                              // 订阅人

	SendId       int `json:"sendId"`                                              // 发送人   发送人给订阅人发消息 系统默认的是 13
	EventId      int `json:"eventId"`                                             // 事件类型： 人 / 文章
	BusinessId   int `json:"businessId"  binding:"required" msg:"采纳对应业务 id 不能未空"` // 业务id  人id / 文章id

其中 SubscriberId EventId BusinessId 做成唯一索引，但是 BusinessId 是业务 id，而订阅的分享会是没有业务id的，订阅的是整个功能，因此需要有一个唯一id进行组成
*/
const MeetingId = 839258394 // 不能改！
