package dao

import "xhyovo.cn/community/server/model"

type OrderDao struct {
}

func (*OrderDao) Save(order model.Orders) {
	model.Order().Save(&order)
}

func (*OrderDao) Page(page, limit int) ([]*model.Orders, int64) {
	var codes []*model.Orders
	tx := model.Order()
	var count int64
	tx.Count(&count)
	tx.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&codes)
	return codes, count
}
