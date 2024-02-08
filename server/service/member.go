package services

import "xhyovo.cn/community/server/model"

type MemberInfoService struct {
}

func (*MemberInfoService) ListMember() []*model.MemberInfos {
	return memberDao.ListMemberInfo()
}

func (*MemberInfoService) SaveMember(member *model.MemberInfos) {
	memberDao.SaveMemberInfo(member)
}
func (*MemberInfoService) DeleteMember(id int) {
	memberDao.DeleteMemberInfo(id)
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
