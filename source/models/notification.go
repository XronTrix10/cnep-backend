package models

import (
	"time"
)

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Type      string    `gorm:"not null" json:"type"`
	Message   string    `gorm:"not null" json:"message"`
	EventLink string    `gorm:"not null" json:"event_link"`
	Read      bool      `gorm:"not null;default:false" json:"read"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
