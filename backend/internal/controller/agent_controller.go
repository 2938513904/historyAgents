package controller

import (
	"fmt"
	"log"
	"net/http"

	"multiagent-chat/internal/model"
	"multiagent-chat/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AgentController 智能体控制器
type AgentController struct {
	agentService *service.AgentService
}

// NewAgentController 创建智能体控制器实例
func NewAgentController(db *gorm.DB) *AgentController {
	return &AgentController{
		agentService: service.NewAgentService(db),
	}
}

// GetAgents 获取所有智能体
func (ac *AgentController) GetAgents(c *gin.Context) {
	agents, err := ac.agentService.GetAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取智能体列表失败"})
		return
	}

	c.JSON(http.StatusOK, agents)
}

// CreateAgent 创建智能体
func (ac *AgentController) CreateAgent(c *gin.Context) {
	var agent model.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.agentService.CreateAgent(&agent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建智能体失败"})
		return
	}

	c.JSON(http.StatusCreated, agent)
}

// UpdateAgent 更新智能体
func (ac *AgentController) UpdateAgent(c *gin.Context) {
	id := c.Param("id")

	var agent model.Agent
	if err := c.ShouldBindJSON(&agent); err != nil {
		log.Printf("JSON绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("更新智能体请求 - ID: %s, 数据: %+v", id, agent)

	if err := ac.agentService.UpdateAgent(id, &agent); err != nil {
		log.Printf("更新智能体失败: %v", err)
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "智能体不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("更新智能体失败: %v", err)})
		return
	}

	// 获取更新后的智能体信息
	var updatedAgent model.Agent
	if err := ac.agentService.GetDB().First(&updatedAgent, "id = ?", id).Error; err != nil {
		log.Printf("获取更新后的智能体信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取更新后的智能体信息失败: %v", err)})
		return
	}

	log.Printf("智能体更新成功 - ID: %s", id)
	c.JSON(http.StatusOK, updatedAgent)
}

// DeleteAgent 删除智能体
func (ac *AgentController) DeleteAgent(c *gin.Context) {
	id := c.Param("id")

	// 检查ID是否为空
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "智能体ID不能为空"})
		return
	}

	if err := ac.agentService.DeleteAgent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除智能体失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "智能体删除成功"})
}

// DeleteEmptyAgents 删除空ID的智能体
func (ac *AgentController) DeleteEmptyAgents(c *gin.Context) {
	if err := ac.agentService.DeleteEmptyAgents(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除空ID智能体失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "空ID智能体删除成功"})
}
