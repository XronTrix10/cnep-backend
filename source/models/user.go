package models

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"unique;not null" json:"email"`
	Phone             string         `json:"phone"`
	Address           string         `json:"address"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	HelpedOthersCount int            `gorm:"default:0" json:"helped_others_count"`
	HelpReceivedCount int            `gorm:"default:0" json:"help_received_count"`
	Rating            int            `gorm:"default:0" json:"rating"`
	Badges            pq.Int64Array  `gorm:"type:integer[]" json:"badges"`
	Designation       string         `json:"designation"`
	Password          string         `gorm:"not null" json:"password"`
	OTP               string         `gorm:"size:6" json:"otp"`
	OTPExpiry         time.Time      `json:"otp_expiry"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}

type UserResponse struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Name              string         `gorm:"not null" json:"name"`
	Email             string         `gorm:"unique;not null" json:"email"`
	Phone             string         `json:"phone"`
	Address           string         `json:"address"`
	Skills            pq.StringArray `gorm:"type:text[]" json:"skills"`
	HelpedOthersCount int            `gorm:"default:0" json:"helped_others_count"`
	HelpReceivedCount int            `gorm:"default:0" json:"help_received_count"`
	Rating            int            `gorm:"default:0" json:"rating"`
	Badges            pq.Int64Array  `gorm:"type:integer[]" json:"badges"`
	Designation       string         `json:"designation"`
	IsVerified        bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt         time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}