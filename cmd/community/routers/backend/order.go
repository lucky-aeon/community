package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitOrderRouters(r *gin.Engine) {

	group := r.Group("/community/admin/order")
	group.GET("", listOrder)

}

func listOrder(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var orderSer = services.OrderServices{}
	orders, count := orderSer.Page(p, limit)
	result.Page(orders, count, nil).Json(ctx)
	return
}
