package dao

import (
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

type Type struct {
}

func (t *Type) Save(types *model.Types) (int, error) {
	d := mysql.GetInstance().Save(&types)
	return types.ID, d.Error
}

func (t *Type) Update(types *model.Types) error {
	return model.Type().Model(types).Save(types).Error
}

func (t *Type) Delete(id int) error {
	d := model.Type().Where("id = ?", id).Delete(&model.Types{})
	if d.Error != nil {
		return d.Error
	}
	d = model.Type().Where("parent_id = ?", id).Delete(&model.Types{})
	return d.Error
}

func (t *Type) List(parentId int) []model.Types {
	var types []model.Types
	model.Type().Where("parent_id = ?", parentId).Find(&types)

	return types
}

func (t *Type) GetById(typeId int) model.Types {
	var typeObject model.Types
	model.Type().Where("id = ?", typeId).Find(&typeObject)
	return typeObject

}
