package dao

import "xhyovo.cn/community/server/model"

type Type struct {
}

func (t *Type) Save(types *model.Types) (int, error) {
	d := model.Type().Create(&types)
	return types.ID, d.Error
}

func (t *Type) Update(types *model.Types) error {
	return model.Type().Model(types).Updates(types).Error
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
