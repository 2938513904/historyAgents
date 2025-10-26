package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"multiagent-chat/internal/model"
	"multiagent-chat/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatRoomController 聊天室控制器
type ChatRoomController struct {
	db *gorm.DB
}

// NewChatRoomController 创建聊天室控制器实例
func NewChatRoomController(db *gorm.DB) *ChatRoomController {
	return &ChatRoomController{
		db: db,
	}
}

// GetChatRooms 获取所有聊天室
func (cc *ChatRoomController) GetChatRooms(c *gin.Context) {
	var chatRooms []model.ChatRoom
	if err := cc.db.Preload("Agents").Find(&chatRooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取聊天室列表失败"})
		return
	}
	c.JSON(http.StatusOK, chatRooms)
}

// CreateChatRoom 创建聊天室
func (cc *ChatRoomController) CreateChatRoom(c *gin.Context) {
	var requestData struct {
		Topic  string   `json:"topic"`
		Agents []string `json:"agents"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 添加日志查看接收到的数据
	// log.Printf("接收到创建聊天室请求 - Topic: %s, Agents: %v", requestData.Topic, requestData.Agents)

	// 创建聊天室
	chatRoom := model.ChatRoom{
		ID:        uuid.New().String(),
		Topic:     requestData.Topic,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// 保存聊天室到数据库
	if err := cc.db.Create(&chatRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建聊天室失败"})
		return
	}

	// 关联智能体
	if len(requestData.Agents) > 0 {
		var agents []model.Agent
		if err := cc.db.Where("id IN ?", requestData.Agents).Find(&agents).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询智能体失败"})
			return
		}

		// 使用Association方法关联智能体
		if err := cc.db.Model(&chatRoom).Association("Agents").Append(agents); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联智能体失败"})
			return
		}
	}

	c.JSON(http.StatusCreated, chatRoom)
}

// GetChatRoom 获取单个聊天室
func (cc *ChatRoomController) GetChatRoom(c *gin.Context) {
	id := c.Param("id")

	var chatRoom model.ChatRoom
	if err := cc.db.Preload("Agents").Preload("Messages").First(&chatRoom, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	c.JSON(http.StatusOK, chatRoom)
}

// StartChatRoom 启动聊天室讨论
func (cc *ChatRoomController) StartChatRoom(c *gin.Context) {
	id := c.Param("id")

	// 获取聊天室
	var chatRoom model.ChatRoom
	if err := cc.db.First(&chatRoom, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	if chatRoom.Status == "running" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "讨论正在进行中"})
		return
	}

	// 更新状态为运行中
	chatRoom.Status = "running"
	if err := cc.db.Save(&chatRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新聊天室状态失败"})
		return
	}

	// 异步开始讨论
	go cc.simulateAgentDiscussion(id)

	c.JSON(http.StatusOK, gin.H{"message": "讨论已开始"})
}

// StopChatRoom 停止聊天室讨论
func (cc *ChatRoomController) StopChatRoom(c *gin.Context) {
	id := c.Param("id")

	var chatRoom model.ChatRoom
	if err := cc.db.First(&chatRoom, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	chatRoom.Status = "stopped"
	if err := cc.db.Save(&chatRoom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新聊天室状态失败"})
		return
	}

	c.JSON(http.StatusOK, chatRoom)
}

// DeleteChatRoom 删除聊天室
func (cc *ChatRoomController) DeleteChatRoom(c *gin.Context) {
	id := c.Param("id")

	// 检查ID是否为空
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "聊天室ID不能为空"})
		return
	}

	// 首先删除聊天室关联的消息
	if err := cc.db.Where("chat_room_id = ?", id).Delete(&model.Message{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除聊天室消息失败"})
		return
	}

	// 删除聊天室与智能体的关联关系
	if err := cc.db.Exec("DELETE FROM chat_room_agents WHERE chat_room_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除聊天室智能体关联失败"})
		return
	}

	// 然后删除聊天室
	if err := cc.db.Delete(&model.ChatRoom{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "聊天室删除成功"})
}

// DeleteEmptyChatRooms 删除空ID的聊天室
func (cc *ChatRoomController) DeleteEmptyChatRooms(c *gin.Context) {
	// 删除所有ID为空的聊天室
	if err := cc.db.Where("id = ?", "").Delete(&model.ChatRoom{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除空ID聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "空ID聊天室删除成功"})
}

// simulateAgentDiscussion 模拟智能体讨论
func (cc *ChatRoomController) simulateAgentDiscussion(chatRoomID string) {
	// 记录开始处理的聊天室ID
	log.Printf("开始处理聊天室讨论，ID: %s", chatRoomID)

	// 获取聊天室并预加载Agents数据
	var chatRoom model.ChatRoom
	if err := cc.db.Preload("Agents").First(&chatRoom, "id = ?", chatRoomID).Error; err != nil {
		log.Printf("获取聊天室失败，ID: %s, 错误: %v", chatRoomID, err)
		return
	}

	log.Printf("成功获取聊天室信息，ID: %s, 主题: %s", chatRoomID, chatRoom.Topic)

	// 立即广播聊天室完整信息，确保前端页面能显示主题、状态和智能体
	roomInfoUpdate := map[string]interface{}{
		"type":      "room_info",
		"chat_room": chatRoom,
	}
	roomInfoJSON, _ := json.Marshal(roomInfoUpdate)
	// 这里需要通过WebSocket广播消息，暂时打印日志
	log.Printf("广播聊天室信息: %s", string(roomInfoJSON))

	// 获取聊天室的最后一条用户消息作为讨论主题
	var lastUserMessage model.Message
	topic := ""
	if err := cc.db.Where("chat_room_id = ? AND type = ?", chatRoomID, "user").
		Order("timestamp desc").First(&lastUserMessage).Error; err == nil {
		// 获取用户输入的最后一条消息作为讨论主题
		topic = lastUserMessage.Content
	} else {
		// 如果没有用户消息，使用聊天室主题作为讨论主题
		topic = chatRoom.Topic

		// 如果聊天室主题也为空，使用默认主题
		if topic == "" {
			topic = "请各位智能体分享你们的观点和见解"
		}

		// 创建一条系统消息作为初始话题
		initialMessage := model.Message{
			ID:         uuid.New().String(),
			Content:    fmt.Sprintf("开始讨论主题：%s", topic),
			Type:       "system",
			Timestamp:  time.Now(),
			ChatRoomID: chatRoomID,
		}

		// 保存初始消息到数据库
		if err := cc.db.Create(&initialMessage).Error; err != nil {
			log.Printf("保存初始消息失败，聊天室ID: %s, 错误: %v", chatRoomID, err)
		}

		// 广播初始消息
		messageJSON, _ := json.Marshal(initialMessage)
		log.Printf("广播初始消息: %s", string(messageJSON))
	}

	// 更新状态为运行中
	chatRoom.Status = "running"
	if err := cc.db.Save(&chatRoom).Error; err != nil {
		log.Printf("更新聊天室状态失败，ID: %s, 错误: %v", chatRoomID, err)
		return
	}

	// 广播状态更新
	statusUpdate := map[string]interface{}{
		"type":   "status_update",
		"status": "running",
	}
	statusJSON, _ := json.Marshal(statusUpdate)
	log.Printf("广播状态更新: %s", string(statusJSON))

	// 获取聊天室关联的智能体
	var agents []model.Agent
	if err := cc.db.Model(&chatRoom).Association("Agents").Find(&agents); err != nil {
		log.Printf("获取智能体失败: %v", err)
		return
	}

	// 为每个智能体生成回复
	for _, agent := range agents {
		// 检查状态
		var currentChatRoom model.ChatRoom
		if err := cc.db.Preload("Agents").First(&currentChatRoom, "id = ?", chatRoomID).Error; err != nil {
			break
		}
		if currentChatRoom.Status != "running" {
			break
		}

		// 构建提示词
		prompt := fmt.Sprintf("你是一个%s，性格特点：%s。当前讨论话题：%s。请发表你的观点，保持简洁有力。",
			agent.Role, agent.Personality, topic)

		// 调用阿里云DashScope API
		log.Printf("调用阿里云DashScope API，智能体: %s, 提示词: %s", agent.Name, prompt)
		content, err := utils.CallDashScopeAPI(prompt)

		if err != nil {
			log.Printf("阿里云API调用失败: %v", err)
			content = fmt.Sprintf("%s暂时无法发言: %v", agent.Name, err)
		}

		// 创建消息
		message := model.Message{
			ID:         uuid.New().String(),
			AgentID:    agent.ID,
			AgentName:  agent.Name,
			Content:    content,
			Type:       "agent",
			Timestamp:  time.Now(),
			ChatRoomID: chatRoomID,
		}

		// 保存消息到数据库
		if err := cc.db.Create(&message).Error; err != nil {
			log.Printf("保存消息失败，聊天室ID: %s, 错误: %v", chatRoomID, err)
			continue
		}

		// 广播消息
		messageJSON, _ := json.Marshal(message)
		log.Printf("广播智能体消息: %s", string(messageJSON))

		// 延迟一下，模拟真实讨论节奏
		time.Sleep(2 * time.Second)
	}

	// 更新状态为完成
	chatRoom.Status = "completed"
	if err := cc.db.Save(&chatRoom).Error; err != nil {
		log.Printf("更新聊天室状态失败，ID: %s, 错误: %v", chatRoomID, err)
		return
	}

	// 广播完成状态
	statusUpdate = map[string]interface{}{
		"type":   "status_update",
		"status": "completed",
	}
	statusJSON, _ = json.Marshal(statusUpdate)
	log.Printf("广播完成状态: %s", string(statusJSON))
}
