package services

import (
	"xhyovo.cn/community/server/model"
)

type TypeService struct {
}

func (s *TypeService) List(parentId uint) []model.Types {
	return typeDao.List(parentId)
}

func (s *TypeService) Save(types *model.Types) uint {
	return typeDao.Save(types)
}

func (s *TypeService) Update(types *model.Types) {
	typeDao.Update(types)
}

func (s *TypeService) Delete(id uint) {
	typeDao.Delete(id)
}
