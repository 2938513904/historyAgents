package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"multiagent-chat/internal/config"
)

// CallDashScopeAPI 调用阿里云DashScope API函数
func CallDashScopeAPI(prompt string) (string, error) {
	// 构造请求体
	requestBody := map[string]interface{}{
		"model": config.HardcodedModel,
		"input": map[string]interface{}{
			"prompt": prompt,
		},
		"parameters": map[string]interface{}{
			"max_tokens":  1500,
			"temperature": 0.8,
		},
	}

	// 将请求体转换为JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构造请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", config.DashScopeEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+config.HardcodedAPIKey)
	req.Header.Set("Content-Type", "application/json")
	// 移除SSE模式，使用标准的JSON响应
	// req.Header.Set("X-DashScope-SSE", "enable")

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 记录响应内容用于调试
	log.Printf("阿里云API响应状态码: %d", resp.StatusCode)
	log.Printf("阿里云API响应内容: %s", string(body))

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API调用失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 检查响应内容是否为空
	if len(body) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	// 检查响应是否为JSON格式（简单检查）
	if len(body) > 0 && body[0] != '{' && body[0] != '[' {
		return "", fmt.Errorf("API返回非JSON格式响应，首字符: %c, 响应内容: %s", body[0], string(body))
	}

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	// 提取生成的文本
	if output, ok := response["output"].(map[string]interface{}); ok {
		if text, ok := output["text"].(string); ok {
			return text, nil
		}
	}

	return "", fmt.Errorf("无法从响应中提取生成的文本")
}
