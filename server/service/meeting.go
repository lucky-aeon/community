package services

import (
	"errors"
	"time"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/server/model"
)

type MeetingService struct {
}

// 保存
func (*MeetingService) Save(meeting *model.Meetings) error {

	// 会议申请时间不低于当前时间
	if !time.Time(meeting.InitiatorTime).Before(time.Now()) {
		return errors.New("会议申请时间不可小于当前时间")
	}

	// 会议状态只一旦非审核中则不可修改
	if meeting.State != constant.Reviewing {
		return errors.New("会议已被锁定,不可修改")
	}
	if err := model.Meeting().Save(&meeting).Error; err != nil {
		return err
	}

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
	} else {
		meeting.JoinUserCount = len(joinUsers)
	}

	return meeting
}

func (*MeetingService) Page(page, limit int) []model.Meetings {
	var meetings []model.Meetings
	model.Meeting().Limit(limit).Offset((page - 1) * limit).Find(&meetings)

	// 收集用户id
	var userIds []int
	for _, meeting := range meetings {
		userIds = append(userIds, meeting.InitiatorId)
	}

	var userS UserService
	nameMap := userS.ListByIdsToMap(userIds)
	for i := range meetings {
		meetings[i].InitiatorName = nameMap[meetings[i].InitiatorId].Name
	}

	return meetings
}
