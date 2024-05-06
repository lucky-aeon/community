package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
)

type Dashboard struct {
	UserCount    int   `json:"userCount"`
	Profit       int   `json:"profit"`
	ArticleCount int64 `json:"articleCount"`
}

func InitDashboardRouters(r *gin.Engine) {
	group := r.Group("/community/admin/dashboard")
	group.GET("", dashboard)
}
func dashboard(ctx *gin.Context) {
	// 查出用户数量
	var codes []model.InviteCodes
	model.InviteCode().Where("state = 1").Select("id", "member_id").Find(&codes)
	var userCount = len(codes)
	// 查出文章数量
	var articleCount int64
	model.Article().Count(&articleCount)
	// 查出盈利
	var Profit int
	var memberInfos []model.MemberInfos
	model.MemberInfo().Select("id", "money").Find(&memberInfos)
	// 查出所有邀请码
	m := make(map[int]int)
	for i := range memberInfos {
		m[memberInfos[i].ID] = memberInfos[i].Money
	}
	for i := range codes {
		Profit += m[codes[i].MemberId]
	}
	d := Dashboard{
		UserCount:    userCount,
		ArticleCount: articleCount,
		Profit:       Profit,
	}
	result.Ok(d, "").Json(ctx)
}
