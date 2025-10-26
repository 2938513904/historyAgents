package model

import (
	"time"
)

// ChatRoom 聊天室模型
type ChatRoom struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Topic     string    `json:"topic" gorm:"type:varchar(200)"`
	Status    string    `json:"status" gorm:"type:varchar(50);default:pending"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages" gorm:"foreignKey:ChatRoomID;references:ID"`
	Agents    []Agent   `json:"agents" gorm:"many2many:chat_room_agents;"`
}
