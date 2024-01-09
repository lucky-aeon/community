package model

import "gorm.io/gorm"

// it's issue or answer
type Article struct {
	gorm.Model
	Title       string
	Description string
	UserId      uint // The id of the author of this article
	IssueId     uint // If issueId is not zero, it is an answer or it is a question
	Solved      bool // The issue status is solved.
	Like        int
}
