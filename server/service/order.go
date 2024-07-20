package services

import (
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type OrderServices struct {
}

func (*OrderServices) Page(page, limit int) ([]*model.Orders, int64) {

	var orderDao = dao.OrderDao{}
	orders, count := orderDao.Page(page, limit)
	// 收集创建者id

	var userIds = mapset.NewSet[int]()
	for _, order := range orders {
		userIds.Add(order.Purchaser)
		userIds.Add(order.Creator)
	}

	var userService = UserService{}
	nameMap := userService.ListByIdsToMap(userIds.ToSlice())
	for _, order := range orders {
		order.PurchaserName = nameMap[order.Purchaser].Name
		order.CreatorName = nameMap[order.Creator].Name
	}

	return orders, count
}

func (*OrderServices) CalculateProfit() int64 {
	// 我现在有一个账单表，需要查出盈利多少，函数名如何取
	var totalProfit int64
	model.Order().Where("acquisition_type = ?", 2).Select("SUM(price)").Scan(&totalProfit)
	return totalProfit
}
