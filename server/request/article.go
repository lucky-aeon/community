package request

type ReqArticle struct {
	ID      int    `gorm:"primarykey" json:"id"`
	Title   string `json:"title" binding:"required" msg:"标题不能未空"`
	Content string `json:"content,omitempty" binding:"required" msg:"描述不能未空"`
	UserId  int    `json:"userId,omitempty"`
	State   int    `json:"state"` // 状态:草稿/发布/待解决/已解决/私密提问
	Type    int    `json:"type"`
	Tags    []int  `json:"tags" gorm:"-"`
}

type TopArticle struct {
	Id        int `json:"id"`
	State     int `json:"state"`
	TopNumber int `json:"topNumber"`
}
