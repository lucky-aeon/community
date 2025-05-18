package services

import (
	"errors"
	"fmt"
	"strings"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

var messageTemplateVar = make(map[string]map[string]string)

func init() {
	messageTemplateVar["user"] = map[string]string{
		"user.id":      "id",
		"user.name":    "name",
		"user.account": "account",
		"user.avatar":  "avatar",
	}
	messageTemplateVar["article"] = map[string]string{
		"article.id":      "id",
		"article.title":   "title",
		"article.content": "content",
		"article.userId":  "user_id",
	}
	messageTemplateVar["comment"] = map[string]string{
		"comment.id":           "id",
		"comment.content":      "content",
		"comment.FromUserName": "from_user_name",
		"comment.ToUserName":   "to_user_name",
		"comment.ArticleTitle": "article_title",
	}
	messageTemplateVar["course"] = map[string]string{
		"course.id":    "id",
		"course.title": "title",
	}
	messageTemplateVar["courses_section"] = map[string]string{
		"courses_section.id":    "id",
		"courses_section.title": "title",
	}
}

// 发送消息中消息模板需要用到的业务id
type SubscribeData struct {
	ArticleId         int
	UserId            int
	CommentId         int
	CurrentBusinessId int // 当前主业务id
	SubscribeId       int // 订阅业务的id(在消息中可以点击跳转的)
	SectionId         int // 章节id
	CourseId          int // 课程id
}

type MessageService struct {
}

func (*MessageService) ListMessageTemplate(page, limit int) ([]*model.MessageTemplates, int64) {
	var count int64
	model.MessageTemplate().Count(&count)
	templates := messageDao.ListMessageTemplate(page, limit)
	eventMap := event.Map()
	for i := range templates {
		templates[i].EventName = eventMap[templates[i].EventId]
	}

	return templates, count
}

func (*MessageService) SaveMessageTemplate(template model.MessageTemplates) error {
	if err := messageDao.SaveMessageTemplate(template); err != nil {
		return errors.New("创建消息模板对应的事件已经存在")
	}
	return nil
}

func (*MessageService) DeleteMessageTemplate(id int) {
	messageDao.DeleteMessageTemplate(id)
}

func (*MessageService) AddMessageLogs(from, types, eventId, businessId int, to []int, content string) {
	var messageLogs []*model.MessageLogs
	for i := range to {
		log := &model.MessageLogs{
			From:      from,
			To:        i,
			Content:   content,
			Type:      types,
			ArticleId: businessId,
			EventId:   eventId,
		}
		messageLogs = append(messageLogs, log)
	}
	messageDao.SaveMessageLogs(messageLogs)
}

func (*MessageService) DeleteMessageLogs(id []int) {
	messageDao.DeleteMessageLogs(id)
}

func (m *MessageService) SendMessage(from, to, types, eventId, businessId int, content string) {

	messageDao.SaveMessage(from, types, eventId, businessId, []int{to}, content)
	// 添加记录
	m.AddMessageLogs(from, types, eventId, businessId, []int{to}, content)
}

func (m *MessageService) SendMessages(from, types, eventId, businessId int, to []int, content string) {

	messageDao.SaveMessage(from, types, eventId, businessId, to, content)
	// 添加记录
	m.AddMessageLogs(from, types, eventId, businessId, to, content)
}

func (*MessageService) ReadMessage2(typee, eventId, businessId, userId int) int64 {
	return messageDao.ReadMessage2(typee, eventId, businessId, userId)
}

func (*MessageService) ReadMessage(id []int, userId int) int64 {
	return messageDao.ReadMessage(id, userId)
}

func (m *MessageService) PageMessage(page, limit, userId, types, state int) (msgs []*model.MessageStates, count int64) {
	msgs = messageDao.ListMessage(page, limit, userId, types, state)
	count = messageDao.CountMessage(userId, types, state)
	return msgs, count
}

// 人：你订阅的 xxx 用户发布了文章，文章标题：xxx
// 文章：你订阅的 xxx 文章，被xxx评论了，评论内容：xxx
func (m *MessageService) GetMsg(template string, b SubscribeData) string {
	BusinessIdMap := businessIdToMap(b)
	for s, v := range messageTemplateVar {
		// 拼接 ${ + s + "." 如果存在则找
		str := fmt.Sprintf("${%s.", s)
		if strings.Contains(template, str) {
			var objet map[string]interface{}
			mysql.GetInstance().Table(s+"s").Where("id = ?", BusinessIdMap[s]).Find(&objet)
			// 遍历 v 从key找template
			for s2 := range v {
				varTemlp := fmt.Sprintf("${%s}", s2)
				if strings.Contains(template, varTemlp) {
					i := objet[v[s2]]
					// 检查值是否为nil
					var valueStr string
					if i != nil {
						valueStr = fmt.Sprintf("%s", i)
					} else {
						valueStr = "" // 如果是nil，替换为空字符串
					}
					template = strings.ReplaceAll(template, varTemlp, valueStr)
				}
			}
		}
	}
	return template
}

func businessIdToMap(b SubscribeData) map[string]int {
	m := map[string]int{
		"user":            b.UserId,
		"article":         b.ArticleId,
		"comment":         b.CommentId,
		"course":          b.CourseId,
		"courses_section": b.SectionId,
	}
	return m
}

func (*MessageService) GetMessageTemplateVar() map[string]map[string]string {
	return messageTemplateVar
}

func (m *MessageService) ClearUnReadMessage(msgType, userId int) {
	model.MessageState().Where("`to` = ? and type = ?", userId, msgType).Update("state", 0)
}

func (m *MessageService) GetUnReadMessageCountByUserId(userId int) (count int64) {
	model.MessageState().Where("`to` = ? and state = 1", userId).Count(&count)
	return
}
