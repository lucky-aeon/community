package data

type ListSortStrategy struct {
	OrderBy   string `json:"orderBy"`
	DescOrder bool   `json:"descOrder"` // true=desc, false=asc
}

type QueryPage struct {
	Page  int
	Limit int
}
