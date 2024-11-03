package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type CoursesSections struct {
	ID          int    `gorm:"primarykey" json:"id"`
	Title       string `json:"title" binding:"required" msg:"标题不能未空"`
	Content     string `json:"content" binding:"required" msg:"内容不能未空"`
	UserId      int    `json:"userId,omitempty"`
	Sort        int    `json:"sort"`
	*UserSimple `json:"user" gorm:"-"`
	CourseId    int            `json:"courseId" binding:"required" msg:"对应课程不能未空"`
	CreatedAt   time.LocalTime `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index,omitempty"`
	PreId       int            `json:"preId" gorm:"-"`
	NextId      int            `json:"nextId" gorm:"-"`
	CourseTitle string         `json:"courseTitle" gorm:"column:courseTitle"`
}

func CoursesSection() *gorm.DB {
	return mysql.GetInstance().Model(&CoursesSections{})
}
