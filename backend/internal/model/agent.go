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
	CreatedAt   time.Time `json:"created_at"`
}
