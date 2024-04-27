package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Courses struct {
	ID          int            `gorm:"primarykey" json:"id"`
	Title       string         `json:"title" binding:"required" msg:"标题不能未空"`
	Desc        string         `json:"desc" binding:"required" msg:"描述不能未空"`
	Technology  string         `json:"technology"`
	TechnologyS []string       `json:"technologys" gorm:"-"`
	Url         string         `json:"url"`
	UserId      int            `json:"userId,omitempty"`
	Money       int            `json:"money"`
	Cover       string         `json:"cover"`
	State       int            `json:"state"`
	CreatedAt   time.LocalTime `json:"createdAt"`
	UpdatedAt   time.LocalTime `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index,omitempty"`
}

func Course() *gorm.DB {
	return mysql.GetInstance().Model(&Courses{})
}
