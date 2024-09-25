package models

import (
	"time"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:current_timestamp" json:"updated_at"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}
