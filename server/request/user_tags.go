package request

type UserTags struct {
	UserId  int   `json:"userId"`
	TagsIds []int `json:"tagsIds"`
}
