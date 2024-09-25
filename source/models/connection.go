package models

import (
	"time"
)

type Connection struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	ReceiverID uint      `gorm:"not null" json:"receiver_id"`
	Status     string    `gorm:"not null;check:status IN ('pending', 'accepted', 'declined', 'deleted')" json:"status"`
	SentAt     time.Time `gorm:"default:current_timestamp" json:"sent_at"`
	UpdatedAt  time.Time `gorm:"default:current_timestamp" json:"updated_at"`
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"`
}
