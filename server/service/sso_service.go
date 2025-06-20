package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/oss"
	"xhyovo.cn/community/server/model"
)

type SsoService struct {
}

// GetApplicationByKey 根据app_key获取应用信息
func (s *SsoService) GetApplicationByKey(appKey string) (*model.SsoApplication, error) {
	var app model.SsoApplication
	err := model.SsoApp().Where("app_key = ? AND status = 1", appKey).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}
	return &app, nil
}

// ValidateRedirectUrl 验证回调地址是否在白名单中
func (s *SsoService) ValidateRedirectUrl(app *model.SsoApplication, redirectUrl string) bool {
	if app.RedirectUrls == "" {
		return false
	}

	allowedUrls := strings.Split(app.RedirectUrls, ",")
	for _, allowedUrl := range allowedUrls {
		allowedUrl = strings.TrimSpace(allowedUrl)
		if strings.HasPrefix(redirectUrl, allowedUrl) {
			return true
		}
	}
	return false
}

// GenerateAuthCode 生成授权码
func (s *SsoService) GenerateAuthCode(appKey string, userId int, redirectUrl string) (string, error) {
	// 生成随机授权码
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	code := hex.EncodeToString(bytes)

	// 保存到数据库，5分钟过期
	authCode := model.SsoAuthCode{
		Code:        code,
		AppKey:      appKey,
		UserId:      userId,
		RedirectUrl: redirectUrl,
		Used:        false,
		ExpiresAt:   time.Now().Add(5 * time.Minute),
	}

	if err := model.SsoAuthCodeModel().Create(&authCode).Error; err != nil {
		return "", err
	}

	return code, nil
}

// ValidateAuthCode 验证授权码并获取用户信息
func (s *SsoService) ValidateAuthCode(appKey, appSecret, code string) (*model.Users, error) {
	// 验证应用
	app, err := s.GetApplicationByKey(appKey)
	if err != nil {
		return nil, err
	}

	if app.AppSecret != appSecret {
		return nil, errors.New("应用密钥错误")
	}

	// 查找授权码
	var authCode model.SsoAuthCode
	err = model.SsoAuthCodeModel().Where("code = ? AND app_key = ? AND used = false", code, appKey).First(&authCode).Error
	if err != nil {
		return nil, errors.New("授权码无效")
	}

	// 检查是否过期
	if time.Now().After(authCode.ExpiresAt) {
		return nil, errors.New("授权码已过期")
	}

	// 标记为已使用
	model.SsoAuthCodeModel().Where("id = ?", authCode.ID).Update("used", true)

	// 获取用户信息
	var user model.Users
	err = model.User().Where("id = ?", authCode.UserId).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// GetUserBasicInfo 获取用户基本信息（脱敏）
func (s *SsoService) GetUserBasicInfo(user *model.Users) map[string]interface{} {
	// 处理头像URL
	avatarUrl := ""
	if user.Avatar != "" {
		// 生成完整的头像URL
		avatarUrl = s.generateAvatarUrl(user.Avatar)
	}
	
	return map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Account,
		"avatar": avatarUrl,
		"desc":   user.Desc,
	}
}

// generateAvatarUrl 生成完整的头像URL
func (s *SsoService) generateAvatarUrl(avatar string) string {
	if avatar == "" {
		return ""
	}
	
	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(avatar, "http://") || strings.HasPrefix(avatar, "https://") {
		return avatar
	}
	
	// 获取OSS配置
	ossConfig := config.GetInstance().OssConfig
	if ossConfig.Cdn != "" {
		return ossConfig.Cdn + "/" + avatar
	}
	
	// 如果没有CDN配置，使用OSS endpoint
	return "https://" + oss.GetInstance().BucketName + "." + oss.GetEndpoint() + "/" + avatar
}
