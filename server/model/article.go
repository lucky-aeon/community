package model

import (
	"gorm.io/gorm"
)

// it's issue or answer
type Article struct {
	gorm.Model
	Title       string
	Description string
	UserId      uint // The id of the author of this article
	State       int  // state code see details:
	Like        int
}
