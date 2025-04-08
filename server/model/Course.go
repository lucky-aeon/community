package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Courses struct {
	ID            int               `gorm:"primarykey" json:"id"`
	Title         string            `json:"title" binding:"required" msg:"标题不能未空"`
	Desc          string            `json:"desc" binding:"required" msg:"描述不能未空"`
	Technology    string            `json:"technology"`
	TechnologyS   []string          `json:"technologys" gorm:"-"`
	Url           string            `json:"url"`
	CustomPageUrl string            `json:"customPageUrl" gorm:"default:''"` // 自定义页面URL，若存在则跳转到对应页面
	UserId        int               `json:"userId,omitempty"`
	Money         int               `json:"money"`
	Cover         string            `json:"cover"`
	Score         int               `json:"score"`
	State         int               `json:"state"`
	DemoUrl       string            `json:"demoUrl"` // 课程演示链接
	CreatedAt     time.LocalTime    `json:"createdAt"`
	UpdatedAt     time.LocalTime    `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt    `json:"deletedAt,omitempty" gorm:"index,omitempty"`
	Sections      []CoursesSections `json:"sections" gorm:"-;"`
	Views         int64             `json:"views" gorm:"-;"`
	Resources     interface{}       `json:"resources" gorm:"-"`        // 课程配套学习资源
	ResourcesJSON string            `json:"-" gorm:"column:resources"` // 资源JSON字符串
}

func Course() *gorm.DB {
	return mysql.GetInstance().Model(&Courses{})
}
