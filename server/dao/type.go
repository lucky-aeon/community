package dao

import "xhyovo.cn/community/server/model"

type Type struct {
}

func (t *Type) Save(types *model.Types) uint {

	model.Type().Create(&types)
	return types.ID
}

func (t *Type) Update(types *model.Types) {
	model.Type().Model(types).Updates(types)
}

func (t *Type) Delete(id uint) {
	model.Type().Where("id = ?", id).Delete(&model.Types{})
}

func (t *Type) List() []model.Types {
	var types []model.Types
	model.Type().Find(&types)
	return types
}
