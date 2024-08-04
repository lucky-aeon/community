package services

import (
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

type CodeService struct {
}

func (*CodeService) PageCodes(page, limit int, code string) (codes []*model.InviteCodes, count int64) {
	codes = codeDao.PageCodes(page, limit, code)
	count = codeDao.GetCount()

	if count == 0 {
		return codes, count
	}

	// set member info
	memberIds := mapset.NewSet[int](len(codes))
	for i := range codes {
		memberIds.Add(codes[i].MemberId)
	}

	var m MemberInfoService
	idNameMap := m.ListNameByIds(memberIds.ToSlice())

	for i := range codes {
		codes[i].MemberName = idNameMap[codes[i].MemberId]
	}

	return codes, count
}
func (*CodeService) GenerateCode(m model.GenerateCode) error {
	memberId := m.MemberId
	var mSer MemberInfoService
	if !mSer.Exist(memberId) {
		return errors.New("对应vip等级不存在")
	}
	number := m.Number
	var flag bool = true
	var code string
	// generate code
	codes := make([]string, 0)
	for ; number > 0; number-- {
		for flag {
			code = utils.GenerateCode(8)
			flag = codeDao.Exist(code)
		}
		codes = append(codes, code)
		flag = true
	}

	codeList := []*model.InviteCodes{}
	for i := range codes {
		c := &model.InviteCodes{
			Code:            codes[i],
			MemberId:        m.MemberId,
			Creator:         m.Creator,
			AcquisitionType: m.AcquisitionType,
		}
		codeList = append(codeList, c)
	}

	codeDao.SaveCodes(codeList)
	return nil
}

func (*CodeService) DestroyCode(code int) error {
	if codeDao.Del(code) == 0 {
		return errors.New("邀请码被使用，无法删除")
	}
	return nil
}

func (*CodeService) SetState(code string) {
	codeDao.SetState(code)
}

func (s *CodeService) CountByMemberId(id int) (count int64) {
	model.InviteCode().Where("member_id = ?", id).Count(&count)
	return count
}

func (s *CodeService) GetByCode(code string) (codeObject model.InviteCodes) {
	model.InviteCode().Where("code = ?", code).Find(&codeObject)
	return
}
