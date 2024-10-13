package models

import (
	"time"
)

type Post struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	IsRequest  bool      `gorm:"not null" json:"is_request"`
	IsUrgent   bool      `gorm:"not null" json:"is_urgent"`
	Status     string    `gorm:"not null;check:status IN ('accepted', 'completed', 'pending')" json:"status"`
	AssignedTo uint      `json:"assigned_to"`
	MediaURL   string    `json:"media_url"`
	MediaType  string    `json:"media_type"`
	Caption    string    `gorm:"not null" json:"caption"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

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

type Reaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null" json:"post_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Reaction  string    `gorm:"not null;check:reaction IN ('like', 'wow', 'love', 'angry', 'sad')" json:"reaction"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"created_at"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}