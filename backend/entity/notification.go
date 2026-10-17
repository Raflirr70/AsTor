package entity

import "time"

type Notifications struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Type       string    `json:"type"`
	IsRead     bool      `json:"is_read"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	EntityType string    `json:"entity_type"`
	EntityID   uint      `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`

	// ForeignKey
	User *Users `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
