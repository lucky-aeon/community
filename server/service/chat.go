package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type ChatService struct {
}

var chatDao = &dao.Chat{}

// GetAIModels 获取所有可用的AI模型
func (s *ChatService) GetAIModels() ([]model.AIModels, error) {
	return chatDao.QueryAIModels()
}

// CreateChatGroup 创建对话分组
func (s *ChatService) CreateChatGroup(userID int64, title string) (*model.ChatGroups, error) {
	group := &model.ChatGroups{
		UserID: userID,
		Title:  title,
	}
	err := chatDao.CreateChatGroup(group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// GetChatGroup 获取对话分组详情
func (s *ChatService) GetChatGroup(id int64) (*model.ChatGroups, error) {
	return chatDao.GetChatGroup(id)
}

// GetUserChatGroups 获取用户的对话分组列表
func (s *ChatService) GetUserChatGroups(userID int64, page, limit int) ([]model.ChatGroups, int64, error) {
	groups, err := chatDao.GetUserChatGroups(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}
	total := chatDao.CountUserChatGroups(userID)
	return groups, total, nil
}

// UpdateChatGroup 更新对话分组
func (s *ChatService) UpdateChatGroup(id int64, userID int64, title string) error {
	group, err := chatDao.GetChatGroup(id)
	if err != nil {
		return err
	}
	if group.UserID != userID {
		return errors.New("无权限修改该对话分组")
	}
	group.Title = title
	return chatDao.UpdateChatGroup(group)
}

// DeleteChatGroup 删除对话分组
func (s *ChatService) DeleteChatGroup(id int64, userID int64) error {
	group, err := chatDao.GetChatGroup(id)
	if err != nil {
		return err
	}
	if group.UserID != userID {
		return errors.New("无权限删除该对话分组")
	}
	return chatDao.DeleteChatGroup(id)
}

// callAIModel 调用AI模型API
func (s *ChatService) callAIModel(aiModel *model.AIModels, messages []model.ChatMessage) (*model.ChatCompletionResponse, error) {
	reqBody := model.ChatCompletionRequest{
		Model:    aiModel.Name,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", aiModel.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", aiModel.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用AI模型失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI模型返回错误: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result model.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &result, nil
}

// SendMessage 发送消息并获取AI回复
func (s *ChatService) SendMessage(groupID int64, userID string, req *model.SendMessageRequest) (*model.ChatMessages, *model.ChatMessages, error) {
	// 验证分组是否存在
	if !chatDao.ExistsByID(groupID) {
		return nil, nil, errors.New("对话分组不存在")
	}

	// 获取AI模型信息
	aiModel, err := chatDao.GetAIModelByID(req.ModelID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取AI模型信息失败: %v", err)
	}

	// 保存用户消息
	userMessage := &model.ChatMessages{
		GroupID: groupID,
		Role:    "user",
		UserID:  userID,
		Content: req.Content,
		FileURL: strings.Join(req.Files, ","), // 多个文件URL用逗号分隔
	}
	if err := chatDao.CreateChatMessage(userMessage); err != nil {
		return nil, nil, fmt.Errorf("保存用户消息失败: %v", err)
	}

	// 获取历史消息作为上下文
	messages, err := chatDao.GetRecentMessages(groupID, math.MaxInt32)
	if err != nil {
		return nil, nil, fmt.Errorf("获取历史消息失败: %v", err)
	}

	// 构建AI请求消息
	var contextMessages []model.ChatMessage
	for _, msg := range messages {
		if aiModel.SupportFile && msg.FileURL != "" {
			// 如果AI模型支持文件且消息包含文件，使用多模态格式
			var contentParts []model.ContentPart
			fileURLs := strings.Split(msg.FileURL, ",")
			for _, fileURL := range fileURLs {
				if fileURL != "" {
					contentParts = append(contentParts, model.ContentPart{
						Type: "image_url",
						ImageURL: &model.ImageURL{
							URL:    fileURL,
							Detail: "high",
						},
					})
				}
			}
			contentParts = append(contentParts, model.ContentPart{
				Type: "text",
				Text: msg.Content,
			})
			contextMessages = append(contextMessages, model.ChatMessage{
				Role:    msg.Role,
				Content: contentParts,
			})
		} else {
			// 如果AI模型不支持文件或消息不包含文件，使用纯文本格式
			contextMessages = append(contextMessages, model.ChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// 添加当前用户消息
	var currentContent interface{}
	if aiModel.SupportFile && len(req.Files) > 0 {
		var contentParts []model.ContentPart
		for _, fileURL := range req.Files {
			contentParts = append(contentParts, model.ContentPart{
				Type: "image_url",
				ImageURL: &model.ImageURL{
					URL:    fileURL,
					Detail: "high",
				},
			})
		}
		contentParts = append(contentParts, model.ContentPart{
			Type: "text",
			Text: req.Content,
		})
		currentContent = contentParts
	} else {
		currentContent = req.Content
	}

	contextMessages = append(contextMessages, model.ChatMessage{
		Role:    "user",
		Content: currentContent,
	})

	// 调用AI模型
	aiResp, err := s.callAIModel(aiModel, contextMessages)
	if err != nil {
		return userMessage, nil, fmt.Errorf("调用AI模型失败: %v", err)
	}

	if len(aiResp.Choices) == 0 {
		return userMessage, nil, errors.New("AI模型未返回有效回复")
	}

	// 保存AI回复
	aiMessage := &model.ChatMessages{
		GroupID: groupID,
		Role:    "assistant",
		UserID:  strconv.FormatInt(req.ModelID, 10),
		Content: aiResp.Choices[0].Message.Content,
	}
	if err := chatDao.CreateChatMessage(aiMessage); err != nil {
		return userMessage, nil, fmt.Errorf("保存AI回复失败: %v", err)
	}

	return userMessage, aiMessage, nil
}

// GetMessages 获取对话消息列表
func (s *ChatService) GetMessages(groupID int64, page, pageSize int) ([]model.ChatMessages, int64, error) {
	// 验证分组是否存在
	if !chatDao.ExistsByID(groupID) {
		return nil, 0, errors.New("对话分组不存在")
	}

	messages, err := chatDao.GetChatMessages(groupID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	total := chatDao.GetTotalMessages(groupID)
	return messages, total, nil
}

// Auth 检查用户是否有权限操作该对话分组
func (s *ChatService) Auth(userID int64, groupID int64) bool {
	group, err := chatDao.GetChatGroup(groupID)
	if err != nil {
		return false
	}
	return group.UserID == userID
}
