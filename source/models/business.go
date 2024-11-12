package models

import (
	"github.com/lib/pq"
	"time"
)

type BusinessPage struct {
	ID          uint          `gorm:"primaryKey" json:"id"`
	OwnerID     uint          `gorm:"not null" json:"owner_id"` // Reference to owner
	User        User          `gorm:"foreignKey:UserID" json:"user"`
	Name        string        `gorm:"not null" json:"name"`
	Badges      pq.Int64Array `gorm:"type:integer[]" json:"badges"`
	Topics      pq.Int64Array `gorm:"type:integer[]" json:"topics"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Logo        string        `json:"logo"`
	CoverImage  string        `json:"cover_image"`
	Contact     string        `json:"contact"`
	Website     string        `json:"website"`
	Location    string        `json:"location"`
	Rating      float32       `gorm:"default:0" json:"rating"`
	Followers   []User        `gorm:"many2many:page_followers;" json:"followers"`
	CreatedAt   time.Time     `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:current_timestamp" json:"updated_at"`
}

// Later we can add more fields to this struct
// for example - business_email, business_phone, business_address, social_media_links, etc.

type Product struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	BusinessPageID uint           `gorm:"not null" json:"business_page_id"`
	Name           string         `gorm:"not null" json:"name"`
	Description    string         `json:"description"`
	Price          float64        `json:"price"`
	Images         pq.StringArray `gorm:"type:text[]" json:"images"`
	Category       string         `json:"category"`
	CreatedAt      time.Time      `gorm:"default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"default:current_timestamp" json:"updated_at"`
}
