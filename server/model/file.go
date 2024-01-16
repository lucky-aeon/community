package model

import "time"

type File struct {
	ID        uint `gorm:"primaryKey"`
	FileKey   string
	Size      int64
	Format    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
