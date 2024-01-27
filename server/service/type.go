package services

import (
	"strings"
	"xhyovo.cn/community/server/model"
)

type TypeService struct {
}

func (s *TypeService) List(parentId uint) []model.Types {
	types := typeDao.List(parentId)
	for i := range types {
		typeObject := types[i]
		typeObject.ArticleStates = strings.Split(typeObject.ArticleState, ",")
	}

	return types
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
