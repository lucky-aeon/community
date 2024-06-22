package request

type OssCallback struct {
	FileKey  string `json:"fileKey"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimeType"`
	UserId   int    `json:"x:userId"`
	Uuid     string `json:"x:uuid"`
}
