package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"multiagent-chat/internal/model"
	"multiagent-chat/internal/utils"
	"multiagent-chat/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ChatRoomController 聊天室控制器
type ChatRoomController struct {
	db         *gorm.DB
	wsManager  *websocket.Manager
}

// NewChatRoomController 创建聊天室控制器实例
func NewChatRoomController(db *gorm.DB, wsManager *websocket.Manager) *ChatRoomController {
	return &ChatRoomController{
		db:        db,
		wsManager: wsManager,
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
	log.Printf("接收到创建聊天室请求 - Topic: %s, Agents: %v", requestData.Topic, requestData.Agents)

	// 创建聊天室
	chatRoom := model.ChatRoom{
		ID:        uuid.New().String(),
		Topic:     requestData.Topic,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// 保存聊天室到数据库
	if err := cc.db.Create(&chatRoom).Error; err != nil {
		log.Printf("创建聊天室失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建聊天室失败"})
		return
	}
	
	log.Printf("聊天室创建成功 - ID: %s, Topic: %s", chatRoom.ID, chatRoom.Topic)

	// 关联智能体
	if len(requestData.Agents) > 0 {
		var agents []model.Agent
		if err := cc.db.Where("id IN ?", requestData.Agents).Find(&agents).Error; err != nil {
			log.Printf("查询智能体失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询智能体失败"})
			return
		}

		// 使用Association方法关联智能体
		if err := cc.db.Model(&chatRoom).Association("Agents").Append(agents); err != nil {
			log.Printf("关联智能体失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联智能体失败"})
			return
		}
		
		log.Printf("智能体关联成功，数量: %d", len(agents))
	}

	// 重新加载聊天室数据以包含关联的智能体
	if err := cc.db.Preload("Agents").First(&chatRoom, "id = ?", chatRoom.ID).Error; err != nil {
		log.Printf("重新加载聊天室数据失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "重新加载聊天室数据失败"})
		return
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
	// 通过WebSocket广播消息
	if cc.wsManager != nil {
		cc.wsManager.BroadcastToRoom(chatRoomID, roomInfoJSON)
	} else {
		log.Printf("WebSocket管理器未初始化，无法广播聊天室信息: %s", string(roomInfoJSON))
	}

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
		if cc.wsManager != nil {
			cc.wsManager.BroadcastToRoom(chatRoomID, messageJSON)
		} else {
			log.Printf("WebSocket管理器未初始化，无法广播初始消息: %s", string(messageJSON))
		}
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
	if cc.wsManager != nil {
		cc.wsManager.BroadcastToRoom(chatRoomID, statusJSON)
	} else {
		log.Printf("WebSocket管理器未初始化，无法广播状态更新: %s", string(statusJSON))
	}

	// 获取聊天室关联的智能体
	var agents []model.Agent
	if err := cc.db.Model(&chatRoom).Association("Agents").Find(&agents); err != nil {
		log.Printf("获取智能体失败: %v", err)
		return
	}

	// 确保至少进行3轮对话
	minRounds := 3
	maxRounds := 5 // 设置最大轮数避免无限循环
	
	log.Printf("开始多轮讨论，智能体数量: %d, 最少轮数: %d", len(agents), minRounds)

	for round := 1; round <= maxRounds; round++ {
		log.Printf("开始第 %d 轮讨论", round)
		
		// 为每个智能体生成回复
		for _, agent := range agents {
			// 检查状态
			var currentChatRoom model.ChatRoom
			if err := cc.db.Preload("Agents").First(&currentChatRoom, "id = ?", chatRoomID).Error; err != nil {
				log.Printf("获取聊天室状态失败: %v", err)
				return
			}
			if currentChatRoom.Status != "running" {
				log.Printf("聊天室状态已改变，停止讨论")
				return
			}

			// 获取之前的对话历史，用于构建更有针对性的提示词
			var previousMessages []model.Message
			cc.db.Where("chat_room_id = ? AND type = ?", chatRoomID, "agent").
				Order("timestamp desc"). // 改为降序，获取最新的消息
				Limit(6). // 获取最近6条消息作为上下文，避免上下文过长
				Find(&previousMessages)

			// 构建包含上下文的提示词
			var contextStr string
			if len(previousMessages) > 0 && round > 1 {
				contextStr = "\n\n最近的讨论内容："
				// 反转消息顺序，按时间正序显示
				for i := len(previousMessages) - 1; i >= 0; i-- {
					msg := previousMessages[i]
					// 避免智能体回应自己的话
					if msg.AgentName != agent.Name {
						contextStr += fmt.Sprintf("\n%s: %s", msg.AgentName, msg.Content)
					}
				}
				contextStr += "\n\n请基于以上讨论内容，"
			}

			// 根据轮数和智能体角色调整提示词，使对话更有针对性
			var roundPrompt string
			switch round {
			case 1:
				roundPrompt = fmt.Sprintf("作为%s，请从你的专业角度发表初步观点，保持简洁有力。", agent.Role)
			case 2:
				if len(previousMessages) > 0 {
					roundPrompt = "请进一步阐述你的观点，或对其他智能体的观点进行回应和补充。"
				} else {
					roundPrompt = "请进一步阐述你的观点，提供更多细节。"
				}
			case 3:
				roundPrompt = "请总结你的核心观点，或提出新的思考角度和建议。"
			default:
				roundPrompt = "请继续深入讨论，提供更多见解、反思或实用建议。"
			}

			// 构建更丰富的提示词，包含角色定位、性格特点、讨论上下文和轮次指导
			prompt := fmt.Sprintf(`你是一个%s，性格特点：%s。
当前讨论话题：%s
%s
%s

请注意：
1. 保持你的角色特色和专业性
2. 回应要有逻辑性和建设性
3. 避免重复之前已经说过的内容
4. 控制回应长度在100-200字之间`,
				agent.Role, agent.Personality, topic, contextStr, roundPrompt)

			// 调用阿里云DashScope API
			log.Printf("第%d轮 - 调用阿里云DashScope API，智能体: %s", round, agent.Name)
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
			if cc.wsManager != nil {
				cc.wsManager.BroadcastToRoom(chatRoomID, messageJSON)
			} else {
				log.Printf("WebSocket管理器未初始化，无法广播智能体消息: %s", string(messageJSON))
			}

			// 延迟一下，模拟真实讨论节奏
			time.Sleep(3 * time.Second)
		}

		log.Printf("第 %d 轮讨论完成", round)
		
		// 在轮次之间稍作停顿
		if round < maxRounds {
			time.Sleep(2 * time.Second)
		}
		
		// 如果已经完成最少轮数，可以考虑提前结束的条件
		// 这里可以根据需要添加更复杂的结束逻辑
		if round >= minRounds {
			// 简单的结束条件：如果智能体数量较少，多进行几轮
			if len(agents) <= 2 && round >= 4 {
				log.Printf("智能体数量较少，已完成 %d 轮讨论，结束讨论", round)
				break
			} else if len(agents) > 2 && round >= minRounds {
				log.Printf("已完成 %d 轮讨论，结束讨论", round)
				break
			}
		}
	}

	log.Printf("多轮讨论结束，聊天室ID: %s", chatRoomID)

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
	if cc.wsManager != nil {
		cc.wsManager.BroadcastToRoom(chatRoomID, statusJSON)
	} else {
		log.Printf("WebSocket管理器未初始化，无法广播完成状态: %s", string(statusJSON))
	}
}
