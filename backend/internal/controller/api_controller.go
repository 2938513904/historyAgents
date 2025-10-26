package controller

import (
	"net/http"

	"multiagent-chat/internal/config"
	"multiagent-chat/internal/utils"

	"github.com/gin-gonic/gin"
)

// APIController API测试控制器
type APIController struct{}

// NewAPIController 创建API测试控制器实例
func NewAPIController() *APIController {
	return &APIController{}
}

// TestAPIConnection 测试API连接
func (ac *APIController) TestAPIConnection(c *gin.Context) {
	// 使用真实的阿里云DashScope API调用进行测试
	// log.Printf("API连接测试使用真实阿里云API，配置信息 - 端点: %s, 模型: %s", config.HardcodedBaseURL, config.HardcodedModel)

	// 构造测试提示词
	testPrompt := "Hello, this is a test message. Please respond with a short greeting."

	// 调用阿里云API
	response, err := utils.CallDashScopeAPI(testPrompt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "API连接失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "API连接测试成功",
		"model":    config.HardcodedModel,
		"provider": "DashScope (阿里云)",
		"response": response,
	})
}

// TestDashScopeAPI 测试阿里云DashScope API
func (ac *APIController) TestDashScopeAPI(c *gin.Context) {
	var requestData struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误: " + err.Error(),
		})
		return
	}

	// 保存原始的模型配置
	originalModel := config.HardcodedModel
	if requestData.Model != "" {
		config.HardcodedModel = requestData.Model
	}

	// 调用阿里云API
	// log.Printf("测试阿里云DashScope API，提示词: %s, 模型: %s", requestData.Prompt, config.HardcodedModel)
	response, err := utils.CallDashScopeAPI(requestData.Prompt)

	// 恢复原始模型配置
	config.HardcodedModel = originalModel

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success":  false,
			"message":  "API调用失败: " + err.Error(),
			"response": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "API调用成功",
		"response": response,
	})
}
