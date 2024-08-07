package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type QaAdoptions struct {
	ID        int            `gorm:"primarykey" json:"id"`
	ArticleId int            `json:"articleId"`
	CommentId int            `json:"commentId" binding:"required" msg:"采纳评论不能未空"`
	CreatedAt time.LocalTime `json:"createdAt"`
}

func QaAdoption() *gorm.DB {

	return mysql.GetInstance().Model(&QaAdoptions{})
}
