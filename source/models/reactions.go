package models

import (
	"time"
)

type Reaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Reaction  string    `gorm:"not null;check:reaction IN ('like', 'wow', 'love', 'angry', 'sad')" json:"reaction"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
