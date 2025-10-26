package model

import (
	"time"
)

// Agent 智能体模型
type Agent struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100)"`
	Role        string    `json:"role" gorm:"type:varchar(200)"`
	Personality string    `json:"personality" gorm:"type:text"`
	Avatar      string    `json:"avatar" gorm:"type:longtext"` // 头像URL或Base64编码，支持大数据
	CozeLink    string    `json:"cozeLink" gorm:"type:varchar(500)"` // coze人物外链
	CreatedAt   time.Time `json:"created_at"`
}
