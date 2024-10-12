package models

import "time"

type Feedback struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"not null" json:"content"`
	Rating     uint8     `gorm:"not null" json:"rating"`
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
}
