package model

import (
	"database/sql"
	"mime/multipart"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// ChatGroups 对话分组表
type ChatGroups struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	UserID    int64          `json:"userId" gorm:"index"`
	Title     string         `json:"title"`
	IsDeleted bool           `json:"isDeleted" gorm:"default:false"`
	CreatedAt time.LocalTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.LocalTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ChatMessages 对话消息表
type ChatMessages struct {
	ID        int64          `gorm:"primarykey" json:"id"`
	GroupID   int64          `json:"groupId" gorm:"index"`
	Role      string         `json:"role"` // user或assistant
	UserID    string         `json:"userId"`
	Content   string         `json:"content"`
	FileURL   sql.NullString `json:"fileUrl"` // 文件URL,可为空
	CreatedAt time.LocalTime `json:"createdAt" gorm:"autoCreateTime"`
}

// AIModels AI模型表
type AIModels struct {
	ID          int64          `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	SupportFile bool           `json:"supportFile" gorm:"default:false"`
	Status      bool           `json:"status" gorm:"default:true"`
	CreatedAt   time.LocalTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.LocalTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest 聊天补全请求
type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatCompletionChunk 流式聊天补全响应片段
type ChatCompletionChunk struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

// Choice 响应选项
type Choice struct {
	Index                int                  `json:"index"`
	Delta                Delta                `json:"delta"`
	FinishReason         *string              `json:"finish_reason"`
	ContentFilterResults ContentFilterResults `json:"content_filter_results"`
}

// Delta 增量内容
type Delta struct {
	Content          *string `json:"content"`
	ReasoningContent string  `json:"reasoning_content"`
	Role             string  `json:"role"`
}

// ContentFilterResults 内容过滤结果
type ContentFilterResults struct {
	Hate     FilterResult `json:"hate"`
	SelfHarm FilterResult `json:"self_harm"`
	Sexual   FilterResult `json:"sexual"`
	Violence FilterResult `json:"violence"`
}

// FilterResult 过滤结果
type FilterResult struct {
	Filtered bool `json:"filtered"`
}

// Usage token使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse 聊天补全响应
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ModelID int64                 `json:"modelId" binding:"required"` // AI模型ID
	Content string                `json:"content" binding:"required"` // 消息内容
	File    *multipart.FileHeader `json:"file"`                       // 可选的文件
}

// ChatGroupDB 获取对话分组表的数据库操作对象
func ChatGroupDB() *gorm.DB {
	return mysql.GetInstance().Model(&ChatGroups{})
}

// ChatMessageDB 获取对话消息表的数据库操作对象
func ChatMessageDB() *gorm.DB {
	return mysql.GetInstance().Model(&ChatMessages{})
}

// AIModelDB 获取AI模型表的数据库操作对象
func AIModelDB() *gorm.DB {
	return mysql.GetInstance().Model(&AIModels{})
}
