package request

type OssCallback struct {
	FileKey  string `json:"fileKey"`
	Size     int64  `json:"size"`
	MineType string `json:"mineType"`
	UserId   int    `json:"x:userId"`
	Uuid     string `json:"x:uuid"`
}
