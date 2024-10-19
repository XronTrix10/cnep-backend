package models

import (
	"time"
)

// ==== User Feedback ====

type Feedback struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `json:"content"`
	Rating     uint8     `gorm:"not null" json:"rating"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

type FeedbackSender struct {
	FeedbackID   uint      `json:"feedback_id"`
	Content      string    `json:"content"`
	Rating       uint8     `json:"rating"`
	SenderID     uint      `json:"sender_id"`
	SenderName   string    `json:"sender_name"`
	SenderEmail  string    `json:"sender_email"`
	SenderRating float32   `json:"sender_rating"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FeedbackWithSender struct {
	FeedbackID uint       `json:"feedback_id"`
	Content    string     `json:"content"`
	Rating     uint8      `json:"rating"`
	Sender     SenderInfo `json:"sender"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ==== User Badges ====

type Badge struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Image       string    `gorm:"not null" json:"image"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
