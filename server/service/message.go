package services

import (
	"fmt"
	"strings"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
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
}

// 发送消息中消息模板需要用到的业务id
type BusinessId struct {
	ArticleId         int
	UserId            int
	CommentId         int
	CurrentBusinessId int // 当前主业务id
}

type MessageService struct {
}

func (*MessageService) ListMessageTemplate(page, limit int) ([]*model.MessageTemplates, int64) {
	var count int64
	model.MessageTemplate().Count(&count)
	return messageDao.ListMessageTemplate(page, limit), count
}

func (*MessageService) SaveMessageTemplate(template model.MessageTemplates) {
	messageDao.SaveMessageTemplate(template)
}

func (*MessageService) DeleteMessageTemplate(id int) {
	messageDao.DeleteMessageTemplate(id)
}

func (*MessageService) ListMessageLogs(page, limit int) []*model.MessageLogs {
	return messageDao.ListMessageLogs(page, limit)
}

func (*MessageService) AddMessageLogs(from, types, businessId int, to []int, content string) {
	var messageLogs []*model.MessageLogs
	for i := range to {
		log := &model.MessageLogs{
			From:      from,
			To:        i,
			Content:   content,
			Type:      types,
			ArticleId: businessId,
		}
		messageLogs = append(messageLogs, log)
	}
	messageDao.SaveMessageLogs(messageLogs)
}

func (*MessageService) DeleteMessageLogs(id []int) {
	messageDao.DeleteMessageLogs(id)
}

func (m *MessageService) SendMessage(from, to, types, businessId int, content string) {

	messageDao.SaveMessage(from, types, businessId, []int{to}, content)
	// 添加记录
	m.AddMessageLogs(from, types, businessId, []int{to}, content)
}

func (m *MessageService) SendMessages(from, types, businessId int, to []int, content string) {

	messageDao.SaveMessage(from, types, businessId, to, content)
	// 添加记录
	m.AddMessageLogs(from, types, businessId, to, content)
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
func (m *MessageService) GetMsg(template string, b BusinessId) string {
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
					template = strings.ReplaceAll(template, varTemlp, fmt.Sprintf("%s", i))
				}
			}
		}
	}
	return template
}

func businessIdToMap(b BusinessId) map[string]int {
	m := map[string]int{
		"user":    b.UserId,
		"article": b.ArticleId,
		"comment": b.CommentId,
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
