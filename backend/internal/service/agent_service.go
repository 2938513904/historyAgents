package service

import (
	"multiagent-chat/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentService 智能体服务
type AgentService struct {
	db *gorm.DB
}

// NewAgentService 创建智能体服务实例
func NewAgentService(db *gorm.DB) *AgentService {
	return &AgentService{
		db: db,
	}
}

// GetDB 获取数据库实例
func (s *AgentService) GetDB() *gorm.DB {
	return s.db
}

// GetAgents 获取所有智能体
func (s *AgentService) GetAgents() ([]model.Agent, error) {
	var agents []model.Agent
	if err := s.db.Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// CreateAgent 创建智能体
func (s *AgentService) CreateAgent(agent *model.Agent) error {
	agent.ID = uuid.New().String()
	// agent.CreatedAt = time.Now() // 由数据库自动设置
	if err := s.db.Create(agent).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAgent 更新智能体
func (s *AgentService) UpdateAgent(id string, agent *model.Agent) error {
	var existingAgent model.Agent
	if err := s.db.First(&existingAgent, "id = ?", id).Error; err != nil {
		return err
	}

	existingAgent.Name = agent.Name
	existingAgent.Role = agent.Role
	existingAgent.Personality = agent.Personality
	existingAgent.Avatar = agent.Avatar
	existingAgent.CozeLink = agent.CozeLink

	if err := s.db.Save(&existingAgent).Error; err != nil {
		return err
	}
	return nil
}

// DeleteAgent 删除智能体
func (s *AgentService) DeleteAgent(id string) error {
	if err := s.db.Delete(&model.Agent{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// DeleteEmptyAgents 删除空ID的智能体
func (s *AgentService) DeleteEmptyAgents() error {
	if err := s.db.Where("id = ?", "").Delete(&model.Agent{}).Error; err != nil {
		return err
	}
	return nil
}
