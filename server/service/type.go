package services

import (
	"strings"

	"xhyovo.cn/community/server/model"
)

type TypeService struct {
}

func (s *TypeService) List(parentId int) []model.Types {
	types := typeDao.List(parentId)
	for i := range types {
		typeObject := types[i]
		typeObject.ArticleStates = strings.Split(typeObject.ArticleState, ",")
	}

	return types
}

func (s *TypeService) Save(types *model.Types) (int, error) {
	return typeDao.Save(types)
}

func (s *TypeService) Update(types *model.Types) error {
	return typeDao.Update(types)
}

func (s *TypeService) Delete(id int) error {
	return typeDao.Delete(id)
}

func (s *TypeService) GetById(id int) (types model.Types) {

	model.Type().Where("id = ?", id).Find(&types)
	return
}

func (s *TypeService) Exist(id int) bool {
	var c int64
	model.Type().Where("id = ?", id).Count(&c)
	return c == 1
}
