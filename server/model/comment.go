package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ParentId   uint
	Content    string
	UserId     uint
	BusinessId uint
	TenantId   uint
}
