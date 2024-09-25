package models

import (
	"time"
)

type Message struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ConversationID uint         `gorm:"not null" json:"conversation_id"`
	SenderID       uint         `gorm:"not null" json:"sender_id"`
	Content        string       `gorm:"not null" json:"content"`
	CreatedAt      time.Time    `gorm:"default:current_timestamp" json:"created_at"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
	Sender         User         `gorm:"foreignKey:SenderID" json:"sender"`
}
