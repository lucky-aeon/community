package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"xhyovo.cn/community/pkg/config"
)

type EmbeddingRequest struct {
	Model string `json:"model"`
	Input struct {
		Texts []string `json:"texts"`
	} `json:"input"`
	Parameters struct {
		TextType string `json:"text_type"`
	} `json:"parameters"`
}

type EmbeddingResponseItem struct {
	TextIndex int       `json:"text_index"`
	Embedding []float32 `json:"embedding"`
}

type Output struct {
	Embeddings []EmbeddingResponseItem `json:"embeddings"`
}

type EmbeddingResponse struct {
	Output Output `json:"output"`
}

type EmbeddingService struct{}

// cleanText 清理 JSON 文本中的非法字符
func cleanText(input string) string {
	// 1. 替换反引号 ` 为单引号 ' 或者直接移除
	input = strings.ReplaceAll(input, "`", "'")

	// 2. 转义双引号和反斜杠
	input = strings.ReplaceAll(input, `\`, `\\`) // 转义反斜杠
	input = strings.ReplaceAll(input, `"`, `\"`) // 转义双引号

	// 3. 处理方括号，避免其影响 JSON 结构
	input = strings.ReplaceAll(input, `[`, `\[`) // 转义左方括号
	input = strings.ReplaceAll(input, `]`, `\]`) // 转义右方括号

	// 4. 删除控制字符（ASCII 范围 0x00 - 0x1F）
	controlChars := regexp.MustCompile(`[\x00-\x1F]`)
	input = controlChars.ReplaceAllString(input, "")

	// 5. 处理其他不合法的字符（可根据需求自定义更多清理规则）

	return input
}

// cleanTexts 批量清理字符串数组
func cleanTexts(texts []string) []string {
	cleanedTexts := make([]string, len(texts))
	for i, text := range texts {
		cleanedTexts[i] = cleanText(text)
	}
	return cleanedTexts
}

func (e *EmbeddingService) GetTextEmbeddings(texts []string) (EmbeddingResponse, error) {
	embeddingConfig := config.GetInstance().EmbeddingConfig

	url := embeddingConfig.Url
	requestBody := EmbeddingRequest{
		Model: embeddingConfig.Model,
		Parameters: struct {
			TextType string `json:"text_type"`
		}{TextType: "query"},
	}
	requestBody.Input.Texts = texts

	body, err := json.Marshal(requestBody)

	var embeddingResponse EmbeddingResponse
	if err != nil {
		return embeddingResponse, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return embeddingResponse, err
	}

	req.Header.Set("Authorization", "Bearer "+embeddingConfig.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	requestDump, err := httputil.DumpRequestOut(req, true) // true 表示包括请求体
	if err != nil {
		return embeddingResponse, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return embeddingResponse, err
	}
	defer resp.Body.Close()

	fmt.Printf("HTTP Request:\n%s\n", string(requestDump)) // 打印到控制台

	if resp.StatusCode != http.StatusOK {
		return embeddingResponse, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&embeddingResponse); err != nil {
		return embeddingResponse, err
	}

	return embeddingResponse, nil
}
