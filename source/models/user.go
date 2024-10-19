package models

import (
	"github.com/lib/pq"
	"time"
)

// Original User struct
type User struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `gorm:"not null" json:"name"`
	Username    string        `gorm:"unique;not null" json:"username"`
	Email       string        `gorm:"unique;not null" json:"email"`
	Password    string        `gorm:"not null" json:"password"`
	Address     string        `json:"address"`
	Designation string        `json:"designation"`
	Phone       string        `json:"phone"`
	Badges      pq.Int64Array `gorm:"type:integer[]" json:"badges"`
	Topics      pq.Int64Array `gorm:"type:integer[]" json:"topics"`
	OTP         string        `gorm:"size:6" json:"otp"`
	OTPExpiry   time.Time     `json:"otp_expiry"`
	IsVerified  bool          `gorm:"default:false" json:"is_verified"`
	Rating      float32       `gorm:"default:0" json:"rating"`
	CreatedAt   time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}

// Excluded sensitive fields from User
type UserResponse struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	Name        string        `gorm:"not null" json:"name"`
	Username    string        `gorm:"unique;not null" json:"username"`
	Email       string        `gorm:"unique;not null" json:"email"`
	Address     string        `json:"address"`
	Designation string        `json:"designation"`
	Phone       string        `json:"phone"`
	Badges      pq.Int64Array `gorm:"type:integer[]" json:"badges"`
	Topics      pq.Int64Array `gorm:"type:integer[]" json:"topics"`
	IsVerified  bool          `gorm:"default:false" json:"is_verified"`
	Rating      float32       `gorm:"default:0" json:"rating"`
	CreatedAt   time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}

type Partner struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Status     string    `gorm:"not null;check:status IN ('pending', 'accepted', 'declined')" json:"status"`
	SentAt     time.Time `gorm:"default:current_timestamp" json:"sent_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}

type SenderInfo struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Rating float32 `json:"rating"`
}

