package models

import (
	"time"
)

type Post struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	IsRequest  bool      `gorm:"not null" json:"is_request"`
	IsUrgent   bool      `gorm:"not null" json:"is_urgent"`
	Status     string    `gorm:"not null;check:status IN ('accepted', 'completed', 'pending', 'urgent')" json:"status"`
	AssignedTo uint      `json:"assigned_to"`
	Content    string    `gorm:"not null" json:"content"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	AssignedUser User    `gorm:"foreignKey:AssignedTo" json:"assigned_user"`
}
