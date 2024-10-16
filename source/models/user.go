package models

import (
	"github.com/lib/pq"
	"time"
)

// Original User struct
type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"unique;not null" json:"email"`
	Phone             string         `json:"phone"`
	Address           string         `json:"address"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	HelpedOthersCount int            `gorm:"default:0" json:"helped_others_count"`
	HelpReceivedCount int            `gorm:"default:0" json:"help_received_count"`
	Rating            float32        `gorm:"default:0" json:"rating"`
	Badges            pq.Int64Array  `gorm:"type:integer[]" json:"badges"`
	Designation       string         `json:"designation"`
	Password          string         `gorm:"not null" json:"password"`
	OTP               string         `gorm:"size:6" json:"otp"`
	OTPExpiry         time.Time      `json:"otp_expiry"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}

// Excluded sensitive fields from User
type UserResponse struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"unique;not null" json:"email"`
	Phone             string         `json:"phone"`
	Address           string         `json:"address"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	HelpedOthersCount int            `gorm:"default:0" json:"helped_others_count"`
	HelpReceivedCount int            `gorm:"default:0" json:"help_received_count"`
	Rating            float32        `gorm:"default:0" json:"rating"`
	Badges            pq.Int64Array  `gorm:"type:integer[]" json:"badges"`
	Designation       string         `json:"designation"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}

type Partner struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Status     string    `gorm:"not null;check:status IN ('pending', 'accepted', 'declined')" json:"status"`
	SentAt     time.Time `gorm:"default:current_timestamp" json:"sent_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

type Feedback struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"not null" json:"content"`
	Rating     uint8     `gorm:"not null" json:"rating"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

// Response struct for feedbacks
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

type SenderInfo struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Rating float32 `json:"rating"`
}

type FeedbackWithSender struct {
	FeedbackID uint       `json:"feedback_id"`
	Content    string     `json:"content"`
	Rating     uint8      `json:"rating"`
	Sender     SenderInfo `json:"sender"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Badge struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Image       string    `gorm:"not null" json:"image"`
	CreatedAt   time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

// TODO: Extract helps in another table from users table
