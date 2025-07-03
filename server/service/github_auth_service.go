package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/config"
	localTime "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

type GitHubAuthService struct {
}

// GetGitHubOAuthURL 生成 GitHub OAuth 授权 URL
func (s *GitHubAuthService) GetGitHubOAuthURL() string {
	githubConfig := config.GetInstance().GitHubConfig
	
	// 生成随机state参数用于防止CSRF攻击
	state := s.generateState()
	
	// 将state存储到缓存中，5分钟过期
	cache.GetInstance().Set("github_oauth_state_"+state, true, 5*time.Minute)
	
	params := url.Values{}
	params.Add("client_id", githubConfig.ClientID)
	params.Add("redirect_uri", githubConfig.RedirectURL)
	params.Add("scope", "user:email")
	params.Add("state", state)
	
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

// HandleOAuthCallback 处理 GitHub OAuth 回调
func (s *GitHubAuthService) HandleOAuthCallback(code, state string) (*model.GitHubUser, error) {
	// 验证state参数
	if !s.validateState(state) {
		return nil, errors.New("无效的state参数")
	}
	
	// 使用授权码获取访问令牌
	accessToken, err := s.getAccessToken(code)
	if err != nil {
		return nil, fmt.Errorf("获取访问令牌失败: %v", err)
	}
	
	// 使用访问令牌获取用户信息
	githubUser, err := s.getGitHubUserInfo(accessToken)
	if err != nil {
		return nil, fmt.Errorf("获取GitHub用户信息失败: %v", err)
	}
	
	return githubUser, nil
}

// FindUserByGitHubID 根据GitHub ID查找已绑定的用户
func (s *GitHubAuthService) FindUserByGitHubID(githubID int64) (*model.Users, error) {
	var binding model.UserGitHubBinding
	err := model.UserGitHubBindingModel().Where("github_id = ?", githubID).First(&binding).Error
	if err != nil {
		return nil, err
	}
	
	var user model.Users
	err = model.User().Where("id = ?", binding.UserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// BindGitHubAccount 绑定 GitHub 账号到用户
func (s *GitHubAuthService) BindGitHubAccount(userID int, githubUser *model.GitHubUser) error {
	// 检查该GitHub账号是否已经绑定其他用户
	var existingBinding model.UserGitHubBinding
	err := model.UserGitHubBindingModel().Where("github_id = ?", githubUser.ID).First(&existingBinding).Error
	if err == nil {
		return errors.New("该GitHub账号已绑定其他用户")
	}
	
	// 检查该用户是否已经绑定了GitHub账号
	err = model.UserGitHubBindingModel().Where("user_id = ?", userID).First(&existingBinding).Error
	if err == nil {
		return errors.New("该用户已绑定GitHub账号")
	}
	
	// 创建新的绑定关系
	binding := model.UserGitHubBinding{
		UserID:         userID,
		GitHubID:       githubUser.ID,
		GitHubUsername: githubUser.Login,
		GitHubEmail:    githubUser.Email,
		GitHubAvatar:   githubUser.AvatarURL,
		BoundAt:        localTime.LocalTime(time.Now()),
	}
	
	return model.UserGitHubBindingModel().Create(&binding).Error
}

// GetGitHubBinding 获取用户的GitHub绑定信息
func (s *GitHubAuthService) GetGitHubBinding(userID int) (*model.UserGitHubBinding, error) {
	var binding model.UserGitHubBinding
	err := model.UserGitHubBindingModel().Where("user_id = ?", userID).First(&binding).Error
	return &binding, err
}

// UnbindGitHubAccount 解绑GitHub账号
func (s *GitHubAuthService) UnbindGitHubAccount(userID int) error {
	return model.UserGitHubBindingModel().Where("user_id = ?", userID).Delete(&model.UserGitHubBinding{}).Error
}

// generateState 生成随机state参数
func (s *GitHubAuthService) generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// validateState 验证state参数
func (s *GitHubAuthService) validateState(state string) bool {
	key := "github_oauth_state_" + state
	_, exists := cache.GetInstance().Get(key)
	if exists {
		// 验证成功后删除state，防止重复使用
		cache.GetInstance().Delete(key)
		return true
	}
	return false
}

// getAccessToken 使用授权码获取访问令牌（带重试机制）
func (s *GitHubAuthService) getAccessToken(code string) (string, error) {
	githubConfig := config.GetInstance().GitHubConfig
	
	data := url.Values{}
	data.Set("client_id", githubConfig.ClientID)
	data.Set("client_secret", githubConfig.ClientSecret)
	data.Set("code", code)
	
	// 重试3次
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
		if err != nil {
			return "", err
		}
		
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if attempt < maxRetries {
				// 等待后重试
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", fmt.Errorf("获取访问令牌失败，已重试%d次: %v", maxRetries, err)
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", err
		}
		
		var tokenResp model.GitHubOAuthResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", err
		}
		
		if tokenResp.AccessToken == "" {
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			return "", errors.New("获取访问令牌失败：响应中无访问令牌")
		}
		
		return tokenResp.AccessToken, nil
	}
	
	return "", errors.New("获取访问令牌失败：超过最大重试次数")
}

// getGitHubUserInfo 使用访问令牌获取GitHub用户信息
func (s *GitHubAuthService) getGitHubUserInfo(accessToken string) (*model.GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var githubUser model.GitHubUser
	if err := json.Unmarshal(body, &githubUser); err != nil {
		return nil, err
	}
	
	// 如果邮箱为空，尝试获取主要邮箱
	if githubUser.Email == "" {
		githubUser.Email, _ = s.getGitHubUserEmail(accessToken)
	}
	
	return &githubUser, nil
}

// getGitHubUserEmail 获取GitHub用户主要邮箱
func (s *GitHubAuthService) getGitHubUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}
	
	// 查找主要邮箱
	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}
	
	// 如果没有主要邮箱，返回第一个邮箱
	if len(emails) > 0 {
		return emails[0].Email, nil
	}
	
	return "", errors.New("未找到邮箱")
}