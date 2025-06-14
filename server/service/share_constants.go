package services

// 分享功能相关常量

// 允许分享的业务类型
const (
	BusinessTypeAiNews = "ai_news" // AI日报
	// BusinessTypeArticle = "article"  // 文章（暂未开放）
	// BusinessTypePost    = "post"     // 帖子（暂未开放）
)

// AllowedShareBusinessTypes 允许分享的业务类型列表
var AllowedShareBusinessTypes = map[string]bool{
	BusinessTypeAiNews: true,
	// 后续可以添加其他业务类型
}

// IsShareableBusinessType 检查业务类型是否允许分享
func IsShareableBusinessType(businessType string) bool {
	return AllowedShareBusinessTypes[businessType]
}

// GetAllowedBusinessTypes 获取所有允许分享的业务类型
func GetAllowedBusinessTypes() []string {
	var types []string
	for businessType := range AllowedShareBusinessTypes {
		types = append(types, businessType)
	}
	return types
}
