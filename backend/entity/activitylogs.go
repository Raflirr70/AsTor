package entity

import "time"

type ActivityLogs struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Action      string    `json:"action"`
	EntityType  string    `json:"entity_type"`
	EntityID    uint      `json:"entity_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	// ForeignKey
	User *Users `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
