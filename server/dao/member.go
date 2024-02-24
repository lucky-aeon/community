package dao

import "xhyovo.cn/community/server/model"

type MemberDao struct {
}

func (*MemberDao) ListMemberInfo() []*model.MemberInfos {
	var members []*model.MemberInfos
	model.MemberInfo().Find(&members)
	return members
}

// save or updated
func (*MemberDao) SaveMemberInfo(memberInfo *model.MemberInfos) {
	if memberInfo.ID == 0 {
		model.MemberInfo().Save(&memberInfo)
	} else {
		model.MemberInfo().Where("id = ?", memberInfo.ID).Updates(&memberInfo)
	}

}

func (*MemberDao) DeleteMemberInfo(id int) {
	model.MemberInfo().Where("id = ?", id).Delete(&model.MemberInfos{})
}

func (*MemberDao) Count(id int) int64 {
	var count int64
	model.MemberInfo().Where("id = ?", id).Count(&count)
	return count
}

func (d *MemberDao) ListByIdsSelectIdAndName(ids []int) []*model.MemberInfos {

	var m []*model.MemberInfos
	model.MemberInfo().Where("id in ?", ids).Find(&m)
	return m
}
