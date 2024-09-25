package models

import (
	"time"
)

type Conversation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	User1ID   uint      `gorm:"not null" json:"user1_id"`
	User2ID   uint      `gorm:"not null" json:"user2_id"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
	User1     User      `gorm:"foreignKey:User1ID" json:"user1"`
	User2     User      `gorm:"foreignKey:User2ID" json:"user2"`
}
