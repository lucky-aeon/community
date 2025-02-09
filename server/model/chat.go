package model

import (
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
	FileURL   string         `json:"fileUrl"` // 文件URL,可为空
	CreatedAt time.LocalTime `json:"createdAt" gorm:"autoCreateTime"`
}

// AIModels AI模型表
type AIModels struct {
	ID          int64          `gorm:"primarykey" json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	SupportFile bool           `json:"supportFile" gorm:"default:false"`
	Status      bool           `json:"status" gorm:"default:true"`
	APIKey      string         `json:"-" gorm:"column:api_key"`  // API密钥，返回给前端时忽略
	BaseURL     string         `json:"-" gorm:"column:base_url"` // API基础URL，返回给前端时忽略
	CreatedAt   time.LocalTime `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.LocalTime `json:"updatedAt" gorm:"autoUpdateTime"`
}

// ContentPart 消息内容部分
type ContentPart struct {
	Type     string    `json:"type"`                // text 或 image_url
	Text     string    `json:"text,omitempty"`      // 当type为text时的文本内容
	ImageURL *ImageURL `json:"image_url,omitempty"` // 当type为image_url时的图片信息
}

// ImageURL 图片URL信息
type ImageURL struct {
	URL    string `json:"url"`    // 图片URL
	Detail string `json:"detail"` // 图片细节级别，如：high, low
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // 可以是string或[]ContentPart
}

// ChatCompletionRequest 聊天补全请求
type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
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
	ID                string             `json:"id"`
	Object            string             `json:"object"`
	Created           int64              `json:"created"`
	Model             string             `json:"model"`
	SystemFingerprint string             `json:"system_fingerprint"`
	Choices           []CompletionChoice `json:"choices"`
	Usage             CompletionUsage    `json:"usage"`
}

// CompletionChoice 补全选项
type CompletionChoice struct {
	Index        int               `json:"index"`
	Message      CompletionMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

// CompletionMessage 补全消息
type CompletionMessage struct {
	Role             string      `json:"role"`
	Content          string      `json:"content"`
	ReasoningContent string      `json:"reasoning_content,omitempty"`
	ToolCalls        interface{} `json:"tool_calls"`
	FunctionCall     interface{} `json:"function_call"`
}

// CompletionUsage token使用统计
type CompletionUsage struct {
	PromptTokens            int         `json:"prompt_tokens"`
	CompletionTokens        int         `json:"completion_tokens"`
	TotalTokens             int         `json:"total_tokens"`
	CompletionTokensDetails interface{} `json:"completion_tokens_details"`
	PromptTokensDetails     interface{} `json:"prompt_tokens_details"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	ModelID int64    `json:"modelId" binding:"required"` // AI模型ID
	Content string   `json:"content" binding:"required"` // 消息内容
	Files   []string `json:"files"`                      // 可选的文件数组
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
