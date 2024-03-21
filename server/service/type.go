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

func (s *TypeService) PageTypes(page, limit int) (types []model.Types, count int64) {
	model.Type().Limit(limit).Offset((page - 1) * limit).Find(&types)
	if len(types) == 0 {
		return types, count
	}
	parentIds := make([]model.Types, 0, len(types))
	for i := range types {
		types[i].ArticleStates = strings.Split(types[i].ArticleState, ",")
		if types[i].ParentId == 0 {
			types[i].Children = make([]model.Types, 0)
			parentIds = append(parentIds, types[i])
		}
	}

	// 根据根分类找子分类
	for i := range types {
		for i2 := range parentIds {
			if types[i].ParentId == parentIds[i2].ID {
				parentIds[i2].Children = append(parentIds[i2].Children, types[i])
			}
		}
	}
	model.Type().Where("parent_id = 0").Count(&count)
	return parentIds, count
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

func (s *TypeService) ListParentTypes() []model.Types {
	var types []model.Types
	model.Type().Where("parent_id = 0").Find(&types)
	return types

}
