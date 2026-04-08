package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client AI大模型调用客户端
type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	HTTP    *http.Client
}

// GenerateResult 返参结构
type GenerateResult struct {
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

// NewClient 创建大模型客户端
func NewClient(baseURL, apiKey, model string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// GenerateSummaryAndTags 生成摘要和标签
func (c *Client) GenerateSummaryAndTags(content string) (*GenerateResult, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("未配置大模型 API Key")
	}

	prompt := "你是一个专业的文章编辑助手。\n" +
		"请根据以下文章正文内容，生成一段简短精炼的摘要（不超过100字），并提取3-5个核心关键词标签。\n" +
		"请强制返回如下格式的纯JSON字符串（不要包含任何Markdown标记符号）：\n" +
		`{"summary": "你的摘要", "tags": ["标签1", "标签2", "标签3"]}` + "\n\n" +
		"正文内容：\n" + content

	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API调用错误状态码: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var respData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &respData); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v", err)
	}

	if len(respData.Choices) == 0 {
		return nil, fmt.Errorf("API没有返回生成的文本")
	}

	contentStr := respData.Choices[0].Message.Content
	contentStr = strings.TrimPrefix(contentStr, "```json\n")
	contentStr = strings.TrimSuffix(contentStr, "\n```")
	contentStr = strings.TrimPrefix(contentStr, "```\n")
	contentStr = strings.TrimSpace(contentStr)

	var result GenerateResult
	if err := json.Unmarshal([]byte(contentStr), &result); err != nil {
		return nil, fmt.Errorf("解析生成的 JSON 失败: %v, 原文: %s", err, contentStr)
	}

	return &result, nil
}
