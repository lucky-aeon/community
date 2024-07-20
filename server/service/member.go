package services

import (
	"errors"
	"xhyovo.cn/community/server/model"
)

type MemberInfoService struct {
}

func (*MemberInfoService) ListMembers(page, limit int) (members []*model.MemberInfos, count int64) {
	db := model.MemberInfo()
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Find(&members)
	return
}

func (*MemberInfoService) SaveMember(member *model.MemberInfos) {
	memberDao.SaveMemberInfo(member)
}
func (*MemberInfoService) DeleteMember(id int) error {
	// 是否被邀请码使用
	var codeS CodeService
	count := codeS.CountByMemberId(id)
	if count > 0 {
		return errors.New("等级被引用,无法删除")
	}
	memberDao.DeleteMemberInfo(id)
	return nil
}

func (*MemberInfoService) Exist(id int) bool {
	return memberDao.Count(id) == 1
}

func (s *MemberInfoService) ListNameByIds(ids []int) map[int]string {

	infos := memberDao.ListByIdsSelectIdAndName(ids)

	m := make(map[int]string, len(infos))

	for i := range infos {
		memberInfos := infos[i]
		m[memberInfos.ID] = memberInfos.Name
	}

	return m
}

func (s *MemberInfoService) GetById(id int) (member model.MemberInfos) {
	model.MemberInfo().Where("id = ?", id).Find(&member)
	return

}
