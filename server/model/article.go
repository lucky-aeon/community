package model

import "time"

// it's issue or answer
type Article struct {
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Title       string
	Description string
	UserId      uint // The id of the author of this article
	State       int  // state code see details:
	Like        int
}
