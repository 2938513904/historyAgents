package model

import (
	"time"
)

// Message 聊天消息模型
type Message struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Type       string    `json:"type" gorm:"type:varchar(20)"`
	Content    string    `json:"content" gorm:"type:text"`
	AgentID    string    `json:"agent_id" gorm:"type:varchar(100)"`
	AgentName  string    `json:"agent_name" gorm:"type:varchar(100)"`
	Timestamp  time.Time `json:"timestamp"`
	ChatRoomID string    `json:"chat_room_id" gorm:"type:varchar(100);index"`
}
