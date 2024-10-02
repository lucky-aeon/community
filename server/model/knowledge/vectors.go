package model

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/postgre"
)

type Vectors struct {
	ID         int             `gorm:"primarykey" json:"id"`
	CreatedAt  time.LocalTime  `json:"createdAt"`
	Content    string          `json:"content" binding:"required"`
	DocumentId int             `json:"documentId" binding:"required"`
	Embedding  pgvector.Vector `gorm:"type:float8[]" json:"embedding"` // 使用 PostgreSQL 数组
}

func Vector() *gorm.DB {
	return postgre.GetInstance().Model(&Vectors{})
}
