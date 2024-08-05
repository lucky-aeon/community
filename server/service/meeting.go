package services

import (
	"errors"
	"fmt"
	"time"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/delay"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/request"

	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

const msgTemp = "用户: %s, 申请了会议，会议标题为: %s,会议描述为: %s"

const msgSuccessTemp = "会议主题：%s 会议已通过审核, 会议截止报名时间为: %v, 开始时间为: %v"

const signupEndTimeTemp = "你参与的 %s 会议报名已截止，你的参会邀请信息：%s"

const startTimeTemp = "你参与的 %s 会议已开始，请及时参会，会议信息为：%s"

type MeetingService struct {
}

func init() {
	// 等待 db 初始化
	go func() {
		time.Sleep(3 * time.Second)
		initMeetingTasks()
	}()

}

func (m *MeetingService) GetJoinMeetingUserSelectAvatar(meetingId int) []string {

	userIds := m.GetJoinUsers(meetingId)
	var avatar []string
	model.User().Where("id in ?", userIds).Select("avatar").Find(&avatar)

	return avatar
}

// 保存
func (m *MeetingService) Save(meeting model.Meetings) error {

	// 会议申请时间不低于当前时间
	if time.Time(meeting.InitiatorTime).Before(time.Now()) {
		return errors.New("会议申请时间不可小于当前时间")
	}

	// 如果是修改
	if meeting.Id > 0 {
		// 会议状态只一旦非审核中则不可修改
		meetingState := m.GetById(meeting.Id).State
		if meetingState != constant.Reviewing {
			return errors.New("会议已被锁定,不可修改")
		}
	}
	meeting.State = constant.Reviewing

	if meeting.Id == 0 {

		if err := model.Meeting().Save(&meeting).Error; err != nil {
			return err
		}
	} else {
		if err := model.Meeting().Where("id = ?", meeting.Id).Updates(&meeting).Error; err != nil {
			return err
		}
	}

	var userS UserService
	user := userS.GetUserById(meeting.InitiatorId)

	// 发送消息给订阅人
	var subS SubscriptionService
	message := fmt.Sprintf(msgTemp, user.Name, meeting.Title, meeting.Description)
	subS.SendMsg(13, event.Meeting, constant.NOTICE, constant.MeetingId, message)
	return nil
}

// 获取
func (*MeetingService) GetById(id int) model.Meetings {
	var meeting model.Meetings
	model.Meeting().Find(&meeting, id)
	if meeting.Id == 0 {
		return meeting
	}

	var userS UserService
	name := userS.GetUserById(meeting.InitiatorId).Name

	meeting.InitiatorName = name

	var joinUsers []model.MeetingJoinUsers
	model.MeetingJoinUser().Where("meeting_id = ?", meeting.Id).Find(&joinUsers)
	// 如果会议已完成则显示昵称，否则只显示人数
	if meeting.State == constant.Completed {
		var userIds []int

		for _, user := range joinUsers {
			userIds = append(userIds, user.UserId)
		}
		nameMap := userS.ListByIdsToMap(userIds)
		for i, _ := range joinUsers {
			joinUsers[i].UserName = nameMap[joinUsers[i].UserId].Name
		}
		meeting.JoinUsers = joinUsers
	} else {
		meeting.JoinUserCount = len(joinUsers)
	}

	return meeting
}

func (*MeetingService) Page(page, limit, userId int) ([]model.Meetings, int64) {
	var meetings []model.Meetings
	var count int64
	tx := model.Meeting()
	tx.Count(&count)

	if count == 0 {
		return meetings, 0
	}

	tx.Limit(limit).Offset((page - 1) * limit).Order("created_at desc")
	if userId != 0 {
		tx.Where("initiator_id = ?", userId)
	}
	tx.Find(&meetings)

	// 收集用户id
	var userIds []int
	for _, meeting := range meetings {
		userIds = append(userIds, meeting.InitiatorId)
	}

	var userS UserService
	nameMap := userS.ListByIdsToMap(userIds)
	for i := range meetings {
		meetings[i].InitiatorName = nameMap[meetings[i].InitiatorId].Name
		meetings[i].InitiatorAvatar = nameMap[meetings[i].InitiatorId].Avatar
	}

	return meetings, count
}

// 管理员删除 传0 即可
func (*MeetingService) DeleteById(id, userId int) error {
	// 只有在审核中可以删除
	var meeting model.Meetings
	model.Meeting().Where("id = ?", id).First(&meeting)
	if meeting.State != constant.Reviewing {
		return errors.New("除了审核状态外均不可删除")
	}
	db := model.Meeting()
	if userId != 0 {
		db.Where("Initiator_id", userId)
	}
	db.Where("id = ?", id).Delete(&model.Meetings{})

	return nil
}

// 加入会议
func (m *MeetingService) JoinMeeting(id, userId int) error {
	meeting := m.GeyByIdSample(id)
	if meeting.Id == 0 {
		return errors.New("会议不存在")
	}
	// 发起者不可申请加入，本身已在
	if meeting.InitiatorId == userId {
		return errors.New("会议发起者不可申请加入，你已在里面")
	}

	// 必须是报名中才可加入
	if meeting.State != constant.Registering {
		return errors.New("参与会议必须是报名中")
	}

	var joinUsers model.MeetingJoinUsers
	joinUsers.UserId = userId
	joinUsers.MeetingId = id
	if err := model.MeetingJoinUser().Save(&joinUsers).Error; err != nil {
		return errors.New("不可重复加入")
	}
	return nil
}

// 退出会议
func (m *MeetingService) QuitJoinMeeting(id, userId int) error {
	meeting := m.GeyByIdSample(id)
	if meeting.Id == 0 {
		return errors.New("会议不存在")
	}

	// 必须是报名中或者筹备中可以退出
	if meeting.State != constant.Registering && meeting.State != constant.Preparing {
		return errors.New("退出会议只能是报名中或者筹备中状态")
	}
	model.MeetingJoinUser().Where("meeting_id = ? and user_id = ?", id, userId).Delete(&model.MeetingJoinUsers{})
	return nil
}

func (*MeetingService) GeyByIdSample(id int) model.Meetings {
	var meeting model.Meetings
	model.Meeting().Where("id = ?", id).First(&meeting)
	return meeting
}

func (*MeetingService) ExistById(id int) bool {
	var count int64
	model.Meeting().Where("id = ? ", id).Count(&count)
	return count == 1
}

func (m *MeetingService) Approve(reqApproveMeeting request.ReqApproveMeeting) error {
	meeting := m.GetById(reqApproveMeeting.Id)
	if meeting.Id == 0 {
		return errors.New("操作会议不存在")
	}

	// 会议报名时间 不能 大于 开始时间，开始时间不能大于结束时间
	signupEndTime := time.Time(*reqApproveMeeting.SignupEndTime)
	startTime := time.Time(*reqApproveMeeting.MeetingStartTime)
	endTime := time.Time(*reqApproveMeeting.MeetingEndTime)

	if meeting.State != constant.Reviewing {
		return errors.New("当前会议状态不可通过")
	}

	if signupEndTime.Before(time.Now()) {
		return errors.New("会议报名时间不能小于当前时间")
	}

	if signupEndTime.After(startTime) {
		return errors.New("会议报名时间不能大于会议开始时间")
	}

	if signupEndTime.After(endTime) {
		return errors.New("会议报名时间不能大于会议结束时间")
	}
	if startTime.After(endTime) {
		return errors.New("会议开始时间不能大于会议结束时间")
	}

	meeting.MeetingStartTime = reqApproveMeeting.MeetingStartTime
	meeting.MeetingEndTime = reqApproveMeeting.MeetingEndTime
	meeting.SignupEndTime = reqApproveMeeting.SignupEndTime
	meeting.MeetingLink = reqApproveMeeting.MeetingLink

	// {name} 会议已通过审核,截止报名时间为：{time}，开始时间为：{time}

	signupMessage := fmt.Sprintf(msgSuccessTemp, meeting.Title, signupEndTime.Format("2006-01-02 15:04:05"), startTime.Format("2006-01-02 15:04:05"))

	// 给订阅人发送邮箱
	var subS SubscriptionService
	subS.SendMsg(13, event.Meeting, constant.NOTICE, constant.MeetingId, signupMessage)

	userIds := m.GetJoinUsers(meeting.Id)

	approveAddTask(meeting, userIds)

	// 修改状态
	meeting.State = constant.Registering

	model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)

	return nil
}

func (m *MeetingService) Pass(reqPassMeeting request.ReqPassMeeting) error {
	meeting := m.GetById(reqPassMeeting.Id)
	if meeting.Id == 0 {
		return errors.New("操作会议不存在")
	}

	if meeting.State != constant.Reviewing {
		return errors.New("会议状态只能是审核中才能被 PASSda∂")
	}

	meeting.State = constant.Pass
	meeting.StateMessage = reqPassMeeting.PassMessage
	model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
	return nil
}
func (m *MeetingService) GetJoinUsers(meetingId int) []int {
	var userIds []int
	model.MeetingJoinUser().Where("meeting_id = ?", meetingId).Select("user_id").Find(&userIds)
	return userIds
}

func (m *MeetingService) Record(reqRecordMeeting request.ReqRecordMeeting) error {
	meeting := m.GetById(reqRecordMeeting.Id)
	if meeting.Id == 0 {
		return errors.New("操作会议不存在")
	}

	// 必须是完成后才可填写
	if meeting.State != constant.Completed {
		return errors.New("会议必须是已完成才可填写会议记录")
	}

	meeting.Record = reqRecordMeeting.Record
	model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
	return nil
}

func (m *MeetingService) InMeetingState(meetingId, userId int) (bool, error) {
	meeting := m.GetById(meetingId)
	if meeting.Id == 0 {
		return false, errors.New("操作会议不存在")
	}
	var count int64
	model.MeetingJoinUser().Where("meeting_id = ? AND user_id = ?", meetingId, userId).Count(&count)
	return count == 1, nil
}

/*
*
初始化定时任务
查出 报名中，筹备中，会议中 的会议状态

报名中状态：添加 报名截止，会议开始，会议结束定时任务
筹备中：添加会议开始，会议结束定时任务
会议中：添加会议结束定时任务
*/
func initMeetingTasks() {
	// 查出在 当前时间 < 报名时间的 数据
	var meetings []model.Meetings
	model.Meeting().Where("signup_end_time > ?", time.Now()).Find(&meetings)
	// 存入队列中
	for _, meeting := range meetings {
		var meetingService MeetingService
		userIds := meetingService.GetJoinUsers(meeting.Id)
		if meeting.State == constant.Registering {
			approveAddTask(meeting, userIds)
		} else if meeting.State == constant.Preparing {
			preparingAddTask(meeting, userIds)
		} else if meeting.State == constant.InMeeting {
			inMeetingAddTask(meeting)
		}

	}
}

func inMeetingAddTask(meeting model.Meetings) {
	endTime := time.Time(*meeting.MeetingEndTime)
	// 加入延迟队列
	delayQueue := delay.GetInstant()

	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	// 会议结束后状态改为会议完成
	delayQueue.Add(meeting.Id, endTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.Completed)
		meeting.State = constant.Completed
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
	})
}

func preparingAddTask(meeting model.Meetings, userIds []int) {
	startTime := time.Time(*meeting.MeetingStartTime)
	endTime := time.Time(*meeting.MeetingEndTime)

	var subS SubscriptionService
	startMessage := fmt.Sprintf(startTimeTemp, meeting.Title, meeting.MeetingLink)
	// 加入延迟队列
	delayQueue := delay.GetInstant()

	// 会议开始后状态改为会议中
	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	delayQueue.Add(meeting.Id, startTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.InMeeting)
		meeting.State = constant.InMeeting
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
		subS.SendMsgByToIds(13, event.Meeting, constant.NOTICE, constant.MeetingId, userIds, startMessage)
	})

	// 会议结束后状态改为会议完成
	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	delayQueue.Add(meeting.Id, endTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.Completed)
		meeting.State = constant.Completed
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
	})
}

func approveAddTask(meeting model.Meetings, userIds []int) {
	signupEndTime := time.Time(*meeting.SignupEndTime)
	startTime := time.Time(*meeting.MeetingStartTime)
	endTime := time.Time(*meeting.MeetingEndTime)

	var subS SubscriptionService
	startMessage := fmt.Sprintf(startTimeTemp, meeting.Title, meeting.MeetingLink)
	// 加入延迟队列
	delayQueue := delay.GetInstant()
	expireTime := signupEndTime

	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	delayQueue.Add(meeting.Id, expireTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.Preparing)
		meeting.State = constant.Preparing
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
		// 发送报名信息
		message := fmt.Sprintf(signupEndTimeTemp, meeting.Title, meeting.MeetingLink)
		subS.SendMsgByToIds(13, event.Meeting, constant.NOTICE, constant.MeetingId, userIds, message)
	})

	// 会议开始后状态改为会议中
	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	delayQueue.Add(meeting.Id, startTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.InMeeting)
		meeting.State = constant.InMeeting
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
		subS.SendMsgByToIds(13, event.Meeting, constant.NOTICE, constant.MeetingId, userIds, startMessage)
	})

	// 会议结束后状态改为会议完成
	log.Infof("延迟队列加入任务：%s", meeting.PrintLog())
	delayQueue.Add(meeting.Id, endTime, func() {
		log.Infof("会议id:%d,会议标题:%s,会议状态:%s,修改会议状态:%s", meeting.Id, meeting.Title, meeting.State, constant.Completed)
		meeting.State = constant.Completed
		model.Meeting().Where("id = ?", meeting.Id).Save(&meeting)
	})
}
